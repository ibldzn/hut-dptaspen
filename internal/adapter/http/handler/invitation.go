package handler

import "net/http"

func (h *Handler) RenderInvitationPage(w http.ResponseWriter, r *http.Request) {
	guest := r.URL.Query().Get("name")
	dressCode := "Bebas Rapih"
	additionalInfo := ""
	table := ""

	if emp, err := h.cfg.EmpService.GetEmployeeByName(r.Context(), guest); err == nil && emp != nil {
		dressCode = "Lihat Surat Edaran"
		additionalInfo = "Catatan: undangan berlaku untuk satu orang"
		if emp.Meja != nil {
			table = *emp.Meja
		}
	} else if g, err := h.cfg.GuestService.GetGuestByName(r.Context(), guest); err == nil && g != nil {
		if g.Meja != nil {
			table = *g.Meja
		}
	}

	data := map[string]any{
		"DressCode":      dressCode,
		"AdditionalInfo": additionalInfo,
		"Table":          table,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := h.templates.ExecuteTemplate(w, "invitation.gohtml", data); err != nil {
		http.Error(w, "failed to render invitation page", http.StatusInternalServerError)
		return
	}
}
