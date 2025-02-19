package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var (
	temporalAddress     = flag.String("temporal-address", "localhost:7233", "Temporal server address")
	namespace           = flag.String("namespace", "benchmark", "Temporal namespace")
	workers             = flag.Int("workers", 10, "Number of concurrent workers")
	numWorkflowsPerIter = flag.Int("num-workflows-per-iter", 100, "Number of workflow executions")
	iterations          = flag.Int("iterations", 1, "Number of iterations")
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
	for range *workers {
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

	total := time.Duration(0)

	for range *iterations {
		wg := sync.WaitGroup{}
		wg.Add(*numWorkflowsPerIter)
		results := make([]time.Duration, *numWorkflowsPerIter)
		startTime := time.Now()
		for i := range *numWorkflowsPerIter {
			ulid := ulid.Make()
			workflowOptions := client.StartWorkflowOptions{
				ID:        fmt.Sprintf("bench-%s", ulid.String()),
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
		// sum the results
		for _, result := range results {
			total += result
		}
		sleepTime := time.Second*60 - time.Since(startTime)
		if sleepTime > 0 {
			time.Sleep(sleepTime)
		}
	}

	log.Printf("Benchmark completed in %v", total)
	log.Printf("Average execution time: %v", total/time.Duration(*numWorkflowsPerIter**iterations))
}
