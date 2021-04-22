#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /app/src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a ./cmd/main.go

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/bin
COPY --from=builder /app/src/main .
COPY --from=builder /app/src/adapters/postgres/schema.sql .
CMD ["./main"]
