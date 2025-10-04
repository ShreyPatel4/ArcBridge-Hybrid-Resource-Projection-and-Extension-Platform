# ArcBridge API Contract (POC)

All endpoints are namespaced under `/api/v1`. Authentication is simulated through static headers in the POC environment.

## POST /api/v1/register
Request body:
```json
{
  "clusterUID": "uuid",
  "bootstrapToken": "string",
  "metadata": {
    "region": "west-europe",
    "tenant": "contoso"
  }
}
```

Response body:
```json
{
  "clusterID": "cluster-123",
  "certificatePEM": "-----BEGIN CERTIFICATE-----...",
  "expiresAt": "2024-12-31T00:00:00Z"
}
```

## PUT /api/v1/inventory
Accepts inventory snapshots using optimistic concurrency via ETags.

## GET /api/v1/desired
Returns desired extension specifications for the calling cluster, filtered by tenant.

## POST /api/v1/status
Streams reconcile progress events. The POC implementation buffers events in memory and logs them for observability exercises.

