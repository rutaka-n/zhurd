package printer

type GetterLister interface {
	Get(int64) (Printer, error)
	List() ([]Printer, error)
}

type QuerySvc struct {
	db GetterLister
}

func NewQuerySvc(db GetterLister) QuerySvc {
	return QuerySvc{db: db}
}

func (svc QuerySvc) Get(printerID int64) (Printer, error) {
	return svc.db.Get(printerID)
}

func (svc QuerySvc) List() ([]Printer, error) {
	return svc.db.List()
}
