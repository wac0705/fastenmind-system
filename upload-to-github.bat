@echo off
echo 📤 FastenMind GitHub 上傳助手
echo ================================
echo.

:: 詢問 GitHub 使用者名稱
set /p GITHUB_USERNAME=請輸入你的 GitHub 使用者名稱: 

if "%GITHUB_USERNAME%"=="" (
    echo ❌ 錯誤：使用者名稱不能為空
    pause
    exit /b 1
)

echo.
echo 🔗 將連接到: https://github.com/%GITHUB_USERNAME%/fastenmind-system.git
echo.

:: 確認是否繼續
set /p CONFIRM=確認要繼續嗎？(Y/N): 
if /i not "%CONFIRM%"=="Y" (
    echo 已取消操作
    pause
    exit /b 0
)

echo.
echo 🚀 開始上傳到 GitHub...
echo.

:: 添加遠端儲存庫
echo 步驟 1/3: 添加遠端儲存庫...
git remote add origin https://github.com/%GITHUB_USERNAME%/fastenmind-system.git

:: 檢查是否成功
if errorlevel 1 (
    echo.
    echo ⚠️  遠端儲存庫可能已存在，嘗試更新...
    git remote set-url origin https://github.com/%GITHUB_USERNAME%/fastenmind-system.git
)

:: 顯示遠端設定
echo.
echo 步驟 2/3: 檢查遠端設定...
git remote -v

:: 推送到 GitHub
echo.
echo 步驟 3/3: 推送程式碼到 GitHub...
echo 這可能需要幾分鐘，請耐心等待...
git push -u origin main

:: 檢查結果
if errorlevel 1 (
    echo.
    echo ❌ 上傳失敗！可能的原因：
    echo 1. GitHub 使用者名稱錯誤
    echo 2. 儲存庫不存在或名稱錯誤
    echo 3. 需要登入 GitHub (會自動彈出登入視窗)
    echo 4. 網路連接問題
    echo.
    echo 💡 建議：
    echo - 確認已在 GitHub 建立 fastenmind-system 儲存庫
    echo - 確認使用者名稱正確
    echo - 如果出現登入視窗，請輸入 GitHub 帳號密碼
    pause
    exit /b 1
)

echo.
echo ✅ 成功上傳到 GitHub！
echo.
echo 🎉 你的專案現在可以在以下網址查看：
echo https://github.com/%GITHUB_USERNAME%/fastenmind-system
echo.
echo 📋 後續建議：
echo 1. 在 GitHub 上添加專案描述和標籤
echo 2. 設定 GitHub Pages (如果需要)
echo 3. 邀請協作者 (如果需要)
echo 4. 設定分支保護規則
echo.
pause