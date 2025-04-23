# ── build stage ──
FROM golang:1.23.8-alpine AS builder

WORKDIR /app

# 拷貝依賴並下載
COPY go.mod go.sum ./
RUN go mod download

# 拷貝整个項目並編譯
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o usersvc cmd/usersvc/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o ordersvc cmd/ordersvc/main.go

# ── final stage ──
FROM alpine:3.17

# 安裝 bash、nginx
RUN apk add --no-cache bash nginx

WORKDIR /app

# 拷貝編譯好的兩支 binary 和腳本
COPY --from=builder /app/usersvc .
COPY --from=builder /app/ordersvc .
COPY run.sh .
COPY nginx.conf /etc/nginx/nginx.conf

# 授權腳本執行
RUN chmod +x run.sh

# 暴露 HTTP 80 （nginx 監聽），内部 services 8000/9000 不必對外
EXPOSE 80 8080

# 啟動腳本
ENTRYPOINT ["./run.sh"]