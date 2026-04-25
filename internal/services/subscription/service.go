package subscription

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fumkaa/go-sub-aggregator-api/internal/domain"
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
	const op = "subscription.SubscriptionManager.CreateSub()"
	sub.ID = uuid.New()

	id, err := m.writer.CreateSub(ctx, sub)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (m *SubscriptionManager) GetSubByID(ctx context.Context, id uuid.UUID) (models.Subscription, error) {
	const op = "subscription.SubscriptionManager.GetSubByID()"

	sub, err := m.reader.GetSubByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrSubNotFound) {
			return models.Subscription{}, err
		}
		return models.Subscription{}, fmt.Errorf("%s: %w", op, err)
	}
	return sub, nil
}

func (m *SubscriptionManager) UpdateSub(ctx context.Context, sub models.Subscription) error {
	const op = "subscription.SubscriptionManager.UpdateSub()"

	err := m.writer.UpdateSub(ctx, sub)
	if err != nil {
		if errors.Is(err, domain.ErrSubNotFound) {
			return err
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (m *SubscriptionManager) DeleteSub(ctx context.Context, id uuid.UUID) error {
	const op = "subscription.SubscriptionManager.DeleteSub()"

	err := m.writer.DeleteSub(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrSubNotFound) {
			return err
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (m *SubscriptionManager) ListSubs(ctx context.Context, userID uuid.UUID) ([]models.Subscription, error) {
	const op = "subscription.SubscriptionManager.ListSubs()"

	subs, err := m.reader.ListSubs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return subs, nil
}

func (m *SubscriptionManager) GetSum(ctx context.Context, params models.ListSubsParams) (int, error) {
	const op = "subscription.SubscriptionManager.GetSum()"

	sum, err := m.aggregator.GetSum(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return sum, nil
}
