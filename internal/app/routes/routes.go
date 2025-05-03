package routes

import (
	"github.com/antoniolopesg/workout-api/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)

	r.Route("/workouts", func(r chi.Router) {
		r.Post("/", app.WorkoutHandler.HandleCreateWorkout)

		r.Get("/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
		r.Put("/{id}", app.WorkoutHandler.HandleUpdateWorkoutByID)
		r.Delete("/{id}", app.WorkoutHandler.HandleDeleteWorkoutByID)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", app.UserHandler.HandleCreateUser)
	})

	return r
}
