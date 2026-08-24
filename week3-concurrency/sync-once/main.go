package main

import (
	"fmt"
	"sync"
)

var once sync.Once

func expensiveSetup() {
	fmt.Println("running expensive setup...")
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			once.Do(expensiveSetup)
			fmt.Printf("goroutine %d done\n", id)
		}(i)
	}

	wg.Wait()
}
