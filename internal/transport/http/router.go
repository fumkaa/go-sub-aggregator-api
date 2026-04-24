package httpserver

import "github.com/go-chi/chi/v5"

func RegisterRoutes(api *serverAPI) *chi.Mux {
	router := chi.NewRouter()

	router.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", api.CreateSub)
		r.Get("/{id}", api.GetSubByID)
		r.Put("/{id}", api.UpdateSub)
		r.Delete("/{id}", api.DeleteSub)
		r.Get("/", api.List)
	})

	return router
}
