package main

import (
	"fmt"
	"sync"
	"time"
)

func source(name string, delay time.Duration, count int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for i := 0; i < count; i++ {
			time.Sleep(delay)
			out <- fmt.Sprintf("%s-item-%d", name, i)
		}
	}()
	return out
}

func merge(channels ...<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan string) {
			defer wg.Done()
			for value := range c {
				out <- value
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	ch1 := source("sensor-A", 300*time.Millisecond, 3)
	ch2 := source("sensor-B", 500*time.Millisecond, 3)
	ch3 := source("sensor-C", 200*time.Millisecond, 3)

	merged := merge(ch1, ch2, ch3)

	for value := range merged {
		fmt.Println("received:", value)
	}
}
