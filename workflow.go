package main

import (
	"context"
	"math/rand"
	"time"

	"go.temporal.io/sdk/workflow"
)

func BenchmarkWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)

	err := workflow.ExecuteActivity(ctx, BenchmarkActivity).Get(ctx, nil)
	if err != nil {
		return err
	}

	var result bool

	for {
		err = workflow.ExecuteActivity(ctx, BenchmarkConditionalActivity).Get(ctx, &result)
		if err != nil {
			return err
		}

		// log the result
		logger.Info("result", "result", result)

		if result {
			break
		}
	}
	return nil
}

func BenchmarkActivity(ctx context.Context) error {
	// Simulate some work
	time.Sleep(100 * time.Millisecond)
	return nil
}

func BenchmarkConditionalActivity(ctx context.Context) (bool, error) {

	time.Sleep(25 * time.Millisecond)
	// random number between 0 and 1
	random := rand.Float64()

	if random < 0.5 {
		return true, nil
	}

	return false, nil
}
