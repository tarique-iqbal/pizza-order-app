FROM golang:1.24-alpine AS builder

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app

COPY identity-service/go.mod identity-service/go.sum ./
RUN go mod download

COPY identity-service/ ./identity-service/

WORKDIR /app/identity-service/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /identity-service .

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache postgresql-client
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /identity-service .

EXPOSE 8080
CMD ["./identity-service"]
