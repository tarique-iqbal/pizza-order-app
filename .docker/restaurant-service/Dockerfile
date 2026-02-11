FROM golang:1.24-alpine AS builder

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app

COPY restaurant-service/go.mod restaurant-service/go.sum ./
RUN go mod download

COPY restaurant-service/ ./restaurant-service/

WORKDIR /app/restaurant-service/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /pizza-api .

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache postgresql-client
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /pizza-api .

EXPOSE 8080
CMD ["./pizza-api"]
