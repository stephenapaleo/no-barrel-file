# Base image for Go
FROM golang:1.23 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go modules and install dependencies
COPY go.mod go.sum ./
COPY vendor/ ./vendor/
COPY Makefile ./Makefile
RUN make vendor

# Copy the entire project into the container
COPY . .

# Build the CLI using the Makefile
RUN make build

# Final stage: minimal runtime image
FROM alpine:latest

# Set up a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory inside the runtime container
WORKDIR /app

# Copy the built CLI binary from the builder stage
COPY --from=builder /app/bin/linux-amd64/no-barrel-file /usr/local/bin/no-barrel-file

# Change ownership and grant execution permissions
RUN chown appuser:appgroup /usr/local/bin/no-barrel-file && chmod +x /usr/local/bin/no-barrel-file

# Switch to the non-root user
USER appuser

# Set the default command to display the help message
ENTRYPOINT ["/usr/local/bin/no-barrel-file"]
