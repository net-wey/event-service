package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "event-service/docs"
	"event-service/internal/service"
)

// NewRouter настраивает все HTTP-маршруты.
func NewRouter(eventSvc *service.EventService, participantSvc *service.ParticipantService, venueSvc *service.VenueService) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.Get("/api/health", HealthCheck)

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	eventHandler := NewEventHandler(eventSvc)
	participantHandler := NewParticipantHandler(participantSvc)
	venueHandler := NewVenueHandler(venueSvc)

	r.Route("/api/events", func(r chi.Router) {
		r.Get("/", eventHandler.List)
		r.Post("/", eventHandler.Create)

		r.Route("/{eventID}", func(r chi.Router) {
			r.Get("/", eventHandler.Get)
			r.Put("/", eventHandler.Update)
			r.Delete("/", eventHandler.Delete)

			r.Route("/participants", func(r chi.Router) {
				r.Get("/", participantHandler.ListByEvent)
				r.Post("/", participantHandler.Register)
				r.Put("/{participantID}", participantHandler.Update)
				r.Delete("/{participantID}", participantHandler.Delete)
			})
		})
	})

	r.Route("/api/venues", func(r chi.Router) {
		r.Get("/", venueHandler.List)
		r.Post("/", venueHandler.Create)

		r.Route("/{venueID}", func(r chi.Router) {
			r.Get("/", venueHandler.Get)
			r.Put("/", venueHandler.Update)
			r.Delete("/", venueHandler.Delete)
		})
	})

	return r
}
