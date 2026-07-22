package main

import (
	"anemone_backend-kanban/internal/api"
	"anemone_backend-kanban/internal/cfg"
	"anemone_backend-kanban/internal/db"
	"anemone_backend-kanban/internal/repo"
	"anemone_backend-kanban/internal/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"
)

func setupCORS(router http.Handler, cfg *cfg.Config) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.CorsDev, cfg.CorsProd},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            false,
	})

	return c.Handler(router)
}

func main() {
	cfg, err := cfg.Load()
	if err != nil {
		log.Fatalf("Critical error: failed to initialize configuration: %v", err)
	}

	database, err := db.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("FATAL: database connection failed: %v", err)
	}
	defer func() {
		log.Println("INFO: Closing database connection...")
		database.Close()
	}()

	repo := repo.NewRepo(database, cfg.DefaultStep)
	src := service.NewService(repo, repo, repo, cfg.DefaultStep)
	handler := api.NewHandler(src, src, src, cfg.AccessSecret)

	r := http.NewServeMux()

	handler.Routes(r)

	handlerWithCORS := setupCORS(r, cfg)

	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: handlerWithCORS,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("INFO: Starting HTTP server on port %s", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("FATAL: failed to start HTTP server: %v", err)
		}
	}()

	log.Println("INFO: All services are running")

	<-quit
	log.Println("INFO: Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("FATAL: Server forced to shutdown: %v", err)
	}

	log.Println("INFO: Server exiting")
}
