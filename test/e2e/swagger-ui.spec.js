// Playwright E2E tests for Swagger UI
// Run with: npx playwright test
const { test, expect } = require('@playwright/test');

test.describe('Hefestus API Swagger UI Tests', () => {
  const BASE_URL = 'http://localhost:8080';

  test.beforeEach(async ({ page }) => {
    // Navigate to Swagger UI
    await page.goto(`${BASE_URL}/swagger/index.html`);
    await page.waitForLoadState('networkidle');
  });

  test('Swagger UI loads successfully', async ({ page }) => {
    // Check that Swagger UI title is present
    await expect(page).toHaveTitle(/Swagger UI/);
    
    // Check for the main Swagger container
    await expect(page.locator('.swagger-ui')).toBeVisible();
    
    // Verify API title is displayed
    await expect(page.locator('.title')).toContainText('Hefestus API');
  });

  test('API models are correctly defined in Swagger', async ({ page }) => {
    // Wait for the schema section to load
    await page.waitForSelector('.models');
    
    // Check that our refactored models are present
    const modelsSection = page.locator('.models');
    await expect(modelsSection).toBeVisible();
    
    // Verify ErrorSolution model with correct field names
    const errorSolutionModel = page.locator('[data-name="ErrorSolution"]');
    await expect(errorSolutionModel).toBeVisible();
    
    // Click to expand the model if needed
    if (await errorSolutionModel.locator('.model-toggle').isVisible()) {
      await errorSolutionModel.locator('.model-toggle').click();
    }
    
    // Check for Portuguese field names (backward compatibility)
    await expect(page.locator('text=causa')).toBeVisible();
    await expect(page.locator('text=solucao')).toBeVisible();
  });

  test('Health check endpoint works via Swagger UI', async ({ page }) => {
    // Find and click the health check endpoint
    const healthEndpoint = page.locator('[data-path="/health"][data-method="get"]');
    await expect(healthEndpoint).toBeVisible();
    await healthEndpoint.click();
    
    // Click "Try it out" button
    const tryItOutBtn = page.locator('button:has-text("Try it out")').first();
    await tryItOutBtn.click();
    
    // Click "Execute" button
    const executeBtn = page.locator('button:has-text("Execute")').first();
    await executeBtn.click();
    
    // Wait for response
    await page.waitForSelector('.responses-wrapper .live-responses-table');
    
    // Check response status
    const responseCode = page.locator('.response-col_status');
    await expect(responseCode).toContainText('200');
    
    // Check response body contains status: ok
    const responseBody = page.locator('.response-col_description .microlight');
    await expect(responseBody).toContainText('"status"');
    await expect(responseBody).toContainText('"ok"');
  });

  test('Error analysis endpoint structure is correct', async ({ page }) => {
    // Find the error analysis endpoint
    const errorEndpoint = page.locator('[data-path="/errors/{domain}"][data-method="post"]');
    await expect(errorEndpoint).toBeVisible();
    await errorEndpoint.click();
    
    // Check endpoint parameters
    const domainParam = page.locator('.parameter .parameter-name:has-text("domain")');
    await expect(domainParam).toBeVisible();
    
    // Check request body schema
    const requestBodySection = page.locator('.request-body');
    await expect(requestBodySection).toBeVisible();
    
    // Verify ErrorRequest model is referenced
    await expect(page.locator('text=ErrorRequest')).toBeVisible();
    
    // Check response schema
    const responsesSection = page.locator('.responses');
    await expect(responsesSection).toBeVisible();
    
    // Verify ErrorResponse model for 200 response
    await expect(page.locator('text=ErrorResponse')).toBeVisible();
    
    // Verify APIError model for error responses
    await expect(page.locator('text=APIError')).toBeVisible();
  });

  test('Models schema validation shows correct types', async ({ page }) => {
    // Navigate to models section
    await page.locator('text=Models').click();
    
    // Check ErrorSolution model details
    const errorSolutionSection = page.locator('[data-name="ErrorSolution"]');
    await errorSolutionSection.click();
    
    // Verify field types
    const causaField = page.locator('.model-property').filter({ hasText: 'causa' });
    await expect(causaField).toContainText('string');
    
    const solucaoField = page.locator('.model-property').filter({ hasText: 'solucao' });
    await expect(solucaoField).toContainText('string');
    
    // Check required fields indicator
    await expect(causaField.locator('.star')).toBeVisible(); // Required field indicator
    await expect(solucaoField.locator('.star')).toBeVisible(); // Required field indicator
  });

  test('API documentation shows correct examples', async ({ page }) => {
    // Expand the error analysis endpoint
    const errorEndpoint = page.locator('[data-path="/errors/{domain}"][data-method="post"]');
    await errorEndpoint.click();
    
    // Check that example values are shown
    await expect(page.locator('text=error_details')).toBeVisible();
    await expect(page.locator('text=context')).toBeVisible();
    
    // Look for example error details
    await expect(page.locator('text=CrashLoopBackOff')).toBeVisible();
    
    // Check response examples
    const responseExample = page.locator('.example-value');
    await expect(responseExample.first()).toBeVisible();
  });

  test('Domain validation shows correct enum values', async ({ page }) => {
    // Expand the error analysis endpoint
    const errorEndpoint = page.locator('[data-path="/errors/{domain}"][data-method="post"]');
    await errorEndpoint.click();
    
    // Click "Try it out" to see parameter details
    const tryItOutBtn = page.locator('button:has-text("Try it out")').first();
    await tryItOutBtn.click();
    
    // Check domain parameter dropdown/validation
    const domainSelect = page.locator('select[data-param-name="domain"]');
    if (await domainSelect.isVisible()) {
      // Check that valid domains are available
      await expect(domainSelect.locator('option[value="kubernetes"]')).toBeVisible();
      await expect(domainSelect.locator('option[value="github"]')).toBeVisible();
      await expect(domainSelect.locator('option[value="argocd"]')).toBeVisible();
    }
  });
});

// Configuration for Playwright
module.exports = {
  testDir: './test/e2e',
  timeout: 30000,
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:8080',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: {
    command: 'go run cmd/server/main.go',
    url: 'http://localhost:8080/api/health',
    reuseExistingServer: !process.env.CI,
    timeout: 120 * 1000,
  },
};