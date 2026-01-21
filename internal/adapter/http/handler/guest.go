package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type guestPayload struct {
	Name string `json:"name"`
}

func (h *Handler) GetGuests(w http.ResponseWriter, r *http.Request) {
	guests, err := h.cfg.GuestService.ListGuests(r.Context())
	if err != nil {
		http.Error(w, "failed to get guests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(guests); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) MarkGuestPresent(w http.ResponseWriter, r *http.Request) {
	var req guestPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		http.Error(w, "guest name is required", http.StatusBadRequest)
		return
	}

	existing, err := h.cfg.GuestService.GetGuestByName(r.Context(), name)
	if err != nil {
		http.Error(w, "failed to get guest", http.StatusInternalServerError)
		return
	}

	if existing != nil && existing.PresentAt != nil {
		http.Error(w, "guest already marked as present", http.StatusBadRequest)
		return
	}

	if err := h.cfg.GuestService.AddGuest(r.Context(), name, time.Now()); err != nil {
		http.Error(w, "failed to mark guest as present", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ResetGuests(w http.ResponseWriter, r *http.Request) {
	if err := h.cfg.GuestService.ResetGuests(r.Context()); err != nil {
		http.Error(w, "failed to reset guests", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
