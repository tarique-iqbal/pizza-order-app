FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY email-service/go.mod email-service/go.sum ./
RUN go mod download

COPY email-service/ ./email-service/

WORKDIR /app/email-service/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /email-service .

FROM alpine:latest
WORKDIR /app

COPY --from=builder /email-service .
COPY --from=builder /app/email-service/internal/infrastructure/email/templates ./internal/infrastructure/email/templates/

CMD ["./email-service"]
