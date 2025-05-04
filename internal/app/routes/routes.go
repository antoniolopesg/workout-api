package routes

import (
	"github.com/antoniolopesg/workout-api/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(app.UserMiddleware.Authenticate)

		r.Post("/users/", app.UserMiddleware.RequireUser(app.WorkoutHandler.HandleCreateWorkout))
		r.Get("/users/{id}", app.UserMiddleware.RequireUser(app.WorkoutHandler.HandleGetWorkoutByID))
		r.Put("/users/{id}", app.UserMiddleware.RequireUser(app.WorkoutHandler.HandleUpdateWorkoutByID))
		r.Delete("/users/{id}", app.UserMiddleware.RequireUser(app.WorkoutHandler.HandleDeleteWorkoutByID))

	})

	r.Get("/health", app.HealthCheck)
	r.Post("/users", app.UserHandler.HandleCreateUser)
	r.Post("/tokens/auth", app.TokenHandler.HandleCreateToken)

	return r
}
