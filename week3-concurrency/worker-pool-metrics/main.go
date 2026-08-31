package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

type Job struct {
	ID    int
	Value int
}

type Result struct {
	JobID  int
	Output int
}

type Metrics struct {
	jobsProcessed     int64
	jobsCancelled     int64
	jobsChannelClosed int64
}

func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup, metrics *Metrics) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				atomic.AddInt64(&metrics.jobsChannelClosed, 1)
				return
			}
			fmt.Printf("worker %d processing job %d\n", id, job.ID)
			results <- Result{JobID: job.ID, Output: job.Value * job.Value}
			atomic.AddInt64(&metrics.jobsProcessed, 1)

		case <-ctx.Done():
			atomic.AddInt64(&metrics.jobsCancelled, 1)
			fmt.Printf("worker %d stopping: %v\n", id, ctx.Err())
			return
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	metrics := &Metrics{}

	jobs := make(chan Job)
	results := make(chan Result)
	var wg sync.WaitGroup

	for i := range 10 {
		wg.Add(1)
		go worker(ctx, i, jobs, results, &wg, metrics)
	}

	go func() {
		for i := range 100 {
			select {
			case jobs <- Job{i, i}:
			case <-ctx.Done():
				close(jobs)
				return
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	count := 0
	for result := range results {
		fmt.Printf("result of job %d is %d\n", result.JobID, result.Output)
		count++
		if count == 20 {
			cancel()
		}
	}

	fmt.Println("--- shutdown report ---")
	fmt.Printf("jobs processed: %d\n", atomic.LoadInt64(&metrics.jobsProcessed))
	fmt.Printf("jobs cancelled before starting: %d\n", atomic.LoadInt64(&metrics.jobsCancelled))
	fmt.Printf("workers found closed channel: %d\n", atomic.LoadInt64(&metrics.jobsChannelClosed))
}
