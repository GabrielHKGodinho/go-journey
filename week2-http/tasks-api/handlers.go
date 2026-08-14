package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

func listarTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func criarTask(w http.ResponseWriter, r *http.Request) {
	var novaTask Task

	if err := json.NewDecoder(r.Body).Decode(&novaTask); err != nil {
		http.Error(w, fmt.Sprintf("corpo inválido: %v", err), http.StatusBadRequest)
		return
	}

	if err := validarTask(novaTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	novaTask.ID = nextID
	nextID++
	tasks = append(tasks, novaTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(novaTask)
}

func buscarTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id inválido, precisa ser um número", http.StatusBadRequest)
		return
	}

	for _, t := range tasks {
		if t.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(t)
			return
		}
	}

	http.Error(w, "task não encontrada", http.StatusNotFound)
}

func atualizarTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id inválido, precisa ser um número", http.StatusBadRequest)
		return
	}

	var novaTask Task
	if err := json.NewDecoder(r.Body).Decode(&novaTask); err != nil {
		http.Error(w, fmt.Sprintf("corpo inválido: %v", err), http.StatusBadRequest)
		return
	}

	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Done = novaTask.Done
			tasks[i].Title = novaTask.Title

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}

	http.Error(w, "task não encontrada", http.StatusNotFound)
}

func excluirTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id inválido, precisa ser um número", http.StatusBadRequest)
		return
	}

	for i, t := range tasks {
		if t.ID == id {
			tasks = slices.Delete(tasks, i, i+1)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "task não encontrada", http.StatusNotFound)
}

func validarTask(t Task) error {
	if strings.TrimSpace(t.Title) == "" {
		return errors.New("title não pode ser vazio")
	}

	if len(t.Title) > 200 {
		return errors.New("title não pode ter mais de 200 caracteres")
	}

	return nil
}
