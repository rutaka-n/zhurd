package httpapi

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"zhurd/internal/label"
	"zhurd/internal/printer"
	pq "zhurd/internal/printingqueue"
)

func New(queue *pq.Pooler) (*mux.Router, error) {
	r := mux.NewRouter()
	r.HandleFunc("/", defaultHandler)
	v1r := r.PathPrefix("/v1").Subrouter()

	// printer
	printerRepo, err := printer.NewMemory()
	if err != nil {
		return nil, err
	}
	printerCommandSvc := printer.NewCommandSvc(printerRepo, queue)
	printerQuerySvc := printer.NewQuerySvc(printerRepo)

	v1r.HandleFunc("/printers", listPrintersHandler(printerQuerySvc)).Methods("GET")
	v1r.HandleFunc("/printers/{printerID}", showPrinterByIDHandler(printerQuerySvc)).Methods("GET")

	v1r.HandleFunc("/printers", createPrinterHandler(printerCommandSvc)).Methods("POST")
	v1r.HandleFunc("/printers/{printerID}", deletePrinterByIDHandler(printerCommandSvc)).Methods("DELETE")

	// label
	labelRepo, err := label.NewMemory()
	if err != nil {
		return nil, err
	}
	labelCommandSvc := label.NewCommandSvc(labelRepo, queue)
	labelQuerySvc := label.NewQuerySvc(labelRepo)

	v1r.HandleFunc("/labels", listLabelsHandler(labelQuerySvc)).Methods("GET")
	v1r.HandleFunc("/labels/{labelID}", showLabelByIDHandler(labelQuerySvc)).Methods("GET")

	v1r.HandleFunc("/labels", createLabelHandler(labelCommandSvc)).Methods("POST")
	v1r.HandleFunc("/labels/{labelID}", deleteLabelByIDHandler(labelCommandSvc)).Methods("DELETE")

	// template
	v1r.HandleFunc("/labels/{labelID}/templates", listTemplatesHandler(labelQuerySvc)).Methods("GET")
	v1r.HandleFunc("/labels/{labelID}/templates/{templateID}", showTemplateByIDHandler(labelQuerySvc)).Methods("GET")

	v1r.HandleFunc("/labels/{labelID}/templates", createTemplateHandler(labelCommandSvc)).Methods("POST")
	v1r.HandleFunc("/labels/{labelID}/templates/{templateID}", deleteTemplateByIDHandler(labelCommandSvc)).Methods("DELETE")

	// enque label
	v1r.HandleFunc("/labels/{labelID}/enqueue", enqueueLabelHandler(labelCommandSvc)).Methods("POST")

	r.Use(loggingMiddleware)

	return r, nil
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("error while reading request body", "error", err)
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		slog.Info("got request",
			"method", r.Method,
			"URI", r.RequestURI,
			"remote addr", r.RemoteAddr,
			"headers", r.Header,
			"body", body,
		)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
