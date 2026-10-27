FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o payment-svc ./cmd

FROM scratch

WORKDIR /app

COPY --from=builder /app/payment-svc .
COPY --from=builder /app/internal/db/migrations ./internal/db/migrations

EXPOSE 8888

ENTRYPOINT [ "./payment-svc" ]
