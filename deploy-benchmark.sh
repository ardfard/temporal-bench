#!/bin/bash
set -e

# build the binary
CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -o temporal-benchmark .

# Build the Docker image
docker build -t ardfard/temporal-benchmark:latest .

# Push the docker image
docker push ardfard/temporal-benchmark:latest

# Create the benchmark namespace if it doesn't exist
kubectl apply -f k8s/benchmark-job.yaml

# Delete the previous job if it exists
kubectl delete job temporal-benchmark -n benchmark --ignore-not-found

# Create the new job
kubectl apply -f k8s/benchmark-job.yaml

# Watch the job progress
echo "Watching job logs..."
kubectl wait --for=condition=complete job/temporal-benchmark -n benchmark --timeout=300s
kubectl logs -f job/temporal-benchmark -n benchmark 
