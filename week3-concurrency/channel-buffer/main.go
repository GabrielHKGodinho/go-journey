package main

import "fmt"

func main() {
	ch := make(chan int, 1)
	ch <- 42
	valor := <-ch
	fmt.Println(valor)
}
