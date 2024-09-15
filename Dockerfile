FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY ./backend/go.mod ./backend/go.sum ./

RUN go mod download

COPY backend/ .

RUN go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./main"]
