package label

import (
	"encoding/base64"
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
	encodedBody := make([]byte, base64.RawStdEncoding.EncodedLen(len(decodedBody)))
	base64.RawStdEncoding.Encode(encodedBody, decodedBody)
	ucs := []struct {
		desc        string
		ct          CreateTemplate
		expectedErr error
	}{
		{
			desc: "happy path",
			ct: CreateTemplate{
				Type: "ZPL",
				Body: encodedBody,
			},
			expectedErr: nil,
		},
		{
			desc: "empty type",
			ct: CreateTemplate{
				Type: "",
				Body: encodedBody,
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
				decoded := make([]byte, base64.RawStdEncoding.DecodedLen(len(us.ct.Body)))
				if _, err := base64.RawStdEncoding.Decode(decoded, us.ct.Body); err != nil {
					t.Fatalf("got error: %s\n", err)
				}
				if !slices.Equal(tmplt.Body, decoded) {
					t.Errorf("expected: %v, got: %v\n", decoded, tmplt.Body)
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
	encodedBody := make([]byte, base64.RawStdEncoding.EncodedLen(len(decodedBody)))
	base64.RawStdEncoding.Encode(encodedBody, decodedBody)
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
		Body: encodedBody,
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
