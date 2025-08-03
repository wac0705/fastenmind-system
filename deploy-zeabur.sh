#!/bin/bash

# FastenMind Zeabur 部署腳本
# 此腳本協助準備和部署到 Zeabur 平台

set -e

echo "🚀 FastenMind Zeabur 部署準備"
echo "================================"

# 檢查必要工具
command -v git >/dev/null 2>&1 || { echo "❌ Git 未安裝. 請先安裝 Git." >&2; exit 1; }

# 檢查是否在 Git 倉庫中
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ 當前目錄不是 Git 倉庫"
    exit 1
fi

# 檢查是否有未提交的更改
if ! git diff-index --quiet HEAD --; then
    echo "⚠️  發現未提交的更改，請先提交所有更改"
    echo "運行以下命令提交更改："
    echo "  git add ."
    echo "  git commit -m 'Deploy to Zeabur'"
    echo "  git push origin main"
    exit 1
fi

# 檢查必要文件
echo "📋 檢查部署文件..."

required_files=(
    "zeabur.yaml"
    "backend/Dockerfile"
    "frontend/Dockerfile"
    "backend/go.mod"
    "frontend/package.json"
)

for file in "${required_files[@]}"; do
    if [ ! -f "$file" ]; then
        echo "❌ 缺少必要文件: $file"
        exit 1
    fi
    echo "✅ $file"
done

# 檢查環境變數範例文件
if [ ! -f ".env.zeabur" ]; then
    echo "❌ 缺少 .env.zeabur 文件，請確保已創建此文件"
    exit 1
fi

echo "✅ 所有必要文件都存在"

# 顯示 Zeabur 配置
echo ""
echo "📝 Zeabur 配置預覽:"
echo "================================"
cat zeabur.yaml

echo ""
echo "🔧 環境變數配置:"
echo "================================"
echo "請在 Zeabur Dashboard 中設置以下環境變數："
echo ""
cat .env.zeabur

echo ""
echo "📖 部署指南:"
echo "================================"
echo "1. 確保代碼已推送到 GitHub"
echo "2. 訪問 https://zeabur.com"
echo "3. 創建新項目並導入此倉庫"
echo "4. 按照 ZEABUR_DEPLOYMENT.md 中的指南操作"
echo ""

# 推送代碼（如果需要）
read -p "🤔 是否要推送當前代碼到 GitHub? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "📤 推送代碼到 GitHub..."
    
    # 獲取當前分支
    current_branch=$(git branch --show-current)
    
    echo "當前分支: $current_branch"
    git push origin "$current_branch"
    
    echo "✅ 代碼已推送到 GitHub"
else
    echo "ℹ️  請手動推送代碼到 GitHub"
fi

echo ""
echo "🎉 部署準備完成！"
echo "================================"
echo "下一步："
echo "1. 訪問 https://zeabur.com"
echo "2. 創建新項目"
echo "3. 導入 GitHub 倉庫"
echo "4. 配置環境變數"
echo "5. 部署服務"
echo ""
echo "詳細步驟請參考 ZEABUR_DEPLOYMENT.md"

# 本地測試選項
echo ""
read -p "🧪 是否要啟動本地 Zeabur 模擬環境進行測試? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🐳 啟動本地測試環境..."
    
    # 檢查 Docker
    command -v docker >/dev/null 2>&1 || { echo "❌ Docker 未安裝. 請先安裝 Docker." >&2; exit 1; }
    command -v docker-compose >/dev/null 2>&1 || { echo "❌ Docker Compose 未安裝. 請先安裝 Docker Compose." >&2; exit 1; }
    
    # 啟動服務
    docker-compose -f docker-compose.zeabur.yml up -d
    
    echo "✅ 本地測試環境已啟動"
    echo ""
    echo "📱 服務地址:"
    echo "  前端: http://localhost:3000"
    echo "  後端: http://localhost:8080"
    echo "  API 文檔: http://localhost:8080/swagger/index.html"
    echo ""
    echo "停止測試環境: docker-compose -f docker-compose.zeabur.yml down"
fi

echo ""
echo "🚀 Happy Deploying!"