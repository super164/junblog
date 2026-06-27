# JunBlog 登录功能自动化测试

本目录包含使用 Playwright 框架编写的登录功能自动化测试。

## 测试文件说明

### 1. login.spec.ts - UI 自动化测试
测试登录页面的前端交互，包括：
- 页面加载测试
- 表单验证测试
- 密码显示/隐藏测试
- 登录按钮状态测试
- 错误处理测试
- 成功登录测试
- 已登录状态测试
- 页面交互测试
- 响应式设计测试

### 2. login-api.spec.ts - API 自动化测试
测试后端登录接口，包括：
- 登录接口测试
- 响应格式验证
- Token 验证
- 受保护 API 访问测试
- Refresh Token 测试

## 前置条件

1. **启动后端服务**
   ```bash
   cd blog_backend
   go run main.go
   ```
   后端服务将在 http://localhost:8080 启动

2. **启动前端服务**（可选，Playwright 会自动启动）
   ```bash
   cd bolg_forntend
   npm run dev
   ```
   前端服务将在 http://localhost:5173 启动

3. **准备测试数据**
   确保数据库中存在测试用户：
   - 用户名: testuser
   - 密码: test123456

## 运行测试

### 运行所有测试
```bash
npm test
```

### 运行特定测试文件
```bash
# 只运行 UI 测试
npx playwright test login.spec.ts

# 只运行 API 测试
npx playwright test login-api.spec.ts
```

### 运行特定测试用例
```bash
# 运行包含 "表单验证" 的测试
npx playwright test -g "表单验证"

# 运行包含 "成功登录" 的测试
npx playwright test -g "成功登录"
```

### 调试模式
```bash
# 使用调试模式运行
npm run test:debug

# 使用 UI 模式运行（推荐）
npm run test:ui
```

### 有界面模式
```bash
# 显示浏览器界面运行
npm run test:headed
```

## 查看测试报告

测试完成后，可以查看详细的测试报告：

```bash
npm run test:report
```

报告将自动在浏览器中打开，包含：
- 测试结果概览
- 每个测试的详细信息
- 失败测试的截图和追踪信息
- 测试执行时间

## 测试结果说明

### 通过的测试 ✅
- 测试用例执行成功
- 所有断言都通过

### 失败的测试 ❌
- 测试用例执行失败
- 查看报告中的错误信息和截图

### 跳过的测试 ⏭️
- 测试被跳过（通常是因为条件不满足）

## 常见问题

### 1. 后端服务未启动
**错误信息**: `Request timeout after 30000ms`

**解决方案**: 确保后端服务已启动并在 http://localhost:8080 运行

### 2. 测试用户不存在
**错误信息**: `登录失败: 用户名或密码错误`

**解决方案**: 在数据库中创建测试用户：
```sql
INSERT INTO users (username, password, email, role, status)
VALUES ('testuser', '$2a$10$...', 'test@example.com', 'user', true);
```

### 3. 浏览器未安装
**错误信息**: `Executable doesn't exist`

**解决方案**: 安装 Playwright 浏览器：
```bash
npx playwright install chromium
```

### 4. 端口被占用
**错误信息**: `listen EADDRINUSE: address already in use :::5173`

**解决方案**: 
- 关闭占用端口的程序
- 或修改 playwright.config.ts 中的端口配置

## 自定义测试

### 添加新的测试用例
在对应的测试文件中添加新的 `test()` 块：

```typescript
test('我的新测试', async ({ page }) => {
  // 测试代码
  await page.goto('/login');
  await expect(page.locator('h2')).toContainText('登录');
});
```

### 修改测试数据
在测试文件顶部修改 `TEST_USER` 对象：

```typescript
const TEST_USER = {
  username: 'your_username',
  password: 'your_password',
};
```

### 添加新的测试描述块
使用 `test.describe()` 组织相关测试：

```typescript
test.describe('我的测试组', () => {
  test('测试1', async ({ page }) => {
    // ...
  });

  test('测试2', async ({ page }) => {
    // ...
  });
});
```

## 测试覆盖率

当前测试覆盖以下功能：

- ✅ 页面加载和渲染
- ✅ 表单字段验证
- ✅ 用户输入交互
- ✅ 密码显示/隐藏
- ✅ 登录按钮状态
- ✅ 错误信息显示
- ✅ 成功登录跳转
- ✅ 本地存储管理
- ✅ 响应式设计
- ✅ API 请求验证
- ✅ Token 生成和验证
- ✅ 受保护资源访问

## 相关链接

- [Playwright 官方文档](https://playwright.dev/)
- [Playwright API 文档](https://playwright.dev/docs/api/class-page)
- [Vue Testing Handbook](https://vue-testing-handbook.com/)
