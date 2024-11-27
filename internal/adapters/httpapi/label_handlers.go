package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"zhurd/internal/label"
)

func listLabelsHandler(svc label.QuerySvc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		labels, err := svc.ListLabels()
		if err != nil {
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
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pr, err := svc.GetLabel(labelID)
		if err != nil {
			if errors.Is(err, label.ErrNotFound) {
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

func createLabelHandler(svc label.CommandSvc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var cp label.CreateLabel
		if err := json.NewDecoder(r.Body).Decode(&cp); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pr, err := svc.CreateLabel(cp)
		if err != nil {
			if errors.Is(err, label.ValidationError) {
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
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := svc.DeleteLabel(labelID); err != nil {
			if errors.Is(err, label.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
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
