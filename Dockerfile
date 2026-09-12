
# === STAGE 1: BUILDER ===
# Use the official Go image to compile the application
FROM golang:1.26-alpine AS builder

# Set necessary environment variables for CGO dependencies if needed, though often not required for pure Go
ENV CGO_ENABLED=0

# Set the working directory inside the container
WORKDIR /app

# Copy dependency files first to leverage Docker caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application statically.
# -o sets the output filename (named 'api')
# -ldflags="-s -w" removes debugging information to keep the binary small
RUN go build -ldflags="-s -w" -o /api main.go


# === STAGE 2: FINAL RUNTIME IMAGE ===
# Use a minimal, secure image for the final deployment
FROM alpine:latest

# Install certificates if your DB uses SSL (good practice, though we disabled it for simplicity)
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /api .

# Expose the port the application listens on
EXPOSE 8080

# Define the command to run when the container starts
CMD ["./api"]