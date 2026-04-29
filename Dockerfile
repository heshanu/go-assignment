FROM golang:1.19 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Copy the books.json file
COPY books.json /app/books.json
RUN CGO_ENABLED=0 GOOS=linux go build -o main .
FROM alpine:latest

# Set the working directory
WORKDIR /root/
COPY --from=builder /app/main .
# Copy the books.json file to the final image
COPY --from=builder /app/books.json /root/books.json
# Expose the port the app runs on
EXPOSE 8081
CMD ["./main"]
