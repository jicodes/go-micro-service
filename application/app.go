package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
	server *http.Server
}

func New() *App {
	app := &App{
		router: loadRoutes(),
		rdb:    redis.NewClient(&redis.Options{}),
	}
	return app
}

func (app *App) Start(ctx context.Context) error {
	err := app.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}

	defer func() {
		if err := app.rdb.Close(); err != nil {
			fmt.Println("failed to close redis connection:", err)
		}
	}()

	fmt.Println("Starting the server on :8080")
	app.server = &http.Server{
		Addr:    ":8080",
		Handler: app.router,
	}

	ch := make(chan error, 1)

	go func() {
		err = app.server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		} 
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		app.server.Shutdown(timeout)
		return nil
	}
}