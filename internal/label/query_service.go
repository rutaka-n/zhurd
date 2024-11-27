package label

type LabelGetterLister interface {
	GetLabel(int64) (Label, error)
	ListLabels() ([]Label, error)
}

type QuerySvc struct {
	db LabelGetterLister
}

func NewQuerySvc(db LabelGetterLister) QuerySvc {
	return QuerySvc{db: db}
}

func (svc QuerySvc) GetLabel(labelID int64) (Label, error) {
	return svc.db.GetLabel(labelID)
}

func (svc QuerySvc) ListLabels() ([]Label, error) {
	return svc.db.ListLabels()
}
