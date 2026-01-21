package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

func (h *Handler) RenderPortalPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "portal.gohtml", nil); err != nil {
		http.Error(w, "failed to render portal page", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) LookupInvitation(w http.ResponseWriter, r *http.Request) {
	nip := strings.TrimSpace(r.URL.Query().Get("nip"))
	if nip == "" {
		http.Error(w, "nip is required", http.StatusBadRequest)
		return
	}

	employee, err := h.cfg.EmpService.GetEmployeeByNIP(r.Context(), nip)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "nip not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to lookup nip", http.StatusInternalServerError)
		return
	}
	if employee == nil {
		http.Error(w, "nip not found", http.StatusNotFound)
		return
	}

	response := struct {
		URL string `json:"url"`
	}{
		URL: "/?name=" + url.QueryEscape(employee.NamaKaryawan),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
