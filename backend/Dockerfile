FROM golang:1.21-alpine

WORKDIR /app

# Install SQLite dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

# Create data directory and set permissions
RUN mkdir -p /app/data
RUN chown -R nobody:nobody /app/data

USER nobody

EXPOSE 8090

CMD ["./main"]
