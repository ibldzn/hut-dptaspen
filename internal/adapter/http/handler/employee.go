package handler

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := h.cfg.EmpService.GetAllEmployees(r.Context())
	if err != nil {
		http.Error(w, "failed to get employees", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

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

func (h *Handler) ExportAttendance(w http.ResponseWriter, r *http.Request) {
	employees, err := h.cfg.EmpService.GetAllEmployees(r.Context())
	if err != nil {
		http.Error(w, "failed to get employees", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=attendance.csv")

	writer := csv.NewWriter(w)
	if err := writer.Write([]string{
		"id",
		"name",
		"position",
		"branch_office",
		"employment_type",
		"present_at",
		"attendance_status",
	}); err != nil {
		http.Error(w, "failed to write csv", http.StatusInternalServerError)
		return
	}

	for _, emp := range employees {
		status := "Belum hadir"
		presentAt := ""
		if emp.PresentAt != nil {
			status = "Hadir"
			presentAt = emp.PresentAt.Format(time.RFC3339)
		}

		record := []string{
			strconv.FormatInt(emp.ID, 10),
			emp.NamaKaryawan,
			emp.Jabatan,
			string(emp.KantorCabang),
			string(emp.JenisKepegawaian),
			presentAt,
			status,
		}

		if err := writer.Write(record); err != nil {
			http.Error(w, "failed to write csv", http.StatusInternalServerError)
			return
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		http.Error(w, "failed to write csv", http.StatusInternalServerError)
		return
	}
}
