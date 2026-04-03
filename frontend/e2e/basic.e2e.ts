// e2e/basic.e2e.ts
import { test, expect } from '@playwright/test';

// 基础端到端测试
test.describe('OpenCMP IAM Module End-to-End Tests', () => {
  
  test('should allow user login', async ({ page }) => {
    // 访问登录页面
    await page.goto('http://localhost:8080/login');
    
    // 输入凭据
    await page.locator('input[type="text"]').fill('admin');
    await page.locator('input[type="password"]').fill('admin123');
    
    // 点击登录按钮
    await page.locator('button').click();
    
    // 验证登录成功
    await expect(page).toHaveURL(/.*dashboard/);
    await expect(page.locator('.user-info .username')).toContainText('admin');
  });

  test('should display user management page', async ({ page }) => {
    // 先登录
    await page.goto('http://localhost:8080/login');
    await page.locator('input[type="text"]').fill('admin');
    await page.locator('input[type="password"]').fill('admin123');
    await page.locator('button').click();
    
    // 导航到用户管理页面
    await page.getByText('IAM').click();
    await page.getByText('用户管理').click();
    
    // 验证页面加载
    await expect(page.locator('.users-page')).toBeVisible();
    await expect(page.locator('.el-table')).toBeVisible();
  });

  test('should allow creating a new user', async ({ page }) => {
    // 先登录
    await page.goto('http://localhost:8080/login');
    await page.locator('input[type="text"]').fill('admin');
    await page.locator('input[type="password"]').fill('admin123');
    await page.locator('button').click();
    
    // 导航到用户管理页面
    await page.getByText('IAM').click();
    await page.getByText('用户管理').click();
    
    // 点击新增用户按钮
    await page.locator('.el-button--primary').first().click();
    
    // 填写用户信息
    await page.locator('[data-testid="user-name-input"]').fill('testuser');
    await page.locator('[data-testid="user-email-input"]').fill('test@example.com');
    await page.locator('[data-testid="user-password-input"]').fill('password123');
    
    // 提交表单
    await page.locator('[data-testid="submit-button"]').click();
    
    // 验证用户创建成功
    await expect(page.locator('.el-message__content')).toContainText('创建成功');
  });

  test('should validate permissions correctly', async ({ page }) => {
    // 测试权限验证功能
    await page.goto('http://localhost:8080/login');
    await page.locator('input[type="text"]').fill('admin');
    await page.locator('input[type="password"]').fill('admin123');
    await page.locator('button').click();
    
    // 导航到权限管理页面
    await page.getByText('IAM').click();
    await page.getByText('权限管理').click();
    
    // 验证权限列表加载
    await expect(page.locator('.permissions-page')).toBeVisible();
    await expect(page.locator('.el-table__body tr')).toHaveCount({ min: 1 });
  });
});