package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sahilrana27582/payment-system/internal/config"
	"github.com/sahilrana27582/payment-system/internal/db"
	"github.com/sahilrana27582/payment-system/internal/handler"
	"github.com/sahilrana27582/payment-system/internal/repository"
	"github.com/sahilrana27582/payment-system/internal/routes"
	"github.com/sahilrana27582/payment-system/internal/service"
)

func main() {
	cfg := config.GetConfig()

	pool, err := db.GetDBConn(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := db.RunMigrations(cfg); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}

	log.Println("Database connected and migrations completed successfully")

	userRepo := repository.NewRepository(pool)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r := routes.NewRouter(userHandler)

	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErr := make(chan error, 1)

	go func() {
		log.Printf("HTTP server listening on %s", srv.Addr)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	stopCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case err := <-serverErr:
		log.Fatalf("HTTP server failed: %v", err)

	case <-stopCtx.Done():
		log.Println("Shutting down HTTP server...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)
		if err := srv.Close(); err != nil {
			log.Printf("Failed to close HTTP server: %v", err)
		}
	}

	log.Println("HTTP server stopped")
}
