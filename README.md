# Temporal Server Benchmark

This project provides a simple benchmark tool for Temporal server performance testing.

## Prerequisites

- Go 1.21 or later
- Running Temporal server
- Docker (optional, for running Temporal server locally)

## Running Temporal Server Locally

If you don't have a Temporal server running, you can start one using Docker:

```bash
docker run --network host temporalio/temporal-server:latest
```

```bash
## Usage

1. Install dependencies:
```bash
go mod tidy
```

2. Run the benchmark:
```bash
go run . --temporal-address=localhost:7233 --namespace=default --workers=10 --iterations=1000
```

Available flags:
- `--temporal-address`: Temporal server address (default: localhost:7233)
- `--namespace`: Temporal namespace (default: default)
- `--workers`: Number of concurrent workers (default: 10)
- `--iterations`: Number of workflow executions (default: 1000)

## Benchmark Results

The tool will output:
- Total execution time
- Average execution time per workflow