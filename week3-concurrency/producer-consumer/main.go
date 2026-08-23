package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 3)

	go producer(ch)
	consumer(ch)
}

func producer(ch chan int) {
	for i := 0; i < 5; i++ {
		fmt.Println("producing:", i)
		ch <- i
	}
	// close(ch)
}

func consumer(ch chan int) {
	timeout := time.NewTimer(2 * time.Second)
	defer timeout.Stop()

	for {
		select {
		case value, ok := <-ch:
			if !ok {
				return // channel fechado, produtor terminou
			}

			if !timeout.Stop() {
				<-timeout.C
			}
			fmt.Println("consumed:", value)
			timeout.Reset(2 * time.Second)

		case <-timeout.C:
			fmt.Println("timeout: no value received in 2 seconds")
			return
		}
	}
}
