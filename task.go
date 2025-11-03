package main

// Task represents the core fields of a Todoist task.
// json tags map Go struct fields to the API's JSON keys.
type Task struct {
	ID          string `json:"id"`
	ProjectID   string `json:"project_id"`
	Content     string `json:"content"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
	Priority    int    `json:"priority"`
	Due         *struct {
		String string `json:"string"`
		Date   string `json:"date"`
	} `json:"due,omitempty"`
}

// CreateTaskRequest defines the payload our MCP service accepts
// for creating a new task.
type CreateTaskRequest struct {
	Content     string `json:"content"`
	ProjectID   string `json:"project_id,omitempty"`
	DueString   string `json:"due_string,omitempty"`
	Description string `json:"description,omitempty"`
	Priority    int    `json:"priority,omitempty"`
}

// UpdateTaskRequest defines the payload our MCP service accepts
// for updating an existing task.
type UpdateTaskRequest struct {
	Content     string `json:"content,omitempty"`
	Description string `json:"description,omitempty"`
	DueString   string `json:"due_string,omitempty"`
	Priority    int    `json:"priority,omitempty"`
}
