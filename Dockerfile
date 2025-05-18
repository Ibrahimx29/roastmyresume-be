# Build stage
FROM golang:1.21 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app

# Final stage
FROM python:3.10-slim

WORKDIR /app

# Install Python dependency for PyMuPDF
RUN pip install --no-cache-dir PyMuPDF

# Copy Go binary and Python script
COPY --from=build /app/app /app/app
COPY extract_text.py /app/extract_text.py

# Optional: copy other dirs
COPY handlers/ /app/handlers/
COPY utils/ /app/utils/
COPY tmp/ /app/tmp/

# Add permissions if needed (useful on Railway)
RUN chmod +x /app/app

EXPOSE 8080

CMD ["./app"]
