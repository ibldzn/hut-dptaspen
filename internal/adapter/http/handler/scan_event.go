package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ibldzn/spinner-hut/internal/model"
)

type scanPayload struct {
	ScannerID int    `json:"scanner_id"`
	Name      string `json:"name"`
}

func (h *Handler) CreateScanEvent(w http.ResponseWriter, r *http.Request) {
	var payload scanPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(payload.Name)
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	if payload.ScannerID < 1 || payload.ScannerID > 3 {
		http.Error(w, "scanner_id must be between 1 and 3", http.StatusBadRequest)
		return
	}

	scannedAt := time.Now()
	created, err := h.ensureAttendance(r.Context(), name, scannedAt)
	if err != nil {
		http.Error(w, "failed to store attendance", http.StatusInternalServerError)
		return
	}
	if !created {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	event := model.ScanEvent{
		ScannerID: payload.ScannerID,
		Name:      name,
		ScannedAt: scannedAt,
	}
	if err := h.cfg.ScanService.AddScanEvent(r.Context(), event); err != nil {
		http.Error(w, "failed to store scan", http.StatusInternalServerError)
		return
	}

	message, err := json.Marshal(event)
	if err == nil {
		h.hub.Broadcast(message)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetRecentScans(w http.ResponseWriter, r *http.Request) {
	result := make(map[string][]model.ScanEvent, 3)
	for scannerID := 1; scannerID <= 3; scannerID++ {
		events, err := h.cfg.ScanService.ListRecentByScanner(r.Context(), scannerID, 3)
		if err != nil {
			http.Error(w, "failed to get scans", http.StatusInternalServerError)
			return
		}
		result[strconv.Itoa(scannerID)] = events
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ensureAttendance(ctx context.Context, name string, scannedAt time.Time) (bool, error) {
	emp, err := h.cfg.EmpService.GetEmployeeByName(ctx, name)
	if err == nil && emp != nil {
		if emp.PresentAt == nil {
			if err := h.cfg.EmpService.UpdateEmployeePresentAt(ctx, emp.ID, scannedAt); err != nil {
				return false, err
			}
			return true, nil
		}
		return false, nil
	}

	guest, guestErr := h.cfg.GuestService.GetGuestByName(ctx, name)
	if guestErr != nil {
		return false, guestErr
	}
	if guest != nil && guest.PresentAt != nil {
		return false, nil
	}
	if err := h.cfg.GuestService.AddGuest(ctx, name, scannedAt); err != nil {
		return false, err
	}
	return true, nil
}
