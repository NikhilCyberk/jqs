# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o jqs ./cmd/main.go

FROM alpine:latest
WORKDIR /app
RUN adduser -D appuser
USER appuser
COPY --from=builder /app/jqs .
EXPOSE 8080
ENV PORT=8080
CMD ["./jqs"] 