import { test, expect } from '@playwright/test';

test.describe('Quote Management', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/login');
    await page.fill('[data-testid="username-input"]', 'engineer1');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');
    await page.waitForURL('**/dashboard');
    
    // Navigate to quotes page
    await page.goto('/quotes');
  });

  test('should display quotes list page', async ({ page }) => {
    // Check page title
    await expect(page.locator('h1')).toContainText('報價管理');
    
    // Check if quote list is visible
    await expect(page.locator('[data-testid="quotes-table"]')).toBeVisible();
    
    // Check if create button is visible
    await expect(page.locator('[data-testid="create-quote-button"]')).toBeVisible();
    
    // Check if search functionality is available
    await expect(page.locator('[data-testid="search-input"]')).toBeVisible();
  });

  test('should filter quotes by status', async ({ page }) => {
    // Wait for quotes to load
    await page.waitForSelector('[data-testid="quotes-table"]');
    
    // Click on status filter
    await page.click('[data-testid="status-filter"]');
    
    // Select 'draft' status
    await page.click('[data-testid="status-draft"]');
    
    // Check if only draft quotes are shown
    const statusBadges = page.locator('[data-testid="quote-status"]');
    const count = await statusBadges.count();
    
    for (let i = 0; i < count; i++) {
      await expect(statusBadges.nth(i)).toContainText('草稿');
    }
  });

  test('should search quotes by quote number', async ({ page }) => {
    // Wait for quotes to load
    await page.waitForSelector('[data-testid="quotes-table"]');
    
    // Search for a specific quote
    await page.fill('[data-testid="search-input"]', 'QUO-2024');
    await page.keyboard.press('Enter');
    
    // Wait for search results
    await page.waitForTimeout(1000);
    
    // Check if results contain the search term
    const quoteNumbers = page.locator('[data-testid="quote-number"]');
    const count = await quoteNumbers.count();
    
    if (count > 0) {
      for (let i = 0; i < count; i++) {
        await expect(quoteNumbers.nth(i)).toContainText('QUO-2024');
      }
    }
  });

  test('should navigate to quote detail page', async ({ page }) => {
    // Wait for quotes to load
    await page.waitForSelector('[data-testid="quotes-table"]');
    
    // Click on first quote row
    const firstQuoteRow = page.locator('[data-testid="quote-row"]').first();
    await firstQuoteRow.click();
    
    // Should navigate to quote detail page
    await page.waitForURL('**/quotes/**');
    
    // Check if quote detail page is loaded
    await expect(page.locator('h1')).toContainText('報價單詳情');
    await expect(page.locator('[data-testid="quote-details"]')).toBeVisible();
  });

  test('should create new quote', async ({ page }) => {
    // Click create quote button
    await page.click('[data-testid="create-quote-button"]');
    
    // Should navigate to create quote page
    await page.waitForURL('**/quotes/new');
    
    // Check if create quote form is visible
    await expect(page.locator('h1')).toContainText('建立報價單');
    await expect(page.locator('[data-testid="quote-form"]')).toBeVisible();
    
    // Fill in quote details
    await page.selectOption('[data-testid="inquiry-select"]', { index: 1 });
    await page.fill('[data-testid="material-cost"]', '1000');
    await page.fill('[data-testid="process-cost"]', '500');
    await page.fill('[data-testid="surface-cost"]', '200');
    await page.fill('[data-testid="packaging-cost"]', '100');
    await page.fill('[data-testid="shipping-cost"]', '300');
    
    // Set rates
    await page.fill('[data-testid="overhead-rate"]', '15');
    await page.fill('[data-testid="profit-rate"]', '20');
    
    // Set delivery and payment terms
    await page.fill('[data-testid="delivery-days"]', '30');
    await page.fill('[data-testid="payment-terms"]', 'T/T 30 days');
    
    // Add notes
    await page.fill('[data-testid="notes"]', 'Test quote creation');
    
    // Submit form
    await page.click('[data-testid="create-quote-submit"]');
    
    // Should redirect to quote detail page
    await page.waitForURL('**/quotes/**');
    
    // Check if quote was created successfully
    await expect(page.locator('[data-testid="success-message"]')).toContainText('報價單建立成功');
  });

  test('should edit quote', async ({ page }) => {
    // Navigate to a draft quote detail page
    await page.goto('/quotes');
    await page.waitForSelector('[data-testid="quotes-table"]');
    
    // Find and click on a draft quote
    const draftQuote = page.locator('[data-testid="quote-row"]').filter({ hasText: '草稿' }).first();
    await draftQuote.click();
    
    await page.waitForURL('**/quotes/**');
    
    // Click edit button
    await page.click('[data-testid="edit-quote-button"]');
    
    // Should navigate to edit page
    await page.waitForURL('**/quotes/**/edit');
    
    // Check if edit form is visible
    await expect(page.locator('h1')).toContainText('編輯報價單');
    await expect(page.locator('[data-testid="quote-edit-form"]')).toBeVisible();
    
    // Update some fields
    await page.fill('[data-testid="material-cost"]', '1200');
    await page.fill('[data-testid="notes"]', 'Updated quote');
    
    // Submit changes
    await page.click('[data-testid="save-quote-button"]');
    
    // Should redirect back to detail page
    await page.waitForURL('**/quotes/**', { timeout: 10000 });
    
    // Check if changes were saved
    await expect(page.locator('[data-testid="success-message"]')).toContainText('報價單更新成功');
  });

  test('should submit quote for approval', async ({ page }) => {
    // Navigate to a draft quote
    await page.goto('/quotes');
    await page.waitForSelector('[data-testid="quotes-table"]');
    
    const draftQuote = page.locator('[data-testid="quote-row"]').filter({ hasText: '草稿' }).first();
    await draftQuote.click();
    
    await page.waitForURL('**/quotes/**');
    
    // Click submit for approval button
    await page.click('[data-testid="submit-approval-button"]');
    
    // Check if status changed
    await expect(page.locator('[data-testid="quote-status"]')).toContainText('待審核');
    
    // Check success message
    await expect(page.locator('[data-testid="success-message"]')).toContainText('已送出審核');
  });

  test('should approve quote (as manager)', async ({ page }) => {
    // Logout and login as manager
    await page.click('[data-testid="user-menu"]');
    await page.click('[data-testid="logout-button"]');
    
    await page.waitForURL('**/login');
    await page.fill('[data-testid="username-input"]', 'manager1');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');
    await page.waitForURL('**/dashboard');
    
    // Navigate to quotes page
    await page.goto('/quotes');
    await page.waitForSelector('[data-testid="quotes-table"]');
    
    // Find a pending approval quote
    const pendingQuote = page.locator('[data-testid="quote-row"]').filter({ hasText: '待審核' }).first();
    await pendingQuote.click();
    
    await page.waitForURL('**/quotes/**');
    
    // Click approve button
    await page.click('[data-testid="approve-button"]');
    
    // Fill approval dialog
    await page.fill('[data-testid="approval-notes"]', 'Approved - looks good');
    await page.click('[data-testid="confirm-approve"]');
    
    // Check if status changed
    await expect(page.locator('[data-testid="quote-status"]')).toContainText('已核准');
  });

  test('should send quote to customer', async ({ page }) => {
    // Navigate to an approved quote
    await page.goto('/quotes');
    await page.waitForSelector('[data-testid="quotes-table"]');
    
    const approvedQuote = page.locator('[data-testid="quote-row"]').filter({ hasText: '已核准' }).first();
    await approvedQuote.click();
    
    await page.waitForURL('**/quotes/**');
    
    // Click send quote button
    await page.click('[data-testid="send-quote-button"]');
    
    // Fill send dialog
    await page.fill('[data-testid="email-message"]', 'Please find attached our quote for your consideration.');
    await page.click('[data-testid="confirm-send"]');
    
    // Check if status changed
    await expect(page.locator('[data-testid="quote-status"]')).toContainText('已發送');
    
    // Check success message
    await expect(page.locator('[data-testid="success-message"]')).toContainText('報價單已發送');
  });

  test('should download quote PDF', async ({ page }) => {
    // Navigate to a quote detail page
    await page.goto('/quotes');
    await page.waitForSelector('[data-testid="quotes-table"]');
    
    const firstQuote = page.locator('[data-testid="quote-row"]').first();
    await firstQuote.click();
    
    await page.waitForURL('**/quotes/**');
    
    // Setup download promise before clicking
    const downloadPromise = page.waitForEvent('download');
    
    // Click download PDF button
    await page.click('[data-testid="download-pdf-button"]');
    
    // Wait for download to start
    const download = await downloadPromise;
    
    // Check if file is downloaded
    expect(download.suggestedFilename()).toMatch(/quote.*\.pdf/i);
  });

  test('should show quote history', async ({ page }) => {
    // Navigate to a quote detail page
    await page.goto('/quotes');
    await page.waitForSelector('[data-testid="quotes-table"]');
    
    const firstQuote = page.locator('[data-testid="quote-row"]').first();
    await firstQuote.click();
    
    await page.waitForURL('**/quotes/**');
    
    // Click on timeline tab
    await page.click('[data-testid="timeline-tab"]');
    
    // Check if timeline is visible
    await expect(page.locator('[data-testid="quote-timeline"]')).toBeVisible();
    
    // Check if timeline has activities
    const activities = page.locator('[data-testid="timeline-activity"]');
    const count = await activities.count();
    expect(count).toBeGreaterThan(0);
  });

  test('should handle quote permissions by role', async ({ page }) => {
    // Test with sales role (should have limited access)
    await page.click('[data-testid="user-menu"]');
    await page.click('[data-testid="logout-button"]');
    
    await page.waitForURL('**/login');
    await page.fill('[data-testid="username-input"]', 'sales1');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');
    await page.waitForURL('**/dashboard');
    
    // Navigate to quotes page
    await page.goto('/quotes');
    
    // Sales should not see create quote button
    await expect(page.locator('[data-testid="create-quote-button"]')).not.toBeVisible();
    
    // Click on a quote to view details
    const firstQuote = page.locator('[data-testid="quote-row"]').first();
    await firstQuote.click();
    
    await page.waitForURL('**/quotes/**');
    
    // Sales should not see edit or approval buttons
    await expect(page.locator('[data-testid="edit-quote-button"]')).not.toBeVisible();
    await expect(page.locator('[data-testid="approve-button"]')).not.toBeVisible();
  });
});