# =========================
# 1. Build stage
# =========================
FROM golang:1.25-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    git \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN --mount=type=secret,id=GITHUB_TOKEN \
    git config --global url."https://x-access-token:$(cat /run/secrets/GITHUB_TOKEN)@github.com/".insteadOf \
    "https://github.com/" && \
    go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -gcflags="all=-l" \
    -o api-gateway ./cmd/main.go


# =========================
# 2. Runtime stage
# =========================
FROM scratch

WORKDIR /app

COPY --from=builder /app/api-gateway /app/api-gateway

ENTRYPOINT ["/app/api-gateway"]