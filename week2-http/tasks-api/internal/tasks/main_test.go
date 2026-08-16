package tasks

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListTasks(t *testing.T) {
	tests := []struct {
		name          string
		initialTasks  []Task
		expectedCount int
	}{
		{name: "empty list", initialTasks: []Task{}, expectedCount: 0},
		{name: "single task", initialTasks: []Task{{ID: 1, Title: "Task A"}}, expectedCount: 1},
		{name: "multiple tasks", initialTasks: []Task{{ID: 1, Title: "A"}, {ID: 2, Title: "B"}}, expectedCount: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks = tt.initialTasks
			nextID = len(tt.initialTasks) + 1

			request := httptest.NewRequest(http.MethodGet, "/tasks", nil)
			response := httptest.NewRecorder()

			ListTasks(response, request)

			var got []Task
			if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
				t.Fatalf("could not decode response body: %v", err)
			}

			if len(got) != tt.expectedCount {
				t.Errorf("expected %d tasks, got %d", tt.expectedCount, len(got))
			}
		})
	}
}

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantTitle  string // vazio = não checar corpo, só status
	}{
		{
			name:       "valid task",
			body:       `{"title":"Learn Go testing","done":false}`,
			wantStatus: http.StatusCreated,
			wantTitle:  "Learn Go testing",
		},
		{
			name:       "malformed json",
			body:       `{"title": "missing closing brace"`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "wrong type for done field",
			body:       `{"title":"Learn Go","done":"yes"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "unknown extra field",
			body:       `{"title":"Learn Go","done":false,"priority":"high"}`,
			wantStatus: http.StatusCreated,
			wantTitle:  "Learn Go",
		},
		{
			name:       "empty title fails validation",
			body:       `{"title":"","done":false}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks = []Task{}
			nextID = 1

			request := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(tt.body))
			response := httptest.NewRecorder()

			CreateTask(response, request)

			if response.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
			}

			if tt.wantTitle != "" {
				var got Task
				if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
					t.Fatalf("could not decode response body: %v", err)
				}
				if got.Title != tt.wantTitle {
					t.Errorf("expected title %q, got %q", tt.wantTitle, got.Title)
				}
				if got.ID == 0 {
					t.Errorf("expected a non-zero ID to be assigned")
				}
			}
		})
	}
}

func TestValidateTask(t *testing.T) {
	tests := []struct {
		name    string
		task    Task
		wantErr bool
	}{
		{name: "valid task", task: Task{Title: "Learn Go"}, wantErr: false},
		{name: "empty title", task: Task{Title: ""}, wantErr: true},
		{name: "whitespace only title", task: Task{Title: "   "}, wantErr: true},
		{name: "title too long", task: Task{Title: strings.Repeat("a", 201)}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTask(tt.task)

			gotErr := err != nil
			if gotErr != tt.wantErr {
				t.Errorf("validarTask(%+v): got error = %v, wantErr = %v", tt.task, gotErr, tt.wantErr)
			}
		})
	}
}

func TestCreateTaskAssignsIncrementingIDs(t *testing.T) {
	tasks = []Task{}
	nextID = 1

	firstRequest := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"First"}`))
	firstResponse := httptest.NewRecorder()
	CreateTask(firstResponse, firstRequest)

	secondRequest := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"Second"}`))
	secondResponse := httptest.NewRecorder()
	CreateTask(secondResponse, secondRequest)

	var first, second Task
	json.NewDecoder(firstResponse.Body).Decode(&first)
	json.NewDecoder(secondResponse.Body).Decode(&second)

	if second.ID <= first.ID {
		t.Errorf("expected second task ID (%d) to be greater than first (%d)", second.ID, first.ID)
	}
}
