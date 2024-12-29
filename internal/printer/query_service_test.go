package printer

import (
	"context"
	"errors"
	"testing"
)

func TestGet(t *testing.T) {
	repo, err := NewMemory()
	if err != nil {
		t.Fatalf("got error: %s\n", err)
	}

	printer := &Printer{
		Addr:    "0.0.0.0:8009",
		Type:    "ZPL",
		Comment: "test printer",
	}
	err = repo.Store(context.Background(), printer)
	if err != nil {
		t.Fatalf("got error: %s\n", err)
	}

	ucs := []struct {
		desc        string
		printerID   int64
		printer     *Printer
		expectedErr error
	}{
		{
			desc:        "happy path",
			printerID:   printer.ID,
			printer:     printer,
			expectedErr: nil,
		},
		{
			desc:        "wrong ID",
			printerID:   -1,
			printer:     nil,
			expectedErr: ErrNotFound,
		},
		{
			desc:        "printer with ID does not exist",
			printerID:   2,
			printer:     nil,
			expectedErr: ErrNotFound,
		},
	}

	for _, us := range ucs {
		us := us
		t.Run(us.desc, func(t *testing.T) {
			svc := NewQuerySvc(repo)

			p, err := svc.Get(context.Background(), us.printerID)
			if !errors.Is(err, us.expectedErr) {
				t.Errorf("expected: %v, got: %v\n", us.expectedErr, err)
			}
			if err == nil && us.printer != nil {
				if p.Addr != us.printer.Addr {
					t.Errorf("expected: %v, got: %v\n", us.printer.Addr, p.Addr)
				}
				if p.Type != us.printer.Type {
					t.Errorf("expected: %v, got: %v\n", us.printer.Type, p.Type)
				}
				if p.Comment != us.printer.Comment {
					t.Errorf("expected: %v, got: %v\n", us.printer.Comment, p.Comment)
				}
			}
		})
	}
}

func TestList(t *testing.T) {
	repo, err := NewMemory()
	if err != nil {
		t.Fatalf("got error: %s\n", err)
	}
	svc := NewQuerySvc(repo)

	// list empty storage
	result, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("got error: %s\n", err)
	}
	if len(result) > 0 {
		t.Fatalf("expected empty result, but got %+v\n", result)
	}

	printers := []Printer{
		{
			Addr:    "0.0.0.0:8009",
			Type:    "ZPL",
			Comment: "test printer",
		},
		{
			Addr:    "0.0.0.0:8019",
			Type:    "ZPL-II",
			Comment: "another test printer",
		},
	}
	for i := range printers {
		printer := &printers[i]
		err = repo.Store(context.Background(), printer)
		if err != nil {
			t.Fatalf("got error: %s\n", err)
		}
	}

	// list all printers in storage
	result, err = svc.List(context.Background())
	if err != nil {
		t.Fatalf("got error: %s\n", err)
	}
	if len(result) != len(printers) {
		t.Fatalf("expected result has same length as stored printers: %d, but got %d\n", len(printers), len(result))
	}
}
