package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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

// doRequest is a helper to create, authenticate, and send HTTP requests.
func (c *APIClient) doRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("Content-Type", "application/json")

	return c.client.Do(req)
}

// createTask makes a request to the Todoist API to create a new task.
func (c *APIClient) createTask(params CreateTaskParams) (Task, error) {
	var task Task

	jsonBody, err := json.Marshal(params)
	if err != nil {
		return task, err
	}

	resp, err := c.doRequest(http.MethodPost, c.baseURL+"/tasks", bytes.NewBuffer(jsonBody))
	if err != nil {
		return task, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return task, err
	}

	return task, nil
}
