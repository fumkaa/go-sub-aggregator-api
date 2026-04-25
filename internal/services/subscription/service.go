package subscription

import (
	"context"
	"log/slog"

	"github.com/fumkaa/go-sub-aggregator-api/internal/domain/models"
	"github.com/google/uuid"
)

type SubscriptionReader interface {
	GetSubByID(ctx context.Context, id uuid.UUID) (models.Subscription, error)
	ListSubs(ctx context.Context, userId uuid.UUID) ([]models.Subscription, error)
}

type SubscriptionWriter interface {
	CreateSub(ctx context.Context, sub models.Subscription) (uuid.UUID, error)
	UpdateSub(ctx context.Context, sub models.Subscription) error
	DeleteSub(ctx context.Context, id uuid.UUID) error
}

type SubscriptionAggregator interface {
	GetSum(ctx context.Context, params models.ListSubsParams) (int, error)
}

type SubscriptionManager struct {
	log        *slog.Logger
	reader     SubscriptionReader
	writer     SubscriptionWriter
	aggregator SubscriptionAggregator
}

func NewSubscriptionManager(log *slog.Logger, reader SubscriptionReader, writer SubscriptionWriter, aggregator SubscriptionAggregator) *SubscriptionManager {
	return &SubscriptionManager{
		log:        log,
		reader:     reader,
		writer:     writer,
		aggregator: aggregator,
	}
}

func (m *SubscriptionManager) CreateSub(ctx context.Context, sub models.Subscription) (uuid.UUID, error) {
	return m.writer.CreateSub(ctx, sub)
}

func (m *SubscriptionManager) GetSubByID(ctx context.Context, id uuid.UUID) (models.Subscription, error) {
	return m.reader.GetSubByID(ctx, id)
}

func (m *SubscriptionManager) UpdateSub(ctx context.Context, sub models.Subscription) error {
	return m.writer.UpdateSub(ctx, sub)
}

func (m *SubscriptionManager) DeleteSub(ctx context.Context, id uuid.UUID) error {
	return m.writer.DeleteSub(ctx, id)
}

func (m *SubscriptionManager) ListSubs(ctx context.Context, userID uuid.UUID) ([]models.Subscription, error) {
	return m.reader.ListSubs(ctx, userID)
}

func (m *SubscriptionManager) GetSum(ctx context.Context, params models.ListSubsParams) (int, error) {
	return m.aggregator.GetSum(ctx, params)
}
