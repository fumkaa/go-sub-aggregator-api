package postgres

import (
	"context"

	"github.com/fumkaa/go-sub-aggregator-api/internal/domain/models"
	"github.com/google/uuid"
)

func (s *Storage) CreateSub(ctx context.Context, sub models.Subscription) (uuid.UUID, error) {
	panic("not implemented")
}

func (s *Storage) UpdateSub(ctx context.Context, sub models.Subscription) error {
	panic("not implemented")
}

func (s *Storage) DeleteSub(ctx context.Context, id uuid.UUID) error {
	panic("not implemented")
}

func (s *Storage) GetSubByID(ctx context.Context, id uuid.UUID) (models.Subscription, error) {
	panic("not implemented")
}

func (s *Storage) ListSubs(ctx context.Context, userId uuid.UUID) ([]models.Subscription, error) {
	panic("not implemented")
}

func (s *Storage) GetSum(ctx context.Context, params models.ListSubsParams) (int, error) {
	panic("not implemented")
}
