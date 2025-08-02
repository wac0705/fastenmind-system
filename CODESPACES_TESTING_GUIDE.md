# 📚 GitHub Codespaces 測試指南

本指南將手把手教你如何在 GitHub Codespaces 中測試 FastenMind 系統。

## 🚀 步驟 1: 開啟 Codespaces

1. 在 GitHub 專案頁面 (https://github.com/wac0705/fastenmind-system)
2. 點擊綠色的 **"Code"** 按鈕
3. 選擇 **"Codespaces"** 標籤
4. 點擊 **"Create codespace on main"**
5. 等待 Codespaces 環境初始化完成（約 2-3 分鐘）

## 🗄️ 步驟 2: 設置 PostgreSQL 資料庫

在 Codespaces 終端中執行以下命令：

### 2.1 安裝並啟動 PostgreSQL

```bash
# 更新套件列表
sudo apt-get update

# 安裝 PostgreSQL
sudo apt-get install -y postgresql postgresql-contrib

# 啟動 PostgreSQL 服務
sudo service postgresql start

# 確認服務狀態
sudo service postgresql status
```

### 2.2 創建資料庫和使用者

```bash
# 切換到 postgres 使用者並創建資料庫
sudo -u postgres psql << EOF
-- 創建使用者
CREATE USER fastenmind WITH PASSWORD 'fastenmind123';

-- 創建資料庫
CREATE DATABASE fastenmind_db OWNER fastenmind;

-- 授予權限
GRANT ALL PRIVILEGES ON DATABASE fastenmind_db TO fastenmind;

-- 顯示創建結果
\l
\du
EOF
```

### 2.3 執行資料庫初始化腳本

```bash
# 進入專案目錄
cd /workspaces/fastenmind-system

# 執行資料庫初始化腳本
PGPASSWORD=fastenmind123 psql -h localhost -U fastenmind -d fastenmind_db < database/init/01_create_tables.sql
PGPASSWORD=fastenmind123 psql -h localhost -U fastenmind -d fastenmind_db < database/init/02_seed_data.sql
PGPASSWORD=fastenmind123 psql -h localhost -U fastenmind -d fastenmind_db < database/init/03_engineer_assignment_tables.sql
PGPASSWORD=fastenmind123 psql -h localhost -U fastenmind -d fastenmind_db < database/init/04_process_cost_tables.sql
PGPASSWORD=fastenmind123 psql -h localhost -U fastenmind -d fastenmind_db < database/init/05_quote_management_tables.sql
```

## ⚙️ 步驟 3: 設置後端環境

### 3.1 準備環境設定檔

```bash
# 進入後端目錄
cd /workspaces/fastenmind-system/backend

# 複製環境設定檔
cp .env.example .env

# 編輯 .env 檔案（使用 nano 或 vim）
nano .env
```

確保 `.env` 檔案中的設定如下：

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

# JWT Configuration
JWT_SECRET_KEY=your-super-secret-jwt-key-for-testing
JWT_ACCESS_TOKEN_EXPIRE=15m
JWT_REFRESH_TOKEN_EXPIRE=7d

# CORS Configuration
CORS_ALLOWED_ORIGINS=*
```

### 3.2 啟動後端服務

```bash
# 給執行檔加上執行權限
chmod +x fastenmind-api

# 執行編譯好的後端服務
./fastenmind-api
```

或者從原始碼執行：

```bash
# 下載依賴
go mod download

# 執行服務
go run ./cmd/server
```

## 🌐 步驟 4: 設置前端環境

開啟新的終端視窗（Terminal → New Terminal）：

```bash
# 進入前端目錄
cd /workspaces/fastenmind-system/frontend

# 安裝依賴
npm install

# 創建 .env.local 檔案
echo "NEXT_PUBLIC_API_URL=http://localhost:8080" > .env.local

# 啟動前端開發伺服器
npm run dev
```

## 🔍 步驟 5: 訪問和測試系統

### 5.1 查看 Ports（端口）

1. 在 Codespaces 中，點擊底部的 **"PORTS"** 標籤
2. 你應該看到：
   - **8080** - 後端 API
   - **3000** - 前端應用

### 5.2 設置端口可見性

右鍵點擊端口，選擇 **"Port Visibility"** → **"Public"**（用於測試）

### 5.3 訪問應用

1. **前端應用**: 點擊 3000 端口旁的地球圖標
2. **API 健康檢查**: 點擊 8080 端口旁的地球圖標，然後加上 `/health`

## 🧪 步驟 6: 基本功能測試

### 6.1 測試 API 端點

```bash
# 健康檢查
curl http://localhost:8080/health

# 獲取系統資訊
curl http://localhost:8080/api/v1/system/info
```

### 6.2 測試登入功能

1. 訪問前端應用 (http://localhost:3000)
2. 使用測試帳號登入：
   - Email: `admin@fastenmind.com`
   - Password: `password123`

### 6.3 測試主要功能

- ✅ 客戶管理
- ✅ 詢價單建立
- ✅ 報價單管理
- ✅ 訂單追蹤

## 🛠️ 快速啟動腳本

為了方便測試，創建一個一鍵啟動腳本：

```bash
# 創建啟動腳本
cat > /workspaces/fastenmind-system/start-test.sh << 'EOF'
#!/bin/bash

echo "🚀 Starting FastenMind Test Environment..."

# Start PostgreSQL
echo "📦 Starting PostgreSQL..."
sudo service postgresql start

# Start Backend
echo "🔧 Starting Backend API..."
cd /workspaces/fastenmind-system/backend
./fastenmind-api &

# Wait for backend to start
sleep 5

# Start Frontend
echo "🎨 Starting Frontend..."
cd /workspaces/fastenmind-system/frontend
npm run dev &

echo "✅ All services started!"
echo "📍 Frontend: http://localhost:3000"
echo "📍 Backend API: http://localhost:8080"
echo "📍 Check the PORTS tab for public URLs"
EOF

# 給腳本執行權限
chmod +x /workspaces/fastenmind-system/start-test.sh
```

使用腳本：
```bash
/workspaces/fastenmind-system/start-test.sh
```

## 🐛 疑難排解

### 問題 1: PostgreSQL 連線失敗
```bash
# 檢查 PostgreSQL 狀態
sudo service postgresql status

# 重啟 PostgreSQL
sudo service postgresql restart

# 檢查連線
PGPASSWORD=fastenmind123 psql -h localhost -U fastenmind -d fastenmind_db -c "SELECT 1;"
```

### 問題 2: 權限被拒絕
```bash
# 確保執行檔有執行權限
chmod +x /workspaces/fastenmind-system/backend/fastenmind-api
```

### 問題 3: 端口被佔用
```bash
# 查看端口使用情況
sudo lsof -i :8080
sudo lsof -i :3000

# 終止佔用的程序
kill -9 <PID>
```

### 問題 4: 前端無法連接後端
確保 `.env.local` 中的 API URL 正確：
```bash
# 檢查環境變數
cat /workspaces/fastenmind-system/frontend/.env.local

# 應該顯示
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## 📊 效能監控

```bash
# 監控資源使用
htop

# 查看日誌
# 後端日誌會直接顯示在終端
# 前端日誌在 npm run dev 的終端中
```

## 🎯 測試檢查清單

- [ ] PostgreSQL 服務正常運行
- [ ] 資料庫表格創建成功
- [ ] 後端 API 健康檢查通過
- [ ] 前端頁面可以訪問
- [ ] 登入功能正常
- [ ] 基本 CRUD 操作正常
- [ ] API 回應時間 < 500ms

## 💡 小技巧

1. **保存 Codespace**: 閒置 30 分鐘後會自動停止，但狀態會被保存
2. **分享測試環境**: 將端口設為 Public 後可以分享 URL 給他人測試
3. **除錯模式**: 在 VS Code 中可以設置斷點進行除錯
4. **查看資料庫**: 安裝 SQLTools 擴充套件可以直接查看資料庫

## 🚀 下一步

測試成功後，你可以：
1. 部署到 Zeabur 生產環境
2. 設置 CI/CD 自動化測試
3. 進行壓力測試和效能優化

---

有任何問題歡迎在 GitHub Issues 中提出！