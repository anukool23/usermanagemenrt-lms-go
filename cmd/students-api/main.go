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
	"github.com/anukool23/usermanagement-lms-go/internal/http/handlers/student"
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

	router.HandleFunc("POST /api/student", student.New(storage))
	router.HandleFunc("GET /api/student/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.ListStudents(storage))
	router.HandleFunc("DELETE /api/student/{id}", student.DeleteById(storage))
	//setup server
	server := http.Server{
		Addr:    cfg.Port,
		Handler: router,
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
