package label

import "encoding/base64"

type GetterLister interface {
	GetLabel(int64) (Label, error)
	ListLabels() ([]Label, error)
	GetTemplate(int64, int64) (Template, error)
	ListTemplates(int64) ([]Template, error)
}

type QuerySvc struct {
	db GetterLister
}

func NewQuerySvc(db GetterLister) QuerySvc {
	return QuerySvc{db: db}
}

func (svc QuerySvc) GetLabel(labelID int64) (Label, error) {
	return svc.db.GetLabel(labelID)
}

func (svc QuerySvc) ListLabels() ([]Label, error) {
	return svc.db.ListLabels()
}

func (svc QuerySvc) GetTemplate(labelID, templateID int64) (Template, error) {
	tmplt, err := svc.db.GetTemplate(labelID, templateID)
	if err != nil {
		return Template{}, err
	}
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(tmplt.Body)))
	base64.StdEncoding.Encode(encoded, tmplt.Body)
	tmplt.Body = encoded
	return tmplt, nil
}

func (svc QuerySvc) ListTemplates(labelID int64) ([]Template, error) {
	tmplts, err :=  svc.db.ListTemplates(labelID)
	if err != nil {
		return nil, err
	}
	for i := range tmplts {
		encoded := make([]byte, base64.StdEncoding.EncodedLen(len(tmplts[i].Body)))
		base64.StdEncoding.Encode(encoded, tmplts[i].Body)
		tmplts[i].Body = encoded
	}

	return tmplts, nil
}
