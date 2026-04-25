package httpserver

import (
	"github.com/fumkaa/go-sub-aggregator-api/internal/domain/models"
	"github.com/fumkaa/go-sub-aggregator-api/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
	middlewareChi "github.com/go-chi/chi/v5/middleware"
)

func RegisterRoutes(api *serverAPI) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middlewareChi.RequestID)

	r.Route("/subscriptions", func(r chi.Router) {
		r.With(middleware.ValidateQuery(api.log, models.ListSubsParams{})).Get("/total", api.GetSum)
		r.With(middleware.ValidateUUID(api.log, "id")).Get("/{id}", api.GetSubByID)
		r.Get("/", api.ListSubs)

		r.With(middleware.BindAndValidate(api.log, models.Subscription{})).Post("/", api.CreateSub)

		r.With(middleware.BindAndValidate(api.log, models.Subscription{})).Put("/{id}", api.UpdateSub)

		r.With(middleware.ValidateUUID(api.log, "id")).Delete("/{id}", api.DeleteSub)

	})

	return r
}
