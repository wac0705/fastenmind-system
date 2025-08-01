# 🚀 FastenMind - 緊固件產業智慧詢報價系統

<div align="center">
  <h3>專為緊固件產業打造的企業級 B2B 詢價報價管理平台</h3>
  <p>
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go" />
    <img src="https://img.shields.io/badge/Next.js-14+-000000?style=for-the-badge&logo=next.js" />
    <img src="https://img.shields.io/badge/PostgreSQL-15+-4169E1?style=for-the-badge&logo=postgresql" />
    <img src="https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker" />
  </p>
</div>

## 📋 目錄

- [系統特色](#-系統特色)
- [技術架構](#-技術架構)
- [快速開始](#-快速開始)
- [功能模組](#-功能模組)
- [部署指南](#-部署指南)
- [開發指南](#-開發指南)
- [授權條款](#-授權條款)

## ✨ 系統特色

### 🏭 產業專屬功能
- **智慧詢價管理** - 自動化詢價單分派與追蹤
- **精準成本計算** - 材料成本、製程成本、關稅計算
- **多階報價審核** - 完整的報價版本控制與審核流程
- **跨國貿易支援** - 多幣別、多語言、關稅管理

### 🛡️ 企業級架構
- **角色權限管理** - 細緻的 RBAC 權限控制
- **多公司架構** - 支援總公司/子公司管理模式
- **API 優先設計** - RESTful API 可整合 ERP/CRM
- **行動裝置支援** - PWA 架構支援離線操作

### 📊 智慧化功能
- **工程師智慧派工** - 基於專長與工作負載自動分派
- **即時匯率更新** - 自動抓取最新匯率資訊
- **自動化工作流程** - N8N 整合支援複雜業務流程
- **數據分析報表** - 詳細的業務分析與洞察

## 🏗️ 技術架構

### 後端技術
- **語言框架**: Go 1.21 + Echo Framework
- **資料庫**: PostgreSQL 15 + Redis 7
- **認證授權**: JWT + RBAC
- **API 文件**: OpenAPI 3.0

### 前端技術
- **框架**: Next.js 14 (App Router)
- **UI 組件**: Radix UI + Tailwind CSS
- **狀態管理**: Zustand + TanStack Query
- **表單驗證**: React Hook Form + Zod

### 基礎設施
- **容器化**: Docker + Docker Compose
- **工作流程**: N8N 自動化平台
- **監控**: Prometheus + Grafana (選配)
- **備份**: 自動化資料庫備份

## 🚀 快速開始

### 前置需求
- Docker & Docker Compose
- Git

### 安裝步驟

```bash
# 1. 複製專案
git clone https://github.com/yourusername/fastenmind-system.git
cd fastenmind-system

# 2. 複製環境設定
cp .env.example .env

# 3. 啟動所有服務
docker-compose up -d

# 4. 初始化資料庫
docker-compose exec backend go run cmd/migrate/main.go up
```

### 預設帳號
- 管理員: `admin` / `admin123`
- 工程師: `engineer1` / `password123`
- 業務員: `sales1` / `password123`

### 存取服務
- 前端應用: http://localhost:3000
- API 文件: http://localhost:8080/swagger
- pgAdmin: http://localhost:8081
- N8N: http://localhost:5678

## 🎯 功能模組

### 📝 詢價管理
- 詢價單建立與管理
- 自動工程師派工
- 詢價狀態追蹤
- 客戶需求分析

### 💰 報價管理
- 多版本報價控制
- 成本自動計算
- 審核流程管理
- PDF 報價單產生

### 📦 訂單管理
- 訂單轉換追蹤
- 生產進度管理
- 交期控制
- 出貨管理

### 💳 財務模組
- 發票管理
- 付款追蹤
- 應收帳款
- 財務報表

### ⚙️ 系統管理
- 用戶權限管理
- 公司資料維護
- 系統參數設定
- 審計日誌

## 🔧 開發指南

### 專案結構
```
fastenmind-system/
├── backend/          # Go 後端 API
├── frontend/         # Next.js 前端
├── database/         # 資料庫腳本
├── docker/           # Docker 配置
├── docs/             # 專案文件
└── scripts/          # 部署腳本
```

### 開發環境設定
```bash
# 後端開發
cd backend
go mod download
go run cmd/server/main.go

# 前端開發
cd frontend
npm install
npm run dev
```

### 測試執行
```bash
# 後端測試
cd backend
go test ./...

# 前端測試
cd frontend
npm test

# E2E 測試
cd e2e
npm run test:e2e
```

## 📚 文件

- [API 文件](./docs/API.md)
- [資料庫設計](./docs/DATABASE.md)
- [部署指南](./PRODUCTION_DEPLOYMENT.md)
- [開發指南](./docs/DEVELOPMENT.md)
- [測試指南](./TESTING_GUIDE.md)

## 🤝 貢獻指南

歡迎貢獻程式碼！請先閱讀 [CONTRIBUTING.md](./CONTRIBUTING.md) 了解貢獻流程。

### 開發流程
1. Fork 專案
2. 建立功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交變更 (`git commit -m 'Add some AmazingFeature'`)
4. 推送分支 (`git push origin feature/AmazingFeature`)
5. 開啟 Pull Request

## 🛡️ 安全性

- 請勿將任何密碼或金鑰提交到版本控制
- 生產環境部署前請參考 [PRODUCTION_DEPLOYMENT.md](./PRODUCTION_DEPLOYMENT.md)
- 發現安全漏洞請私下回報至 security@fastenmind.com

## 📝 授權條款

本專案採用 MIT 授權條款 - 詳見 [LICENSE](./LICENSE) 檔案

## 🙏 致謝

- [Echo Framework](https://echo.labstack.com/)
- [Next.js](https://nextjs.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [N8N](https://n8n.io/)

---

<div align="center">
  <p>Made with ❤️ by FastenMind Team</p>
  <p>
    <a href="https://github.com/yourusername/fastenmind-system/issues">回報問題</a>
    ·
    <a href="https://github.com/yourusername/fastenmind-system/discussions">討論區</a>
  </p>
</div>