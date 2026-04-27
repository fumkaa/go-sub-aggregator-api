package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/fumkaa/go-sub-aggregator-api/internal/domain"
	"github.com/fumkaa/go-sub-aggregator-api/internal/domain/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) CreateSub(ctx context.Context, sub models.Subscription) (uuid.UUID, error) {
	const op = "postgres.Storage.CreateSub()"

	query := `
		INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, TO_DATE($5, 'MM-YYYY'), TO_DATE($6, 'MM-YYYY'))
	`

	_, err := s.dbpool.Exec(ctx, query,
		sub.ID,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return sub.ID, nil
}

func (s *Storage) UpdateSub(ctx context.Context, sub models.Subscription) error {
	const op = "postgres.Storage.UpdateSub()"

	query := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, user_id = $3, start_date = TO_DATE($4, 'MM-YYYY'), end_date = TO_DATE($5, 'MM-YYYY')
		WHERE id = $6
	`

	res, err := s.dbpool.Exec(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res.RowsAffected() == 0 {
		return domain.ErrSubNotFound
	}

	return nil
}

func (s *Storage) DeleteSub(ctx context.Context, id uuid.UUID) error {
	const op = "postgres.Storage.DeleteSub()"

	query := `
		DELETE FROM subscriptions
		WHERE id = $1
	`

	res, err := s.dbpool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res.RowsAffected() == 0 {
		return domain.ErrSubNotFound
	}

	return nil
}

func (s *Storage) GetSubByID(ctx context.Context, id uuid.UUID) (models.Subscription, error) {
	const op = "postgres.Storage.GetSubByID()"

	query := `
		SELECT id, service_name, price, user_id,
			TO_CHAR(start_date, 'MM-YYYY'),
			TO_CHAR(end_date, 'MM-YYYY')
		FROM subscriptions
		WHERE id = $1
	`

	var sub models.Subscription
	err := s.dbpool.QueryRow(ctx, query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Subscription{}, domain.ErrSubNotFound
	} else if err != nil {
		return models.Subscription{}, fmt.Errorf("%s: %w", op, err)
	}

	return sub, nil
}

func (s *Storage) ListSubs(ctx context.Context, userId uuid.UUID, limit, offset int) ([]models.Subscription, error) {
	const op = "postgres.Storage.ListSubs()"

	query := `
		SELECT id, service_name, price, user_id,
			TO_CHAR(start_date, 'MM-YYYY'),
			TO_CHAR(end_date, 'MM-YYYY')
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY start_date DESC
        LIMIT $2 OFFSET $3
	`

	var subs []models.Subscription
	rows, err := s.dbpool.Query(ctx, query, userId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return subs, nil
}

func (s *Storage) GetSum(ctx context.Context, params models.ListSubsParams) (int, error) {
	const op = "postgres.Storage.GetSum()"

	query := `
		SELECT COALESCE(SUM(
            price * (
                EXTRACT(YEAR FROM age(
                    LEAST(end_date, COALESCE(TO_DATE($2, 'MM-YYYY'), end_date)), 
                    GREATEST(start_date, COALESCE(TO_DATE($3, 'MM-YYYY'), start_date))
                )) * 12 +
                EXTRACT(MONTH FROM age(
                    LEAST(end_date, COALESCE(TO_DATE($2, 'MM-YYYY'), end_date)), 
                    GREATEST(start_date, COALESCE(TO_DATE($3, 'MM-YYYY'), start_date))
                )) + 1
            )
        ), 0)::INTEGER
        FROM subscriptions
        WHERE user_id = $1
	`
	var pStart, pEnd *string
	if params.StartDate != "" {
		pStart = &params.StartDate
	}
	if params.EndDate != "" {
		pEnd = &params.EndDate
	}

	args := []any{params.UserID, pEnd, pStart}

	if params.ServiceName != "" {
		query += " AND service_name = $4"
		args = append(args, params.ServiceName)
	}

	query += " AND start_date <= COALESCE(TO_DATE($2, 'MM-YYYY'), end_date)"
	query += " AND end_date >= COALESCE(TO_DATE($3, 'MM-YYYY'), start_date)"

	var sum int
	err := s.dbpool.QueryRow(ctx, query, args...).Scan(&sum)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return sum, nil
}
