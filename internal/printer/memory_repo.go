package printer

import (
	"encoding/json"
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("record not found")
)

type Memory struct {
	m      map[int64][]byte
	nextID int64
	mu     sync.RWMutex
}

func NewMemory() (*Memory, error) {
	return &Memory{
		m:      make(map[int64][]byte),
		nextID: 1,
		mu:     sync.RWMutex{},
	}, nil
}

func (m *Memory) Store(p *Printer) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	if p.ID != 0 {
		m.m[p.ID] = data
		return nil
	}

	if _, ok := m.m[m.nextID]; ok {
		panic("could not generate unique ID for printer")
	}
	p.ID = m.nextID
	m.nextID += 1
	m.m[p.ID] = data

	return nil
}

func (m *Memory) Get(id int64) (Printer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data, ok := m.m[id]
	if !ok {
		return Printer{}, ErrNotFound
	}
	var p Printer
	if err := json.Unmarshal(data, &p); err != nil {
		return Printer{}, err
	}
	return p, nil
}
