package label

import (
	"errors"
	"testing"
)

func TestRegister(t *testing.T) {
	ucs := []struct {
		desc        string
		cl          CreateLabel
		expectedErr error
	}{
		{
			desc: "happy path",
			cl: CreateLabel{
				Name:    "my label",
				Comment: "test label",
			},
			expectedErr: nil,
		},
		{
			desc: "empty name",
			cl: CreateLabel{
				Name:    "",
				Comment: "test label",
			},
			expectedErr: ValidationError,
		},
	}

	for _, us := range ucs {
		us := us
		t.Run(us.desc, func(t *testing.T) {
			repo, err := NewMemory()
			if err != nil {
				t.Fatalf("got error: %s\n", err)
			}
			svc := NewCommandSvc(repo)

			p, err := svc.CreateLabel(us.cl)
			if !errors.Is(err, us.expectedErr) {
				t.Errorf("expected: %v, got: %v\n", us.expectedErr, err)
			}
			if err == nil {
				if p.Name != us.cl.Name {
					t.Errorf("expected: %v, got: %v\n", us.cl.Name, p.Name)
				}
				if p.Comment != us.cl.Comment {
					t.Errorf("expected: %v, got: %v\n", us.cl.Comment, p.Comment)
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	repo, err := NewMemory()
	if err != nil {
		t.Fatalf("got error: %s\n", err)
	}
	svc := NewCommandSvc(repo)
	label := &Label{
		Name:    "new label",
		Comment: "test label",
	}

	if err := repo.StoreLabel(label); err != nil {
		t.Fatalf("got error: %s\n", err)
	}

	if err := svc.DeleteLabel(label.ID); err != nil {
		t.Fatalf("got error: %s\n", err)
	}

	_, err = repo.GetLabel(label.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected: %v, got: %v\n", ErrNotFound, err)
	}
}
