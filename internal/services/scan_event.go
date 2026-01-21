package services

import (
	"context"

	"github.com/ibldzn/spinner-hut/internal/model"
	"github.com/ibldzn/spinner-hut/internal/repository"
)

type ScanEventService struct {
	Repo *repository.ScanEventRepository
}

func NewScanEventService(repo *repository.ScanEventRepository) *ScanEventService {
	return &ScanEventService{Repo: repo}
}

func (s *ScanEventService) AddScanEvent(ctx context.Context, event model.ScanEvent) error {
	return s.Repo.AddScanEvent(ctx, event)
}

func (s *ScanEventService) ListRecentByScanner(ctx context.Context, scannerID int, limit int) ([]model.ScanEvent, error) {
	return s.Repo.ListRecentByScanner(ctx, scannerID, limit)
}
