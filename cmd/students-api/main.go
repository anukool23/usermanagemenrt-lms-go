package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anukool23/usermanagement-lms-go/internal/config"
	"github.com/anukool23/usermanagement-lms-go/internal/http/middleware"
	"github.com/anukool23/usermanagement-lms-go/internal/http/routes"
	"github.com/anukool23/usermanagement-lms-go/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		slog.Error("Failed to initialize storage: ", slog.String("error", err.Error()))
	}
	slog.Info("Storage initialized successfully", slog.String("env", cfg.Env))

	//setup router
	router := http.NewServeMux()
	routes.Register(router, storage)

	allowedSecretKeys := []string{
		"7b8d9f7b8d9f7",
		"7b8d9f7b8d9f8",
	}

	handler := middleware.Chain(
		router,
		middleware.SecretKeyAuth(allowedSecretKeys),
	)

	//setup server
	server := http.Server{
		Addr:              cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}


	log.Printf("Server started on %s", cfg.Port)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server: ", err)
		}
	}()
	<-done

	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown server: ", slog.String("error", err.Error()))
	}
	slog.Info("Server gracefully stopped")
}
