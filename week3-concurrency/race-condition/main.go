package main

import (
	"fmt"
	"sync"
)

// var counter int

// func main() {
// 	var wg sync.WaitGroup

// 	for i := 0; i < 1000; i++ {
// 		wg.Add(1)
// 		go increment(&wg)
// 	}

// 	wg.Wait()
// 	fmt.Println("final counter:", counter)
// }

// func increment(wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	counter++
// }

// ====================== FIXED VERSION =================================
var counter int
var mu sync.Mutex

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go increment(&wg)
	}

	wg.Wait()
	fmt.Println("final counter:", counter)
}

func increment(wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	counter++
	mu.Unlock()
}
