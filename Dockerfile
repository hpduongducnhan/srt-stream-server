# Stage 1: Build Go application
FROM ubuntu:20.04 AS builder

ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /build

# Install Go and build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    wget git ca-certificates gcc libc6-dev \
    && rm -rf /var/lib/apt/lists/*

# Install Go 1.25
RUN wget -q https://go.dev/dl/go1.25.0.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.25.0.linux-amd64.tar.gz \
    && rm go1.25.0.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

# Copy go mod files first for better caching
COPY mgmt/go.mod mgmt/go.sum* ./
RUN go mod download

# Copy source code
COPY mgmt/ ./

# Build the application with CGO enabled for sqlite3
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /build/mgmt-app .


# Stage 2: Final image based on ossrs/srs:6
FROM ossrs/srs:6

LABEL maintainer="NhanDD"
LABEL description="StreamServer - SRS + Management App with Supervisord"

# Install supervisord and other utilities
USER root
RUN apt-get update && apt-get install -y --no-install-recommends \
    supervisor \
    curl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean

# Create directories
RUN mkdir -p /var/log/supervisor /app /app/data /etc/supervisor/conf.d

# Copy the built Go application
COPY --from=builder /build/mgmt-app /app/mgmt-app
RUN chmod +x /app/mgmt-app

# Copy SRT configuration (can be overridden by volume mount)
COPY srt.conf /usr/local/srs/conf/srt.conf

# Copy supervisord configuration
COPY supervisord.conf /etc/supervisor/supervisord.conf

# ============================================
# Default Environment Variables
# ============================================
# Discord Bot Token (build-time default, can be overridden at runtime)

# SRS API endpoint (internal communication)
ENV SRS_API_URL="http://127.0.0.1:1985"

# Database path
ENV NDD_DB_PATH="/app/data/stream-server.db"

# Candidate IP for WebRTC (0.0.0.0 means auto-detect)
ENV CANDIDATE="0.0.0.0"

# ============================================
# Expose ports
# ============================================
# RTMP port
EXPOSE 1935
# HTTP API port (SRS)
EXPOSE 1985
# HTTP Server port (HTTP-FLV, HLS)
EXPOSE 8080
# SRT port (UDP)
EXPOSE 10080/udp
# Management API port
EXPOSE 10081

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD curl -f http://127.0.0.1:1985/api/v1/versions || exit 1

# Start supervisord
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]
