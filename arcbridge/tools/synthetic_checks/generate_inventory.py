#!/usr/bin/env python3
"""Generate synthetic cluster inventory payloads for testing APIs."""

import json
import random
import sys
from datetime import datetime

CLUSTERS = ["cluster-alpha", "cluster-beta", "cluster-gamma"]
EXTENSIONS = ["nginx-ingress", "logging", "policy"]

def generate_payload(cluster: str) -> dict:
    extensions = {name: f"{random.randint(0,2)}.{random.randint(0,9)}.{random.randint(0,9)}" for name in EXTENSIONS}
    return {
        "clusterID": cluster,
        "nodes": random.randint(1, 10),
        "pods": random.randint(20, 200),
        "extensions": extensions,
        "timestamp": datetime.utcnow().isoformat() + "Z",
    }


def main() -> None:
    payloads = [generate_payload(cluster) for cluster in CLUSTERS]
    json.dump(payloads, sys.stdout, indent=2)


if __name__ == "__main__":
    main()
