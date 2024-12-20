package printer

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var ValidationError = errors.New("Validation error")

type StorerDeleter interface {
	Store(*Printer) error
	Delete(int64) error
}

type Queue interface {
	Add(printer Printer)
	Delete(id int64) error
}

type CreatePrinter struct {
	Addr    string `json:"addr" validate:"required,hostname_port"`
	Type    string `json:"type" validate:"required"`
	Comment string `json:"comment"`
}

type CommandSvc struct {
	db       StorerDeleter
	queue    Queue
	validate *validator.Validate
}

func NewCommandSvc(db StorerDeleter, queue Queue) CommandSvc {
	return CommandSvc{
		db:       db,
		queue:    queue,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (svc CommandSvc) Create(cp CreatePrinter) (Printer, error) {
	if err := svc.validate.Struct(cp); err != nil {
		return Printer{}, fmt.Errorf("%w: %w", ValidationError, err)
	}
	p := New(cp.Type, cp.Addr, cp.Comment)

	if err := svc.db.Store(&p); err != nil {
		return Printer{}, err
	}

	svc.queue.Add(p)
	return p, nil
}

func (svc CommandSvc) Delete(printerID int64) error {
	return svc.db.Delete(printerID)
}
