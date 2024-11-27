package printingqueue

import (
	"context"
    "sync"
    "log/slog"
)

type Pooler struct {
    bufferSize int
    wg sync.WaitGroup
    queues map[int64]Queue
    addCh chan Queue
    deleteCh chan int64

}

func NewPooler(bufferSize int) *Pooler {
	return &Pooler{
        bufferSize: bufferSize,
        queues: map[int64]Queue{},
    }
}

func (p *Pooler) Add(id int64, addr string) {
    p.addCh <- New(id, addr, p.bufferSize)
}

func (p *Pooler) Delete(id int64) error {
    p.deleteCh <- id
    q, ok := p.queues[id];
    if !ok {
        slog.Warn("trying to delete queue that does not exist, ingored", "printerID", id)
        return nil
    }
    delete(p.queues, id)
    return q.Close()
}

func (p *Pooler) Run(ctx context.Context) {
    for {
        select {
        case q := <-p.addCh:
            slog.Debug("pooler: got new queue to run", "printerID", q.PrinterID)
            if _, ok := p.queues[q.PrinterID]; ok {
                slog.Warn("adding existing queue, ignored", "printerID", q.PrinterID)
                continue
            }
            p.queues[q.PrinterID] = q
            p.wg.Add(1)
            go func () {
                defer p.wg.Done()
                q.Process(ctx)
            }()
        case id := <-p.deleteCh:
            slog.Debug("pooler: got command to stop queue", "printerID", id)
            q, ok := p.queues[id];
            if !ok {
                slog.Warn("trying to delete queue that does not exist, ingored", "printerID", id)
                continue
            }
            delete(p.queues, id)
            err := q.Close()
            if err != nil {
                slog.Error("closing queue", "printerID", q.PrinterID, "error", err)
            }
        case <-ctx.Done():
            slog.Debug("pooler is done")
            for _, q := range p.queues {
                err := q.Close()
                if err != nil {
                    slog.Error("closing queue", "printerID", q.PrinterID, "error", err)
                }
            }
            p.wg.Wait()
            return
        }
    }
}
