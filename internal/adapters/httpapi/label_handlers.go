package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"zhurd/internal/label"

	"github.com/gorilla/mux"
)

func listLabelsHandler(svc label.QuerySvc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		labels, err := svc.ListLabels(r.Context())
		if err != nil {
			slog.Error("cannot list labels", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(labels)
	}
}

func showLabelByIDHandler(svc label.QuerySvc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		labelID, err := getLabelID(r)
		if err != nil {
			slog.Error("cannot parse labelID", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pr, err := svc.GetLabel(r.Context(), labelID)
		if err != nil {
			if errors.Is(err, label.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			slog.Error("cannot get label", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(pr)
	}
}

func createLabelHandler(svc label.CommandSvc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var cp label.CreateLabel
		if err := json.NewDecoder(r.Body).Decode(&cp); err != nil {
			slog.Error("cannot parse request", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pr, err := svc.CreateLabel(r.Context(), cp)
		if err != nil {
			if errors.Is(err, label.ValidationError) {
				slog.Error("validation error", "error", err)
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

func deleteLabelByIDHandler(svc label.CommandSvc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		labelID, err := getLabelID(r)
		if err != nil {
			slog.Error("cannot get labelID", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := svc.DeleteLabel(r.Context(), labelID); err != nil {
			if errors.Is(err, label.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			slog.Error("cannot delete label", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func enqueueLabelHandler(svc label.CommandSvc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		labelID, err := getLabelID(r)
		if err != nil {
			slog.Error("cannot get labelID", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var enqueueLabel label.EnqueueLabel
		if err := json.NewDecoder(r.Body).Decode(&enqueueLabel); err != nil {
			slog.Error("cannot parse request", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := svc.Enqueue(r.Context(), labelID, enqueueLabel); err != nil {
			if errors.Is(err, label.ValidationError) {
				slog.Error("validation error", "error", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			slog.Error("cannot enqueue label", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func getLabelID(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	val := vars["labelID"]
	return strconv.ParseInt(val, 10, 64)
}
