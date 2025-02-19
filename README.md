# Temporal Server Benchmark

This project provides a simple benchmark tool for Temporal server performance testing.

## Prerequisites

- Go 1.21 or later
- Docker
- Running Temporal server in Kubernetes
- kubectl configured to access your cluster

## Running the Benchmark in Kubernetes

1. Build and deploy the benchmark:

```bash
chmod +x deploy-benchmark.sh
./deploy-benchmark.sh
```

2. Check the job status:

```bash
kubectl get jobs -n benchmark
```

3. View the benchmark results:

```bash
kubectl logs job/temporal-benchmark -n benchmark
```

## Local Development

If you want to run the benchmark locally:

1. Install dependencies:

```bash
go mod tidy
```

2. Run the benchmark:

```bash
go run . --temporal-address=localhost:7233 --namespace=default --workers=10 --iterations=10 --num-workflows-per-iter
```

## Configuration

Available flags:

- `--temporal-address`: Temporal server address (default: localhost:7233)
- `--namespace`: Temporal namespace (default: benchmark)
- `--workers`: Number of concurrent workers (default: 10)
- `--num-workflows-per-iter`: Number of workflow executions (default: 100)
- '--iterations': Number of iterations (default: 10)

## Benchmark Results

The tool will output:

- Total execution time
- Average execution time per workflow
