package label

import (
	"context"
	"encoding/base64"
)

type GetterLister interface {
	GetLabel(context.Context, int64) (Label, error)
	ListLabels(context.Context) ([]Label, error)
	GetTemplate(context.Context, int64, int64) (Template, error)
	ListTemplates(context.Context, int64) ([]Template, error)
}

type QuerySvc struct {
	db GetterLister
}

func NewQuerySvc(db GetterLister) QuerySvc {
	return QuerySvc{db: db}
}

func (svc QuerySvc) GetLabel(ctx context.Context, labelID int64) (Label, error) {
	return svc.db.GetLabel(ctx, labelID)
}

func (svc QuerySvc) ListLabels(ctx context.Context) ([]Label, error) {
	return svc.db.ListLabels(ctx)
}

func (svc QuerySvc) GetTemplate(ctx context.Context, labelID, templateID int64) (Template, error) {
	tmplt, err := svc.db.GetTemplate(ctx, labelID, templateID)
	if err != nil {
		return Template{}, err
	}
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(tmplt.Body)))
	base64.StdEncoding.Encode(encoded, tmplt.Body)
	tmplt.Body = encoded
	return tmplt, nil
}

func (svc QuerySvc) ListTemplates(ctx context.Context, labelID int64) ([]Template, error) {
	tmplts, err := svc.db.ListTemplates(ctx, labelID)
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
