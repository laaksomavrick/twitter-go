# Build
FROM golang:1.12-alpine AS builder
WORKDIR /app
# https://github.com/docker-library/golang/issues/209
RUN apk add --no-cache git
COPY go.mod go.sum ./
COPY services/common services/common
COPY services/ready services/ready
RUN go build -ldflags="-s -w" -o app services/ready/cmd/main.go

# App
FROM alpine AS final
WORKDIR /app
COPY --from=builder /app /app
EXPOSE 3000
ENTRYPOINT ./app