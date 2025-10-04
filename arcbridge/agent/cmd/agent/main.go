package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"arcbridge/arcbridge/agent/pkg/informers"
	"arcbridge/arcbridge/agent/pkg/reconcile"
	"arcbridge/arcbridge/agent/pkg/security"
	"arcbridge/arcbridge/telemetry/setup"
)

func main() {
	logger := log.New(os.Stdout, "agent | ", log.LstdFlags)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	registration, err := security.LoadRegistration()
	if err != nil {
		logger.Fatalf("failed to load registration: %v", err)
	}

	shutdownTracing := setup.InitTracing(logger)
	defer func() {
		_ = shutdownTracing(context.Background())
	}()

	informer := informers.NewClusterInformer(logger)
	reconciler := reconcile.NewExtensionReconciler(logger)

	queue := make(chan informers.ClusterEvent, 10)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				queue <- informers.ClusterEvent{ClusterID: registration.ClusterID, Generation: rand.Intn(1000)}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case evt := <-queue:
				status := reconciler.Reconcile(ctx, evt)
				logger.Printf("reconcile complete: %+v", status)
			}
		}
	}()

	informer.Start(ctx)
	wg.Wait()
	logger.Println("agent shutting down")
}
