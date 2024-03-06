package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb *redis.Client
}

func New() *App {
	app := &App{
		router: loadRoutes(),
		rdb: redis.NewClient(&redis.Options{}),
	}
	return app
}

func (app *App) Start(ctx context.Context) error {

	err := app.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}

	fmt.Println("Starting the server on :8080")

	err = http.ListenAndServe(":8080" , app.router)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}