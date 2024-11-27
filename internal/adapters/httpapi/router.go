package httpapi

import (
    "fmt"
	"net/http"

	"github.com/gorilla/mux"
	"zhurd/internal/printer"
)

func New() (*mux.Router, error) {
    r := mux.NewRouter()
    r.HandleFunc("/", defaultHandler)
    v1r := r.PathPrefix("/v1").Subrouter()

    printerRepo, err := printer.NewMemory()
    if err != nil {
        return nil, err
    }
    printerCommandSvc := printer.NewCommandSvc(printerRepo)
    printerQuerySvc := printer.NewQuerySvc(printerRepo)

    v1r.HandleFunc("/printers", listPrintersHandler(printerQuerySvc)).Methods("GET")
    v1r.HandleFunc("/printers/{printerID}", showPrinterByIDHandler(printerQuerySvc)).Methods("GET")

    v1r.HandleFunc("/printers", createPrinterHandler(printerCommandSvc)).Methods("POST")
    v1r.HandleFunc("/printers/{printerID}", deletePrinterByIDHandler(printerCommandSvc)).Methods("DELETE")

    r.Use(loggingMiddleware)

    return r, nil
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Println(r.RequestURI)
        // Call the next handler, which can be another middleware in the chain, or the final handler.
        next.ServeHTTP(w, r)
    })
}
