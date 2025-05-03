package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/antoniolopesg/workout-api/internal/api"
	"github.com/antoniolopesg/workout-api/internal/store"
	"github.com/antoniolopesg/workout-api/migrations"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	pgDb, err := store.Open()
	if err != nil {
		return nil, err
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	err = store.MigrateFS(pgDb, migrations.FS, ".")

	if err != nil {
		panic(err)
	}

	// stores
	workoutStore := store.NewPostgresWorkoutStore(pgDb)
	userStore := store.NewPostgresUserStore(pgDb)

	// handlers
	userHandler := api.NewUserHandler(userStore, logger)
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHandler,
		DB:             pgDb,
	}
	return app, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
