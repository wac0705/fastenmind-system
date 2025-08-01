# 🧪 FastenMind 測試指南

本指南說明如何運行 FastenMind 系統的各種測試。

## 📋 測試架構概覽

### 測試層級

1. **單元測試 (Unit Tests)**
   - 測試個別函數和組件
   - 位置：`backend/tests/unit/`、`frontend/tests/`
   - 工具：Go testing、Jest + React Testing Library

2. **整合測試 (Integration Tests)**
   - 測試 API 和服務間的整合
   - 位置：`backend/tests/integration/`
   - 工具：Go testing + testify、真實資料庫

3. **端到端測試 (E2E Tests)**
   - 測試完整用戶流程
   - 位置：`e2e/tests/`
   - 工具：Playwright

## 🚀 快速開始

### 前置需求

```bash
# 安裝 Go (後端測試)
go version  # 需要 1.21+

# 安裝 Node.js (前端測試)
node --version  # 需要 18+

# 安裝 Docker (整合測試資料庫)
docker --version
```

### 設定測試環境

```bash
# 1. 後端測試環境
cd backend
go mod download

# 2. 前端測試環境
cd frontend
npm install

# 3. E2E 測試環境
cd e2e
npm install
npx playwright install
```

## 🧪 運行測試

### 後端測試

```bash
cd backend

# 運行所有測試
make test

# 只運行單元測試
make test-unit

# 運行整合測試（需要測試資料庫）
make test-setup  # 啟動測試資料庫
make test-integration
make test-cleanup  # 清理測試資料庫

# 生成覆蓋率報告
make test-coverage

# 運行特定測試
go test -v ./tests/unit/service/auth_service_test.go
go test -v ./tests/integration/auth_integration_test.go
```

### 前端測試

```bash
cd frontend

# 運行所有測試
npm test

# 運行測試並監聽變更
npm run test:watch

# 生成覆蓋率報告
npm run test:coverage

# 運行特定測試
npm test Button.test.tsx
npm test LoginPage.test.tsx
```

### E2E 測試

```bash
cd e2e

# 運行所有 E2E 測試
npm test

# 運行測試（顯示瀏覽器）
npm run test:headed

# 除錯模式
npm run test:debug

# 使用 UI 模式
npm run test:ui

# 運行特定瀏覽器
npm run test:chrome
npm run test:firefox
npm run test:safari

# 運行手機測試
npm run test:mobile

# 運行特定測試
npm run test:auth
npm run test:quotes

# 查看測試報告
npm run report
```

## 📊 測試覆蓋率

### 目標覆蓋率

- **後端**：80% 以上
- **前端**：75% 以上
- **E2E**：主要用戶流程 100%

### 查看覆蓋率報告

```bash
# 後端覆蓋率
cd backend && make test-coverage
open coverage.html

# 前端覆蓋率
cd frontend && npm run test:coverage
open coverage/lcov-report/index.html

# E2E 測試報告
cd e2e && npm run report
```

## 🏗️ 測試結構

### 後端測試結構

```
backend/tests/
├── unit/                    # 單元測試
│   ├── service/            # 服務層測試
│   │   ├── auth_service_test.go
│   │   └── quote_service_test.go
│   └── handler/            # 處理器測試
│       └── auth_handler_test.go
├── integration/            # 整合測試
│   ├── auth_integration_test.go
│   └── quote_integration_test.go
├── testutils/              # 測試工具
│   └── database.go
└── Makefile               # 測試腳本
```

### 前端測試結構

```
frontend/tests/
├── components/             # 組件測試
│   ├── ui/                # UI 組件測試
│   │   ├── Button.test.tsx
│   │   └── Badge.test.tsx
│   └── pages/             # 頁面組件測試
│       └── LoginPage.test.tsx
├── utils/                 # 測試工具
│   └── test-utils.tsx
├── jest.config.js         # Jest 配置
└── jest.setup.js          # 測試環境設定
```

### E2E 測試結構

```
e2e/
├── tests/                 # 測試檔案
│   ├── auth.spec.ts      # 認證流程測試
│   └── quotes.spec.ts    # 報價管理測試
├── global-setup.ts       # 全域設定
├── global-teardown.ts    # 全域清理
├── playwright.config.ts  # Playwright 配置
└── package.json          # 依賴和腳本
```

## 🎯 編寫測試的最佳實踐

### 通用原則

1. **AAA 模式**：Arrange（準備）、Act（執行）、Assert（驗證）
2. **測試命名**：清楚描述測試的行為和預期結果
3. **獨立性**：每個測試應該能獨立運行
4. **可重複性**：測試結果應該是一致的

### 後端測試

```go
func TestAuthService_Login_Success(t *testing.T) {
    // Arrange
    mockRepo := new(MockAccountRepository)
    authService := service.NewAuthService(mockRepo)
    // ... 設定 mock 和測試資料
    
    // Act
    result, err := authService.Login(context.Background(), "testuser", "password123")
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, expectedAccount.ID, result.Account.ID)
    mockRepo.AssertExpectations(t)
}
```

### 前端測試

```typescript
test('should login successfully with valid credentials', async () => {
  // Arrange
  mockServices.auth.login.mockResolvedValue(mockApiResponses.login)
  render(<LoginPage />)
  
  // Act
  fireEvent.change(screen.getByLabelText('帳號'), { target: { value: 'testuser' } })
  fireEvent.change(screen.getByLabelText('密碼'), { target: { value: 'password123' } })
  fireEvent.click(screen.getByRole('button', { name: '登入' }))
  
  // Assert
  await waitFor(() => {
    expect(mockServices.auth.login).toHaveBeenCalledWith({
      username: 'testuser',
      password: 'password123',
    })
  })
})
```

### E2E 測試

```typescript
test('should login successfully with valid credentials', async ({ page }) => {
  // Arrange
  await page.goto('/login')
  
  // Act
  await page.fill('[data-testid="username-input"]', 'admin')
  await page.fill('[data-testid="password-input"]', 'password123')
  await page.click('[data-testid="login-button"]')
  
  // Assert
  await page.waitForURL('**/dashboard')
  await expect(page.locator('h1')).toContainText('儀表板')
})
```

## 🤖 CI/CD 整合

### GitHub Actions

```yaml
name: Tests

on: [push, pull_request]

jobs:
  backend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      - name: Run backend tests
        run: |
          cd backend
          make test-coverage

  frontend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Run frontend tests
        run: |
          cd frontend
          npm ci
          npm run test:coverage

  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Install Playwright
        run: |
          cd e2e
          npm ci
          npx playwright install --with-deps
      - name: Run E2E tests
        run: |
          cd e2e
          npm run test:ci
```

## 🐛 除錯測試

### 後端測試除錯

```bash
# 使用 Delve 除錯器
cd backend
dlv test ./tests/unit/service/auth_service_test.go

# 增加詳細輸出
go test -v -run TestAuthService_Login_Success ./tests/unit/service/
```

### 前端測試除錯

```bash
# 使用 VS Code 除錯
# 在測試檔案中設定中斷點，然後使用 "Debug Jest Tests" 配置

# 在瀏覽器中除錯
npm test -- --debug
```

### E2E 測試除錯

```bash
# 視覺化除錯
npm run test:debug

# 使用 UI 模式
npm run test:ui

# 生成測試程式碼
npm run codegen
```

## 📈 測試報告和監控

### 測試報告

- **後端**：HTML 覆蓋率報告 (`coverage.html`)
- **前端**：Jest 覆蓋率報告 (`coverage/lcov-report/`)
- **E2E**：Playwright HTML 報告

### 持續監控

1. 設定覆蓋率門檻
2. 監控測試執行時間
3. 追蹤測試失敗率
4. 定期檢視測試報告

## ❓ 常見問題

### Q: 測試資料庫連線失敗
A: 確認 Docker 已啟動且測試資料庫容器正在運行：
```bash
docker ps | grep postgres
```

### Q: E2E 測試在 CI 中失敗
A: 檢查是否有足夠的等待時間，並確認測試環境配置正確：
```typescript
await page.waitForLoadState('networkidle')
```

### Q: 前端測試 Mock 不生效
A: 確認 Mock 在測試前已正確設定：
```typescript
beforeEach(() => {
  jest.clearAllMocks()
})
```

## 🆘 取得協助

- 查看測試日誌和錯誤訊息
- 檢查相關文檔和範例
- 在團隊聊天室詢問
- 建立 Issue 回報問題

---

保持測試的完整性和品質是確保系統穩定性的關鍵！🚀