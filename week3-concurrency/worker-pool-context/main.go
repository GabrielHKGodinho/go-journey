package main

import (
	"context"
	"fmt"
	"sync"
)

type Job struct {
	ID    int
	Value int
}

type Result struct {
	JobID  int
	Output int
}

func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			fmt.Printf("worker %d processing job %d\n", id, job.ID)
			results <- Result{JobID: job.ID, Output: job.Value * job.Value}

		case <-ctx.Done():
			fmt.Printf("worker %d stopping: %v\n", id, ctx.Err())
			return
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	jobs := make(chan Job)
	results := make(chan Result)
	var wg sync.WaitGroup

	for i := range 10 {
		wg.Add(1)
		go worker(ctx, i, jobs, results, &wg)
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
			cancel() // cancela depois de processar só 20 resultados, de propósito
		}
	}
}
