FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .

RUN go build -o /app/migrator ./cmd/migrator/main.go

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/migrator .
COPY --from=builder /app/migrations ./migrations

ENTRYPOINT ["./migrator"]