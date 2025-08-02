# 🚀 FastenMind 快速開始指南 (Codespaces)

## 超快速開始（30 秒）

在 GitHub Codespaces 中，只需執行一個命令：

```bash
# 給腳本執行權限並執行
chmod +x /workspaces/fastenmind-system/scripts/*.sh
/workspaces/fastenmind-system/scripts/start-codespaces.sh
```

就這樣！系統會自動：
- ✅ 安裝並設置 PostgreSQL
- ✅ 初始化資料庫
- ✅ 啟動後端 API
- ✅ 啟動前端應用

## 🎯 快速測試檢查清單

### 1️⃣ 基本連線測試
```bash
# 檢查服務狀態
/workspaces/fastenmind-system/check-status.sh

# API 健康檢查
curl http://localhost:8080/health

# 前端首頁
curl -I http://localhost:3000
```

### 2️⃣ 登入測試
1. 開啟瀏覽器訪問 http://localhost:3000
2. 使用測試帳號：
   - Email: `admin@fastenmind.com`
   - Password: `password123`

### 3️⃣ 功能測試清單

#### 客戶管理
- [ ] 新增客戶
- [ ] 編輯客戶資料
- [ ] 搜尋客戶
- [ ] 刪除客戶

#### 詢價管理
- [ ] 建立詢價單
- [ ] 自動派工給工程師
- [ ] 上傳圖檔
- [ ] 查看詢價狀態

#### 報價管理
- [ ] 從詢價單建立報價
- [ ] 成本計算
- [ ] 多版本管理
- [ ] PDF 輸出

#### 訂單管理
- [ ] 從報價單轉訂單
- [ ] 訂單狀態追蹤
- [ ] 出貨管理
- [ ] 付款記錄

### 4️⃣ API 測試範例

```bash
# 取得客戶列表
curl -X GET http://localhost:8080/api/v1/customers \
  -H "Authorization: Bearer YOUR_TOKEN"

# 建立新詢價單
curl -X POST http://localhost:8080/api/v1/inquiries \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "customer_id": "UUID",
    "product_name": "Test Product",
    "quantity": 100,
    "required_date": "2024-12-31"
  }'
```

## 🛠️ 實用命令

### 查看日誌
```bash
# 後端日誌
tail -f /tmp/backend.log

# 前端日誌  
tail -f /tmp/frontend.log

# PostgreSQL 日誌
sudo tail -f /var/log/postgresql/postgresql-*.log
```

### 重啟服務
```bash
# 停止所有服務
/workspaces/fastenmind-system/stop-services.sh

# 重新啟動
/workspaces/fastenmind-system/scripts/start-codespaces.sh
```

### 資料庫操作
```bash
# 連接資料庫
PGPASSWORD=fastenmind123 psql -h localhost -U fastenmind -d fastenmind_db

# 查看所有表格
\dt

# 查看表格結構
\d companies

# 執行查詢
SELECT * FROM accounts LIMIT 5;
```

## 🐛 常見問題

### 1. "Permission denied" 錯誤
```bash
chmod +x /workspaces/fastenmind-system/backend/fastenmind-api
chmod +x /workspaces/fastenmind-system/scripts/*.sh
```

### 2. 資料庫連線失敗
```bash
# 重啟 PostgreSQL
sudo service postgresql restart

# 重新初始化資料庫
/workspaces/fastenmind-system/scripts/init-database.sh
```

### 3. 端口被佔用
```bash
# 查看佔用程序
sudo lsof -i :8080
sudo lsof -i :3000

# 強制停止
pkill -f fastenmind-api
pkill -f "next dev"
```

### 4. 前端無法連接後端
```bash
# 確認環境變數
cat /workspaces/fastenmind-system/frontend/.env.local

# 應該包含
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## 📊 效能測試

### 簡單壓力測試
```bash
# 安裝 Apache Bench
sudo apt-get install -y apache2-utils

# 測試 API 效能
ab -n 100 -c 10 http://localhost:8080/health
```

### 監控資源使用
```bash
# 即時監控
htop

# 查看記憶體使用
free -h

# 查看磁碟使用
df -h
```

## 🎨 自定義設定

### 修改環境變數
```bash
# 後端設定
nano /workspaces/fastenmind-system/backend/.env

# 前端設定
nano /workspaces/fastenmind-system/frontend/.env.local
```

### 修改資料庫連線
編輯 `.env` 中的資料庫設定：
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=fastenmind
DB_PASSWORD=fastenmind123
DB_NAME=fastenmind_db
```

## 🚀 下一步

1. **探索 API 文件**
   - 訪問 http://localhost:8080/swagger (如果有設置)

2. **執行測試**
   ```bash
   # 後端測試
   cd /workspaces/fastenmind-system/backend
   go test ./...

   # 前端測試
   cd /workspaces/fastenmind-system/frontend
   npm test
   ```

3. **部署準備**
   - 設置生產環境變數
   - 優化資料庫索引
   - 設置 SSL 憑證

## 💡 專業提示

1. **使用 VS Code 除錯器**
   - 在程式碼中設置斷點
   - 使用 F5 啟動除錯模式

2. **資料庫 GUI 工具**
   - 安裝 SQLTools 擴充套件
   - 使用內建的資料庫瀏覽器

3. **API 測試工具**
   - 使用 Thunder Client 擴充套件
   - 匯入 Postman 集合

---

需要幫助？查看 [完整測試指南](CODESPACES_TESTING_GUIDE.md) 或在 GitHub Issues 中提問！