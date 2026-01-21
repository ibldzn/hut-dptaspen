package handler

import "net/http"

func (h *Handler) RenderScanPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "scan.gohtml", nil); err != nil {
		http.Error(w, "failed to render scan page", http.StatusInternalServerError)
		return
	}
}
