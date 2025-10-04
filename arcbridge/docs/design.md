# ArcBridge Platform Design (POC)

This document captures the proof-of-concept design for the ArcBridge hybrid resource projection and extension platform. The goal is to provide enough scaffolding for experimentation while leaving clear placeholders for production-grade implementations.

## High-Level Goals
- Connect any Kubernetes cluster to a central control plane and manage it as a first-class resource.
- Enable platform teams to declare extensions once and roll them out safely across many clusters with drift detection and rollback.
- Offer strong health signals and runbooks so a designated responder can restore service quickly.

## Architecture Overview
1. **Agent**: Runs in each customer cluster, watches the `ArcBridgeCluster` custom resource, and reconciles desired extensions.
2. **Control Plane**: Exposes REST APIs for registration, desired-state publication, and status reporting. A lightweight gRPC side-channel is reserved for streaming heartbeats and inventories.
3. **Storage**: Postgres for durable cluster and extension state, Redis for hot-path queues and rate limiting.
4. **Observability**: OpenTelemetry traces and Prometheus metrics across agent and control plane.
5. **Security**: Azure AD OIDC for users and services, multi-tenant namespaces, and mTLS between agents and the control plane.

## Data Flow Summary
1. A new cluster posts to `/api/v1/register` with a bootstrap token and receives a signed client certificate.
2. The agent stores registration metadata in the `ArcBridgeCluster` CR and begins streaming inventory via gRPC.
3. The control plane publishes desired extension specs via `/api/v1/desired`.
4. The agent reconciles desired state using Helm or Kustomize, verifies readiness, and reports status events.

## Drift and Rollback Strategy
- Each extension status includes a reason code, message, and last transition time.
- The control plane tracks the last known good spec per extension, enabling single-click rollback.

## Open Questions
- **Secret Distribution**: Prototype uses local environment variables; production will leverage Azure Key Vault / AWS Secrets Manager via CSI drivers.
- **GitOps Integration**: Flux integration is stubbed; future work should solidify a Git-based desired state pipeline.
- **Fleet Pacing**: Token bucket limits are configurable but static in the POC. Real systems may require adaptive rate control.

