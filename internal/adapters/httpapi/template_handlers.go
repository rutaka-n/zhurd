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

func listTemplatesHandler(svc label.QuerySvc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		labelID, err := getLabelID(r)
		if err != nil {
			slog.Error("cannot get labelID", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		templates, err := svc.ListTemplates(r.Context(), labelID)
		if err != nil {
			slog.Error("cannot list templates", "error", err)
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
			slog.Error("cannot get templateID", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		labelID, err := getLabelID(r)
		if err != nil {
			slog.Error("cannot get labelID", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pr, err := svc.GetTemplate(r.Context(), labelID, templateID)
		if err != nil {
			if errors.Is(err, label.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			slog.Error("cannot get template", "error", err)
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

		var ct label.CreateTemplate
		if err := json.NewDecoder(r.Body).Decode(&ct); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			slog.Debug("cannot decode request body", "error", err)
			return
		}

		labelID, err := getLabelID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			slog.Debug("cannot find label ID", "error", err)
			return
		}
		ct.LabelID = labelID

		pr, err := svc.CreateTemplate(r.Context(), ct)
		if err != nil {
			if errors.Is(err, label.ValidationError) {
				slog.Error("template validation error", "error", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			slog.Error("cannot create template", "error", err)
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
			slog.Error("cannot get labelID", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		templateID, err := getTemplateID(r)
		if err != nil {
			slog.Error("cannot get templateID", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := svc.DeleteTemplate(r.Context(), labelID, templateID); err != nil {
			if errors.Is(err, label.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			slog.Error("cannot delete template", "error", err)
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
