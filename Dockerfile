# Use Go to build your app
FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

# Final image
FROM debian:bullseye

WORKDIR /app

# Copy Go binary
COPY --from=builder /app/main .

# Copy Python-built binary
COPY --from=builder /app/dist/extract_text ./extract_text

# Make sure it’s executable
RUN chmod +x ./extract_text

EXPOSE 8080
CMD ["./main"]
