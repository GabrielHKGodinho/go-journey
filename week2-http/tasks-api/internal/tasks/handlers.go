package tasks

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (s *TaskStore) ListTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query("SELECT id, title, done FROM tasks")
	if err != nil {
		http.Error(w, fmt.Sprintf("could not query tasks: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	result := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done); err != nil {
			http.Error(w, fmt.Sprintf("could not scan task: %v", err), http.StatusInternalServerError)
			return
		}
		result = append(result, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *TaskStore) CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task

	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, fmt.Sprintf("invalid body: %v", err), http.StatusBadRequest)
		return
	}

	if err := validateTask(newTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	row := s.db.QueryRow(
		"INSERT INTO tasks (title, done) VALUES ($1, $2) RETURNING id",
		newTask.Title, newTask.Done,
	)

	if err := row.Scan(&newTask.ID); err != nil {
		http.Error(w, fmt.Sprintf("could not create task: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func (s *TaskStore) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseTaskID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var t Task
	row := s.db.QueryRow("SELECT id, title, done FROM tasks WHERE id = $1", id)

	if err := row.Scan(&t.ID, &t.Title, &t.Done); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("could not get task: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (s *TaskStore) UpdateTask(w http.ResponseWriter, r *http.Request) {
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

	if err := validateTask(updatedTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := s.db.Exec(
		"UPDATE tasks SET title = $1, done = $2 WHERE id = $3",
		updatedTask.Title, updatedTask.Done, id,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not update task: %v", err), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf("could not check update result: %v", err), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	updatedTask.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

func (s *TaskStore) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseTaskID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := s.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not delete task: %v", err), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf("could not check delete result: %v", err), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
