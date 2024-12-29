package printer

import (
	"context"
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

func (m *Memory) Store(ctx context.Context, p *Printer) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if p.ID == 0 {
		if _, ok := m.m[m.nextID]; ok {
			panic("could not generate unique ID for printer")
		}
		p.ID = m.nextID
		m.nextID += 1
	}

	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	m.m[p.ID] = data

	return nil
}

func (m *Memory) List(ctx context.Context) ([]Printer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	printers := make([]Printer, 0, len(m.m))

	for _, val := range m.m {
		var p Printer
		if err := json.Unmarshal(val, &p); err != nil {
			return nil, err
		}
		printers = append(printers, p)
	}
	return printers, nil
}

func (m *Memory) Get(ctx context.Context, id int64) (Printer, error) {
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

func (m *Memory) Delete(ctx context.Context, id int64) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if _, ok := m.m[id]; !ok {
		return ErrNotFound
	}
	delete(m.m, id)
	return nil
}
