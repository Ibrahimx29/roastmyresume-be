# Use Go image with Linux base
FROM golang:1.21

# Install Python
RUN apt-get update && apt-get install -y python3 python3-pip

# Set working directory
WORKDIR /app

# Copy everything to container
COPY . .

# Build Go app
RUN go build -o app

# Run the Go binary
CMD ["./app"]
