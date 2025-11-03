# --- Build Stage ---
# Use the official Go image as the builder.
FROM golang:1.25.3-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app, creating a static binary.
# CGO_ENABLED=0 is crucial for static linking.
# -w -s strips debug symbols, reducing binary size.
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o go-todoist-mcp .

# --- Final Stage ---
# Use a minimal base image.
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the 'builder' stage.
COPY --from=builder /app/go-todoist-mcp .

# The config.yaml will be mounted as a volume at runtime.
# We do NOT copy it into the image to protect the API token.

# Expose port 8080 (which our service listens on)
EXPOSE 8080

# The command to run when the container starts.
ENTRYPOINT ["./go-todoist-mcp"]
