@echo off
echo 🚀 準備 FastenMind 專案上傳 GitHub...

echo.
echo 📋 檢查 Git 安裝...
git --version
if errorlevel 1 (
    echo ❌ 錯誤: Git 未安裝，請先安裝 Git
    echo 下載連結: https://git-scm.com/download/win
    pause
    exit /b 1
)

echo.
echo 🔧 初始化 Git 儲存庫...
git init

echo.
echo 📝 設定 Git 使用者資訊 (如需要)...
REM git config user.name "Your Name"
REM git config user.email "your.email@example.com"

echo.
echo 🔒 檢查敏感檔案...
if exist .env (
    echo ⚠️  發現 .env 檔案，確認已加入 .gitignore
)
if exist docker-compose.override.yml (
    echo ⚠️  發現 docker-compose.override.yml，確認已加入 .gitignore
)

echo.
echo 📦 新增檔案到 Git...
git add .

echo.
echo 📊 檢視 Git 狀態...
git status

echo.
echo ✅ 準備工作完成！
echo.
echo 📌 接下來的步驟：
echo 1. 在 GitHub 建立新的儲存庫
echo 2. 執行以下命令：
echo    git commit -m "Initial commit: FastenMind - Fastener Industry B2B Platform"
echo    git branch -M main
echo    git remote add origin https://github.com/YOUR_USERNAME/fastenmind-system.git
echo    git push -u origin main
echo.
echo 💡 提示：記得將 YOUR_USERNAME 替換為你的 GitHub 使用者名稱
echo.
pause