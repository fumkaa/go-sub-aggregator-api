package httpserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/fumkaa/go-sub-aggregator-api/internal/domain/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	CreateSub(ctx context.Context, sub models.Subscription) (uuid.UUID, error)
	GetSubByID(ctx context.Context, id uuid.UUID) (models.Subscription, error)
	UpdateSub(ctx context.Context, sub models.Subscription) error
	DeleteSub(ctx context.Context, id uuid.UUID) error
	ListSubs(ctx context.Context, userId uuid.UUID) ([]models.Subscription, error)
	GetSum(ctx context.Context, params models.ListSubsParams) (int, error)
}

func (s *serverAPI) CreateSub(w http.ResponseWriter, r *http.Request) {
	const op = "httpserver.serverAPI.CreateSub()"
	log := s.log.With(
		slog.String("op", op),
	)

	var sub models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		log.ErrorContext(r.Context(), "failed to decode request body", slog.String("error", err.Error()))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	sub.ID = uuid.New()

	log.InfoContext(r.Context(), "attempting to create subscription",
		slog.String("user_id", sub.UserID.String()),
		slog.String("service", sub.ServiceName),
	)

	id, err := s.subService.CreateSub(r.Context(), sub)
	if err != nil {
		log.ErrorContext(r.Context(), "failed to create subscription", slog.String("error", err.Error()))
		http.Error(w, "Failed to create subscription", http.StatusInternalServerError)
		return
	}

	log.InfoContext(r.Context(), "subscription created", slog.String("id", id.String()))

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]uuid.UUID{"id": id})
}

func (s *serverAPI) GetSubByID(w http.ResponseWriter, r *http.Request) {
	const op = "httpserver.serverAPI.GetSubByID()"
	log := s.log.With(
		slog.String("op", op),
	)

	id := chi.URLParam(r, "id")

	log.InfoContext(r.Context(), "getting subscription by ID", slog.String("id", id))

	idUUID, err := uuid.Parse(id)
	if err != nil {
		log.ErrorContext(r.Context(), "failed to parse subscription ID", slog.String("error", err.Error()))
		http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
		return
	}

	sub, err := s.subService.GetSubByID(r.Context(), idUUID)
	if err != nil {
		log.ErrorContext(r.Context(), "failed to get subscription", slog.String("error", err.Error()))
		http.Error(w, "Failed to get subscription", http.StatusInternalServerError)
		return
	}

	log.InfoContext(r.Context(), "subscription retrieved", slog.String("id", sub.ID.String()))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sub)
}

func (s *serverAPI) UpdateSub(w http.ResponseWriter, r *http.Request) {
	const op = "httpserver.serverAPI.UpdateSub()"
	log := s.log.With(
		slog.String("op", op),
	)

	id := chi.URLParam(r, "id")

	log.InfoContext(r.Context(), "updating subscription", slog.String("id", id))

	idUUID, err := uuid.Parse(id)
	if err != nil {
		log.ErrorContext(r.Context(), "failed to parse subscription ID", slog.String("error", err.Error()))
		http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
		return
	}

	var sub models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		log.ErrorContext(r.Context(), "failed to decode request body", slog.String("error", err.Error()))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	sub.ID = idUUID

	err = s.subService.UpdateSub(r.Context(), sub)
	if err != nil {
		log.ErrorContext(r.Context(), "failed to update subscription", slog.String("error", err.Error()))
		http.Error(w, "Failed to update subscription", http.StatusInternalServerError)
		return
	}

	log.InfoContext(r.Context(), "subscription updated", slog.String("id", sub.ID.String()))

	w.WriteHeader(http.StatusOK)
}

func (s *serverAPI) DeleteSub(w http.ResponseWriter, r *http.Request) {
	const op = "httpserver.serverAPI.DeleteSub()"
	log := s.log.With(
		slog.String("op", op),
	)

	id := chi.URLParam(r, "id")

	log.InfoContext(r.Context(), "deleting subscription", slog.String("id", id))

	idUUID, err := uuid.Parse(id)
	if err != nil {
		log.ErrorContext(r.Context(), "failed to parse subscription ID", slog.String("error", err.Error()))
		http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
		return
	}

	err = s.subService.DeleteSub(r.Context(), idUUID)
	if err != nil {
		log.ErrorContext(r.Context(), "failed to delete subscription", slog.String("error", err.Error()))
		http.Error(w, "Failed to delete subscription", http.StatusInternalServerError)
		return
	}

	log.InfoContext(r.Context(), "subscription deleted", slog.String("id", id))

	w.WriteHeader(http.StatusOK)
}

func (s *serverAPI) ListSubs(w http.ResponseWriter, r *http.Request) {
	const op = "httpserver.serverAPI.ListSubs()"
	log := s.log.With(
		slog.String("op", op),
	)

	userId := r.URL.Query().Get("user_id")

	log.InfoContext(r.Context(), "listing subscriptions", slog.String("user_id", userId))

	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		log.ErrorContext(r.Context(), "failed to parse user ID", slog.String("error", err.Error()))
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	subs, err := s.subService.ListSubs(r.Context(), userIdUUID)
	if err != nil {
		log.ErrorContext(r.Context(), "failed to list subscriptions", slog.String("error", err.Error()))
		http.Error(w, "Failed to list subscriptions", http.StatusInternalServerError)
		return
	}

	log.InfoContext(r.Context(), "subscriptions listed", slog.Int("count", len(subs)))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(subs)
}

func (s *serverAPI) GetSum(w http.ResponseWriter, r *http.Request) {
	const op = "httpserver.serverAPI.GetSum()"
	log := s.log.With(
		slog.String("op", op),
	)

	userId := chi.URLParam(r, "user_id")
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		log.ErrorContext(r.Context(), "failed to parse user ID", slog.String("error", err.Error()))
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	params := models.ListSubsParams{
		UserID:      userIdUUID,
		ServiceName: chi.URLParam(r, "service_name"),
		StartDate:   chi.URLParam(r, "start_date"),
		EndDate:     chi.URLParam(r, "end_date"),
	}

	log.InfoContext(r.Context(), "calculating total cost",
		slog.Any("params", params),
	)

	sum, err := s.subService.GetSum(r.Context(), params)
	if err != nil {
		log.ErrorContext(r.Context(), "failed to get sum", slog.String("error", err.Error()))
		http.Error(w, "Failed to get sum", http.StatusInternalServerError)
		return
	}

	log.InfoContext(r.Context(), "sum calculated successfully", slog.Int("result", sum))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"sum": sum})
}
