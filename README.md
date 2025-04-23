# Go 微服務範例專案

一個結合 **Gin**、**Nginx** 反向代理與 **Docker** 的 Go 微服務示範，包含兩個獨立服務：

| 服務名稱    | 說明                             | 預設埠號 |
|-----------|----------------------------------|--------|
| **usersvc**  | 提供使用者資料 API (`/users/:id`)  | 8000   |
| **ordersvc** | 提供訂單資料 API (`/orders/:id`) | 9000   |
| **nginx**    | 反向代理並統一對外監聽            |(8080) |

## 架構圖
```
Client ↔ Nginx(80/8080) ↔ { usersvc(8000), ordersvc(9000) }
```

---
## 目錄結構
```
├── cmd/
│   ├── usersvc/           # User Service 主程式
│   └── ordersvc/          # Order Service 主程式
├── internal/              # 服務共用程式碼 (Handler、Service)
├── run.sh                 # 啟動兩個服務與 Nginx
├── nginx.conf             # Nginx 反向代理設定
├── Dockerfile             # 多階段 Build + 最終 Image
└── README.md              # 本檔案
``` 

---
## 本機開發 & 測試
1. **Go 環境**：需安裝 Go 1.23+。
2. **執行 User Service**：
   ```bash
   cd cmd/usersvc
   go run main.go            # 監聽 :8000
   ```
3. **執行 Order Service**：
   ```bash
   cd cmd/ordersvc
   go run main.go            # 監聽 :9000，並呼叫 http://localhost:8000/users/123
   ```
4. **測試**：
   ```bash
   curl http://localhost:8000/users/123
   curl http://localhost:9000/orders/abc
   ```

---
## Docker 化部署
1. **Build Image**：
   ```bash
   docker build -t go-microservices .
   ```
2. **Run Container**：
   ```bash
   docker run -d \
     --name micro-svcs \
     -p 8080:8080  # (若需映射平台預設埠) \
     go-microservices
   ```
3. **驗證**：
   ```bash
   curl http://localhost:8080/users/123
   curl http://localhost:8080/orders/abc
   ```

---
## 主要檔案說明
- **run.sh**：
  ```bash
  #!/usr/bin/env bash
  ./usersvc &> usersvc.log &    # 背景啟動 User Service
  ./ordersvc &> ordersvc.log &  # 背景啟動 Order Service
  exec nginx -g "daemon off;"   # 前台啟 Nginx
  ```
- **nginx.conf**：
  ```nginx
  events {}

  http {
    upstream usersvc { server 127.0.0.1:8000; }
    upstream ordersvc { server 127.0.0.1:9000; }

    server {
      listen 8080;                 # 若需支援平台預設埠
      server_name microgo.zeabur.app;

      location /users/  { proxy_pass http://usersvc; }
      location /orders/ { proxy_pass http://ordersvc; }
      location / { return 200 'Go 微服務 API Gateway'; }
    }
  }
  ```
- **Dockerfile**：多階段建置，最終 image 包含二進位檔、啟動腳本、Nginx。

---
## 平台部署注意事項
- **容器暴露埠**：`EXPOSE 8080`。
- **平台 Web Port**：設定成 `8080` (或要對應的 8080)。
- **域名綁定**：如 `microgo.zeabur.app` → 對應到容器內的 8080。

---
## 常見問題
- **502 Bad Gateway**：通常因為平台路由到錯誤埠（如 8080 → Nginx 未監聽）或 Go 服務未啟動。
- **Host Header**：Nginx 會根據 `server_name` 或 `default_server` 決定要 proxy 哪個 `upstream`。

---
## 授權
MIT © 2025
