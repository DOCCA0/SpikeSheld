# Backend Dockerfile
FROM golang:1.21-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download || true

# Copy source code
COPY . .

# Build the application
RUN go build -o main .

# Run the application
CMD ["./main"]
