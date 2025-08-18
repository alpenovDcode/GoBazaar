# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payment ./cmd/payment

# Final stage
FROM alpine:3.21
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/payment .
EXPOSE 8084
CMD ["./payment"] 