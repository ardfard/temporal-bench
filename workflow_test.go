package main

import (
	"context"
	"testing"
)

func TestWorkflow(t *testing.T) {

	for i := 0; i < 10; i++ {
		for {
			res, _ := BenchmarkConditionalActivity(context.Background())
			t.Logf("rand: %v\n", res)
			if res {
				break
			}
		}
	}
}
