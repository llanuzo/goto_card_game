# STAGE 1: Build the binary
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/card-game

# STAGE 2: Runtime
FROM alpine:latest
RUN adduser -D appuser
USER appuser
WORKDIR /app
COPY ./config.yml ./config.yml
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
