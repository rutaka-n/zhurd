package label

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
	return svc.db.GetTemplate(labelID, templateID)
}

func (svc QuerySvc) ListTemplates(labelID int64) ([]Template, error) {
	return svc.db.ListTemplates(labelID)
}
