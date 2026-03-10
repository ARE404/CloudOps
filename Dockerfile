# Stage 1: Build Go binary
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .
RUN swag init -g main.go --output docs
RUN CGO_ENABLED=0 GOOS=linux go build -o cloudops main.go

# Stage 2: Minimal runtime image
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/cloudops .
COPY --from=builder /app/config ./config

EXPOSE 8080
CMD ["./cloudops"]
