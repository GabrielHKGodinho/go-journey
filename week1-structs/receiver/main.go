package main

import "fmt"

type ContaBancaria struct {
	saldo  int
	Numero int
	Banco  string
}

func (c ContaBancaria) Saldo() int {
	return c.saldo
}

func (c ContaBancaria) changeBancoValor(novoBanco string) {
	c.Banco = novoBanco
}

func (c *ContaBancaria) changeBancoReferencia(novoBanco string) {
	c.Banco = novoBanco
}

func main() {
	conta := ContaBancaria{
		saldo:  1000,
		Numero: 1415,
		Banco:  "Santander",
	}

	fmt.Println(conta.Banco)
	conta.changeBancoValor("Itaú")
	fmt.Println(conta.Banco)
	conta.changeBancoReferencia("Caixa")
	fmt.Println(conta.Banco)
}
