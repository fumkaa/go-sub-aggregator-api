FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o /app/sub-aggregator ./cmd/sub-aggregator/main.go

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/sub-aggregator .
COPY --from=builder /app/configs ./configs

CMD ["./sub-aggregator"]