package application

import (
	"context"
	"fmt"
	"net/http"
)

type App struct {
	router http.Handler
}

func New() *App {
	app := &App{
		router: loadRoutes(),
	}
	return app
}

func (app *App) Start(ctx context.Context) error {
	err := http.ListenAndServe(":8080" , app.router)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}