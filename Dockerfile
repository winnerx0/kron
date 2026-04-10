FROM golang:alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN cd cmd/kron && CGO_ENABLED=0 GOOS=linux go build -o kron . && mv kron ../../

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/kron /app/

ENTRYPOINT ["./kron"]