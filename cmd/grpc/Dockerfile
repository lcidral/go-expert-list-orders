FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o grpc-server ./cmd/grpc/main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/grpc-server .
EXPOSE 50051
ENTRYPOINT ["./grpc-server"]