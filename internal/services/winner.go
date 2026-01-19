package services

import (
	"context"

	"github.com/ibldzn/spinner-hut/internal/model"
	"github.com/ibldzn/spinner-hut/internal/repository"
)

type WinnerService struct {
	WinnerRepository *repository.WinnerRepository
}

func NewWinnerService(repo *repository.WinnerRepository) *WinnerService {
	return &WinnerService{
		WinnerRepository: repo,
	}
}

func (s *WinnerService) AddWinners(ctx context.Context, winners []model.Winner) error {
	return s.WinnerRepository.AddWinners(ctx, winners)
}

func (s *WinnerService) GetWinners(ctx context.Context) ([]model.Winner, error) {
	return s.WinnerRepository.ListWinners(ctx)
}

func (s *WinnerService) GetWinnersByType(ctx context.Context, prizeType string) ([]model.Winner, error) {
	return s.WinnerRepository.ListWinnersByType(ctx, prizeType)
}
