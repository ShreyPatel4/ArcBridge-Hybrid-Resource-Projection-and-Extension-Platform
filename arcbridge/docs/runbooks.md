# ArcBridge Runbooks (POC)

These runbooks describe the procedures for the most common incidents encountered in the proof-of-concept environment. They assume access to the local Kind-based fleet and placeholder telemetry.

## Agent Cannot Register
1. Verify the control plane is running: `make run-control-plane`.
2. Inspect control plane logs for `register` errors.
3. Confirm the bootstrap token in `deploy/kind/bootstrap-token.yaml` matches the agent config.
4. Regenerate local certificates with `tools/synthetic_checks/gen_cert.sh` and restart the agent deployment.

## Agent Fails to Install Extensions
1. Check the agent pod logs: `kubectl logs -n arcbridge agent -c agent`.
2. Run `kubectl get arcbridgeclusters -o yaml` to view the CR status and reason fields.
3. Retry reconciliation using `kubectl annotate arcbridgecluster <name> arcbridge.io/force-reconcile=$(date +%s)`.
4. If Helm install fails, run `helm template` locally with the same values file stored in `deploy/helm/charts/<extension>`.

## Control Plane Backlog High
1. View queue depth metrics via `telemetry/prom/rules/queue_depth.rules.yaml` or Grafana dashboard `telemetry/grafana/dashboards/fleet.json`.
2. Scale the control plane deployment replicas in `deploy/kustomize/controlplane/`.
3. Verify Redis is reachable using `tools/synthetic_checks/redis_ping.sh`.
4. Reduce rollout speed by lowering the token bucket rates in `controlplane/pkg/handlers/ratelimit_config.go`.

## Per-Tenant Saturation
1. Examine rate-limit counters via the `/metrics` endpoint.
2. Temporarily disable new rollouts by toggling the feature flag in `controlplane/pkg/storage/featureflags.go`.
3. Notify affected tenants using the contact list in `docs/api.md` (placeholder).

