package main

import (
	"context"
	"fmt"
)

func leaky() {
	ch := make(chan int)
	go func() {
		value := <-ch // presa pra sempre, se nada escrever em ch
		fmt.Println(value)
	}()
}

func notLeaky(ctx context.Context) {
	ch := make(chan int)
	go func() {
		select {
		case value := <-ch:
			fmt.Println(value)
		case <-ctx.Done():
			return // a "porta de saída" que garante que a goroutine sempre pode terminar
		}
	}()
}

func firstToRespond(sources ...<-chan int) int {
	out := make(chan int)
	for _, src := range sources {
		go func(s <-chan int) {
			out <- <-s // manda o primeiro valor que essa fonte der
		}(src)
	}
	return <-out // só pega UM valor, e retorna
}

func main() {
	leaky()
}
