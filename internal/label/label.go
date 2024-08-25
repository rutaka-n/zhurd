package label

import (
	"fmt"
)

type Label struct {
	ID           int64
	Comment      string
	Templates    map[string]Template
	Placeholders map[string]string
}

func (l Label) Print(pType string) ([]byte, error) {
	tplt, ok := l.Templates[pType]
	if !ok {
		return nil, fmt.Errorf("label has no template with type: %s", pType)
	}
	return tplt.Print(l.Placeholders)
}
