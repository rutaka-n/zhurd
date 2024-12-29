package printer

import "context"

type GetterLister interface {
	Get(context.Context, int64) (Printer, error)
	List(context.Context) ([]Printer, error)
}

type QuerySvc struct {
	db GetterLister
}

func NewQuerySvc(db GetterLister) QuerySvc {
	return QuerySvc{db: db}
}

func (svc QuerySvc) Get(ctx context.Context, printerID int64) (Printer, error) {
	return svc.db.Get(ctx, printerID)
}

func (svc QuerySvc) List(ctx context.Context) ([]Printer, error) {
	return svc.db.List(ctx)
}
