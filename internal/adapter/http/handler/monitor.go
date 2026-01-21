package handler

import "net/http"

func (h *Handler) RenderMonitorPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "monitor.gohtml", nil); err != nil {
		http.Error(w, "failed to render monitor page", http.StatusInternalServerError)
		return
	}
}
