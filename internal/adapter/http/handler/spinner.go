package handler

import "net/http"

func (h *Handler) RenderSpinnerPage(w http.ResponseWriter, r *http.Request) {
	if err := h.templates.ExecuteTemplate(w, "spinner.gohtml", nil); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
		return
	}
}
