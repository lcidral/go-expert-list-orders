FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o graphql-server ./cmd/graphql/main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/graphql-server .
EXPOSE 8081
ENTRYPOINT ["./graphql-server"]