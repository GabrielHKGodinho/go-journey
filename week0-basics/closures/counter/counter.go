package main

import "fmt"

func criarContador() func() int {
	contagem := 0

	return func() int {
		contagem++
		return contagem
	}
}

func main() {
	funcao := criarContador
	proximo := funcao()
	fmt.Println(proximo()) // 1
	fmt.Println(proximo()) // 2
	fmt.Println(proximo()) // 3

	outroContador := criarContador()
	fmt.Println(outroContador()) // 1 — variável independente!
}
