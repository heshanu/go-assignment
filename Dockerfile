# Use the official Golang image to create a build artifact.
FROM golang:1.19 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Copy the books.json file
COPY books.json /app/books.json

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a minimal base image for the final stage
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy the books.json file to the final image
COPY --from=builder /app/books.json /root/books.json

# Expose the port the app runs on
EXPOSE 8081

# Command to run the executable
CMD ["./main"]
