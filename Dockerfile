# Build Stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o web-app-analyzer ./cmd/app/main.go

# Final Stage
FROM alpine:3.18

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/web-app-analyzer .

# ✅ Copy templates directory so it is available in the container
COPY --from=builder /app/internal/templates /internal/templates

EXPOSE 8080

CMD ["/web-app-analyzer"]