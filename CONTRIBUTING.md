# 🤝 貢獻指南

感謝您有興趣為 FastenMind 專案做出貢獻！

## 📋 貢獻流程

### 1. 準備工作
- Fork 本專案到您的 GitHub 帳號
- Clone 到本地開發環境
- 設定上游遠端倉庫

```bash
git clone https://github.com/yourusername/fastenmind-system.git
cd fastenmind-system
git remote add upstream https://github.com/original/fastenmind-system.git
```

### 2. 建立分支
```bash
# 同步最新程式碼
git fetch upstream
git checkout main
git merge upstream/main

# 建立功能分支
git checkout -b feature/your-feature-name
```

### 3. 開發規範

#### 提交訊息格式
```
<type>(<scope>): <subject>

<body>

<footer>
```

類型 (type):
- `feat`: 新功能
- `fix`: 修復錯誤
- `docs`: 文件更新
- `style`: 程式碼格式調整
- `refactor`: 重構
- `test`: 測試相關
- `chore`: 維護性工作

範例:
```
feat(auth): add two-factor authentication

- Implement TOTP-based 2FA
- Add QR code generation for authenticator apps
- Update user model with 2FA fields

Closes #123
```

#### 程式碼風格

**Go 後端**:
- 使用 `gofmt` 格式化程式碼
- 遵循 [Effective Go](https://golang.org/doc/effective_go.html)
- 執行 `golint` 檢查

**TypeScript/React 前端**:
- 使用 ESLint + Prettier
- 遵循 React Hooks 規則
- 元件使用 PascalCase 命名

### 4. 測試要求
- 新功能必須包含單元測試
- 修復錯誤應包含回歸測試
- 確保所有測試通過

```bash
# 後端測試
cd backend && go test ./...

# 前端測試
cd frontend && npm test

# E2E 測試
cd e2e && npm run test:e2e
```

### 5. 提交 Pull Request

#### PR 檢查清單
- [ ] 程式碼已格式化
- [ ] 所有測試通過
- [ ] 更新相關文件
- [ ] 無敏感資訊洩露
- [ ] PR 描述清晰完整

#### PR 模板
```markdown
## 描述
簡要描述這個 PR 的內容

## 變更類型
- [ ] 錯誤修復
- [ ] 新功能
- [ ] 破壞性變更
- [ ] 文件更新

## 測試
描述您如何測試這些變更

## 相關 Issue
Closes #(issue)

## 截圖 (如適用)
```

## 🐛 回報問題

### Issue 模板
```markdown
## 問題描述
清楚描述遇到的問題

## 重現步驟
1. 前往 '...'
2. 點擊 '....'
3. 滾動到 '....'
4. 看到錯誤

## 預期行為
描述預期應該發生什麼

## 實際行為
描述實際發生了什麼

## 環境資訊
- OS: [e.g. Windows 10]
- Browser: [e.g. Chrome 91]
- Version: [e.g. 1.0.0]

## 附加資訊
任何有助於解決問題的資訊
```

## 💬 開發討論

- 使用 GitHub Discussions 進行功能討論
- 加入我們的 Discord 社群 (連結待定)
- 訂閱開發者郵件列表

## 🏆 貢獻者行為準則

### 我們的承諾
- 營造友善、包容的環境
- 尊重不同觀點和經驗
- 優雅地接受建設性批評
- 專注於對社群最有利的事

### 不可接受的行為
- 使用性暗示語言或圖像
- 人身攻擊或政治攻擊
- 公開或私下騷擾
- 未經許可發布他人資訊

## 📜 授權協議

提交程式碼即表示您同意將程式碼以 MIT 授權條款授權。

## 🙋 需要幫助？

- 查看 [文件](./docs)
- 在 [Discussions](https://github.com/yourusername/fastenmind-system/discussions) 提問
- 聯繫維護者

感謝您的貢獻！ 🎉