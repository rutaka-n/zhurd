package label

import (
	"errors"
	"slices"
	"testing"
	"time"
	"zhurd/internal/printer"
)

type TestQueue struct {
	enqueued int
}

func (q *TestQueue) Enqueue(printerID int64, document printer.Printable, quantity int, timeout time.Duration) {
	q.enqueued++
}

func TestRegisterLabel(t *testing.T) {
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
			q := &TestQueue{}
			svc := NewCommandSvc(repo, q)

			l, err := svc.CreateLabel(us.cl)
			if !errors.Is(err, us.expectedErr) {
				t.Errorf("expected: %v, got: %v\n", us.expectedErr, err)
			}
			if err == nil {
				if l.Name != us.cl.Name {
					t.Errorf("expected: %v, got: %v\n", us.cl.Name, l.Name)
				}
				if l.Comment != us.cl.Comment {
					t.Errorf("expected: %v, got: %v\n", us.cl.Comment, l.Comment)
				}
			}
		})
	}
}

func TestDeleteLabel(t *testing.T) {
	repo, err := NewMemory()
	if err != nil {
		t.Fatalf("got error: %s\n", err)
	}
	q := &TestQueue{}
	svc := NewCommandSvc(repo, q)
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

func TestRegisterTemplate(t *testing.T) {
	decodedBody := []byte(`^XA
^FX Third section with bar code.
^BY5,2,270
^FO100,550^BC^FD12345678^FS
^XZ
`)
	ucs := []struct {
		desc        string
		ct          CreateTemplate
		expectedErr error
	}{
		{
			desc: "happy path",
			ct: CreateTemplate{
				Type: "ZPL",
				Body: decodedBody,
			},
			expectedErr: nil,
		},
		{
			desc: "empty type",
			ct: CreateTemplate{
				Type: "",
				Body: decodedBody,
			},
			expectedErr: ValidationError,
		},
		{
			desc: "empty body",
			ct: CreateTemplate{
				Type: "ZPL",
				Body: nil,
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
			label := &Label{
				Name:    "new label",
				Comment: "test label",
			}

			if err := repo.StoreLabel(label); err != nil {
				t.Fatalf("got error: %s\n", err)
			}

			us.ct.LabelID = label.ID
			tmplt, err := svc.CreateTemplate(us.ct)
			if !errors.Is(err, us.expectedErr) {
				t.Errorf("expected: %v, got: %v\n", us.expectedErr, err)
			}
			if err == nil {
				if tmplt.Type != us.ct.Type {
					t.Errorf("expected: %v, got: %v\n", us.ct.Type, tmplt.Type)
				}
				if !slices.Equal(tmplt.Body, us.ct.Body) {
					t.Errorf("expected: %v, got: %v\n", us.ct.Body, tmplt.Body)
				}
			}
		})
	}
}

func TestDeleteTemplate(t *testing.T) {
	decodedBody := []byte(`^XA
^FX Third section with bar code.
^BY5,2,270
^FO100,550^BC^FD12345678^FS
^XZ
`)
	repo, err := NewMemory()
	if err != nil {
		t.Fatalf("got error: %s\n", err)
	}
	label := &Label{
		Name: "label",
	}
	repo.StoreLabel(label)
	q := &TestQueue{}
	svc := NewCommandSvc(repo, q)
	template := &Template{
		LabelID: label.ID,
		Type:    "ZPL",
		Body:    decodedBody,
	}

	if err := repo.StoreTemplate(template); err != nil {
		t.Fatalf("got error: %s\n", err)
	}

	if err := svc.DeleteTemplate(template.LabelID, template.ID); err != nil {
		t.Fatalf("got error: %s\n", err)
	}

	_, err = repo.GetTemplate(template.LabelID, template.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected: %v, got: %v\n", ErrNotFound, err)
	}
}

func TestEnqueue(t *testing.T) {
	decodedBody := []byte(`^XA
^FX Third section with bar code.
^BY5,2,270
^FO100,550^BC^FD12345678^FS
^XZ
`)
	repo, err := NewMemory()
	if err != nil {
		t.Fatalf("got error: %s\n", err)
	}
	label := &Label{
		Name: "label",
	}
	repo.StoreLabel(label)
	q := &TestQueue{}
	svc := NewCommandSvc(repo, q)
	template := &Template{
		LabelID: label.ID,
		Type:    "ZPL",
		Body:    decodedBody,
	}

	if err := repo.StoreTemplate(template); err != nil {
		t.Fatalf("got error: %s\n", err)
	}

	if q.enqueued != 0 {
		t.Fatalf("expect enqueued tasks: %d, got: %d\n", 0, q.enqueued)
	}

	enc := EnqueueLabel{
		PrinterID: 1,
		Quantity: 1,
		Timeout: time.Second,
		Placeholders: []Placeholder{},
	}

	svc.Enqueue(label.ID, enc)
	if q.enqueued != 1 {
		t.Fatalf("expect enqueued tasks: %d, got: %d\n", 1, q.enqueued)
	}
}
