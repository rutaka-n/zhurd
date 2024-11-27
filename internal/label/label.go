package label

import (
	"fmt"
)

type Label struct {
    ID           int64 `json:"id"`
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
