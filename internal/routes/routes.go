package routes

import (
	"github.com/AryanBhatDev/Workout-App/internal/app"
	"github.com/go-chi/chi/v5"
)

func Setup(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)

		r.Group(func(r chi.Router) {
			r.Use(app.Middleware.RequireUser)

			r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
			r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
			r.Put("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
			r.Delete("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkoutByID)
		})
	})

	r.Get("/health", app.HealthCheck)

	r.Post("/users", app.UserHandler.HandleRegisterUser)
	r.Post("/tokens/authentication", app.TokenHandler.HandleCreateToken)

	return r
}
