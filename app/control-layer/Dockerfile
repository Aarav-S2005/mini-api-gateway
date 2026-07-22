FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

WORKDIR /app/control-layer

RUN CGO_ENABLED=0 GOOS=linux go build -o control-layer ./cmd/server

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/control-layer/control-layer .

EXPOSE 3001

CMD ["./control-layer"]