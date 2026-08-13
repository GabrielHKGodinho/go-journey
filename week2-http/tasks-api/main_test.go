package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListTasks(t *testing.T) {
	// reset the global state so this test doesn't depend on whatever
	// manual curl/Invoke-WebRequest calls left behind
	tasks = []Task{}
	nextID = 1

	request := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	response := httptest.NewRecorder()

	listarTasks(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var got []Task
	if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
		t.Fatalf("could not decode response body: %v", err)
	}

	if len(got) != 0 {
		t.Errorf("expected an empty task list, got %d tasks", len(got))
	}
}

func TestCreateTask(t *testing.T) {
	tasks = []Task{}
	nextID = 1

	requestBody := strings.NewReader(`{"title":"Learn Go testing","done":false}`)
	request := httptest.NewRequest(http.MethodPost, "/tasks", requestBody)
	response := httptest.NewRecorder()

	criarTask(response, request)

	if response.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, response.Code)
	}

	var got Task
	if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
		t.Fatalf("could not decode response body: %v", err)
	}

	if got.Title != "Learn Go testing" {
		t.Errorf("expected title %q, got %q", "Learn Go testing", got.Title)
	}

	if got.ID != 1 {
		t.Errorf("expected ID 1, got %d", got.ID)
	}

	if len(tasks) != 1 {
		t.Errorf("expected 1 task stored, got %d", len(tasks))
	}
}
