package printingqueue

import (
	"context"
	"log/slog"
	"sync"

	"zhurd/internal/printer"
)

type Pooler struct {
	bufferSize int
	wg         sync.WaitGroup
	queues     map[int64]*Queue
	addCh      chan *Queue
	deleteCh   chan int64
}

func NewPooler(bufferSize int) *Pooler {
	return &Pooler{
		bufferSize: bufferSize,
		queues:     map[int64]*Queue{},
		addCh:      make(chan *Queue),
		deleteCh:   make(chan int64),
	}
}

func (p *Pooler) Add(printer printer.Printer) {
	p.addCh <- New(printer, p.bufferSize)
}

func (p *Pooler) Delete(id int64) error {
	p.deleteCh <- id
	q, ok := p.queues[id]
	if !ok {
		slog.Warn("trying to delete queue that does not exist, ingored", "printerID", id)
		return nil
	}
	delete(p.queues, id)
	return q.Close()
}

func (p *Pooler) Run(ctx context.Context) {
	slog.Debug("pooler started")
	for {
		select {
		case q := <-p.addCh:
			slog.Debug("pooler: got new queue to run", "printerID", q.printer.ID)
			if _, ok := p.queues[q.printer.ID]; ok {
				slog.Warn("adding existing queue, ignored", "printerID", q.printer.ID)
				continue
			}
			p.queues[q.printer.ID] = q
			p.wg.Add(1)
			go func() {
				defer p.wg.Done()
				q.Process(ctx)
			}()
		case id := <-p.deleteCh:
			slog.Debug("pooler: got command to stop queue", "printerID", id)
			q, ok := p.queues[id]
			if !ok {
				slog.Warn("trying to delete queue that does not exist, ingored", "printerID", id)
				continue
			}
			delete(p.queues, id)
			err := q.Close()
			if err != nil {
				slog.Error("closing queue", "printerID", q.printer.ID, "error", err)
			}
		case <-ctx.Done():
			for _, q := range p.queues {
				err := q.Close()
				if err != nil {
					slog.Error("closing queue", "printerID", q.printer.ID, "error", err)
				}
			}
			slog.Debug("pooler is done")
			p.wg.Wait()
			return
		}
	}
}
