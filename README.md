# ArcBridge-Hybrid-Resource-Projection-and-Extension-Platform
Manage many Kubernetes clusters anywhere and project their state into a central Azure like control plane with secure registration, extension lifecycle, and production grade observability.



Goals
• Provide a simple way to connect any cluster to a central control plane and see it as a first class resource
• Let platform teams declare extensions once and roll them out safely across many clusters with drift detection and rollback
• Offer strong health signals and runbooks so a DRI can restore service quickly

Architecture overview
• Agent inside each customer cluster watches for a small CRD named ArcBridgeCluster that holds registration info and desired extensions
• Control plane service exposes REST APIs to issue tokens, accept heartbeats, receive inventory, and publish desired state per cluster
• Reconciler in the agent installs extensions using Helm and Kustomize, verifies readiness, and reports detailed status with reasons and last transition time
• Observability layer uses OpenTelemetry traces and Prometheus metrics emitted by both agent and control plane so you can follow a request across the system
• Identity and security rely on Azure AD OIDC for user and service auth, per tenant namespaces, and mTLS between agents and control plane

Key components
• Agent written in Go with client go informers, work queues, and a backoff aware reconcile loop
• Control plane written in Go with Gin or Chi for REST and a small gRPC side channel for streaming heartbeats and inventory
• Storage Postgres for cluster and extension state Redis for hot queues and rate limiting
• GitOps optional path with Flux where the control plane writes desired state to a Git repo and agents pull and apply
• Policy Open Policy Agent admission checks for extension parameters and namespace safety

Interesting problems and solutions
• Duplicate registration and stale tokens
Implement a one time registration flow with a short lived bootstrap token that swaps to a long lived client cert tied to cluster UID
• Large fleet rollout noise and thundering herds
Use a token bucket per region and per tenant to pace extension installs and a global rate cap to protect the control plane
• Partial failure and drift
Record a per extension status with an explicit reason field and last good spec so rollback is a single click in the control plane
• Secure secrets
Keep credentials in Azure Key Vault or AWS Secrets Manager and mount them on demand through CSI rather than copying into cluster specs

APIs and contracts
• POST register returns a signed client certificate and a cluster id
• PUT inventory accepts nodes pods and extension version info with etags to prevent races
• GET desired returns target extension specs scoped to cluster and tenant
• POST status streams reconcile progress as small events suitable for real time dashboards

Telemetry and SRE practices
• Metrics request rate error rate latency histograms queue depths reconcile attempts retries token issuance failures
• Traces wrap registration inventory fetch apply and health checks so on call can follow a single cluster through the pipeline
• Alerts p99 registration time above three minutes extension error rate above two percent heartbeat gap over two minutes reconcile retry storms
• Runbooks agent cannot register agent repeatedly fails install control plane backlog high per tenant saturation

Validation and results
• Local fleet of fifty clusters using Kind across five virtual regions registers in under two minutes median with full inventory in under four minutes
• Rolling out a basic NGINX ingress and a logging extension across fifty clusters completes in about eight minutes with less than one percent failure rate and automatic retries
• Fault tests include network partitions certificate expiry and bad chart versions with clean rollback in each case
• Load tests with K6 show the control plane sustains three thousand register or status events per second on a small three node cluster

Repository layout
arcbridge
docs design md runbooks md api md
agent cmd controller pkg reconcile pkg informers pkg security
controlplane cmd apiserver pkg handlers pkg storage pkg grpc
telemetry otel exporters prom rules grafana dashboards
deploy kind tilt kustomize helm charts
tools loadgen synthetic checks trace viewer

Academic tie in
This project grew from an academic study on portable control planes for edge and data center clusters. The study compared pull based and push based reconciliation, evaluated informer cache consistency, and measured the tradeoff between fine grained and batched status updates. ArcBridge uses a pull based approach with batched status to lower control plane load while keeping operator feedback responsive.
