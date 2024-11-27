package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"zhurd/internal/printer"
)

func listPrintersHandler(svc printer.QuerySvc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		printers, err := svc.List()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(printers)
	}
}

func showPrinterByIDHandler(svc printer.QuerySvc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		printerID, err := getPrinterID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pr, err := svc.Get(printerID)
		if err != nil {
			if errors.Is(err, printer.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(pr)
	}
}

func createPrinterHandler(svc printer.CommandSvc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var cp printer.CreatePrinter
		if err := json.NewDecoder(r.Body).Decode(&cp); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pr, err := svc.Create(cp)
		if err != nil {
			if errors.Is(err, printer.ValidationError) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(pr)
	}
}

func deletePrinterByIDHandler(svc printer.CommandSvc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		printerID, err := getPrinterID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := svc.Delete(printerID); err != nil {
			if errors.Is(err, printer.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func getPrinterID(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	val := vars["printerID"]
	return strconv.ParseInt(val, 10, 64)
}
