package informers

import (
	"context"
	"log"
	"time"
)

// ClusterEvent represents a reconciliation trigger for a given cluster generation.
type ClusterEvent struct {
	ClusterID  string
	Generation int
}

// ClusterInformer is a placeholder watcher that simulates CRD events.
type ClusterInformer struct {
	logger *log.Logger
}

// NewClusterInformer constructs a new informer instance.
func NewClusterInformer(logger *log.Logger) *ClusterInformer {
	return &ClusterInformer{logger: logger}
}

// Start simulates informer startup and logs readiness.
func (c *ClusterInformer) Start(ctx context.Context) {
	c.logger.Println("cluster informer starting")
	<-ctx.Done()
	c.logger.Println("cluster informer stopping")
}

// Sync simulates cache synchronization.
func (c *ClusterInformer) Sync() {
	time.Sleep(500 * time.Millisecond)
	c.logger.Println("cluster informer synced")
}
