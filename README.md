# FastenMind 緊固件詢報價系統

## 🚀 快速開始

### 前置需求
- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Docker & Docker Compose
- N8N (選用)

### 專案結構
```
fastenmind-system/
├── backend/        # Go + Echo 後端 API
├── frontend/       # Next.js 前端應用
├── database/       # 資料庫設計與遷移腳本
├── n8n/           # N8N 工作流程設定
└── docs/          # 系統文件
```

### 快速啟動
```bash
# 1. 啟動資料庫
cd database
docker-compose up -d

# 2. 啟動後端
cd ../backend
go mod download
go run cmd/server/main.go

# 3. 啟動前端
cd ../frontend
npm install
npm run dev
```

### 開發環境設定
詳見各子專案的 README.md 文件。

## 📚 相關文件
- [系統架構設計](docs/architecture.md)
- [API 文件](backend/docs/api.md)
- [資料庫設計](database/schema.md)
- [部署指南](docs/deployment.md)