FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source code
COPY backend/ ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Runtime stage
FROM alpine:latest

# Install runtime dependencies including LaTeX
RUN apk add --no-cache \
    texlive \
    texlive-xetex \
    texlive-latex-extra \
    texmf-dist-latexextra \
    texmf-dist-fontsextra \
    ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server /app/server

# Create directories
RUN mkdir -p /app/uploads /app/outputs

# Expose port
EXPOSE 8080

# Run
CMD ["/app/server"]
