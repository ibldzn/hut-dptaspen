package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

func (h *Handler) GetPresentEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := h.cfg.EmpService.GetPresentEmployees(r.Context())
	if err != nil {
		http.Error(w, "failed to get present employees", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) MarkEmployeePresent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	emp, err := h.cfg.EmpService.GetEmployeeByName(r.Context(), req.Name)
	if err != nil {
		http.Error(w, "failed to get employee", http.StatusInternalServerError)
		return
	}

	if emp == nil {
		http.Error(w, "employee not found", http.StatusNotFound)
		return
	}

	if emp.PresentAt != nil {
		http.Error(w, "employee already marked as present", http.StatusBadRequest)
		return
	}

	if err := h.cfg.EmpService.UpdateEmployeePresentAt(r.Context(), emp.ID, time.Now()); err != nil {
		http.Error(w, "failed to mark employee as present", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
