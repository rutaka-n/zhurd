package printingqueue

import (
	"context"
	"log/slog"
	"time"

	"zhurd/internal/printer"
)

type Printable interface {
	Print(pType string) ([]byte, error)
}

type Task struct {
	Quantity int
	Timeout  time.Duration
	Document Printable
}

type Queue struct {
	printer printer.Printer
	q       chan Task
	cancel  context.CancelFunc
}

func New(printer printer.Printer, size int) *Queue {
	return &Queue{
		printer: printer,
		q:       make(chan Task, size),
		cancel:  func() {}, // noop cancel func
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
	slog.Debug("start queue processing for printer", "printerID", q.printer.ID)
	ctx, cancel := context.WithCancel(ctx)
	q.cancel = cancel
	if err := q.printer.Connect(); err != nil {
		slog.Error("cannot connect to printer", "printerID", q.printer.ID, "addr", q.printer.Addr, "error", err)
	}
	for {
		select {
		case task := <-q.q:
			if !q.printer.IsConnected() {
				slog.Debug("printer is not connectied, try to connect", "printerID", q.printer.ID)
				if err := q.printer.Connect(); err != nil {
					slog.Error("cannot connect to printer", "printerID", q.printer.ID, "addr", q.printer.Addr, "error", err)
					continue
				}
			}
			for i := 0; i < task.Quantity; i++ {
				err := q.printer.Enqueue(task.Document)
				if err != nil {
					slog.Error("queue: printing failed", "printerID", q.printer.ID, "error", err)
				}
				time.Sleep(task.Timeout)
			}
		case <-ctx.Done():
			slog.Debug("processing queue for printer is done", "printerID", q.printer.ID)
			return nil
		}
	}
}
