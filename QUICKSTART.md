# 🚀 FastenMind 快速啟動指南

## 📋 前置需求

- Docker Desktop
- Node.js 18+ (選用，用於本地開發)
- Go 1.21+ (選用，用於本地開發)

## 🎯 一鍵啟動（推薦）

```bash
# 1. 複製環境變數設定
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env

# 2. 啟動所有服務
docker-compose up -d

# 3. 等待服務啟動（約 30 秒）
docker-compose ps

# 4. 檢查服務健康狀態
docker-compose logs
```

## 🌐 服務訪問地址

| 服務 | 地址 | 說明 |
|------|------|------|
| 前端應用 | http://localhost:3000 | Next.js 前端介面 |
| 後端 API | http://localhost:8080 | Go Echo API 服務 |
| pgAdmin | http://localhost:8081 | 資料庫管理介面 |
| N8N | http://localhost:5678 | 工作流程自動化平台 |

## 👤 測試帳號

| 帳號 | 密碼 | 角色 | 權限說明 |
|------|------|------|----------|
| admin | password123 | 系統管理員 | 完整系統權限 |
| manager1 | password123 | 業務主管 | 審核報價、查看報表 |
| engineer1 | password123 | 工程師 | 處理詢價、建立報價 |
| sales1 | password123 | 業務人員 | 建立詢價、查看報價 |

## 🛠️ 本地開發模式

### 後端開發
```bash
cd backend
go mod download
go run cmd/server/main.go
```

### 前端開發
```bash
cd frontend
npm install
npm run dev
```

### 資料庫管理
```bash
cd database
docker-compose up -d
```

## 📦 功能模組

### 已實作功能
- ✅ JWT 雙令牌認證系統
- ✅ 多租戶公司管理
- ✅ 使用者角色權限控制
- ✅ 資料庫架構與初始資料
- ✅ Docker 容器化部署

### 開發中功能
- 🚧 詢價單管理介面
- 🚧 工程師分派系統
- 🚧 製程成本計算
- 🚧 報價審核流程
- 🚧 N8N 工作流程整合

## 🔧 常用指令

```bash
# 停止所有服務
docker-compose down

# 重新建構並啟動
docker-compose up -d --build

# 查看後端日誌
docker-compose logs -f backend

# 查看前端日誌
docker-compose logs -f frontend

# 進入後端容器
docker exec -it fastenmind_backend sh

# 進入資料庫
docker exec -it fastenmind_postgres psql -U fastenmind -d fastenmind_db
```

## 🐛 疑難排解

### 1. 端口被佔用
```bash
# 檢查端口使用
netstat -ano | findstr :3000
netstat -ano | findstr :8080

# 修改 docker-compose.yml 中的端口映射
```

### 2. 資料庫連線失敗
```bash
# 確認資料庫服務正常
docker-compose ps postgres
docker-compose logs postgres

# 重新啟動資料庫
docker-compose restart postgres
```

### 3. 前端無法連線後端
- 檢查 frontend/.env 中的 NEXT_PUBLIC_API_URL
- 確認後端服務正常運行
- 檢查 CORS 設定

## 📚 下一步

1. 訪問 http://localhost:3000 開始使用系統
2. 使用測試帳號登入體驗功能
3. 查看 `/docs` 目錄了解詳細技術文件
4. 開始開發新功能或客製化

## 💡 提示

- 開發環境的資料會保存在 Docker volumes 中
- 首次啟動可能需要較長時間下載 Docker 映像
- 建議使用 Chrome 或 Firefox 瀏覽器獲得最佳體驗

---

需要協助？查看 [完整文件](./docs) 或提交 Issue。