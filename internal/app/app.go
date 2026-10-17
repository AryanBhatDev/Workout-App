package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/AryanBhatDev/Workout-App/internal/api"
)

type Application struct {
	Logger         *slog.Logger
	WorkoutHandler *api.WorkoutHandler
}

func NewApplication() (*Application, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	workoutHandler := api.NewWorkoutHandler()

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
	}

	return app, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
