# ArcBridge Hybrid Resource Projection and Extension Platform (POC)

This repository contains a proof-of-concept implementation of ArcBridge, a hybrid management plane that connects Kubernetes clusters to a central control plane and safely delivers platform extensions.

## Repository Structure
- `arcbridge/docs/` – Design notes, API contract, and incident runbooks.
- `arcbridge/agent/` – Go-based agent that simulates watching the `ArcBridgeCluster` CRD and reconciling extensions.
- `arcbridge/controlplane/` – Go-based control plane exposing REST endpoints for registration, inventory, desired state, and status events.
- `arcbridge/telemetry/` – OpenTelemetry collector config, Prometheus rules, and Grafana dashboards.
- `arcbridge/deploy/` – Kind cluster specs, Kustomize overlays, and Helm charts for sample extensions.
- `arcbridge/tools/` – Synthetic data generators, load tests, and trace exploration helpers.

## Getting Started
1. **Install dependencies**
   - Go 1.21+
   - Node.js (for k6) and Python 3 (for synthetic scripts)

2. **Run the control plane locally**
   ```bash
   go run ./arcbridge/controlplane/cmd/apiserver
   ```

3. **Simulate the agent**
   ```bash
   go run ./arcbridge/agent/cmd/agent
   ```

4. **Generate synthetic inventory**
   ```bash
   python arcbridge/tools/synthetic_checks/generate_inventory.py
   ```

5. **Load test registration endpoint**
   ```bash
   k6 run arcbridge/tools/loadgen/register_scenario.js
   ```

## Telemetry
- Launch the collector and Jaeger locally using Docker Compose from `tools/trace_viewer/`.
- Metrics can be scraped from the collector's Prometheus exporter on port 9464.

## Notes
- All credentials and certificates are placeholders; integrate with Azure Key Vault or AWS Secrets Manager for production use.
- GitOps integration via Flux is not implemented, but hooks are present for future extension.
