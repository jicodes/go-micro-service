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
	config Config
}

func New(config Config) *App {
	app := &App{
		rdb:    redis.NewClient(&redis.Options{
			Addr: config.RedisAddr,
		}),
		config: config,
	}

	app.loadRoutes()

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

	fmt.Printf("Starting the server on: %d\n", app.config.ServerPort)

	app.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.ServerPort),
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