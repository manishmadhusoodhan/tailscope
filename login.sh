#!/bin/sh
set -e

echo "Checking Kubernetes login..."

kubectl get ns >/dev/null

echo "Starting TailScope..."
docker compose up -d