package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		go printNumber(i)
	}

	time.Sleep(100 * time.Millisecond) // vou explicar essa linha logo abaixo
}

func printNumber(n int) {
	fmt.Println(n)
}
