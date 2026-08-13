package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := setupRoutes()

	fmt.Println("Servidor rodando em :8080")
	http.ListenAndServe(":8080", mux)
}
