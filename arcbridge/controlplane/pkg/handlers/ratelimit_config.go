package handlers

// RateLimitConfig defines rollout pacing for the POC.
type RateLimitConfig struct {
	PerTenant int `json:"perTenant"`
	PerRegion int `json:"perRegion"`
	Global    int `json:"global"`
}

// DefaultRateLimits returns a conservative configuration suitable for Kind.
func DefaultRateLimits() RateLimitConfig {
	return RateLimitConfig{
		PerTenant: 5,
		PerRegion: 10,
		Global:    25,
	}
}
