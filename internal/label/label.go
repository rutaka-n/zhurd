package label

import (
	"fmt"
)

type Label struct {
    ID           int64 `json:"id"`
    Name         string `json:"name"`
    Comment      string `json:"comment"`
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
