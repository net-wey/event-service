# Этап 1: Сборка приложения
FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app main.go

# Этап 2: Запуск приложения
FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app .

EXPOSE 8080

CMD ["./app"]
