package services

import (
	"context"
	"time"

	"github.com/ibldzn/spinner-hut/internal/model"
	"github.com/ibldzn/spinner-hut/internal/repository"
)

type GuestService struct {
	GuestRepository *repository.GuestRepository
}

func NewGuestService(repo *repository.GuestRepository) *GuestService {
	return &GuestService{
		GuestRepository: repo,
	}
}

func (s *GuestService) ListGuests(ctx context.Context) ([]model.Guest, error) {
	return s.GuestRepository.ListGuests(ctx)
}

func (s *GuestService) GetGuestByName(ctx context.Context, name string) (*model.Guest, error) {
	return s.GuestRepository.GetGuestByName(ctx, name)
}

func (s *GuestService) AddGuest(ctx context.Context, name string, presentAt time.Time) error {
	return s.GuestRepository.AddGuest(ctx, name, presentAt)
}

func (s *GuestService) ResetGuests(ctx context.Context) error {
	return s.GuestRepository.ResetGuests(ctx)
}
