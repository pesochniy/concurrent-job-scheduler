package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pesochniy/concurrent-job-scheduler/handlers"
	"github.com/pesochniy/concurrent-job-scheduler/internal/job"
	"github.com/pesochniy/concurrent-job-scheduler/internal/scheduler"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	HTTPAddr string
}

func loadConfig() Config {
	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8080"
	}

	return Config{
		HTTPAddr: httpAddr,
	}
}

func main() {
	cfg := loadConfig()
	mux := http.NewServeMux()

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/jobsdb"
	}

	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	//store := job.NewMemoryStore()
	store := job.NewPostgresStore(db)
	sched := scheduler.NewScheduler(store, 4)

	h := handlers.NewHandler(store, sched)
	handlers.Register(mux, h)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: mux,
	}
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()
	log.Printf("server started on %s", cfg.HTTPAddr)

	go func() {
		if err := sched.Start(ctx); err != nil {
			log.Printf("scheduler stopped: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
