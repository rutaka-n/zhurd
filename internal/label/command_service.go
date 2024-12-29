package label

import (
	"context"
	"errors"
	"fmt"
	"time"

	"zhurd/internal/printer"

	"github.com/go-playground/validator/v10"
)

var ValidationError = errors.New("Validation error")

type StorerDeleter interface {
	StoreLabel(context.Context, *Label) error
	DeleteLabel(context.Context, int64) error
	GetLabel(context.Context, int64) (Label, error)
	StoreTemplate(context.Context, *Template) error
	DeleteTemplate(context.Context, int64, int64) error
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

func (svc CommandSvc) CreateLabel(ctx context.Context, cl CreateLabel) (Label, error) {
	if err := svc.validate.Struct(cl); err != nil {
		return Label{}, fmt.Errorf("%w: %w", ValidationError, err)
	}
	l := Label{
		Name:    cl.Name,
		Comment: cl.Comment,
	}

	if err := svc.db.StoreLabel(ctx, &l); err != nil {
		return Label{}, err
	}

	return l, nil
}

func (svc CommandSvc) DeleteLabel(ctx context.Context, labelID int64) error {
	return svc.db.DeleteLabel(ctx, labelID)
}

func (svc CommandSvc) CreateTemplate(ctx context.Context, ct CreateTemplate) (Template, error) {
	if err := svc.validate.Struct(ct); err != nil {
		return Template{}, fmt.Errorf("%w: %w", ValidationError, err)
	}
	if _, err := svc.db.GetLabel(ctx, ct.LabelID); err != nil {
		return Template{}, err
	}
	t, err := NewTemplate(ct.LabelID, ct.Type, ct.Body)
	if err != nil {
		return Template{}, fmt.Errorf("%w: %w", ValidationError, err)
	}

	if err := svc.db.StoreTemplate(ctx, &t); err != nil {
		return Template{}, err
	}

	return t, nil
}

func (svc CommandSvc) DeleteTemplate(ctx context.Context, labelID, templateID int64) error {
	return svc.db.DeleteTemplate(ctx, labelID, templateID)
}

func (svc CommandSvc) Enqueue(ctx context.Context, labelID int64, enqueueLabel EnqueueLabel) error {
	label, err := svc.db.GetLabel(ctx, labelID)
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
