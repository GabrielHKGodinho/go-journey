package main

import "net/http"

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", listarTasks)
	mux.HandleFunc("POST /tasks", criarTask)
	mux.HandleFunc("GET /tasks/{id}", buscarTask)
	mux.HandleFunc("PUT /tasks/{id}", atualizarTask)
	mux.HandleFunc("DELETE /tasks/{id}", excluirTask)

	return mux
}
