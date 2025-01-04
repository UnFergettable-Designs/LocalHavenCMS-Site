FROM golang:1.22-alpine

WORKDIR /app

# Install SQLite dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

# Create data directory with proper permissions
RUN mkdir -p /app/data && \
    adduser -D -u 1000 appuser && \
    chown -R appuser:appuser /app && \
    chmod -R 755 /app && \
    chmod -R 777 /app/data

USER appuser

EXPOSE 8090

# Use shell form to ensure directory exists and has proper permissions
CMD sh -c "mkdir -p /app/data && chmod 777 /app/data && ./main"
