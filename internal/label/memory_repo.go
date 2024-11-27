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
