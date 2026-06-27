# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: interaction.spec.ts >> 未登录状态测试 >> 未登录点击点赞应提示登录
- Location: tests\interaction.spec.ts:20:3

# Error details

```
Test timeout of 30000ms exceeded.
```

```
Error: locator.click: Test timeout of 30000ms exceeded.
Call log:
  - waiting for locator('.interact-btn').first()
    - locator resolved to <button data-v-b3b8ad5c="" class="interact-btn">🤍 1</button>
  - attempting click action
    - waiting for element to be visible, enabled and stable
    - element is visible, enabled and stable
    - scrolling into view if needed
    - done scrolling
    - performing click action

```

# Test source

```ts
  1   | import { test, expect } from '@playwright/test';
  2   | 
  3   | // 测试账号的 token（通过 API 获取）
  4   | let authToken = '';
  5   | let refreshToken = '';
  6   | let testUser = { id: 3, username: 'admin', role: 'admin' };
  7   | 
  8   | test.beforeAll(async ({ request }) => {
  9   |   // 获取测试账号的 token
  10  |   const loginRes = await request.post('http://localhost:8080/api/v1/auth/login', {
  11  |     data: { username: 'admin', password: 'admin123' },
  12  |   });
  13  |   const loginBody = await loginRes.json();
  14  |   authToken = loginBody.data.token;
  15  |   refreshToken = loginBody.data.refresh_token;
  16  |   testUser = loginBody.data.user;
  17  | });
  18  | 
  19  | test.describe('未登录状态测试', () => {
  20  |   test('未登录点击点赞应提示登录', async ({ page }) => {
  21  |     const dialogPromise = page.waitForEvent('dialog')
  22  | 
  23  |     await page.goto('/articles/go-vue-blog')
  24  |     await page.waitForSelector('.interaction-bar')
  25  | 
  26  |     const likeBtn = page.locator('.interact-btn').first()
> 27  |     await likeBtn.click()
      |                   ^ Error: locator.click: Test timeout of 30000ms exceeded.
  28  | 
  29  |     const dialog = await dialogPromise
  30  |     expect(dialog.message()).toContain('请先登录')
  31  |     await dialog.accept()
  32  |   })
  33  | 
  34  |   test('未登录应显示评论登录提示', async ({ page }) => {
  35  |     await page.goto('/articles/go-vue-blog')
  36  |     await page.waitForSelector('.comment-section')
  37  | 
  38  |     await expect(page.locator('.login-hint')).toBeVisible()
  39  |     await expect(page.locator('.login-hint')).toContainText('登录')
  40  |   })
  41  | })
  42  | 
  43  | test.describe('已登录状态测试', () => {
  44  |   test.beforeEach(async ({ page }) => {
  45  |     // 直接注入 token 到 localStorage，模拟已登录状态
  46  |     await page.goto('/')
  47  |     await page.evaluate(({ token, refresh, user }) => {
  48  |       localStorage.setItem('junblog_token', token)
  49  |       localStorage.setItem('junblog_refresh_token', refresh)
  50  |       localStorage.setItem('junblog_user', JSON.stringify(user))
  51  |     }, { token: authToken, refresh: refreshToken, user: testUser })
  52  |   })
  53  | 
  54  |   test('登录后应显示评论输入框', async ({ page }) => {
  55  |     await page.goto('/articles/go-vue-blog')
  56  |     await page.waitForSelector('.comment-section')
  57  | 
  58  |     await expect(page.locator('.comment-form')).toBeVisible()
  59  |     await expect(page.locator('textarea')).toBeVisible()
  60  |   })
  61  | 
  62  |   test('点赞后计数应更新', async ({ page }) => {
  63  |     await page.goto('/articles/go-vue-blog')
  64  |     await page.waitForSelector('.interaction-bar')
  65  | 
  66  |     const likeBtn = page.locator('.interact-btn').first()
  67  |     const initialText = await likeBtn.textContent()
  68  |     const initialCount = parseInt(initialText?.match(/\d+/)?.[0] || '0')
  69  | 
  70  |     await likeBtn.click()
  71  |     await page.waitForTimeout(1000)
  72  | 
  73  |     const newText = await likeBtn.textContent()
  74  |     const newCount = parseInt(newText?.match(/\d+/)?.[0] || '0')
  75  |     expect(newCount).toBe(initialCount + 1)
  76  |   })
  77  | 
  78  |   test('收藏后计数应更新', async ({ page }) => {
  79  |     await page.goto('/articles/go-vue-blog')
  80  |     await page.waitForSelector('.interaction-bar')
  81  | 
  82  |     const favBtn = page.locator('.interact-btn').last()
  83  |     const initialText = await favBtn.textContent()
  84  |     const initialCount = parseInt(initialText?.match(/\d+/)?.[0] || '0')
  85  | 
  86  |     await favBtn.click()
  87  |     await page.waitForTimeout(1000)
  88  | 
  89  |     const newText = await favBtn.textContent()
  90  |     const newCount = parseInt(newText?.match(/\d+/)?.[0] || '0')
  91  |     expect(newCount).toBe(initialCount + 1)
  92  |   })
  93  | 
  94  |   test('提交评论后应显示在列表中', async ({ page }) => {
  95  |     await page.goto('/articles/go-vue-blog')
  96  |     await page.waitForSelector('.comment-section')
  97  | 
  98  |     const initialComments = await page.locator('.comment-item').count()
  99  | 
  100 |     const testComment = 'Playwright 测试 ' + Date.now()
  101 |     await page.fill('textarea', testComment)
  102 |     await page.click('button', { hasText: '发表评论' })
  103 |     await page.waitForTimeout(2000)
  104 | 
  105 |     const newComments = await page.locator('.comment-item').count()
  106 |     expect(newComments).toBe(initialComments + 1)
  107 |   })
  108 | })
  109 | 
```