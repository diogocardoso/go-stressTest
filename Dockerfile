FROM golang:1.21.3-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o loadtest

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/loadtest .

ENTRYPOINT ["/app/loadtest"] 