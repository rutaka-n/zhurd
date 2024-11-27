package label

import (
	"errors"
	"testing"
)

func TestGet(t *testing.T) {
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
		labelID   int64
		label     *Label
		expectedErr error
	}{
		{
			desc:        "happy path",
			labelID:   label.ID,
			label:     label,
			expectedErr: nil,
		},
		{
			desc:        "wrong ID",
			labelID:   -1,
			label:     nil,
			expectedErr: ErrNotFound,
		},
		{
			desc:        "label with ID does not exist",
			labelID:   2,
			label:     nil,
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

func TestList(t *testing.T) {
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
