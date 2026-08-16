package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

func ListTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task

	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, fmt.Sprintf("invalid body: %v", err), http.StatusBadRequest)
		return
	}

	if err := validateTask(newTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTask.ID = nextID
	nextID++
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseTaskID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, t := range tasks {
		if t.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(t)
			return
		}
	}

	http.Error(w, "task not found", http.StatusNotFound)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseTaskID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updatedTask Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, fmt.Sprintf("invalid body: %v", err), http.StatusBadRequest)
		return
	}

	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Done = updatedTask.Done
			tasks[i].Title = updatedTask.Title

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}

	http.Error(w, "task not found", http.StatusNotFound)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseTaskID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, t := range tasks {
		if t.ID == id {
			tasks = slices.Delete(tasks, i, i+1)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "task not found", http.StatusNotFound)
}

// parseTaskID extracts and validates the "id" path parameter,
// shared by every handler that operates on a single task.
func parseTaskID(r *http.Request) (int, error) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return 0, errors.New("invalid id, must be a number")
	}
	return id, nil
}

func validateTask(t Task) error {
	if strings.TrimSpace(t.Title) == "" {
		return errors.New("title cannot be empty")
	}

	if len(t.Title) > 200 {
		return errors.New("title cannot be longer than 200 characters")
	}

	return nil
}
