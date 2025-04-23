#!/usr/bin/env bash
set -e

# 啟動 Auth Service
./authsvc &

# 啟動 User Service
./usersvc &

# 啟動 Order Service
./ordersvc &

# 最後啟動 Nginx（前面都 background，nginx 用 foreground 方式）
nginx -g "daemon off;" &> nginx.log

# 等待所有子程序結束
wait