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
	handlers.Register(mux)

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

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
