# =============================================================================
# Stage 1: Build the Go binary
# =============================================================================
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /src

# Cache dependency downloads.
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /bin/incharge ./cmd/server/

# =============================================================================
# Stage 2: Minimal runtime image
# =============================================================================
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

# Create a non-root user.
RUN addgroup -S incharge && adduser -S incharge -G incharge

WORKDIR /app

# Copy binary and templates from builder.
COPY --from=builder /bin/incharge .
COPY --from=builder /src/templates ./templates

# Use non-root user.
USER incharge

EXPOSE 8080

ENTRYPOINT ["./incharge"]
