package handler

import "net/http"

func (h *Handler) RenderInvitationPage(w http.ResponseWriter, r *http.Request) {
	guest := r.URL.Query().Get("name")
	dressCode := "Bebas Rapih"

	if emp, err := h.cfg.EmpService.GetEmployeeByName(r.Context(), guest); err == nil && emp != nil {
		dressCode = "Lihat Surat Edaran"
	}

	data := map[string]any{
		"DressCode": dressCode,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := h.templates.ExecuteTemplate(w, "invitation.gohtml", data); err != nil {
		http.Error(w, "failed to render invitation page", http.StatusInternalServerError)
		return
	}
}
