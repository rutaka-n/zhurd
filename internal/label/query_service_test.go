package label

import (
	"errors"
	"slices"
	"testing"
)

func TestGetLabel(t *testing.T) {
	repo, err := NewMemory()
	if err != nil {
		t.Fatalf("got error: %s\n", err)
	}

	label := &Label{
		Name:    "label",
		Comment: "test label",
	}
	err = repo.StoreLabel(label)
	if err != nil {
		t.Fatalf("got error: %s\n", err)
	}

	ucs := []struct {
		desc        string
		labelID     int64
		label       *Label
		expectedErr error
	}{
		{
			desc:        "happy path",
			labelID:     label.ID,
			label:       label,
			expectedErr: nil,
		},
		{
			desc:        "wrong ID",
			labelID:     -1,
			label:       nil,
			expectedErr: ErrNotFound,
		},
		{
			desc:        "label with ID does not exist",
			labelID:     2,
			label:       nil,
			expectedErr: ErrNotFound,
		},
	}

	for _, us := range ucs {
		us := us
		t.Run(us.desc, func(t *testing.T) {
			svc := NewQuerySvc(repo)

			p, err := svc.GetLabel(us.labelID)
			if !errors.Is(err, us.expectedErr) {
				t.Errorf("expected: %v, got: %v\n", us.expectedErr, err)
			}
			if err == nil && us.label != nil {
				if p.Name != us.label.Name {
					t.Errorf("expected: %v, got: %v\n", us.label.Name, p.Name)
				}
				if p.Comment != us.label.Comment {
					t.Errorf("expected: %v, got: %v\n", us.label.Comment, p.Comment)
				}
			}
		})
	}
}

func TestListLabel(t *testing.T) {
	repo, err := NewMemory()
	if err != nil {
		t.Fatalf("got error: %s\n", err)
	}
	svc := NewQuerySvc(repo)

	// list empty storage
	result, err := svc.ListLabels()
	if err != nil {
		t.Fatalf("got error: %s\n", err)
	}
	if len(result) > 0 {
		t.Fatalf("expected empty result, but got %+v\n", result)
	}

	labels := []Label{
		{
			Name:    "label 1",
			Comment: "test label",
		},
		{
			Name:    "label 2",
			Comment: "another test label",
		},
	}
	for i := range labels {
		label := &labels[i]
		err = repo.StoreLabel(label)
		if err != nil {
			t.Fatalf("got error: %s\n", err)
		}
	}

	// list all labels in storage
	result, err = svc.ListLabels()
	if err != nil {
		t.Fatalf("got error: %s\n", err)
	}
	if len(result) != len(labels) {
		t.Fatalf("expected result has same length as stored labels: %d, but got %d\n", len(labels), len(result))
	}
}

func TestGetTemplate(t *testing.T) {
	repo, err := NewMemory()
	if err != nil {
		t.Fatalf("got error: %s\n", err)
	}
	label := &Label{
		Name: "label",
	}
	repo.StoreLabel(label)

	template := &Template{
		LabelID: label.ID,
		Type:    "ZPL-II",
		Body: []byte(`^XA
^FX Third section with bar code.
^BY5,2,270
^FO100,550^BC^FD12345678^FS
^XZ
`),
	}
	err = repo.StoreTemplate(template)
	if err != nil {
		t.Fatalf("got error: %s\n", err)
	}

	ucs := []struct {
		desc        string
		templateID  int64
		template    *Template
		expectedErr error
	}{
		{
			desc:        "happy path",
			templateID:  template.ID,
			template:    template,
			expectedErr: nil,
		},
		{
			desc:        "wrong ID",
			templateID:  -1,
			template:    nil,
			expectedErr: ErrNotFound,
		},
		{
			desc:        "template with ID does not exist",
			templateID:  2,
			template:    nil,
			expectedErr: ErrNotFound,
		},
	}

	for _, us := range ucs {
		us := us
		t.Run(us.desc, func(t *testing.T) {
			svc := NewQuerySvc(repo)

			tmplt, err := svc.GetTemplate(template.LabelID, us.templateID)
			if !errors.Is(err, us.expectedErr) {
				t.Errorf("expected: %v, got: %v\n", us.expectedErr, err)
			}
			if err == nil && us.template != nil {
				if tmplt.Type != us.template.Type {
					t.Errorf("expected: %s, got: %s\n", us.template.Type, tmplt.Type)
				}
				if !slices.Equal(tmplt.Body, us.template.Body) {
					t.Errorf("expected: %s, got: %s\n", us.template.Body, tmplt.Body)
				}
			}
		})
	}
}

func TestListTemplate(t *testing.T) {
	repo, err := NewMemory()
	if err != nil {
		t.Fatalf("got error: %s\n", err)
	}
	label := &Label{
		Name: "label",
	}
	repo.StoreLabel(label)
	svc := NewQuerySvc(repo)

	// list empty storage
	result, err := svc.ListTemplates(1)
	if err != nil {
		t.Fatalf("got error: %s\n", err)
	}
	if len(result) > 0 {
		t.Fatalf("expected empty result, but got %+v\n", result)
	}

	templates := []Template{
		{
			LabelID: label.ID,
			Type:    "ZPL",
			Body: []byte(`^XA
^FX Third section with bar code.
^BY5,2,270
^FO100,550^BC^FD12345678^FS
^XZ
`),
		},
		{
			LabelID: label.ID,
			Type:    "ZPL-II",
			Body: []byte(`^XA
^FX Third section with bar code.
^BY5,2,270
^FO100,550^BC^FD12345678^FS
^XZ
`),
		},
	}
	for i := range templates {
		template := &templates[i]
		err = repo.StoreTemplate(template)
		if err != nil {
			t.Fatalf("got error: %s\n", err)
		}
	}

	// list all templates in storage
	result, err = svc.ListTemplates(label.ID)
	if err != nil {
		t.Fatalf("got error: %s\n", err)
	}
	if len(result) != len(templates) {
		t.Fatalf("expected result has same length as stored templates: %d, but got %d\n", len(templates), len(result))
	}
}
