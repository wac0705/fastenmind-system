import { test, expect } from '@playwright/test';

test.describe('Finance Management', () => {
  test.beforeEach(async ({ page }) => {
    // Login as finance user
    await page.goto('/login');
    await page.fill('[data-testid="username-input"]', 'finance1');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');
    await page.waitForURL('**/dashboard');
  });

  test('should display finance dashboard', async ({ page }) => {
    // Navigate to finance page
    await page.goto('/finance');
    
    // Check page title
    await expect(page.locator('h1')).toContainText('財務管理');
    
    // Check if main sections are visible
    await expect(page.locator('[data-testid="revenue-summary"]')).toBeVisible();
    await expect(page.locator('[data-testid="expenses-summary"]')).toBeVisible();
    await expect(page.locator('[data-testid="receivables-summary"]')).toBeVisible();
    await expect(page.locator('[data-testid="cashflow-summary"]')).toBeVisible();
  });

  test('should create invoice from order', async ({ page }) => {
    // Go to orders page first
    await page.goto('/orders');
    await page.waitForSelector('[data-testid="orders-table"]');
    
    // Find a completed order
    const completedOrder = page.locator('[data-testid="order-row"]').filter({ hasText: '已完成' }).first();
    const orderCount = await completedOrder.count();
    
    if (orderCount === 0) {
      test.skip(); // Skip if no completed orders
      return;
    }
    
    await completedOrder.click();
    await page.waitForURL('**/orders/**');
    
    // Click create invoice button
    await page.click('[data-testid="create-invoice-button"]');
    
    // Should navigate to create invoice page
    await page.waitForURL('**/finance/invoices/new**');
    
    // Fill invoice details
    await page.fill('[data-testid="invoice-no"]', 'INV-E2E-001');
    
    // Tax rate should be pre-filled, verify it
    const taxRateValue = await page.inputValue('[data-testid="tax-rate"]');
    expect(parseFloat(taxRateValue)).toBeGreaterThanOrEqual(0);
    
    // Set due date (30 days from now)
    const futureDate = new Date();
    futureDate.setDate(futureDate.getDate() + 30);
    const dateString = futureDate.toISOString().split('T')[0];
    await page.fill('[data-testid="due-date"]', dateString);
    
    // Add notes
    await page.fill('[data-testid="notes"]', 'E2E test invoice');
    
    // Submit invoice
    await page.click('[data-testid="create-invoice-submit"]');
    
    // Should redirect to invoice detail page
    await page.waitForURL('**/finance/invoices/**');
    
    // Check if invoice was created successfully
    await expect(page.locator('[data-testid="success-message"]')).toContainText('發票建立成功');
    await expect(page.locator('[data-testid="invoice-number"]')).toContainText('INV-E2E-001');
  });

  test('should process payment for invoice', async ({ page }) => {
    // Navigate to invoices page
    await page.goto('/finance/invoices');
    await page.waitForSelector('[data-testid="invoices-table"]');
    
    // Find a pending invoice
    const pendingInvoice = page.locator('[data-testid="invoice-row"]').filter({ hasText: '待付款' }).first();
    const invoiceCount = await pendingInvoice.count();
    
    if (invoiceCount === 0) {
      test.skip(); // Skip if no pending invoices
      return;
    }
    
    await pendingInvoice.click();
    await page.waitForURL('**/finance/invoices/**');
    
    // Click process payment button
    await page.click('[data-testid="process-payment-button"]');
    
    // Fill payment details
    await page.fill('[data-testid="payment-amount"]', '5000');
    await page.selectOption('[data-testid="payment-method"]', 'bank_transfer');
    await page.fill('[data-testid="payment-reference"]', 'TXN-E2E-123456');
    await page.fill('[data-testid="payment-notes"]', 'E2E test payment');
    
    // Submit payment
    await page.click('[data-testid="process-payment-submit"]');
    
    // Check if payment was processed
    await expect(page.locator('[data-testid="success-message"]')).toContainText('付款處理成功');
    
    // Check if invoice status updated
    await expect(page.locator('[data-testid="invoice-status"]')).toContainText('部分付款');
  });

  test('should submit expense for approval', async ({ page }) => {
    // Navigate to expenses page
    await page.goto('/finance/expenses');
    
    // Check if create expense button is visible
    await expect(page.locator('[data-testid="create-expense-button"]')).toBeVisible();
    
    // Click create expense
    await page.click('[data-testid="create-expense-button"]');
    
    // Should navigate to create expense page
    await page.waitForURL('**/finance/expenses/new');
    
    // Fill expense details
    await page.selectOption('[data-testid="expense-category"]', 'travel');
    await page.fill('[data-testid="expense-amount"]', '500');
    await page.fill('[data-testid="expense-description"]', 'E2E test business trip');
    
    // Set expense date
    const today = new Date().toISOString().split('T')[0];
    await page.fill('[data-testid="expense-date"]', today);
    
    // Upload receipt (simulate)
    await page.fill('[data-testid="receipt-url"]', 'https://example.com/receipt.jpg');
    
    // Submit expense
    await page.click('[data-testid="submit-expense-button"]');
    
    // Should redirect to expenses list
    await page.waitForURL('**/finance/expenses');
    
    // Check if expense was created
    await expect(page.locator('[data-testid="success-message"]')).toContainText('費用申請已提交');
  });

  test('should approve expense (as manager)', async ({ page }) => {
    // Logout and login as manager
    await page.click('[data-testid="user-menu"]');
    await page.click('[data-testid="logout-button"]');
    
    await page.waitForURL('**/login');
    await page.fill('[data-testid="username-input"]', 'manager1');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');
    await page.waitForURL('**/dashboard');
    
    // Navigate to expenses page
    await page.goto('/finance/expenses');
    await page.waitForSelector('[data-testid="expenses-table"]');
    
    // Find a pending expense
    const pendingExpense = page.locator('[data-testid="expense-row"]').filter({ hasText: '待審核' }).first();
    const expenseCount = await pendingExpense.count();
    
    if (expenseCount === 0) {
      test.skip(); // Skip if no pending expenses
      return;
    }
    
    await pendingExpense.click();
    await page.waitForURL('**/finance/expenses/**');
    
    // Click approve button
    await page.click('[data-testid="approve-expense-button"]');
    
    // Fill approval notes
    await page.fill('[data-testid="approval-notes"]', 'Approved for reimbursement');
    await page.click('[data-testid="confirm-approval"]');
    
    // Check if expense was approved
    await expect(page.locator('[data-testid="expense-status"]')).toContainText('已核准');
    await expect(page.locator('[data-testid="success-message"]')).toContainText('費用已核准');
  });

  test('should view receivables with aging analysis', async ({ page }) => {
    // Navigate to receivables page
    await page.goto('/finance/receivables');
    
    // Check page title
    await expect(page.locator('h1')).toContainText('應收帳款');
    
    // Check if receivables table is visible
    await expect(page.locator('[data-testid="receivables-table"]')).toBeVisible();
    
    // Check if aging analysis is displayed
    await expect(page.locator('[data-testid="aging-analysis"]')).toBeVisible();
    
    // Check aging categories
    await expect(page.locator('[data-testid="current-amount"]')).toBeVisible();
    await expect(page.locator('[data-testid="overdue-1-30"]')).toBeVisible();
    await expect(page.locator('[data-testid="overdue-31-60"]')).toBeVisible();
    await expect(page.locator('[data-testid="overdue-61-90"]')).toBeVisible();
    await expect(page.locator('[data-testid="overdue-over-90"]')).toBeVisible();
  });

  test('should filter receivables by status', async ({ page }) => {
    // Navigate to receivables page
    await page.goto('/finance/receivables');
    await page.waitForSelector('[data-testid="receivables-table"]');
    
    // Filter by overdue status
    await page.click('[data-testid="status-filter"]');
    await page.click('[data-testid="status-overdue"]');
    
    // Wait for filtered results
    await page.waitForTimeout(1000);
    
    // All visible receivables should be overdue
    const statusBadges = page.locator('[data-testid="receivable-status"]');
    const count = await statusBadges.count();
    
    for (let i = 0; i < count; i++) {
      await expect(statusBadges.nth(i)).toContainText('逾期');
    }
  });

  test('should view financial summary with charts', async ({ page }) => {
    // Navigate to finance dashboard
    await page.goto('/finance');
    
    // Check if charts are visible
    await expect(page.locator('[data-testid="revenue-chart"]')).toBeVisible();
    await expect(page.locator('[data-testid="expense-chart"]')).toBeVisible();
    await expect(page.locator('[data-testid="cashflow-chart"]')).toBeVisible();
    
    // Check summary cards
    await expect(page.locator('[data-testid="total-revenue"]')).toBeVisible();
    await expect(page.locator('[data-testid="total-expenses"]')).toBeVisible();
    await expect(page.locator('[data-testid="net-profit"]')).toBeVisible();
    
    // Check if amounts are displayed correctly (should be numbers)
    const totalRevenue = await page.locator('[data-testid="total-revenue-amount"]').textContent();
    expect(totalRevenue).toMatch(/[\d,]+/);
  });

  test('should export financial reports', async ({ page }) => {
    // Navigate to finance dashboard
    await page.goto('/finance');
    
    // Setup download promise before clicking
    const downloadPromise = page.waitForEvent('download');
    
    // Click export button
    await page.click('[data-testid="export-financial-report"]');
    
    // Select report type
    await page.selectOption('[data-testid="report-type"]', 'summary');
    
    // Set date range
    await page.selectOption('[data-testid="date-range"]', 'current_month');
    
    // Confirm export
    await page.click('[data-testid="confirm-export"]');
    
    // Wait for download to start
    const download = await downloadPromise;
    
    // Check if file is downloaded
    expect(download.suggestedFilename()).toMatch(/financial.*\.(pdf|xlsx)/i);
  });

  test('should handle finance permissions by role', async ({ page }) => {
    // Test with different roles
    const roles = [
      { username: 'admin', role: 'admin', canCreateInvoice: true, canApproveExpense: true },
      { username: 'finance1', role: 'finance', canCreateInvoice: true, canApproveExpense: true },
      { username: 'manager1', role: 'manager', canCreateInvoice: true, canApproveExpense: true },
      { username: 'sales1', role: 'sales', canCreateInvoice: false, canApproveExpense: false },
      { username: 'engineer1', role: 'engineer', canCreateInvoice: false, canApproveExpense: false },
    ];
    
    for (const roleTest of roles) {
      // Logout current user
      await page.click('[data-testid="user-menu"]');
      await page.click('[data-testid="logout-button"]');
      
      // Login with test role
      await page.waitForURL('**/login');
      await page.fill('[data-testid="username-input"]', roleTest.username);
      await page.fill('[data-testid="password-input"]', 'password123');
      await page.click('[data-testid="login-button"]');
      await page.waitForURL('**/dashboard');
      
      // Navigate to finance page
      await page.goto('/finance');
      
      // Check invoice creation permission
      const createInvoiceVisible = await page.locator('[data-testid="create-invoice-button"]').isVisible().catch(() => false);
      expect(createInvoiceVisible).toBe(roleTest.canCreateInvoice);
      
      // Check expense approval permission
      if (roleTest.canApproveExpense) {
        await page.goto('/finance/expenses');
        await expect(page.locator('[data-testid="approve-expense-button"]')).toBeVisible();
      }
    }
  });

  test('should display payment history', async ({ page }) => {
    // Navigate to payments page
    await page.goto('/finance/payments');
    
    // Check page title
    await expect(page.locator('h1')).toContainText('付款管理');
    
    // Check if payments table is visible
    await expect(page.locator('[data-testid="payments-table"]')).toBeVisible();
    
    // Check table headers
    await expect(page.locator('th')).toContainText('付款日期');
    await expect(page.locator('th')).toContainText('發票號碼');
    await expect(page.locator('th')).toContainText('付款金額');
    await expect(page.locator('th')).toContainText('付款方式');
    await expect(page.locator('th')).toContainText('狀態');
  });

  test('should validate payment amounts', async ({ page }) => {
    // Navigate to invoices and try to overpay
    await page.goto('/finance/invoices');
    await page.waitForSelector('[data-testid="invoices-table"]');
    
    const partialInvoice = page.locator('[data-testid="invoice-row"]').filter({ hasText: '部分付款' }).first();
    const invoiceCount = await partialInvoice.count();
    
    if (invoiceCount === 0) {
      test.skip(); // Skip if no partial invoices
      return;
    }
    
    await partialInvoice.click();
    await page.waitForURL('**/finance/invoices/**');
    
    // Get remaining amount
    const remainingText = await page.locator('[data-testid="remaining-amount"]').textContent();
    const remainingAmount = parseFloat(remainingText?.replace(/[^\d.]/g, '') || '0');
    
    // Click process payment
    await page.click('[data-testid="process-payment-button"]');
    
    // Try to pay more than remaining amount
    await page.fill('[data-testid="payment-amount"]', String(remainingAmount + 1000));
    await page.selectOption('[data-testid="payment-method"]', 'bank_transfer');
    
    // Submit payment - should show error
    await page.click('[data-testid="process-payment-submit"]');
    
    // Should show validation error
    await expect(page.locator('[data-testid="error-message"]')).toContainText('付款金額超過剩餘金額');
  });
});