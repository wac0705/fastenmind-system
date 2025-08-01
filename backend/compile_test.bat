@echo off
echo FastenMind Backend Compilation Test
echo ===================================
echo.

cd /d "%~dp0"

echo Checking Go installation...
go version
if errorlevel 1 (
    echo Error: Go is not installed or not in PATH
    pause
    exit /b 1
)

echo.
echo Building the application...
go build -v ./... 2>&1

if errorlevel 1 (
    echo.
    echo Build failed! Please check the errors above.
) else (
    echo.
    echo Build successful!
)

echo.
echo Press any key to exit...
pause > nul