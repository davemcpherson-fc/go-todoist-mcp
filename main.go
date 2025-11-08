package main

import (
	"log"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Load configuration from config.yaml
	cfg := LoadConfig()

	// Create a new API client
	client := NewAPIClient(cfg.Todoist.APIToken, cfg.Todoist.BaseURL)

	// Create a new MCP server
	s := server.NewMCPServer(
		"Dave's ToDoIst MCP Server",
		"0.0.1",
		server.WithToolCapabilities(true),
	)

	// Register the tools
	s.AddTools(getTools(s, client)...)

	// Run the server
	log.Println("Starting MCP server with stdio transport")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server exited with error: %s", err)
	}
}
