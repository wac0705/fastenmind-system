# 🚀 FastenMind 生產環境部署指南

## ⚠️ 重要警告
**目前的開發配置不適合直接部署到生產環境！**
請務必按照此指南完成安全配置後再部署。

## 🛡️ 部署前檢查清單

### ✅ **必須完成 (關鍵安全項目)**
- [ ] 更改所有預設密碼
- [ ] 設定強密碼政策
- [ ] 配置 HTTPS/SSL 憑證
- [ ] 關閉不必要的端口暴露
- [ ] 設定防火牆規則
- [ ] 配置備份機制

### ✅ **建議完成 (最佳實務)**
- [ ] 設定監控告警
- [ ] 配置日誌收集
- [ ] 實作速率限制
- [ ] 設定資料庫連線池
- [ ] 配置負載平衡器

## 🔧 快速部署步驟

### 1. 環境準備
```bash
# 複製生產環境配置
cp .env.production.example .env.production

# 編輯配置檔案，填入安全的密碼
nano .env.production
```

### 2. 生成安全密鑰
```bash
# 生成 JWT 密鑰
openssl rand -base64 32

# 生成資料庫密碼
openssl rand -base64 16

# 生成 Redis 密碼
openssl rand -base64 16
```

### 3. 執行部署腳本
```bash
# 給予執行權限
chmod +x scripts/deploy.sh

# 執行部署
./scripts/deploy.sh
```

## 🔐 安全設定詳解

### 密碼要求
- **長度**: 最少 16 字元
- **複雜性**: 包含大小寫字母、數字、特殊符號
- **唯一性**: 不使用預設或常見密碼

### SSL 憑證設定
```bash
# 方法1: 使用 Let's Encrypt (免費)
certbot --nginx -d yourdomain.com

# 方法2: 使用自簽憑證 (僅測試)
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout nginx/ssl/key.pem \
  -out nginx/ssl/cert.pem
```

## 📊 監控與維護

### 服務狀態檢查
```bash
# 檢查所有服務狀態
docker-compose -f docker-compose.production.yml ps

# 查看服務日誌
docker-compose -f docker-compose.production.yml logs -f backend

# 檢查資源使用
docker stats
```

### 備份與恢復
```bash
# 手動備份
docker exec fastenmind_postgres pg_dump -U $DB_USER $DB_NAME > backup.sql

# 恢復備份
docker exec -i fastenmind_postgres psql -U $DB_USER $DB_NAME < backup.sql
```

## 🚨 常見問題與解決

### 1. 服務無法啟動
```bash
# 檢查日誌
docker-compose logs service_name

# 檢查端口占用
netstat -tulpn | grep :80
```

### 2. 資料庫連接失敗
```bash
# 檢查資料庫狀態
docker exec fastenmind_postgres pg_isready

# 測試連接
docker exec -it fastenmind_postgres psql -U $DB_USER -d $DB_NAME
```

### 3. SSL 憑證問題
```bash
# 檢查憑證有效性
openssl x509 -in nginx/ssl/cert.pem -text -noout

# 測試 HTTPS 連接
curl -k https://localhost/api/health
```

## 🔧 進階配置

### 負載平衡設定
```nginx
upstream backend {
    server backend1:8080;
    server backend2:8080;
    server backend3:8080;
}
```

### 資料庫讀寫分離
```yaml
postgres-master:
  image: postgres:15-alpine
  
postgres-slave:
  image: postgres:15-alpine
  command: postgres -c wal_level=replica
```

## 📈 效能優化建議

### 1. 資料庫優化
- 定期 VACUUM 和 ANALYZE
- 調整 shared_buffers 和 work_mem
- 啟用查詢計畫快取

### 2. 應用程式優化
- 啟用 Gzip 壓縮
- 設定適當的快取標頭
- 使用 CDN 加速靜態資源

### 3. 系統優化
- 調整檔案描述符限制
- 優化 TCP 參數
- 設定適當的 swap 大小

## 🛡️ 安全最佳實務

### 1. 網路安全
- 使用 VPN 或跳板機存取管理介面
- 實作 IP 白名單
- 定期更新系統和套件

### 2. 應用程式安全
- 實作 API 速率限制
- 加強輸入驗證
- 定期安全掃描

### 3. 資料安全
- 加密敏感資料
- 定期備份測試
- 實作存取日誌記錄

## 📞 技術支援

如遇到部署問題，請檢查：
1. 系統需求是否滿足
2. 環境變數是否正確設定
3. 防火牆和網路設定
4. 服務日誌錯誤訊息

---
**⚠️ 重要提醒**: 生產環境部署涉及重要的安全和穩定性考量，建議尋求專業的 DevOps 工程師協助。