package label

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var ValidationError = errors.New("Validation error")

type LabelStorerDeleter interface {
	StoreLabel(*Label) error
	DeleteLabel(int64) error
}

type CreateLabel struct {
	Name    string `json:"name" validate:"required"`
	Comment string `json:"comment"`
}

type CommandSvc struct {
	db       LabelStorerDeleter
	validate *validator.Validate
}

func NewCommandSvc(db LabelStorerDeleter) CommandSvc {
	return CommandSvc{
		db:       db,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (svc CommandSvc) CreateLabel(cl CreateLabel) (Label, error) {
	if err := svc.validate.Struct(cl); err != nil {
		return Label{}, fmt.Errorf("%w: %w", ValidationError, err)
	}
	l := Label{
        Name: cl.Name,
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
