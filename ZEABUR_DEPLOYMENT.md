# FastenMind Zeabur 部署指南

## 快速部署

### 1. 準備工作

1. 註冊 [Zeabur](https://zeabur.com) 帳號
2. 安裝 Zeabur CLI (可選)
3. 確保代碼已推送到 GitHub

### 2. 一鍵部署

[![Deploy on Zeabur](https://zeabur.com/button.svg)](https://zeabur.com/new?template=fastenmind-system)

或者點擊上方按鈕進行一鍵部署。

### 3. 手動部署步驟

#### 3.1 創建項目
1. 登入 Zeabur Dashboard
2. 點擊 "New Project"
3. 選擇 "Import from GitHub"
4. 選擇此倉庫

#### 3.2 部署服務

##### 後端服務
1. 添加服務 → 選擇 "backend" 資料夾
2. 選擇 Dockerfile 模式
3. 設置環境變數（見下方配置）

##### 前端服務  
1. 添加服務 → 選擇 "frontend" 資料夾
2. 選擇 Node.js (Next.js) 模式
3. 設置環境變數（見下方配置）

##### 資料庫服務
1. 添加服務 → 選擇 "PostgreSQL"
2. 版本選擇 15
3. 記錄連接資訊

##### Redis 服務
1. 添加服務 → 選擇 "Redis"
2. 版本選擇 7
3. 記錄連接資訊

### 4. 環境變數配置

#### 後端環境變數
```bash
# 伺服器配置
SERVER_PORT=8080
SERVER_ENV=production

# 資料庫配置
DB_PRIMARY_HOST={POSTGRES_HOST}
DB_PRIMARY_PORT={POSTGRES_PORT}
DB_PRIMARY_NAME={POSTGRES_DATABASE}
DB_PRIMARY_USER={POSTGRES_USERNAME}
DB_PRIMARY_PASSWORD={POSTGRES_PASSWORD}
DB_PRIMARY_SSL_MODE=require

# Redis 配置
REDIS_HOST={REDIS_HOST}
REDIS_PORT={REDIS_PORT}
REDIS_PASSWORD={REDIS_PASSWORD}

# JWT 配置
JWT_SECRET_KEY={你的安全密鑰}

# CORS 配置
CORS_ALLOWED_ORIGINS={前端域名}
```

#### 前端環境變數
```bash
# Next.js 配置
NODE_ENV=production
NEXT_TELEMETRY_DISABLED=1

# API 配置
NEXT_PUBLIC_API_URL={後端域名}/api/v1
NEXT_PUBLIC_APP_URL={前端域名}
```

### 5. 域名配置

1. 在 Zeabur Dashboard 中為各服務設置域名
2. 更新環境變數中的域名配置
3. 確保 CORS 設置正確

### 6. 數據庫初始化

部署完成後，後端服務會自動運行數據庫遷移。如需手動初始化：

```bash
# 連接到後端容器
zeabur exec [service-id] -- /bin/sh

# 運行遷移（如果有）
./main migrate
```

### 7. 健康檢查

部署完成後，檢查以下端點：

- 後端健康檢查: `{BACKEND_URL}/health`
- 前端頁面: `{FRONTEND_URL}`
- API 文檔: `{BACKEND_URL}/swagger/index.html`

### 8. 監控和日誌

1. 在 Zeabur Dashboard 查看服務狀態
2. 檢查各服務的日誌輸出
3. 設置告警（如需要）

### 9. 故障排除

#### 常見問題

1. **數據庫連接失敗**
   - 檢查數據庫服務狀態
   - 確認環境變數配置正確
   - 檢查 SSL 模式設置

2. **前端無法訪問後端**
   - 檢查 CORS 配置
   - 確認 API URL 正確
   - 檢查網路連通性

3. **編譯失敗**
   - 檢查依賴版本
   - 查看編譯日誌
   - 確認 Dockerfile 配置

#### 日誌查看
```bash
# 查看服務日誌
zeabur logs [service-id]

# 實時日誌
zeabur logs [service-id] --follow
```

### 10. 自動部署

設置 GitHub Webhook 實現自動部署：

1. 在 Zeabur 項目設置中啟用自動部署
2. 選擇要監聽的分支（建議 main）
3. 推送代碼到對應分支觸發部署

### 11. 環境管理

建議設置多個環境：

- **開發環境**: 用於功能開發和測試
- **預發環境**: 用於發布前驗證  
- **生產環境**: 正式運行環境

每個環境使用不同的環境變數配置。

### 12. 安全注意事項

1. 使用強密碼和密鑰
2. 定期更新依賴
3. 啟用 HTTPS
4. 限制資料庫訪問權限
5. 監控異常活動

### 13. 成本優化

1. 根據使用量選擇合適的服務規格
2. 設置自動縮放規則
3. 定期檢查資源使用情況
4. 優化數據庫查詢和 API 響應

## 支援

如有部署問題，請參考：
- [Zeabur 官方文檔](https://zeabur.com/docs)
- 項目 Issues 頁面
- 聯繫開發團隊