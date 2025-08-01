# FastenMind 資料庫設計

## 快速啟動

```bash
# 啟動資料庫服務
docker-compose up -d

# 檢查服務狀態
docker-compose ps

# 查看日誌
docker-compose logs -f postgres
```

## 服務連線資訊

### PostgreSQL
- Host: localhost
- Port: 5432
- Database: fastenmind_db
- User: fastenmind
- Password: fastenmind123

### Redis
- Host: localhost
- Port: 6379

### pgAdmin
- URL: http://localhost:8081
- Email: admin@fastenmind.com
- Password: admin123

## 測試帳號

| 使用者名稱 | 密碼 | 角色 | 說明 |
|------------|------|------|------|
| admin | password123 | admin | 系統管理員 |
| manager1 | password123 | manager | 業務主管 |
| engineer1 | password123 | engineer | 工程師 |
| sales1 | password123 | sales | 業務人員 |

## 資料庫架構

### 核心資料表
- `companies` - 公司資料
- `accounts` - 使用者帳號
- `customers` - 客戶資料
- `inquiries` - 詢價單
- `quotes` - 報價單

### 製程相關
- `process_categories` - 製程分類
- `equipment_master` - 設備主檔
- `process_cost_config` - 製程成本設定

### 貿易相關
- `global_tariff_rates` - 全球關稅稅率
- `assignment_rules` - 工程師分派規則

## 維護指令

```bash
# 停止服務
docker-compose down

# 清除資料（謹慎使用）
docker-compose down -v

# 備份資料庫
docker exec fastenmind_postgres pg_dump -U fastenmind fastenmind_db > backup.sql

# 還原資料庫
docker exec -i fastenmind_postgres psql -U fastenmind fastenmind_db < backup.sql
```