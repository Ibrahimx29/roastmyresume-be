FROM golang:1.24-alpine

# Install Python and required dependencies
RUN apk add --no-cache python3 py3-pip

# Create symbolic link for python command
RUN ln -sf /usr/bin/python3 /usr/bin/python

# Install PyMuPDF
RUN pip3 install PyMuPDF

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the Go application
RUN go build -o main .

# Make sure the Python script is executable
RUN chmod +x ./extract_text.py

# Set the PORT environment variable
ENV PORT=8080

# Expose the port
EXPOSE 8080

# Command to run the application
CMD ["./main"]