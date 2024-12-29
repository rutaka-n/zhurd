package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"zhurd/internal/printer"
)

func listPrintersHandler(svc printer.QuerySvc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		printers, err := svc.List(r.Context())
		if err != nil {
			slog.Error("cannot list printers", "error", err)
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
			slog.Error("cannot get printerID", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pr, err := svc.Get(r.Context(), printerID)
		if err != nil {
			if errors.Is(err, printer.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			slog.Error("cannot get printer", "error", err)
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
			slog.Error("cannot parse request", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pr, err := svc.Create(r.Context(), cp)
		if err != nil {
			if errors.Is(err, printer.ValidationError) {
				slog.Error("validation error", "error", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			slog.Error("cannot create printer", "error", err)
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

		if err := svc.Delete(r.Context(), printerID); err != nil {
			if errors.Is(err, printer.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			slog.Error("cannot delete printer", "error", err)
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
