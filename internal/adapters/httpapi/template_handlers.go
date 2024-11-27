package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"zhurd/internal/label"
)

func listTemplatesHandler(svc label.QuerySvc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		labelID, err := getLabelID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		templates, err := svc.ListTemplates(labelID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(templates)
	}
}

func showTemplateByIDHandler(svc label.QuerySvc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		templateID, err := getTemplateID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		labelID, err := getLabelID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pr, err := svc.GetTemplate(labelID, templateID)
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

func createTemplateHandler(svc label.CommandSvc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var cp label.CreateTemplate
		if err := json.NewDecoder(r.Body).Decode(&cp); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		labelID, err := getLabelID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		cp.LabelID = labelID

		pr, err := svc.CreateTemplate(cp)
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

func deleteTemplateByIDHandler(svc label.CommandSvc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		labelID, err := getLabelID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		templateID, err := getTemplateID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := svc.DeleteTemplate(labelID, templateID); err != nil {
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

func getTemplateID(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	val := vars["templateID"]
	return strconv.ParseInt(val, 10, 64)
}
