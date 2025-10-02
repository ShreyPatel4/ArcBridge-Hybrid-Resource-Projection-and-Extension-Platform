# ArcBridge-Hybrid-Resource-Projection-and-Extension-Platform
Manage many Kubernetes clusters anywhere and project their state into a central Azure like control plane with secure registration, extension lifecycle, and production grade observability.


# ArcBridge

Portable control plane for fleets of Kubernetes clusters with safe multi cluster rollout, strong health signals, and fast restore paths.

## Goals

* Provide a simple way to connect any cluster to a central control plane and see it as a first class resource
* Let platform teams declare extensions once and roll them out safely across many clusters with drift detection and rollback
* Offer strong health signals and runbooks so a DRI can restore service quickly

## Architecture overview

* Agent inside each customer cluster watches a small CRD named **ArcBridgeCluster** that holds registration info and desired extensions
* Control plane service exposes REST APIs to issue tokens, accept heartbeats, receive inventory, and publish desired state per cluster
* Reconciler in the agent installs extensions using Helm and Kustomize, verifies readiness, and reports detailed status with reasons and last transition time
* Observability layer uses OpenTelemetry traces and Prometheus metrics emitted by both agent and control plane so you can follow a request across the system
* Identity and security rely on Azure AD OIDC for user and service auth, per tenant namespaces, and mTLS between agents and control plane

## Key components

* **Agent** written in Go with client go informers, work queues, and a backoff aware reconcile loop
* **Control plane** written in Go with Gin or Chi for REST and a small gRPC side channel for streaming heartbeats and inventory
* **Storage** Postgres for cluster and extension state, Redis for hot queues and rate limiting
* **GitOps optional path** with Flux where the control plane writes desired state to a Git repo and agents pull and apply
* **Policy** Open Policy Agent admission checks for extension parameters and namespace safety

## Interesting problems and solutions

* **Duplicate registration and stale tokens**  
  Implement a one time registration flow with a short lived bootstrap token that swaps to a long lived client cert tied to cluster UID

* **Large fleet rollout noise and thundering herds**  
  Use a token bucket per region and per tenant to pace extension installs and a global rate cap to protect the control plane

* **Partial failure and drift**  
  Record a per extension status with an explicit reason field and last good spec so rollback is a single click in the control plane

* **Secure secrets**  
  Keep credentials in Azure Key Vault or AWS Secrets Manager and mount them on demand through CSI rather than copying into cluster specs

