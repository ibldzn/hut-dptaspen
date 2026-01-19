package handler

import "net/http"

func (h *Handler) RenderAdminPage(w http.ResponseWriter, r *http.Request) {
	if err := h.templates.ExecuteTemplate(w, "admin.gohtml", nil); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
		return
	}
}
