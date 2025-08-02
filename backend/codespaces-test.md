# 在 GitHub Codespaces 上測試 FastenMind API

## 快速開始

1. **在 Codespaces 中開啟專案**
   - 在 GitHub 上點擊 "Code" -> "Codespaces" -> "Create codespace on main"

2. **設定環境變數**
   ```bash
   cd backend
   cp .env.example .env
   # 編輯 .env 檔案，設定適當的值
   ```

3. **安裝 PostgreSQL (如果需要)**
   ```bash
   # 在 Codespaces 中安裝 PostgreSQL
   sudo apt-get update
   sudo apt-get install -y postgresql postgresql-contrib
   
   # 啟動 PostgreSQL
   sudo service postgresql start
   
   # 創建資料庫和使用者
   sudo -u postgres psql -c "CREATE USER fastenmind WITH PASSWORD 'fastenmind123';"
   sudo -u postgres psql -c "CREATE DATABASE fastenmind_db OWNER fastenmind;"
   ```

4. **執行程式**
   ```bash
   # 直接執行編譯好的 Linux 執行檔
   ./fastenmind-api
   
   # 或者從原始碼執行
   go run ./cmd/server
   ```

## 環境設定建議

在 Codespaces 中測試時，建議修改 `.env` 檔案：

```env
# Server Configuration
SERVER_PORT=8080
SERVER_ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=fastenmind
DB_PASSWORD=fastenmind123
DB_NAME=fastenmind_db
DB_SSL_MODE=disable

# Redis Configuration (可選，如果不使用可以留空)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# CORS Configuration
CORS_ALLOWED_ORIGINS=*
```

## 測試 API

1. **健康檢查**
   ```bash
   curl http://localhost:8080/health
   ```

2. **查看 API 文件**
   - 如果有 Swagger UI，訪問：http://localhost:8080/swagger

## 注意事項

1. **Port Forwarding**
   - Codespaces 會自動轉發端口，你可以在 "Ports" 標籤頁中查看
   - 點擊地球圖標可以將端口設為公開（用於測試）

2. **資料庫持久性**
   - Codespaces 中的資料庫數據在關閉後可能會丟失
   - 建議使用外部資料庫服務進行生產測試

3. **記憶體限制**
   - 免費版 Codespaces 有資源限制
   - 如果遇到問題，可以升級到更大的機器類型

## 疑難排解

1. **執行權限問題**
   ```bash
   chmod +x fastenmind-api
   ```

2. **資料庫連線失敗**
   - 確認 PostgreSQL 服務正在運行：`sudo service postgresql status`
   - 檢查連線設定是否正確

3. **端口被佔用**
   - 修改 .env 中的 SERVER_PORT 為其他端口

## 部署到 Zeabur

測試成功後，可以直接使用編譯好的 `fastenmind-api` 部署到 Zeabur：

1. 確保 `.env` 設定適合生產環境
2. 上傳 `fastenmind-api` 執行檔
3. 設定環境變數
4. 啟動服務