@echo off
echo Fixing import paths in Go files...

cd backend

:: Fix all files with wrong import paths
powershell -Command "(Get-Content internal\repositories\advanced.repository.go) -replace 'fastenmind-system/internal', 'github.com/fastenmind/fastener-api/internal' | Set-Content internal\repositories\advanced.repository.go"
powershell -Command "(Get-Content internal\repositories\integration.repository.go) -replace 'fastenmind-system/internal', 'github.com/fastenmind/fastener-api/internal' | Set-Content internal\repositories\integration.repository.go"
powershell -Command "(Get-Content internal\repositories\trade.repository.go) -replace 'fastenmind-system/internal', 'github.com/fastenmind/fastener-api/internal' | Set-Content internal\repositories\trade.repository.go"
powershell -Command "(Get-Content internal\services\integration.service.go) -replace 'fastenmind-system/internal', 'github.com/fastenmind/fastener-api/internal' | Set-Content internal\services\integration.service.go"
powershell -Command "(Get-Content internal\services\trade.service.go) -replace 'fastenmind-system/internal', 'github.com/fastenmind/fastener-api/internal' | Set-Content internal\services\trade.service.go"

echo Done!
cd ..