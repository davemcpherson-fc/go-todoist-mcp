package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestCreateTaskHandler verifies the POST /tasks endpoint.
func TestCreateTaskHandler(t *testing.T) {
	// 1. Create a mock Todoist API server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the auth token was passed
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer fake-token" {
			t.Errorf("Expected auth header 'Bearer fake-token', got '%s'", authHeader)
		}

		// Check the method and URL
		if r.Method != http.MethodPost || r.URL.Path != "/tasks" {
			t.Errorf("Mock server expected POST /tasks, got %s %s", r.Method, r.URL.Path)
		}

		// Send back a mock task response
		w.WriteHeader(http.StatusCreated)
		mockTask := Task{
			ID:      "12345",
			Content: "Test task from mock",
		}
		json.NewEncoder(w).Encode(mockTask)
	}))
	defer mockServer.Close()

	// 2. Create our APIClient, pointing its baseURL to our mock server
	client := &APIClient{
		client:  mockServer.Client(), // Use the mock server's client
		token:   "fake-token",
		baseURL: mockServer.URL, // Point to the mock server
	}

	// 3. Craft a fake incoming request for our MCP service
	taskPayload := `{"content":"Test task from mock"}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(taskPayload))

	// 4. Use a ResponseRecorder to capture the handler's response
	rr := httptest.NewRecorder()

	// 5. Call the handler directly
	client.handleTasks(rr, req)

	// 6. Assert the results
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body
	var createdTask Task
	if err := json.NewDecoder(rr.Body).Decode(&createdTask); err != nil {
		t.Fatalf("Could not decode response body: %v", err)
	}

	if createdTask.ID != "12345" || createdTask.Content != "Test task from mock" {
		t.Errorf("Handler returned unexpected body: got %v", createdTask)
	}
}
