# Build stage
FROM --platform=$BUILDPLATFORM golang:alpine AS builder

ARG TARGETOS
ARG TARGETARCH

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application for target platform
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o server ./cmd/server

# Final stage
FROM alpine:3.21

RUN apk add --no-cache git openssh-client ca-certificates tzdata && \
    adduser -D -h /app app && \
    chmod 1777 /tmp

WORKDIR /app

# Copy the binary from builder
COPY --from=builder --chown=app:app /app/server .

USER app

# Expose HTTP port
EXPOSE 3000

# Run the application
ENTRYPOINT ["./server"]
