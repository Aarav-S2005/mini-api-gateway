FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

WORKDIR /app/gateway

RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./cmd/server

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/gateway/gateway .

EXPOSE 3002

CMD ["./gateway"]