# ==================== builder ====================
FROM golang:1.23.0-alpine AS builder
WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build cmd/podlogger/main.go

# ==================== runner ====================
FROM alpine:latest AS runner
WORKDIR /app

COPY --from=builder /app/main ./main

CMD ["./main"]
