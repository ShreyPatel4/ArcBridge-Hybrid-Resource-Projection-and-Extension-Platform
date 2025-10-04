package storage

import "sync"

// FeatureFlags stores toggleable features for experimentation.
type FeatureFlags struct {
	mu     sync.RWMutex
	values map[string]bool
}

// NewFeatureFlags creates a feature flag store with defaults.
func NewFeatureFlags() *FeatureFlags {
	return &FeatureFlags{
		values: map[string]bool{
			"rollouts.enabled": true,
		},
	}
}

// Enabled returns whether a flag is active.
func (f *FeatureFlags) Enabled(name string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.values[name]
}

// Set updates a flag value.
func (f *FeatureFlags) Set(name string, value bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.values[name] = value
}
