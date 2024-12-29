package printingqueue

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"zhurd/internal/printer"
)

type Pooler struct {
	bufferSize int
	wg         sync.WaitGroup
	queues     map[int64]*Queue
	addCh      chan *Queue
	deleteCh   chan int64
	tasksCh    chan Task
}

func NewPooler(bufferSize int) *Pooler {
	return &Pooler{
		bufferSize: bufferSize,
		queues:     map[int64]*Queue{},
		addCh:      make(chan *Queue),
		deleteCh:   make(chan int64),
		tasksCh:    make(chan Task),
	}
}

func (p *Pooler) Add(printers ...printer.Printer) {
	for i := range printers {
		p.addCh <- New(printers[i], p.bufferSize)
	}
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

func (p *Pooler) Enqueue(printerID int64, document printer.Printable, quantity int, timeout time.Duration) {
	task := Task{
		printerID: printerID,
		Quantity:  quantity,
		Timeout:   timeout,
		Document:  document,
	}
	p.tasksCh <- task
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
		case task := <-p.tasksCh:
			slog.Debug("pooler: got task to enqueue", "task", task)
			q, ok := p.queues[task.printerID]
			if !ok {
				slog.Warn("trying to enqueue document for printer that are not exists, ignoring", "printerID", task.printerID)
				continue
			}
			if err := q.Enqueue(task); err != nil {
				slog.Warn("cannot enqueue task", "error", err)
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
