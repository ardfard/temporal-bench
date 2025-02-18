package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var (
	temporalAddress = flag.String("temporal-address", "localhost:7233", "Temporal server address")
	namespace       = flag.String("namespace", "benchmark", "Temporal namespace")
	workers         = flag.Int("workers", 10, "Number of concurrent workers")
	iterations      = flag.Int("iterations", 1000, "Number of workflow executions")
)

func main() {
	flag.Parse()

	// Create Temporal client
	c, err := client.Dial(client.Options{
		HostPort:  *temporalAddress,
		Namespace: *namespace,
	})
	if err != nil {
		log.Fatalf("Failed to create Temporal client: %v", err)
	}
	defer c.Close()

	// Start workers
	for i := 0; i < *workers; i++ {
		w := worker.New(c, "benchmark-taskqueue", worker.Options{})
		w.RegisterWorkflow(BenchmarkWorkflow)
		w.RegisterActivity(BenchmarkActivity)
		w.RegisterActivity(BenchmarkConditionalActivity)

		err = w.Start()
		if err != nil {
			log.Fatalf("Failed to start worker: %v", err)
		}
	}

	// Start benchmark
	wg := sync.WaitGroup{}

	wg.Add(*iterations)
	results := make([]time.Duration, *iterations)
	startTime := time.Now()

	for i := 0; i < *iterations; i++ {
		workflowOptions := client.StartWorkflowOptions{
			ID:        fmt.Sprintf("benchmark-workflow-%d", i),
			TaskQueue: "benchmark-taskqueue",
		}

		go func(i int) {
			defer wg.Done()
			startTime := time.Now()
			res, err := c.ExecuteWorkflow(context.Background(), workflowOptions, BenchmarkWorkflow)
			if err != nil {
				log.Printf("Failed to execute workflow: %v", err)
				return
			}

			err = res.Get(context.Background(), nil)
			if err != nil {
				log.Printf("Failed to get workflow result: %v", err)
				return
			}

			results[i] = time.Since(startTime)
		}(i)
	}
	wg.Wait()

	duration := time.Since(startTime)
	// sum the results
	total := time.Duration(0)
	for _, result := range results {
		total += result
	}

	log.Printf("Benchmark completed in %v", duration)
	log.Printf("Average execution time: %v", total/time.Duration(*iterations))
}
