package tasks

import "net/http"

func SetupRoutes(store *TaskStore) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", store.ListTasks)
	mux.HandleFunc("POST /tasks", store.CreateTask)
	mux.HandleFunc("GET /tasks/{id}", store.GetTask)
	mux.HandleFunc("PUT /tasks/{id}", store.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", store.DeleteTask)

	return mux
}
