FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files and download dependencies
COPY backend/go.mod ./
RUN go mod download || true

# Copy source code
COPY backend/ ./

# Tidy and download dependencies
RUN go mod tidy && go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Runtime stage
FROM alpine:latest

# Install runtime dependencies including LaTeX and PDF tools
RUN apk add --no-cache \
    texlive \
    texlive-xetex \
    texmf-dist-latexextra \
    texmf-dist-fontsextra \
    texmf-dist-langchinese \
    font-noto-cjk \
    poppler-utils \
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
