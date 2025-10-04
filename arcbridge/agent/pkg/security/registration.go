package security

import (
	"encoding/json"
	"errors"
	"os"
)

// Registration holds cluster identity data returned by the control plane.
type Registration struct {
	ClusterID string `json:"clusterID"`
	CertPath  string `json:"certPath"`
	KeyPath   string `json:"keyPath"`
}

// LoadRegistration loads registration data from a local file.
func LoadRegistration() (*Registration, error) {
	path := os.Getenv("ARCBRIDGE_REGISTRATION_FILE")
	if path == "" {
		path = "./deploy/kind/registration.json"
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Registration{ClusterID: "cluster-poc", CertPath: "", KeyPath: ""}, nil
	} else if err != nil {
		return nil, err
	}

	var reg Registration
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, err
	}
	if reg.ClusterID == "" {
		reg.ClusterID = "cluster-poc"
	}
	return &reg, nil
}
