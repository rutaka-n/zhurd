package printingqueue

import (
	"context"
	"log/slog"
	"time"
)

type Printable interface {
	Print(pType string) ([]byte, error)
}

type Printer interface {
	Connect() error
	Close() error
	IsConnected() bool
	Enqueue(label Printable) error
}

type Task struct {
	Quantity int
	Timeout  time.Duration
	Document Printable
}

type Queue struct {
	PrinterID int64
	printer   Printer
	q         chan Task
	cancel    context.CancelFunc
}

func New(printerID int64, printer Printer, size int) Queue {
	return Queue{
		PrinterID: printerID,
		printer:   printer,
		q:         make(chan Task, size),
		cancel:    func() {}, // noop cancel func
	}
}

func (q *Queue) Enqueue(task Task) {
	q.q <- task
}

func (q *Queue) Close() error {
	defer close(q.q)
	defer q.cancel()
	return q.printer.Close()
}

func (q *Queue) Process(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	q.cancel = cancel
	if err := q.printer.Connect(); err != nil {
		return err
	}
	for {
		select {
		case task := <-q.q:
			for i := 0; i < task.Quantity; i++ {
				err := q.printer.Enqueue(task.Document)
				if err != nil {
					slog.Error("queue: printing failed", "printerID", q.printerID, "error", err)
				}
				time.Sleep(task.Timeout)
			}
		case <-ctx.Done():
			slog.Debug("processing queue for printer is done", "printerID", q.PrinterID)
			return nil
		}
	}
}
