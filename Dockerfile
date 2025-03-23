# Use the official Golang image as the base image
FROM golang:1.19

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files (if you have a go.mod file)
COPY go.mod .
COPY go.sum .

# Download and install dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Command to run the application
CMD ["./main"]