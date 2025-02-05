FROM golang:1.21.3-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o stresstest

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/stresstest .

ENTRYPOINT ["/app/stresstest"] 