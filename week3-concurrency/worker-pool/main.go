package main

import (
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

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("worker %d processing job %d\n", id, job.ID)
		results <- Result{JobID: job.ID, Output: job.Value * job.Value}
	}
}

func main() {
	jobs := make(chan Job)
	results := make(chan Result)
	var wg sync.WaitGroup

	for i := range 10 {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	go func() {
		defer close(jobs)
		for i := range 100 {
			jobs <- Job{i, i}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Printf("result of job %d is %d\n", result.JobID, result.Output)
	}
}
