package main

import (
	"anemone_backend-microservices/internal/mail/api"
	"anemone_backend-microservices/internal/mail/config"
	"anemone_backend-microservices/internal/mail/database"
	"anemone_backend-microservices/internal/mail/repository"
	"anemone_backend-microservices/internal/mail/services"
	"anemone_backend-microservices/internal/mail/smtp"
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func setupCORS(router http.Handler) http.Handler {
	cfg := config.Load()

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
	cfg := config.Load()

	db, err := database.NewConnection(cfg.DatabaseURL, migrationFS)
	if err != nil {
		log.Fatalf("FATAL: database connection failed: %v", err)
	}
	defer db.Close()
	log.Println("INFO: Database connection successful")

	repo := repository.New(db)
	service := services.New(repo, cfg.DomainName)
	handler := api.NewMailHandler(service, repo)

	r := mux.NewRouter()
	handler.Register(r, cfg.JWTSecret)

	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: setupCORS(r),
	}

	smtpServer := smtp_server.NewServer(cfg, repo)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("INFO: Starting SMTP server")
		smtpServer.Start()
	}()

	go func() {
		log.Printf("INFO: Starting HTTP server on port %s", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("FATAL: failed to start HTTP server: %v", err)
		}
	}()

	log.Println("INFO: All services are running")

	<-quit
	log.Println("INFO: Shutting down services...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("ERROR: HTTP server forced to shutdown: %v", err)
	}

	log.Println("INFO: Server exited properly")
}
