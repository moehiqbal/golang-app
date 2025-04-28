FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Use a minimal alpine image for the final container
FROM alpine:latest

# Add ca-certificates for any HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./app"]