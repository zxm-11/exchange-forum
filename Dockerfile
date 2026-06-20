# ============================================
# 阶段1: 编译
# ============================================
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /build

COPY exchangeapp/go.mod exchangeapp/go.sum ./
RUN go mod download

COPY exchangeapp/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /build/server .

# ============================================
# 阶段2: 运行
# ============================================
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/server .
COPY exchangeapp/config/config.yml ./config/config.yml

ENV GIN_MODE=release \
    TZ=Asia/Shanghai

RUN adduser -D -H -h /app appuser && \
    chown -R appuser:appuser /app
USER appuser

EXPOSE 3000

CMD ["./server"]
