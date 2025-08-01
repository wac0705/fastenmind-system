import { test, expect } from '@playwright/test';

test.describe('Order Management', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/login');
    await page.fill('[data-testid="username-input"]', 'sales1');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');
    await page.waitForURL('**/dashboard');
    
    // Navigate to orders page
    await page.goto('/orders');
  });

  test('should display orders list page', async ({ page }) => {
    // Check page title
    await expect(page.locator('h1')).toContainText('訂單管理');
    
    // Check if order list is visible
    await expect(page.locator('[data-testid="orders-table"]')).toBeVisible();
    
    // Check if create button is visible (for sales role)
    await expect(page.locator('[data-testid="create-order-button"]')).toBeVisible();
    
    // Check if search functionality is available
    await expect(page.locator('[data-testid="search-input"]')).toBeVisible();
  });

  test('should filter orders by status', async ({ page }) => {
    // Wait for orders to load
    await page.waitForSelector('[data-testid="orders-table"]');
    
    // Click on status filter
    await page.click('[data-testid="status-filter"]');
    
    // Select 'confirmed' status
    await page.click('[data-testid="status-confirmed"]');
    
    // Check if only confirmed orders are shown
    const statusBadges = page.locator('[data-testid="order-status"]');
    const count = await statusBadges.count();
    
    for (let i = 0; i < count; i++) {
      await expect(statusBadges.nth(i)).toContainText('已確認');
    }
  });

  test('should search orders by order number', async ({ page }) => {
    // Wait for orders to load
    await page.waitForSelector('[data-testid="orders-table"]');
    
    // Search for a specific order
    await page.fill('[data-testid="search-input"]', 'ORD-2024');
    await page.keyboard.press('Enter');
    
    // Wait for search results
    await page.waitForTimeout(1000);
    
    // Check if results contain the search term
    const orderNumbers = page.locator('[data-testid="order-number"]');
    const count = await orderNumbers.count();
    
    if (count > 0) {
      for (let i = 0; i < count; i++) {
        await expect(orderNumbers.nth(i)).toContainText('ORD-2024');
      }
    }
  });

  test('should navigate to order detail page', async ({ page }) => {
    // Wait for orders to load
    await page.waitForSelector('[data-testid="orders-table"]');
    
    // Click on first order row
    const firstOrderRow = page.locator('[data-testid="order-row"]').first();
    await firstOrderRow.click();
    
    // Should navigate to order detail page
    await page.waitForURL('**/orders/**');
    
    // Check if order detail page is loaded
    await expect(page.locator('h1')).toContainText('訂單詳情');
    await expect(page.locator('[data-testid="order-details"]')).toBeVisible();
  });

  test('should create new order from accepted quote', async ({ page }) => {
    // First need to have an accepted quote
    await page.goto('/quotes');
    await page.waitForSelector('[data-testid="quotes-table"]');
    
    // Find an accepted quote and click on it
    const acceptedQuote = page.locator('[data-testid="quote-row"]').filter({ hasText: '已接受' }).first();
    
    // If no accepted quote exists, create one first
    const quoteCount = await acceptedQuote.count();
    if (quoteCount === 0) {
      // Skip this test if no accepted quotes exist
      test.skip();
      return;
    }
    
    await acceptedQuote.click();
    await page.waitForURL('**/quotes/**');
    
    // Click create order button
    await page.click('[data-testid="create-order-button"]');
    
    // Should navigate to create order page
    await page.waitForURL('**/orders/new**');
    
    // Check if create order form is visible
    await expect(page.locator('h1')).toContainText('建立訂單');
    await expect(page.locator('[data-testid="order-form"]')).toBeVisible();
    
    // Fill in order details
    await page.fill('[data-testid="po-number"]', 'PO-2024-TEST-001');
    
    // Set delivery date (30 days from now)
    const futureDate = new Date();
    futureDate.setDate(futureDate.getDate() + 30);
    const dateString = futureDate.toISOString().split('T')[0];
    await page.fill('[data-testid="delivery-date"]', dateString);
    
    // Select delivery method
    await page.selectOption('[data-testid="delivery-method"]', '海運');
    
    // Fill shipping address
    await page.fill('[data-testid="shipping-address"]', '123 Test Street, Test City, Test Country');
    
    // Set payment terms
    await page.fill('[data-testid="payment-terms"]', 'T/T 30 days');
    
    // Set down payment
    await page.fill('[data-testid="down-payment"]', '3000');
    
    // Add notes
    await page.fill('[data-testid="notes"]', 'E2E test order creation');
    
    // Submit form
    await page.click('[data-testid="create-order-submit"]');
    
    // Should redirect to order detail page
    await page.waitForURL('**/orders/**');
    
    // Check if order was created successfully
    await expect(page.locator('[data-testid="success-message"]')).toContainText('訂單建立成功');
    await expect(page.locator('[data-testid="po-number"]')).toContainText('PO-2024-TEST-001');
  });

  test('should update order status through workflow', async ({ page }) => {
    // Navigate to an order detail page
    await page.goto('/orders');
    await page.waitForSelector('[data-testid="orders-table"]');
    
    // Find a pending order
    const pendingOrder = page.locator('[data-testid="order-row"]').filter({ hasText: '待確認' }).first();
    const orderCount = await pendingOrder.count();
    
    if (orderCount === 0) {
      test.skip(); // Skip if no pending orders exist
      return;
    }
    
    await pendingOrder.click();
    await page.waitForURL('**/orders/**');
    
    // Should see update status button
    await expect(page.locator('[data-testid="update-status-button"]')).toBeVisible();
    
    // Click update status
    await page.click('[data-testid="update-status-button"]');
    
    // Select new status (confirmed)
    await page.selectOption('[data-testid="status-select"]', 'confirmed');
    
    // Add notes
    await page.fill('[data-testid="status-notes"]', 'Order confirmed by customer');
    
    // Submit status update
    await page.click('[data-testid="confirm-status-update"]');
    
    // Check if status was updated
    await expect(page.locator('[data-testid="order-status"]')).toContainText('已確認');
    await expect(page.locator('[data-testid="success-message"]')).toContainText('訂單狀態已更新');
  });

  test('should view order items', async ({ page }) => {
    // Navigate to an order detail page
    await page.goto('/orders');
    await page.waitForSelector('[data-testid="orders-table"]');
    
    const firstOrder = page.locator('[data-testid="order-row"]').first();
    await firstOrder.click();
    await page.waitForURL('**/orders/**');
    
    // Click on items tab
    await page.click('[data-testid="items-tab"]');
    
    // Check if items table is visible
    await expect(page.locator('[data-testid="order-items-table"]')).toBeVisible();
    
    // Check if items are displayed with required columns
    await expect(page.locator('th')).toContainText('料號');
    await expect(page.locator('th')).toContainText('描述');
    await expect(page.locator('th')).toContainText('數量');
    await expect(page.locator('th')).toContainText('單價');
    await expect(page.locator('th')).toContainText('總價');
  });

  test('should view order timeline', async ({ page }) => {
    // Navigate to an order detail page
    await page.goto('/orders');
    await page.waitForSelector('[data-testid="orders-table"]');
    
    const firstOrder = page.locator('[data-testid="order-row"]').first();
    await firstOrder.click();
    await page.waitForURL('**/orders/**');
    
    // Click on timeline tab
    await page.click('[data-testid="timeline-tab"]');
    
    // Check if timeline is visible
    await expect(page.locator('[data-testid="order-timeline"]')).toBeVisible();
    
    // Check if timeline has activities
    const activities = page.locator('[data-testid="timeline-activity"]');
    const count = await activities.count();
    expect(count).toBeGreaterThan(0);
  });

  test('should export orders', async ({ page }) => {
    // Wait for orders to load
    await page.waitForSelector('[data-testid="orders-table"]');
    
    // Setup download promise before clicking
    const downloadPromise = page.waitForEvent('download');
    
    // Click export button
    await page.click('[data-testid="export-orders-button"]');
    
    // Wait for download to start
    const download = await downloadPromise;
    
    // Check if file is downloaded
    expect(download.suggestedFilename()).toMatch(/orders.*\.csv/i);
  });

  test('should handle order permissions by role', async ({ page }) => {
    // Test with manager role (should have full access)
    await page.click('[data-testid="user-menu"]');
    await page.click('[data-testid="logout-button"]');
    
    await page.waitForURL('**/login');
    await page.fill('[data-testid="username-input"]', 'manager1');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');
    await page.waitForURL('**/dashboard');
    
    // Navigate to orders page
    await page.goto('/orders');
    
    // Manager should see create order button
    await expect(page.locator('[data-testid="create-order-button"]')).toBeVisible();
    
    // Click on an order to view details
    const firstOrder = page.locator('[data-testid="order-row"]').first();
    await firstOrder.click();
    await page.waitForURL('**/orders/**');
    
    // Manager should see edit and status update buttons
    await expect(page.locator('[data-testid="edit-order-button"]')).toBeVisible();
    await expect(page.locator('[data-testid="update-status-button"]')).toBeVisible();
  });

  test('should restrict access for viewer role', async ({ page }) => {
    // Test with viewer role (should have limited access)
    await page.click('[data-testid="user-menu"]');
    await page.click('[data-testid="logout-button"]');
    
    await page.waitForURL('**/login');
    await page.fill('[data-testid="username-input"]', 'viewer1');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');
    await page.waitForURL('**/dashboard');
    
    // Navigate to orders page
    await page.goto('/orders');
    
    // Viewer should not see create order button
    await expect(page.locator('[data-testid="create-order-button"]')).not.toBeVisible();
    
    // Click on an order to view details
    const firstOrder = page.locator('[data-testid="order-row"]').first();
    await firstOrder.click();
    await page.waitForURL('**/orders/**');
    
    // Viewer should not see edit or update buttons
    await expect(page.locator('[data-testid="edit-order-button"]')).not.toBeVisible();
    await expect(page.locator('[data-testid="update-status-button"]')).not.toBeVisible();
  });

  test('should handle order editing', async ({ page }) => {
    // Navigate to an order in draft/confirmed status
    await page.goto('/orders');
    await page.waitForSelector('[data-testid="orders-table"]');
    
    const editableOrder = page.locator('[data-testid="order-row"]').filter({ hasText: '已確認' }).first();
    const orderCount = await editableOrder.count();
    
    if (orderCount === 0) {
      test.skip(); // Skip if no editable orders exist
      return;
    }
    
    await editableOrder.click();
    await page.waitForURL('**/orders/**');
    
    // Click edit button
    await page.click('[data-testid="edit-order-button"]');
    
    // Should navigate to edit page
    await page.waitForURL('**/orders/**/edit');
    
    // Check if edit form is visible
    await expect(page.locator('h1')).toContainText('編輯訂單');
    await expect(page.locator('[data-testid="order-edit-form"]')).toBeVisible();
    
    // Update some fields
    await page.fill('[data-testid="shipping-address"]', 'Updated shipping address');
    await page.fill('[data-testid="notes"]', 'Updated order notes');
    
    // Submit changes
    await page.click('[data-testid="save-order-button"]');
    
    // Should redirect back to detail page
    await page.waitForURL('**/orders/**');
    
    // Check if changes were saved
    await expect(page.locator('[data-testid="success-message"]')).toContainText('訂單更新成功');
  });

  test('should validate order status transitions', async ({ page }) => {
    // Navigate to a completed order
    await page.goto('/orders');
    await page.waitForSelector('[data-testid="orders-table"]');
    
    const completedOrder = page.locator('[data-testid="order-row"]').filter({ hasText: '已完成' }).first();
    const orderCount = await completedOrder.count();
    
    if (orderCount === 0) {
      test.skip(); // Skip if no completed orders exist
      return;
    }
    
    await completedOrder.click();
    await page.waitForURL('**/orders/**');
    
    // Update status button should not be visible for completed orders
    await expect(page.locator('[data-testid="update-status-button"]')).not.toBeVisible();
  });

  test('should show payment information', async ({ page }) => {
    // Navigate to an order detail page
    await page.goto('/orders');
    await page.waitForSelector('[data-testid="orders-table"]');
    
    const firstOrder = page.locator('[data-testid="order-row"]').first();
    await firstOrder.click();
    await page.waitForURL('**/orders/**');
    
    // Check if payment information is displayed
    await expect(page.locator('[data-testid="payment-status"]')).toBeVisible();
    await expect(page.locator('[data-testid="total-amount"]')).toBeVisible();
    await expect(page.locator('[data-testid="down-payment"]')).toBeVisible();
    await expect(page.locator('[data-testid="paid-amount"]')).toBeVisible();
    
    // Check payment terms
    await expect(page.locator('[data-testid="payment-terms"]')).toBeVisible();
  });

  test('should handle mobile responsive design', async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });
    
    // Navigate to orders page
    await page.goto('/orders');
    await page.waitForSelector('[data-testid="orders-table"]');
    
    // On mobile, table should be scrollable or cards should be used
    const isMobileView = await page.locator('[data-testid="mobile-order-cards"]').isVisible().catch(() => false);
    const isTableScrollable = await page.locator('[data-testid="orders-table"]').isVisible();
    
    // Either mobile cards or scrollable table should be present
    expect(isMobileView || isTableScrollable).toBeTruthy();
    
    // Navigation should be accessible
    await expect(page.locator('[data-testid="mobile-nav"]')).toBeVisible();
  });
});