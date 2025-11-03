package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// APIClient holds the http client and auth token.
type APIClient struct {
	client  *http.Client
	token   string
	baseURL string
}

// NewAPIClient creates a new client for our service.
func NewAPIClient(token, baseURL string) *APIClient {
	return &APIClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		token:   token,
		baseURL: baseURL,
	}
}

// main is the entry point for our service.
func main() {
	// Load configuration from config.yaml
	cfg := LoadConfig()

	client := NewAPIClient(cfg.Todoist.APIToken, cfg.Todoist.BaseURL)

	// Use http.ServeMux for a simple router.
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", client.handleTasks)       // Handles GET (all) and POST
	mux.HandleFunc("/tasks/", client.handleTaskByID) // Handles GET (one), PUT, DELETE

	listenAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting MCP service on http://localhost%s", listenAddr)
	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}

// handleTasks routes requests for the /tasks endpoint.
func (c *APIClient) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.getTasks(w, r)
	case http.MethodPost:
		c.createTask(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// handleTaskByID routes requests for /tasks/{id}.
func (c *APIClient) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if id == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		c.getTask(w, r, id)
	case http.MethodPut:
		c.updateTask(w, r, id) // Our MCP accepts PUT for updates
	case http.MethodDelete:
		c.deleteTask(w, r, id)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// --- CRUD Function Implementations ---

// getTasks (Read All) fetches all active tasks.
func (c *APIClient) getTasks(w http.ResponseWriter, r *http.Request) {
	resp, err := c.doRequest(http.MethodGet, c.baseURL+"/tasks", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	proxyResponse(w, resp)
}

// getTask (Read One) fetches a single task by its ID.
func (c *APIClient) getTask(w http.ResponseWriter, r *http.Request, id string) {
	resp, err := c.doRequest(http.MethodGet, c.baseURL+"/tasks/"+id, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	proxyResponse(w, resp)
}

// createTask (Create) creates a new task.
func (c *APIClient) createTask(w http.ResponseWriter, r *http.Request) {
	var reqBody CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	jsonBody, _ := json.Marshal(reqBody)
	resp, err := c.doRequest(http.MethodPost, c.baseURL+"/tasks", bytes.NewBuffer(jsonBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	proxyResponse(w, resp)
}

// updateTask (Update) updates an existing task.
// This translates our local PUT to Todoist's POST.
func (c *APIClient) updateTask(w http.ResponseWriter, r *http.Request, id string) {
	var reqBody UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	jsonBody, _ := json.Marshal(reqBody)
	// Note: Todoist API uses POST for updates.
	resp, err := c.doRequest(http.MethodPost, c.baseURL+"/tasks/"+id, bytes.NewBuffer(jsonBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	proxyResponse(w, resp)
}

// deleteTask (Delete) deletes a task.
func (c *APIClient) deleteTask(w http.ResponseWriter, r *http.Request, id string) {
	resp, err := c.doRequest(http.MethodDelete, c.baseURL+"/tasks/"+id, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// On success, Todoist returns 204 No Content.
	w.WriteHeader(resp.StatusCode)
}

// --- Helper Functions ---

// doRequest is a helper to create, authenticate, and send HTTP requests.
func (c *APIClient) doRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Add the authentication header
	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("Content-Type", "application/json")

	return c.client.Do(req)
}

// proxyResponse copies the status code, headers, and body from the
// Todoist API response directly to our service's response writer.
func proxyResponse(w http.ResponseWriter, resp *http.Response) {
	for k, v := range resp.Header {
		w.Header().Set(k, strings.Join(v, ","))
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
