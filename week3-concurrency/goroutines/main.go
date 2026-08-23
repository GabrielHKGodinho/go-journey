package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go printNumber(i, &wg)
	}

	wg.Wait()
	fmt.Println("all goroutines finished")
}

func printNumber(n int, wg *sync.WaitGroup) {
	// defer wg.Done()
	fmt.Println(n)
}
