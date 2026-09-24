FROM golang:alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kron ./cmd/kron

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/kron /app/

EXPOSE 5000

ENTRYPOINT ["./kron"]