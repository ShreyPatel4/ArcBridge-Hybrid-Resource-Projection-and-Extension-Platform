package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"arcbridge/arcbridge/controlplane/pkg/handlers"
	"arcbridge/arcbridge/telemetry/setup"
)

func main() {
	logger := log.New(os.Stdout, "controlplane | ", log.LstdFlags)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	shutdownTracing := setup.InitTracing(logger)
	defer func() {
		_ = shutdownTracing(context.Background())
	}()

	router := chi.NewRouter()
	state := handlers.NewStateStore()

	router.Post("/api/v1/register", handlers.RegisterHandler(logger, state))
	router.Put("/api/v1/inventory", handlers.InventoryHandler(logger, state))
	router.Get("/api/v1/desired", handlers.DesiredHandler(logger, state))
	router.Post("/api/v1/status", handlers.StatusHandler(logger, state))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	logger.Println("control plane listening on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("server error: %v", err)
	}
	wg.Wait()
	logger.Println("control plane shutdown complete")

	// Dump state snapshot for quick visibility in POC.
	snapshot, _ := json.MarshalIndent(state.DebugSnapshot(), "", "  ")
	logger.Printf("final state: %s", snapshot)
}
