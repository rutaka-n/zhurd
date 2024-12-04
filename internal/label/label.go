package label

import (
	"fmt"
	"time"
)

type Label struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Comment      string `json:"comment"`
	templates    map[string]Template
	placeholders map[string]string
}

func (l Label) Print(pType string) ([]byte, error) {
	tplt, ok := l.templates[pType]
	if !ok {
		return nil, fmt.Errorf("label has no template with type: %s", pType)
	}
	return tplt.Print(l.placeholders)
}

type Placeholder struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type EnqueueLabel struct {
	PrinterID    int64         `json:"printer_id"`
	Quantity     int           `json:"quantity"`
	Timeout      time.Duration `json:"timeout"`
	Placeholders []Placeholder `json:"placeholders"`
}
