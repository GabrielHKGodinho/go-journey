package main

import "fmt"

// fibonacci retorna uma função que retorna
// o próximo número de Fibonacci a cada chamada.
func fibonacci() func() int {
	a, b := 0, 1

	return func() int {
		resultado := a
		a, b = b, a+b
		return resultado
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
