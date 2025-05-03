package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/antoniolopesg/workout-api/internal/app"
	"github.com/antoniolopesg/workout-api/internal/app/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "go backend server ort")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	defer app.DB.Close()

	r := routes.SetupRoutes(app)

	app.Logger.Println("Application Structure Instantiated Successfully!")

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      r,
	}

	app.Logger.Printf("HTTP Initializing on port %d\n", port)

	err = server.ListenAndServe()

	if err != nil {
		app.Logger.Fatal(err)
	}
}
