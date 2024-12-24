FROM golang:1.21-alpine

WORKDIR /app

# Install SQLite dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

# Create data directory with proper permissions
RUN mkdir -p /app/data && \
    chown -R nobody:nobody /app && \
    chmod 777 /app/data

USER nobody

EXPOSE 8090

CMD ["./main"]
