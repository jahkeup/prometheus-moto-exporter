# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang alpine image
FROM golang:alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go files
COPY go.mod go.sum ./

# Download and verify dependencies.
RUN go mod download && go mod verify

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the app
RUN go build -ldflags="-w -s" ./cmd/prometheus-moto-exporter

# Build a small image
FROM alpine

COPY --from=builder /app/prometheus-moto-exporter /go/bin/prometheus-moto-exporter

# Set default values
ENV MOTO_USERNAME=admin
ENV MOTO_PASSWORD=motorola

# Expose port 9731
EXPOSE 9731

# Run the executable
ENTRYPOINT ["/go/bin/prometheus-moto-exporter", "--bind", "0.0.0.0:9731"]
