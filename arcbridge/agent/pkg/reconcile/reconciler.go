package reconcile

import (
	"context"
	"log"
	"math/rand"
	"time"
)

// ExtensionStatus describes the result of an extension reconcile attempt.
type ExtensionStatus struct {
	ExtensionName      string
	LastTransition     time.Time
	Ready              bool
	Reason             string
	ObservedGeneration int
}

// ExtensionReconciler simulates applying extensions via Helm or Kustomize.
type ExtensionReconciler struct {
	logger *log.Logger
}

// NewExtensionReconciler constructs a placeholder reconciler.
func NewExtensionReconciler(logger *log.Logger) *ExtensionReconciler {
	return &ExtensionReconciler{logger: logger}
}

// Reconcile simulates installation, readiness checks, and emits pseudo status.
func (r *ExtensionReconciler) Reconcile(ctx context.Context, evt interface{}) ExtensionStatus {
	clusterEvent, _ := evt.(struct {
		ClusterID  string
		Generation int
	})

	delay := time.Duration(rand.Intn(400)+100) * time.Millisecond
	select {
	case <-ctx.Done():
		return ExtensionStatus{ExtensionName: "nginx-ingress", Ready: false, Reason: "context cancelled", LastTransition: time.Now(), ObservedGeneration: 0}
	case <-time.After(delay):
	}

	status := ExtensionStatus{
		ExtensionName:      "nginx-ingress",
		Ready:              rand.Intn(10) > 0,
		Reason:             "applied",
		LastTransition:     time.Now(),
		ObservedGeneration: clusterEvent.Generation,
	}

	if !status.Ready {
		status.Reason = "transient-error"
	}

	r.logger.Printf("cluster %s generation %d reconcile delay %s", clusterEvent.ClusterID, clusterEvent.Generation, delay)
	return status
}
