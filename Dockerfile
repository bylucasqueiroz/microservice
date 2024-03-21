# Use an official Golang runtime as the base image
FROM golang:1.22-alpine

# Set the working directory in the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go application
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 5000

# Command to run the executable
CMD ["./main"]