package handlers

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// StateStore simulates persistent storage for the POC.
type StateStore struct {
	mu        sync.RWMutex
	clusters  map[string]*ClusterRecord
	inventory map[string]Inventory
}

// ClusterRecord tracks registration and desired state for a cluster.
type ClusterRecord struct {
	ClusterID string            `json:"clusterID"`
	Tenant    string            `json:"tenant"`
	Region    string            `json:"region"`
	Desired   map[string]string `json:"desired"`
}

// Inventory captures a simplified cluster inventory snapshot.
type Inventory struct {
	Nodes      int               `json:"nodes"`
	Pods       int               `json:"pods"`
	Extensions map[string]string `json:"extensions"`
	UpdatedAt  time.Time         `json:"updatedAt"`
}

// DebugSnapshot returns the internal state for runbook inspection.
func (s *StateStore) DebugSnapshot() map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshot := make(map[string]any)
	for id, record := range s.clusters {
		snapshot[id] = map[string]any{
			"record":    record,
			"inventory": s.inventory[id],
		}
	}
	return snapshot
}

// NewStateStore creates a new in-memory store.
func NewStateStore() *StateStore {
	return &StateStore{
		clusters:  make(map[string]*ClusterRecord),
		inventory: make(map[string]Inventory),
	}
}

// RegisterHandler handles cluster registration.
func RegisterHandler(logger *log.Logger, store *StateStore) http.HandlerFunc {
	type request struct {
		ClusterUID     string `json:"clusterUID"`
		BootstrapToken string `json:"bootstrapToken"`
		Metadata       struct {
			Region string `json:"region"`
			Tenant string `json:"tenant"`
		} `json:"metadata"`
	}
	type response struct {
		ClusterID      string `json:"clusterID"`
		CertificatePEM string `json:"certificatePEM"`
		ExpiresAt      string `json:"expiresAt"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		store.mu.Lock()
		defer store.mu.Unlock()

		clusterID := "cluster-" + req.ClusterUID
		record := &ClusterRecord{
			ClusterID: clusterID,
			Tenant:    req.Metadata.Tenant,
			Region:    req.Metadata.Region,
			Desired: map[string]string{
				"nginx-ingress": "1.0.0",
				"logging":       "0.2.1",
			},
		}
		store.clusters[clusterID] = record

		resp := response{
			ClusterID:      clusterID,
			CertificatePEM: "-----BEGIN CERTIFICATE-----POC-----END CERTIFICATE-----",
			ExpiresAt:      time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		}

		logger.Printf("registered cluster %s tenant=%s region=%s", clusterID, record.Tenant, record.Region)
		_ = json.NewEncoder(w).Encode(resp)
	}
}

// InventoryHandler receives inventory updates and stores them.
func InventoryHandler(logger *log.Logger, store *StateStore) http.HandlerFunc {
	type request struct {
		ClusterID  string            `json:"clusterID"`
		Nodes      int               `json:"nodes"`
		Pods       int               `json:"pods"`
		Extensions map[string]string `json:"extensions"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		store.mu.Lock()
		store.inventory[req.ClusterID] = Inventory{
			Nodes:      req.Nodes,
			Pods:       req.Pods,
			Extensions: req.Extensions,
			UpdatedAt:  time.Now(),
		}
		store.mu.Unlock()

		logger.Printf("inventory update for %s nodes=%d pods=%d", req.ClusterID, req.Nodes, req.Pods)
		w.WriteHeader(http.StatusAccepted)
	}
}

// DesiredHandler returns desired extension specs for a cluster.
func DesiredHandler(logger *log.Logger, store *StateStore) http.HandlerFunc {
	type response struct {
		Extensions map[string]string `json:"extensions"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		clusterID := r.URL.Query().Get("clusterID")
		store.mu.RLock()
		record, ok := store.clusters[clusterID]
		store.mu.RUnlock()

		if !ok {
			http.Error(w, "cluster not found", http.StatusNotFound)
			return
		}

		logger.Printf("served desired state for %s", clusterID)
		_ = json.NewEncoder(w).Encode(response{Extensions: record.Desired})
	}
}

// StatusHandler accepts reconcile status events.
func StatusHandler(logger *log.Logger, store *StateStore) http.HandlerFunc {
	type request struct {
		ClusterID string `json:"clusterID"`
		Name      string `json:"name"`
		Ready     bool   `json:"ready"`
		Reason    string `json:"reason"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		logger.Printf("status event cluster=%s extension=%s ready=%t reason=%s", req.ClusterID, req.Name, req.Ready, req.Reason)

		// Simulate drift detection by randomly adjusting desired state.
		if !req.Ready && rand.Intn(5) == 0 {
			store.mu.Lock()
			if record, ok := store.clusters[req.ClusterID]; ok {
				record.Desired[req.Name] = "rollback-" + time.Now().Format("150405")
				logger.Printf("triggered rollback for %s extension=%s", req.ClusterID, req.Name)
			}
			store.mu.Unlock()
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
