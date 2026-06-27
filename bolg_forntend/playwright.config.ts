import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  // 测试文件目录
  testDir: './tests',

  // 测试结果输出目录
  outputDir: './test-results',

  // 每个测试的超时时间（30秒）
  timeout: 30000,

  // expect断言的超时时间
  expect: {
    timeout: 5000,
  },

  // 完全并行运行测试
  fullyParallel: true,

  // CI环境中禁止.only
  forbidOnly: !!process.env.CI,

  // 失败重试次数
  retries: process.env.CI ? 2 : 0,

  // 并发 worker 数量
  workers: process.env.CI ? 1 : undefined,

  // 报告器配置
  reporter: [
    ['html', { outputFolder: 'playwright-report', open: 'never' }],
    ['list'],
  ],

  // 全局测试配置
  use: {
    // 基础URL（前端开发服务器）
    baseURL: 'http://localhost:5173',

    // 收集失败时的追踪信息
    trace: 'on-first-retry',

    // 截图配置
    screenshot: 'only-on-failure',

    // 视频配置
    video: 'retain-on-failure',

    // 浏览器视口大小
    viewport: { width: 1280, height: 720 },

    // 忽略HTTPS错误（开发环境）
    ignoreHTTPSErrors: true,
  },

  // 浏览器项目配置
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],

  // 本地开发服务器配置（自动启动前端服务）
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
    timeout: 120000,
  },
});
