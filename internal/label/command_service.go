package label

import (
	"errors"
	"fmt"
	"time"

	"zhurd/internal/printer"

	"github.com/go-playground/validator/v10"
)

var ValidationError = errors.New("Validation error")

type StorerDeleter interface {
	StoreLabel(*Label) error
	DeleteLabel(int64) error
	GetLabel(int64) (Label, error)
	StoreTemplate(*Template) error
	DeleteTemplate(int64, int64) error
}

type Queue interface {
	Enqueue(printerID int64, document printer.Printable, quantity int, timeout time.Duration)
}

type CreateLabel struct {
	Name    string `json:"name" validate:"required"`
	Comment string `json:"comment"`
}

type CreateTemplate struct {
	LabelID int64
	Type    string `json:"type" validate:"required"`
	Body    []byte `json:"body" validate:"required"`
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

func (svc CommandSvc) CreateLabel(cl CreateLabel) (Label, error) {
	if err := svc.validate.Struct(cl); err != nil {
		return Label{}, fmt.Errorf("%w: %w", ValidationError, err)
	}
	l := Label{
		Name:    cl.Name,
		Comment: cl.Comment,
	}

	if err := svc.db.StoreLabel(&l); err != nil {
		return Label{}, err
	}

	return l, nil
}

func (svc CommandSvc) DeleteLabel(labelID int64) error {
	return svc.db.DeleteLabel(labelID)
}

func (svc CommandSvc) CreateTemplate(ct CreateTemplate) (Template, error) {
	if err := svc.validate.Struct(ct); err != nil {
		return Template{}, fmt.Errorf("%w: %w", ValidationError, err)
	}
	if _, err := svc.db.GetLabel(ct.LabelID); err != nil {
		return Template{}, err
	}
	t, err := NewTemplate(ct.LabelID, ct.Type, ct.Body)
	if err != nil {
		return Template{}, fmt.Errorf("%w: %w", ValidationError, err)
	}

	if err := svc.db.StoreTemplate(&t); err != nil {
		return Template{}, err
	}

	return t, nil
}

func (svc CommandSvc) DeleteTemplate(labelID, templateID int64) error {
	return svc.db.DeleteTemplate(labelID, templateID)
}

func (svc CommandSvc) Enqueue(labelID int64, enqueueLabel EnqueueLabel) error {
	label, err := svc.db.GetLabel(labelID)
	if err != nil {
		return err
	}
	label.placeholders = make(map[string]string, len(enqueueLabel.Placeholders))
	for _, ph := range enqueueLabel.Placeholders {
		label.placeholders[ph.Name] = ph.Value
	}
	svc.queue.Enqueue(enqueueLabel.PrinterID, label, enqueueLabel.Quantity, enqueueLabel.Timeout)
	return nil
}
