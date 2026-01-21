package handler

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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

	name := strings.TrimSpace(req.Name)
	if name == "" {
		http.Error(w, "employee name is required", http.StatusBadRequest)
		return
	}

	emp, err := h.cfg.EmpService.GetEmployeeByName(r.Context(), name)
	if err == nil && emp != nil {
		if emp.PresentAt != nil {
			http.Error(w, "employee already marked as present", http.StatusBadRequest)
			return
		}

		if err := h.cfg.EmpService.UpdateEmployeePresentAt(r.Context(), emp.ID, time.Now()); err != nil {
			http.Error(w, "failed to mark employee as present", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	guest, guestErr := h.cfg.GuestService.GetGuestByName(r.Context(), name)
	if guestErr != nil {
		http.Error(w, "failed to get guest", http.StatusInternalServerError)
		return
	}

	if guest != nil && guest.PresentAt != nil {
		http.Error(w, "guest already marked as present", http.StatusBadRequest)
		return
	}

	if err := h.cfg.GuestService.AddGuest(r.Context(), name, time.Now()); err != nil {
		http.Error(w, "failed to mark guest as present", http.StatusInternalServerError)
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

	guests, err := h.cfg.GuestService.ListGuests(r.Context())
	if err != nil {
		http.Error(w, "failed to get guests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=attendance.csv")

	writer := csv.NewWriter(w)
	if err := writer.Write([]string{
		"person_type",
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
			"employee",
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

	for _, guest := range guests {
		presentAt := ""
		if guest.PresentAt != nil {
			presentAt = guest.PresentAt.Format(time.RFC3339)
		}
		record := []string{
			"guest",
			strconv.FormatInt(guest.ID, 10),
			guest.NamaTamu,
			"",
			"",
			"",
			presentAt,
			"Hadir",
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

func (h *Handler) ResetAllAttendances(w http.ResponseWriter, r *http.Request) {
	if err := h.cfg.EmpService.ResetAllAttendances(r.Context()); err != nil {
		http.Error(w, "failed to reset attendance", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
