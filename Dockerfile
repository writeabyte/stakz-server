# Stage 1: Build the Go binary
FROM golang:alpine AS builder

WORKDIR /app

# Copy the Go source code into the container
COPY main.go .
COPY go.mod .

# Build the Go binary
RUN go build -o myapp

# Stage 2: Create a minimal image with only the binary
FROM alpine

WORKDIR /app

# Copy the binary from the previous stage
COPY --from=builder /app/myapp .

# Expose the port your application will run on
EXPOSE 3001

# Command to run your application
CMD ["./myapp"]

