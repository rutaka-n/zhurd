package printer

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var ValidationError = errors.New("Validation error")

type Storer interface {
	Store(*Printer) error
}

type CreatePrinter struct {
	Addr    string `json:"addr" validate:"required,hostname_port"`
	Type    string `json:"type" validate:"required"`
	Comment string `json:"comment"`
}

type CommandSvc struct {
	db       Storer
	validate *validator.Validate
}

func NewCommandSvc(db Storer) CommandSvc {
	return CommandSvc{
		db:       db,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (svc CommandSvc) Call(cp CreatePrinter) (Printer, error) {
	if err := svc.validate.Struct(cp); err != nil {
		return Printer{}, fmt.Errorf("%w: %w", ValidationError, err)
	}
	p := New(cp.Type, cp.Addr, cp.Comment)

	if err := svc.db.Store(&p); err != nil {
		return Printer{}, err
	}

	return p, nil
}
