import { test, expect, type Page } from '@playwright/test';

// 测试数据
const TEST_USER = {
  username: 'testuser',
  password: 'test123456',
};

const INVALID_USER = {
  username: 'nonexistentuser',
  password: 'wrongpassword',
};

// 辅助函数：等待页面加载完成
async function waitForPageLoad(page: Page) {
  await page.waitForLoadState('networkidle');
}

// 辅助函数：清空 localStorage
async function clearAuthStorage(page: Page) {
  await page.evaluate(() => {
    localStorage.removeItem('junblog_token');
    localStorage.removeItem('junblog_refresh_token');
    localStorage.removeItem('junblog_user');
  });
}

test.describe('登录功能测试', () => {
  test.beforeEach(async ({ page }) => {
    // 每个测试前清空认证状态
    await page.goto('/login');
    await waitForPageLoad(page);
    await clearAuthStorage(page);
    await page.reload();
    await waitForPageLoad(page);
  });

  test.describe('页面加载测试', () => {
    test('应该正确加载登录页面', async ({ page }) => {
      // 验证页面标题
      await expect(page).toHaveTitle(/JunBlog/);

      // 验证登录表单元素存在
      await expect(page.locator('h2')).toContainText('登录');
      await expect(page.locator('input[type="text"]')).toBeVisible();
      await expect(page.locator('input[type="password"]')).toBeVisible();
      await expect(page.locator('button[type="submit"]')).toBeVisible();
    });

    test('应该显示返回首页链接', async ({ page }) => {
      const backHomeLink = page.locator('a.back-home');
      await expect(backHomeLink).toBeVisible();
      await expect(backHomeLink).toContainText('← 返回首页');
    });

    test('应该显示注册链接', async ({ page }) => {
      const registerLink = page.locator('a[href="/register"]');
      await expect(registerLink).toBeVisible();
      await expect(registerLink).toContainText('立即注册');
    });
  });

  test.describe('表单验证测试', () => {
    test('空用户名应该显示错误提示', async ({ page }) => {
      // 只输入密码
      await page.fill('input[type="password"]', TEST_USER.password);
      await page.click('button[type="submit"]');

      // 验证错误提示
      await expect(page.locator('.field-error')).toContainText('请输入用户名');
    });

    test('用户名少于2位应该显示错误提示', async ({ page }) => {
      // 输入太短的用户名
      await page.fill('input[type="text"]', 'a');
      await page.fill('input[type="password"]', TEST_USER.password);
      await page.click('button[type="submit"]');

      // 验证错误提示
      await expect(page.locator('.field-error')).toContainText('用户名至少2位');
    });

    test('空密码应该显示错误提示', async ({ page }) => {
      // 只输入用户名
      await page.fill('input[type="text"]', TEST_USER.username);
      await page.click('button[type="submit"]');

      // 验证错误提示
      await expect(page.locator('.field-error')).toContainText('请输入密码');
    });

    test('密码少于6位应该显示错误提示', async ({ page }) => {
      // 输入太短的密码
      await page.fill('input[type="text"]', TEST_USER.username);
      await page.fill('input[type="password"]', '12345');
      await page.click('button[type="submit"]');

      // 验证错误提示
      await expect(page.locator('.field-error')).toContainText('密码至少6位');
    });

    test('输入时应该清除对应的错误提示', async ({ page }) => {
      // 先触发错误
      await page.click('button[type="submit"]');
      await expect(page.locator('.field-error')).toBeVisible();

      // 输入用户名，应该清除用户名的错误
      await page.fill('input[type="text"]', TEST_USER.username);
      await expect(page.locator('.field:has(input[type="text"]) .field-error')).not.toBeVisible();
    });
  });

  test.describe('密码显示/隐藏测试', () => {
    test('点击眼睛按钮应该切换密码可见性', async ({ page }) => {
      const passwordInput = page.locator('input[type="password"]');
      const eyeButton = page.locator('.eye-btn');

      // 初始状态：密码隐藏
      await expect(passwordInput).toHaveAttribute('type', 'password');

      // 点击眼睛按钮，显示密码
      await eyeButton.click();
      await expect(passwordInput).toHaveAttribute('type', 'text');

      // 再次点击，隐藏密码
      await eyeButton.click();
      await expect(passwordInput).toHaveAttribute('type', 'password');
    });
  });

  test.describe('登录按钮状态测试', () => {
    test('点击登录按钮应该显示加载状态', async ({ page }) => {
      // 填写表单
      await page.fill('input[type="text"]', TEST_USER.username);
      await page.fill('input[type="password"]', TEST_USER.password);

      // 监听登录请求
      const loginResponsePromise = page.waitForResponse(
        (response) => response.url().includes('/api/v1/auth/login')
      );

      // 点击登录
      await page.click('button[type="submit"]');

      // 验证按钮显示加载状态
      await expect(page.locator('button[type="submit"]')).toContainText('登录中...');
      await expect(page.locator('.spinner')).toBeVisible();

      // 等待请求完成
      await loginResponsePromise;
    });
  });

  test.describe('错误处理测试', () => {
    test('用户名不存在应该显示错误信息', async ({ page }) => {
      // 填写不存在的用户名
      await page.fill('input[type="text"]', INVALID_USER.username);
      await page.fill('input[type="password"]', INVALID_USER.password);

      // 点击登录
      await page.click('button[type="submit"]');

      // 验证错误提示显示
      await expect(page.locator('.form-alert.error')).toBeVisible();
      await expect(page.locator('.form-alert.error')).toContainText('登录失败');
    });

    test('密码错误应该显示错误信息', async ({ page }) => {
      // 填写错误的密码
      await page.fill('input[type="text"]', TEST_USER.username);
      await page.fill('input[type="password"]', 'wrongpassword');

      // 点击登录
      await page.click('button[type="submit"]');

      // 验证错误提示显示
      await expect(page.locator('.form-alert.error')).toBeVisible();
    });
  });

  test.describe('成功登录测试', () => {
    test('输入正确的用户名和密码应该登录成功', async ({ page }) => {
      // 填写正确的登录信息
      await page.fill('input[type="text"]', TEST_USER.username);
      await page.fill('input[type="password"]', TEST_USER.password);

      // 点击登录
      await page.click('button[type="submit"]');

      // 验证成功提示显示
      await expect(page.locator('.toast.success')).toBeVisible();
      await expect(page.locator('.toast.success')).toContainText('登录成功');

      // 等待跳转（根据角色跳转到不同页面）
      await page.waitForURL(/(\/admin|\/)/, { timeout: 5000 });

      // 验证 localStorage 中保存了认证信息
      const token = await page.evaluate(() => localStorage.getItem('junblog_token'));
      expect(token).toBeTruthy();
    });

    test('登录成功后应该保存用户信息到 localStorage', async ({ page }) => {
      // 填写正确的登录信息
      await page.fill('input[type="text"]', TEST_USER.username);
      await page.fill('input[type="password"]', TEST_USER.password);

      // 点击登录
      await page.click('button[type="submit"]');

      // 等待登录完成
      await page.waitForURL(/(\/admin|\/)/, { timeout: 5000 });

      // 验证 localStorage 中的用户信息
      const user = await page.evaluate(() => {
        const userStr = localStorage.getItem('junblog_user');
        return userStr ? JSON.parse(userStr) : null;
      });

      expect(user).toBeTruthy();
      expect(user.username).toBe(TEST_USER.username);
    });
  });

  test.describe('已登录状态测试', () => {
    test('已登录用户应该自动跳转', async ({ page }) => {
      // 先模拟登录状态
      await page.evaluate(() => {
        localStorage.setItem('junblog_token', 'mock-token');
        localStorage.setItem('junblog_user', JSON.stringify({
          id: 1,
          username: 'testuser',
          role: 'user',
        }));
      });

      // 访问登录页
      await page.goto('/login');
      await waitForPageLoad(page);

      // 应该自动跳转到首页
      await page.waitForURL('/', { timeout: 5000 });
    });
  });

  test.describe('页面交互测试', () => {
    test('点击返回首页链接应该跳转', async ({ page }) => {
      await page.click('a.back-home');
      await page.waitForURL('/');
    });

    test('点击注册链接应该跳转到注册页', async ({ page }) => {
      await page.click('a[href="/register"]');
      await page.waitForURL('/register');
    });
  });

  test.describe('响应式设计测试', () => {
    test('移动端应该正确显示', async ({ page }) => {
      // 设置移动端视口
      await page.setViewportSize({ width: 375, height: 812 });

      // 验证页面元素仍然可见
      await expect(page.locator('h2')).toContainText('登录');
      await expect(page.locator('input[type="text"]')).toBeVisible();
      await expect(page.locator('button[type="submit"]')).toBeVisible();
    });
  });
});
