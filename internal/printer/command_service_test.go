package printer

import (
	"errors"
	"testing"
)

type TestQueue struct {
	added int
	deleted int
}
func (q *TestQueue) Add(printer Printer) {
	q.added++
}
func (q *TestQueue) Delete(id int64) error {
	q.deleted++
	return nil
}

func TestRegister(t *testing.T) {
	ucs := []struct {
		desc        string
		cp          CreatePrinter
		expectedErr error
	}{
		{
			desc: "happy path",
			cp: CreatePrinter{
				Addr:    "0.0.0.0:8009",
				Type:    "ZPL",
				Comment: "test printer",
			},
			expectedErr: nil,
		},
		{
			desc: "invalid port",
			cp: CreatePrinter{
				Addr:    "0.0.256.0:128009",
				Type:    "ZPL",
				Comment: "test printer",
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
			q := &TestQueue{}
			svc := NewCommandSvc(repo, q)

			p, err := svc.Create(us.cp)
			if !errors.Is(err, us.expectedErr) {
				t.Errorf("expected: %v, got: %v\n", us.expectedErr, err)
			}
			if err == nil {
				if p.Addr != us.cp.Addr {
					t.Errorf("expected: %v, got: %v\n", us.cp.Addr, p.Addr)
				}
				if p.Type != us.cp.Type {
					t.Errorf("expected: %v, got: %v\n", us.cp.Type, p.Type)
				}
				if p.Comment != us.cp.Comment {
					t.Errorf("expected: %v, got: %v\n", us.cp.Comment, p.Comment)
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
	q := &TestQueue{}
	svc := NewCommandSvc(repo, q)
	printer := &Printer{
		Addr: "0.0.0.0:8009",
		Type: "ZPL",
	}

	if err := repo.Store(printer); err != nil {
		t.Fatalf("got error: %s\n", err)
	}

	if err := svc.Delete(printer.ID); err != nil {
		t.Fatalf("got error: %s\n", err)
	}

    _, err = repo.Get(printer.ID)
    if !errors.Is(err, ErrNotFound) {
        t.Errorf("expected: %v, got: %v\n", ErrNotFound, err)
    }
}
