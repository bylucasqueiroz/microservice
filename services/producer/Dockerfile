# Use an official Golang runtime as the base image
FROM golang:1.22-alpine AS builder

# Install GCC and development files for the C standard library
RUN apk --no-cache add ca-certificates gcc
RUN 

# Set necessary environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory in the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Kafka configs
RUN apk -U add ca-certificates
RUN apk update && apk upgrade && apk add pkgconf git bash build-base sudo
RUN git clone https://github.com/confluentinc/librdkafka.git && cd librdkafka && ./configure --prefix /usr && make && make install

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -tags musl --ldflags "-extldflags -static" -o main .

# Start a new stage from scratch
FROM alpine:latest

# Set necessary environment variables
ENV PORT=3000

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/main /app/main

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["/app/main"]