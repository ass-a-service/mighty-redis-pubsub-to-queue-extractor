# Stage 1 (Build)
FROM golang:1.16-alpine AS builder

RUN apk add --update --no-cache git make upx
WORKDIR /app/
COPY go.mod go.sum /app/
RUN go mod download
COPY . /app/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X github.com/ass-a-service/mighty-redis-pubsub-to-queue-extractor" \
    -v \
    -trimpath \
    -o redis_queue_extractor \
    main.go
RUN upx redis_queue_extractor
RUN echo "ID=\"distroless\"" > /etc/os-release

# Stage 2 (Final)
FROM gcr.io/distroless/static:latest
COPY --from=builder /etc/os-release /etc/os-release

COPY --from=builder /app/redis_queue_extractor /usr/bin/
CMD [ "/usr/bin/redis_queue_extractor" ]