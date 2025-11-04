# Go Todoist MCP Service

This is a local microservice (MCP) that acts as a wrapper for the Todoist v2 REST API. It provides simple, local REST endpoints for agentic workflows.

---

## Prerequisites

1.  **Docker:** You must have Docker installed and running.
2.  **Go:** (Optional, for local testing) Go 1.21+
3.  **Todoist API Token:** A valid Todoist API token.

---

## Configuration

1.  Duplicate `config_defaults.yaml` and rename to `config.yaml' in this directory.
2.  Open the `config.yaml` file.
3.  Replace the placeholder `REPLACE_WITH_YOUR_TODOIST_API_TOKEN` with your actual Todoist API token.

```yaml
# config.yaml
todoist:
  api_token: "your-real-token-goes-here"
  base_url: "[https://api.todoist.com/rest/v2](https://api.todoist.com/rest/v2)"
server:
  port: "8080"
```

-----

## Running with Docker

These instructions will build the Docker image and run it as a detached container.

### 1\. Build the Image

From the root of the project directory (`~/projects/go-todoist-mcp`), run:

```bash
docker build -t go-todoist-mcp .
```

### 2\. Run the Container

This command will start the container, map port 8080 to your local machine, and—most importantly—mount your local `config.yaml` file into the container in read-only mode.

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/config.yaml:/app/config.yaml:ro \
  --name todoist-mcp \
  go-todoist-mcp
```

  * `-d`: Run in detached mode (in the background).
  * `-p 8080:8080`: Map your local port 8080 to the container's port 8080.
  * `-v $(pwd)/config.yaml:/app/config.yaml:ro`: Mounts your local config file into the container at `/app/config.yaml` (where the app expects it). The `:ro` makes it read-only.
  * `--name todoist-mcp`: Gives your container a memorable name.

-----

## Service Management

### Confirming Access

You can confirm the service is running by sending a `curl` request from your local machine:

```bash
curl http://localhost:8080/tasks
```

If it's working, you will receive a JSON array of your Todoist tasks (or an empty array `[]`).

### Viewing Logs

The service logs (including startup messages) are sent to `stdout` inside the container. You can view them using the `docker logs` command.

**To see all logs:**

```bash
docker logs todoist-mcp
```

**To follow the logs in real-time (like `tail -f`):**

```bash
docker logs -f todoist-mcp
```

### Stopping and Removing

To stop the service:

```bash
docker stop todoist-mcp
```

To remove the stopped container (e.g., to run it again):

```bash
docker rm todoist-mcp
```

-----

## API Endpoints

The service provides the following local endpoints:

  * **Create Task:** `POST http://localhost:8080/tasks`
  * **Get All Tasks:** `GET http://localhost:8080/tasks`
  * **Get One Task:** `GET http://localhost:8080/tasks/{id}`
  * **Update Task:** `PUT http://localhost:8080/tasks/{id}`
  * **Delete Task:** `DELETE http://localhost:8080/tasks/{id}`
