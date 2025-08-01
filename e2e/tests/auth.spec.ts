import { test, expect } from '@playwright/test';

test.describe('Authentication', () => {
  test.beforeEach(async ({ page }) => {
    // Go to login page
    await page.goto('/login');
  });

  test('should display login form', async ({ page }) => {
    // Check if login form is visible
    await expect(page.locator('[data-testid="login-form"]')).toBeVisible();
    await expect(page.locator('[data-testid="username-input"]')).toBeVisible();
    await expect(page.locator('[data-testid="password-input"]')).toBeVisible();
    await expect(page.locator('[data-testid="login-button"]')).toBeVisible();
    
    // Check page title
    await expect(page).toHaveTitle(/FastenMind/);
    
    // Check login header
    await expect(page.locator('h1')).toContainText('登入 FastenMind');
  });

  test('should show validation errors for empty fields', async ({ page }) => {
    // Click login button without filling fields
    await page.click('[data-testid="login-button"]');
    
    // Check for validation errors
    await expect(page.locator('text=請輸入帳號')).toBeVisible();
    await expect(page.locator('text=請輸入密碼')).toBeVisible();
  });

  test('should show error for invalid credentials', async ({ page }) => {
    // Fill in invalid credentials
    await page.fill('[data-testid="username-input"]', 'invalid_user');
    await page.fill('[data-testid="password-input"]', 'wrong_password');
    
    // Click login button
    await page.click('[data-testid="login-button"]');
    
    // Wait for error message
    await expect(page.locator('text=帳號或密碼錯誤')).toBeVisible();
  });

  test('should login successfully with valid credentials', async ({ page }) => {
    // Fill in valid credentials
    await page.fill('[data-testid="username-input"]', 'admin');
    await page.fill('[data-testid="password-input"]', 'password123');
    
    // Click login button
    await page.click('[data-testid="login-button"]');
    
    // Wait for redirect to dashboard
    await page.waitForURL('**/dashboard');
    
    // Check if we're on dashboard
    await expect(page.locator('h1')).toContainText('儀表板');
    
    // Check if user menu is visible
    await expect(page.locator('[data-testid="user-menu"]')).toBeVisible();
  });

  test('should handle password visibility toggle', async ({ page }) => {
    const passwordInput = page.locator('[data-testid="password-input"]');
    const toggleButton = page.locator('[data-testid="password-toggle"]');
    
    // Initially password should be hidden
    await expect(passwordInput).toHaveAttribute('type', 'password');
    
    // Click toggle to show password
    await toggleButton.click();
    await expect(passwordInput).toHaveAttribute('type', 'text');
    
    // Click toggle to hide password again
    await toggleButton.click();
    await expect(passwordInput).toHaveAttribute('type', 'password');
  });

  test('should redirect authenticated user away from login', async ({ page }) => {
    // First login
    await page.fill('[data-testid="username-input"]', 'admin');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');
    
    // Wait for redirect to dashboard
    await page.waitForURL('**/dashboard');
    
    // Try to go back to login page
    await page.goto('/login');
    
    // Should be redirected to dashboard
    await expect(page).toHaveURL(/\/dashboard/);
  });

  test('should handle logout', async ({ page }) => {
    // Login first
    await page.fill('[data-testid="username-input"]', 'admin');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');
    
    await page.waitForURL('**/dashboard');
    
    // Open user menu and logout
    await page.click('[data-testid="user-menu"]');
    await page.click('[data-testid="logout-button"]');
    
    // Should be redirected to login page
    await page.waitForURL('**/login');
    await expect(page.locator('[data-testid="login-form"]')).toBeVisible();
  });

  test('should handle session expiry', async ({ page }) => {
    // Login first
    await page.fill('[data-testid="username-input"]', 'admin');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');
    
    await page.waitForURL('**/dashboard');
    
    // Simulate session expiry by clearing localStorage
    await page.evaluate(() => {
      localStorage.clear();
      sessionStorage.clear();
    });
    
    // Try to navigate to a protected route
    await page.goto('/quotes');
    
    // Should be redirected to login
    await page.waitForURL('**/login');
    await expect(page.locator('[data-testid="login-form"]')).toBeVisible();
  });

  test('should remember login state after page refresh', async ({ page }) => {
    // Login first
    await page.fill('[data-testid="username-input"]', 'admin');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');
    
    await page.waitForURL('**/dashboard');
    
    // Refresh the page
    await page.reload();
    
    // Should still be on dashboard
    await expect(page).toHaveURL(/\/dashboard/);
    await expect(page.locator('h1')).toContainText('儀表板');
  });

  test('should handle keyboard navigation', async ({ page }) => {
    // Focus should start on username input
    await page.keyboard.press('Tab');
    await expect(page.locator('[data-testid="username-input"]')).toBeFocused();
    
    // Tab to password input
    await page.keyboard.press('Tab');
    await expect(page.locator('[data-testid="password-input"]')).toBeFocused();
    
    // Tab to login button
    await page.keyboard.press('Tab');
    await expect(page.locator('[data-testid="login-button"]')).toBeFocused();
    
    // Enter should submit form
    await page.fill('[data-testid="username-input"]', 'admin');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.keyboard.press('Enter');
    
    await page.waitForURL('**/dashboard');
    await expect(page.locator('h1')).toContainText('儀表板');
  });
});