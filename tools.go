package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// CreateTaskParams defines the input schema for the create_task tool.
type CreateTaskParams struct {
	Content     string   `json:"content"`
	Description string   `json:"description,omitempty"`
	ProjectID   string   `json:"project_id,omitempty"`
	SectionID   string   `json:"section_id,omitempty"`
	ParentID    string   `json:"parent_id,omitempty"`
	Order       int      `json:"order,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	Priority    int      `json:"priority,omitempty"`
	DueDate     string   `json:"due_date,omitempty"`
	DueLang     string   `json:"due_lang,omitempty"`
	AssigneeID  string   `json:"assignee_id,omitempty"`
}

// GetTaskParams defines the input schema for the get_task tool.
type GetTaskParams struct {
	ID string `json:"id"`
}

// UpdateTaskParams defines the input schema for the update_task tool.
type UpdateTaskParams struct {
	ID          string   `json:"id"`
	Content     *string  `json:"content,omitempty"`
	Description *string  `json:"description,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	Priority    *int     `json:"priority,omitempty"`
	DueDate     *string  `json:"due_date,omitempty"`
	DueLang     *string  `json:"due_lang,omitempty"`
	AssigneeID  *string  `json:"assignee_id,omitempty"`
}

// DeleteTaskParams defines the input schema for the delete_task tool.
type DeleteTaskParams struct {
	ID string `json:"id"`
}

func getTools(s *server.MCPServer, client *APIClient) []server.ServerTool {
	return []server.ServerTool{
		{
			Tool: mcp.Tool{
				Name:        "create_task",
				Description: "Creates a new task in Todoist.",
			},
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				args := request.GetArguments()

				var p CreateTaskParams
				if v, ok := args["content"].(string); ok {
					p.Content = v
				}
				if v, ok := args["description"].(string); ok {
					p.Description = v
				}
				if v, ok := args["project_id"].(string); ok {
					p.ProjectID = v
				}
				if v, ok := args["section_id"].(string); ok {
					p.SectionID = v
				}
				if v, ok := args["parent_id"].(string); ok {
					p.ParentID = v
				}
				if v, ok := args["order"].(float64); ok {
					p.Order = int(v)
				}
				if v, ok := args["labels"].([]interface{}); ok {
					for _, li := range v {
						if s, ok := li.(string); ok {
							p.Labels = append(p.Labels, s)
						}
					}
				}
				if v, ok := args["priority"].(float64); ok {
					p.Priority = int(v)
				}
				if v, ok := args["due_date"].(string); ok {
					p.DueDate = v
				}
				if v, ok := args["due_lang"].(string); ok {
					p.DueLang = v
				}
				if v, ok := args["assignee_id"].(string); ok {
					p.AssigneeID = v
				}

				task, err := client.createTask(p)
				if err != nil {
					return mcp.NewToolResultErrorFromErr("failed to create task", err), nil
				}

				return mcp.NewToolResultText(fmt.Sprintf("Successfully created task with ID: %s", task.ID)), nil
			},
		},
		{
			Tool: mcp.Tool{Name: "get_all_tasks",
				Description: "Gets all tasks from Todoist.",
			},
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return mcp.NewToolResultText("get_all_tasks not implemented"), nil
			},
		},
		{
			Tool: mcp.Tool{
				Name:        "get_task",
				Description: "Gets a single task from Todoist by its ID.",
			},
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return mcp.NewToolResultText("get_task not implemented"), nil
			},
		},
		{
			Tool: mcp.Tool{
				Name:        "update_task",
				Description: "Updates a task in Todoist.",
			},
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return mcp.NewToolResultText("update_task not implemented"), nil
			},
		},
		{
			Tool: mcp.Tool{
				Name:        "delete_task",
				Description: "Deletes a task from Todoist.",
			},
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return mcp.NewToolResultText("delete_task not implemented"), nil
			},
		},
	}
}
