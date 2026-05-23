# =========================
# 1. Build stage
# =========================
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=secret,id=GITHUB_TOKEN \
    git config --global url."https://x-access-token:$(cat /run/secrets/GITHUB_TOKEN)@github.com/".insteadOf \
    "https://github.com/" && \
    go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o api-gateway ./cmd/main.go


# =========================
# 2. Runtime stage
# =========================
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/api-gateway /app/api-gateway

ENTRYPOINT ["/app/api-gateway"]