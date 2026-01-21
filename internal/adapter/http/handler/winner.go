package handler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ibldzn/spinner-hut/internal/model"
)

type winnerPayload struct {
	RoundID    string        `json:"round_id"`
	RoundLabel string        `json:"round_label"`
	PrizeType  string        `json:"prize_type"`
	Winners    []winnerInput `json:"winners"`
}

type winnerInput struct {
	EmployeeID     string `json:"employee_id"`
	Name           string `json:"name"`
	Position       string `json:"position"`
	Branch         string `json:"branch"`
	EmploymentType string `json:"employment_type"`
}

func (h *Handler) AddWinners(w http.ResponseWriter, r *http.Request) {
	var payload winnerPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	prizeType := normalizePrizeType(payload.PrizeType)
	if prizeType == "" {
		http.Error(w, "invalid prize type", http.StatusBadRequest)
		return
	}

	if payload.RoundLabel == "" || payload.RoundID == "" {
		http.Error(w, "round metadata is required", http.StatusBadRequest)
		return
	}

	if len(payload.Winners) == 0 {
		http.Error(w, "winners list is empty", http.StatusBadRequest)
		return
	}

	now := time.Now()
	winners := make([]model.Winner, 0, len(payload.Winners))
	for _, entry := range payload.Winners {
		if strings.TrimSpace(entry.Name) == "" {
			http.Error(w, "winner name is required", http.StatusBadRequest)
			return
		}
		winners = append(winners, model.Winner{
			EmployeeID:     strings.TrimSpace(entry.EmployeeID),
			Name:           strings.TrimSpace(entry.Name),
			Position:       strings.TrimSpace(entry.Position),
			Branch:         strings.TrimSpace(entry.Branch),
			EmploymentType: strings.TrimSpace(entry.EmploymentType),
			PrizeType:      prizeType,
			RoundID:        strings.TrimSpace(payload.RoundID),
			RoundLabel:     strings.TrimSpace(payload.RoundLabel),
			WonAt:          now,
		})
	}

	if err := h.cfg.WinnerService.AddWinners(r.Context(), winners); err != nil {
		http.Error(w, "failed to store winners", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetWinners(w http.ResponseWriter, r *http.Request) {
	prizeType := normalizePrizeType(r.URL.Query().Get("type"))
	var (
		winners []model.Winner
		err     error
	)

	if prizeType == "" {
		winners, err = h.cfg.WinnerService.GetWinners(r.Context())
	} else {
		winners, err = h.cfg.WinnerService.GetWinnersByType(r.Context(), prizeType)
	}

	if err != nil {
		http.Error(w, "failed to get winners", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(winners); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ExportWinners(w http.ResponseWriter, r *http.Request) {
	rawType := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("type")))
	var (
		winners   []model.Winner
		err       error
		filename  string
		prizeType string
	)

	if rawType == "" || rawType == "all" {
		winners, err = h.cfg.WinnerService.GetWinners(r.Context())
		filename = "winners-all.csv"
	} else {
		prizeType = normalizePrizeType(rawType)
		if prizeType == "" {
			http.Error(w, "type must be door or grand", http.StatusBadRequest)
			return
		}
		winners, err = h.cfg.WinnerService.GetWinnersByType(r.Context(), prizeType)
		filename = fmt.Sprintf("winners-%s.csv", prizeType)
	}
	if err != nil {
		http.Error(w, "failed to get winners", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	writer := csv.NewWriter(w)
	if err := writer.Write([]string{
		"employee_id",
		"name",
		"position",
		"branch",
		"employment_type",
		"prize_type",
		"round_id",
		"round_label",
		"won_at",
	}); err != nil {
		http.Error(w, "failed to write csv", http.StatusInternalServerError)
		return
	}

	for _, winner := range winners {
		record := []string{
			winner.EmployeeID,
			winner.Name,
			winner.Position,
			winner.Branch,
			winner.EmploymentType,
			winner.PrizeType,
			winner.RoundID,
			winner.RoundLabel,
			winner.WonAt.Format(time.RFC3339),
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

func (h *Handler) ResetWinners(w http.ResponseWriter, r *http.Request) {
	if err := h.cfg.WinnerService.ResetWinners(r.Context()); err != nil {
		http.Error(w, "failed to reset winners", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func normalizePrizeType(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "door", "doorprize":
		return "door"
	case "grand", "grandprize":
		return "grand"
	default:
		return ""
	}
}
