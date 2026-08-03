package main

import "fmt"

func main() {
	s := make([]int, 0)
	for i := 0; i < 10; i++ {
		s = append(s, i)
		fmt.Printf("len=%d cap=%d\n", len(s), cap(s))
	}
}
