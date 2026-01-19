package repository

import (
	"context"
	"sync"

	"github.com/ibldzn/spinner-hut/internal/model"
)

type WinnerRepository struct {
	mu      sync.RWMutex
	winners []model.Winner
	index   map[string]struct{}
}

func NewWinnerRepository() *WinnerRepository {
	return &WinnerRepository{
		winners: make([]model.Winner, 0),
		index:   make(map[string]struct{}),
	}
}

func (r *WinnerRepository) AddWinners(ctx context.Context, winners []model.Winner) error {
	if len(winners) == 0 {
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, winner := range winners {
		key := winnerKey(winner)
		if _, exists := r.index[key]; exists {
			continue
		}
		r.index[key] = struct{}{}
		r.winners = append(r.winners, winner)
	}

	return nil
}

func (r *WinnerRepository) ListWinners(ctx context.Context) ([]model.Winner, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]model.Winner, len(r.winners))
	copy(out, r.winners)
	return out, nil
}

func (r *WinnerRepository) ListWinnersByType(ctx context.Context, prizeType string) ([]model.Winner, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filtered := make([]model.Winner, 0, len(r.winners))
	for _, winner := range r.winners {
		if winner.PrizeType == prizeType {
			filtered = append(filtered, winner)
		}
	}
	return filtered, nil
}

func winnerKey(winner model.Winner) string {
	return winner.EmployeeID + "|" + winner.RoundID + "|" + winner.PrizeType
}
