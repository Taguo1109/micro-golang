#!/usr/bin/env bash
set -e

# 先啟動 User Service
./usersvc &
# 再啟動 Order Service
./ordersvc &
# 啟動 Nginx
nginx -g "daemon off;" &> nginx.log
# 等待所有子程序
wait