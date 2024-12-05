package label

import (
	"encoding/json"
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("record not found")
)

type Memory struct {
	labels         map[int64][]byte
	templates      map[int64][]byte
	nextLabelID    int64
	nextTemplateID int64
	mu             sync.RWMutex
}

func NewMemory() (*Memory, error) {
	return &Memory{
		labels:         make(map[int64][]byte),
		templates:      make(map[int64][]byte),
		nextLabelID:    1,
		nextTemplateID: 1,
		mu:             sync.RWMutex{},
	}, nil
}

func (m *Memory) StoreLabel(l *Label) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if l.ID == 0 {
		if _, ok := m.labels[m.nextLabelID]; ok {
			panic("could not generate unique ID for label")
		}
		l.ID = m.nextLabelID
		m.nextLabelID += 1
	}

	data, err := json.Marshal(l)
	if err != nil {
		return err
	}
	m.labels[l.ID] = data

	return nil
}

func (m *Memory) ListLabels() ([]Label, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	labels := make([]Label, 0, len(m.labels))

	for _, val := range m.labels {
		var l Label
		if err := json.Unmarshal(val, &l); err != nil {
			return nil, err
		}
		labels = append(labels, l)
	}
	return labels, nil
}

func (m *Memory) GetLabel(id int64) (Label, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data, ok := m.labels[id]
	if !ok {
		return Label{}, ErrNotFound
	}
	var l Label
	if err := json.Unmarshal(data, &l); err != nil {
		return Label{}, err
	}
	tmplts, err := m.ListTemplates(id)
	if err != nil {
		return Label{}, err
	}
	l.templates = make(map[string]Template, len(tmplts))
	for _, tmplt := range tmplts {
		l.templates[tmplt.Type] = tmplt
	}
	return l, nil
}

func (m *Memory) DeleteLabel(id int64) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if _, ok := m.labels[id]; !ok {
		return ErrNotFound
	}
	delete(m.labels, id)
	return nil
}

func (m *Memory) StoreTemplate(t *Template) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if t.ID == 0 {
		if _, ok := m.templates[m.nextTemplateID]; ok {
			panic("could not generate unique ID for template")
		}
		t.ID = m.nextTemplateID
		m.nextTemplateID += 1
	}

	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	m.templates[t.ID] = data

	return nil
}

func (m *Memory) ListTemplates(labelID int64) ([]Template, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	templates := make([]Template, 0, len(m.templates))

	for _, val := range m.templates {
		var t Template
		if err := json.Unmarshal(val, &t); err != nil {
			return nil, err
		}
		if t.LabelID == labelID {
			templates = append(templates, t)
		}
	}
	return templates, nil
}

func (m *Memory) GetTemplate(labelID, id int64) (Template, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data, ok := m.templates[id]
	if !ok {
		return Template{}, ErrNotFound
	}
	var t Template
	if err := json.Unmarshal(data, &t); err != nil {
		return Template{}, err
	}
	if t.LabelID != labelID {
		return Template{}, ErrNotFound
	}

	return t, nil
}

func (m *Memory) DeleteTemplate(labelID, templateID int64) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.templates[templateID]
	if !ok {
		return ErrNotFound
	}

	var t Template
	if err := json.Unmarshal(val, &t); err != nil {
		return err
	}
	if t.LabelID != labelID {
		return ErrNotFound
	}

	delete(m.templates, templateID)
	return nil
}
