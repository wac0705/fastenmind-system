@echo off

echo Building the application...

REM Build the application
go build -o server.exe cmd/api/main.go

if %ERRORLEVEL% EQU 0 (
    echo Build successful!
) else (
    echo Build failed!
    exit /b 1
)