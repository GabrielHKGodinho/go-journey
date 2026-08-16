package tasks

import "net/http"

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", ListTasks)
	mux.HandleFunc("POST /tasks", CreateTask)
	mux.HandleFunc("GET /tasks/{id}", GetTask)
	mux.HandleFunc("PUT /tasks/{id}", UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", DeleteTask)

	return mux
}
