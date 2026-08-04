package main

import "fmt"

type Radio struct {
	Nome       string
	Tipo       string
	Frequencia float64
}

type Cliente struct {
	Nome string
	CNPJ string
}

type Pacote struct {
	Nome                string
	MateriasContratadas int
	Cliente             Cliente
}

func test() {
	radio := Radio{
		Nome:       "Radio Tupi",
		Tipo:       "AM",
		Frequencia: 97.90,
	}

	cliente := Cliente{
		Nome: "Apple",
		CNPJ: "102.341.983/0001",
	}

	pacote := Pacote{
		Nome:                "Verao",
		MateriasContratadas: 4,
		Cliente:             cliente,
	}

	fmt.Printf("%+v\n", radio)
	fmt.Println(cliente.CNPJ)
	fmt.Println(pacote.Cliente.Nome)
}
