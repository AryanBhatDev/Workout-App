package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/AryanBhatDev/Workout-App/internal/app"
	"github.com/AryanBhatDev/Workout-App/internal/routes"
)

func main() {

	var port int
	flag.IntVar(&port, "port", 8080, "port address")
	flag.Parse()

	application, err := app.NewApplication()

	if err != nil {
		panic(err)
	}

	application.Logger.Info("App running")

	r := routes.Setup(application)

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ErrorLog:     slog.NewLogLogger(application.Logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = server.ListenAndServe()

	if err != nil {
		application.Logger.Error(err.Error())
		os.Exit(1)
	}
}
