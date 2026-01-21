package handler

import "net/http"

func (h *Handler) RenderInvitationPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "invitation.gohtml", nil); err != nil {
		http.Error(w, "failed to render invitation page", http.StatusInternalServerError)
		return
	}
}
