package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/AryanBhatDev/Workout-App/internal/api"
	"github.com/AryanBhatDev/Workout-App/internal/store"
	"github.com/AryanBhatDev/Workout-App/migrations"
)

type Application struct {
	Logger         *slog.Logger
	WorkoutHandler *api.WorkoutHandler
	DB             *sql.DB
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
}

func NewApplication() (*Application, error) {

	pgDB, err := store.OpenDB()

	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")

	if err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)

	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		DB:             pgDB,
		UserHandler:    userHandler,
		TokenHandler:   tokenHandler,
	}

	return app, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
