import { chromium, FullConfig } from '@playwright/test';

async function globalSetup(config: FullConfig) {
  console.log('🚀 Starting global setup...');
  
  const { baseURL } = config.projects[0].use;
  
  // Start browser for setup
  const browser = await chromium.launch();
  const context = await browser.newContext();
  const page = await context.newPage();
  
  try {
    // Wait for the application to be ready
    console.log('⏳ Waiting for application to be ready...');
    await page.goto(`${baseURL}/login`, { waitUntil: 'networkidle' });
    
    // Check if the application is running
    await page.waitForSelector('[data-testid="login-form"]', { timeout: 30000 });
    console.log('✅ Application is ready!');
    
    // Setup test data if needed
    await setupTestData(page, baseURL);
    
    // Login and save authentication state
    await loginAndSaveAuth(page, baseURL);
    
  } catch (error) {
    console.error('❌ Global setup failed:', error);
    throw error;
  } finally {
    await browser.close();
  }
  
  console.log('✅ Global setup completed!');
}

async function setupTestData(page: any, baseURL: string) {
  console.log('📊 Setting up test data...');
  
  // You can call API endpoints directly to set up test data
  // or use the page to interact with the UI
  
  try {
    // Example: Create test users, companies, etc.
    const response = await page.request.post(`${baseURL}/api/test/setup`, {
      data: {
        action: 'create_test_data'
      }
    });
    
    if (response.ok()) {
      console.log('✅ Test data created successfully');
    }
  } catch (error) {
    console.log('⚠️  Test data setup skipped (API not available)');
  }
}

async function loginAndSaveAuth(page: any, baseURL: string) {
  console.log('🔐 Setting up authentication...');
  
  try {
    await page.goto(`${baseURL}/login`);
    
    // Login with test user
    await page.fill('[data-testid="username-input"]', 'admin');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');
    
    // Wait for successful login
    await page.waitForURL('**/dashboard', { timeout: 10000 });
    
    // Save authentication state
    await page.context().storageState({ path: 'auth-state.json' });
    
    console.log('✅ Authentication state saved!');
  } catch (error) {
    console.log('⚠️  Authentication setup skipped:', error.message);
  }
}

export default globalSetup;