import { chromium, FullConfig } from '@playwright/test';

async function globalTeardown(config: FullConfig) {
  console.log('🧹 Starting global teardown...');
  
  const { baseURL } = config.projects[0].use;
  
  // Start browser for teardown
  const browser = await chromium.launch();
  const context = await browser.newContext();
  const page = await context.newPage();
  
  try {
    // Clean up test data
    await cleanupTestData(page, baseURL);
    
    console.log('✅ Global teardown completed!');
  } catch (error) {
    console.error('❌ Global teardown failed:', error);
  } finally {
    await browser.close();
  }
}

async function cleanupTestData(page: any, baseURL: string) {
  console.log('🗑️  Cleaning up test data...');
  
  try {
    // Call cleanup API endpoint
    const response = await page.request.post(`${baseURL}/api/test/cleanup`, {
      data: {
        action: 'cleanup_test_data'
      }
    });
    
    if (response.ok()) {
      console.log('✅ Test data cleaned up successfully');
    }
  } catch (error) {
    console.log('⚠️  Test data cleanup skipped (API not available)');
  }
}

export default globalTeardown;