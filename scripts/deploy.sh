#!/bin/bash

# FastenMind 生產環境部署腳本
# 使用方法: ./scripts/deploy.sh

set -e  # 遇到錯誤立即停止

echo "🚀 開始 FastenMind 生產環境部署..."

# 檢查必要檔案
if [ ! -f ".env.production" ]; then
    echo "❌ 錯誤: .env.production 檔案不存在"
    echo "請複製 .env.production.example 並填入正確的設定值"
    exit 1
fi

# 檢查 Docker
if ! command -v docker &> /dev/null; then
    echo "❌ 錯誤: Docker 未安裝"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ 錯誤: Docker Compose 未安裝"
    exit 1
fi

# 建立必要目錄
echo "📁 建立必要目錄..."
mkdir -p nginx/ssl
mkdir -p backups
mkdir -p logs
mkdir -p monitoring

# 生成自簽 SSL 憑證 (僅用於測試，生產環境請使用正式憑證)
if [ ! -f "nginx/ssl/cert.pem" ]; then
    echo "🔐 生成自簽 SSL 憑證..."
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout nginx/ssl/key.pem \
        -out nginx/ssl/cert.pem \
        -subj "/C=TW/ST=Taiwan/L=Taipei/O=FastenMind/CN=localhost"
fi

# 建立 Nginx 配置
cat > nginx/nginx.conf << 'EOF'
events {
    worker_connections 1024;
}

http {
    upstream backend {
        server backend:8080;
    }
    
    upstream frontend {
        server frontend:3000;
    }

    # HTTP 重導向到 HTTPS
    server {
        listen 80;
        server_name _;
        return 301 https://$host$request_uri;
    }

    # HTTPS 設定
    server {
        listen 443 ssl http2;
        server_name _;

        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;
        
        # 安全標頭
        add_header X-Frame-Options SAMEORIGIN;
        add_header X-Content-Type-Options nosniff;
        add_header X-XSS-Protection "1; mode=block";
        add_header Strict-Transport-Security "max-age=31536000";

        # API 代理
        location /api/ {
            proxy_pass http://backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # 前端代理
        location / {
            proxy_pass http://frontend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }
}
EOF

# 檢查環境變數
echo "🔍 檢查環境變數..."
source .env.production

required_vars=("DB_PASSWORD" "REDIS_PASSWORD" "JWT_SECRET_KEY")
for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo "❌ 錯誤: 環境變數 $var 未設定"
        exit 1
    fi
done

# 停止現有服務
echo "🛑 停止現有服務..."
docker-compose -f docker-compose.production.yml down

# 建立映像檔
echo "🔨 建立 Docker 映像檔..."
docker-compose -f docker-compose.production.yml build --no-cache

# 啟動服務
echo "▶️  啟動生產服務..."
docker-compose -f docker-compose.production.yml up -d

# 等待服務啟動
echo "⏳ 等待服務啟動..."
sleep 30

# 健康檢查
echo "🏥 執行健康檢查..."
if curl -f -k https://localhost/api/health > /dev/null 2>&1; then
    echo "✅ 後端服務正常"
else
    echo "❌ 後端服務異常"
    exit 1
fi

if curl -f -k https://localhost > /dev/null 2>&1; then
    echo "✅ 前端服務正常"
else
    echo "❌ 前端服務異常"
    exit 1
fi

# 顯示服務狀態
echo "📊 服務狀態:"
docker-compose -f docker-compose.production.yml ps

echo ""
echo "🎉 部署完成!"
echo "🌐 應用程式網址: https://localhost"
echo "📊 監控面板: http://localhost:9090 (如果啟用)"
echo ""
echo "📋 重要提醒:"
echo "1. 預設使用自簽憑證，瀏覽器會顯示安全警告"
echo "2. 生產環境請使用正式的 SSL 憑證"
echo "3. 定期檢查 ./backups 目錄中的資料庫備份"
echo "4. 監控 ./logs 目錄中的應用程式日誌"