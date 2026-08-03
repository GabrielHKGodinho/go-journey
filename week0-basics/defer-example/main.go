package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	//abre o arquivo
	file, err := os.Open("../painel-diario-go.html")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}

	//fecha o arquivo no lifo
	defer file.Close()

	//le cada linha
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
	}
}
