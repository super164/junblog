# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: interaction.spec.ts >> 已登录状态测试 >> 提交评论后应显示在列表中
- Location: tests\interaction.spec.ts:94:3

# Error details

```
Error: expect(received).toBe(expected) // Object.is equality

Expected: 3
Received: 0
```

# Page snapshot

```yaml
- generic [ref=e3]:
  - banner [ref=e4]:
    - generic [ref=e5]:
      - link "J JunBlog" [ref=e6] [cursor=pointer]:
        - /url: /
        - generic [ref=e7]: J
        - generic [ref=e8]: JunBlog
      - navigation [ref=e9]:
        - link "首页" [ref=e10] [cursor=pointer]:
          - /url: /
        - link "文章" [ref=e11] [cursor=pointer]:
          - /url: /articles
        - link "关于" [ref=e12] [cursor=pointer]:
          - /url: /about
      - link "登录" [ref=e14] [cursor=pointer]:
        - /url: /login
  - main [ref=e15]:
    - generic [ref=e16]:
      - text: 文
      - generic [ref=e17]:
        - generic [ref=e18]:
          - generic [ref=e19]: 个人博客项目
          - heading "用简洁的设计， 承载持续表达的内容" [level=1] [ref=e20]:
            - text: 用简洁的设计，
            - text: 承载持续表达的内容
          - paragraph [ref=e21]: 一个前后台分离的博客系统，前台供读者浏览，后台管理文章、分类、标签和评论。
          - generic [ref=e22]:
            - link "浏览文章" [ref=e23] [cursor=pointer]:
              - /url: /articles
            - link "了解作者" [ref=e24] [cursor=pointer]:
              - /url: /about
        - link "— Featured 技术笔记 · 3天前 使用 Go + Vue 搭建前后台分离博客系统 这是一篇关于使用 Go 和 Vue 3 搭建前后台分离博客系统的文章。" [ref=e26] [cursor=pointer]:
          - /url: /articles/go-vue-blog
          - generic [ref=e27]:
            - paragraph [ref=e28]: — Featured
            - generic [ref=e29]:
              - generic [ref=e30]: 技术笔记
              - generic [ref=e31]: · 3天前
            - heading "使用 Go + Vue 搭建前后台分离博客系统" [level=3] [ref=e32]
            - paragraph [ref=e33]: 这是一篇关于使用 Go 和 Vue 3 搭建前后台分离博客系统的文章。
    - generic [ref=e34]:
      - generic [ref=e35]:
        - generic [ref=e36]:
          - paragraph [ref=e37]: — Recent Posts
          - heading "最新文章" [level=2] [ref=e38]
        - link "查看全部 →" [ref=e39] [cursor=pointer]:
          - /url: /articles
      - generic [ref=e40]:
        - 'link "测试文章 No.01 技术笔记 · 3小时前 测试文章 11111 # Go # Vue" [ref=e41] [cursor=pointer]':
          - /url: /articles/测试文章-1780882807
          - img "测试文章" [ref=e43]
          - generic [ref=e44]:
            - generic [ref=e45]: No.01
            - generic [ref=e46]:
              - generic [ref=e47]: 技术笔记
              - generic [ref=e48]: · 3小时前
            - heading "测试文章" [level=3] [ref=e49]
            - paragraph [ref=e50]: "11111"
            - generic [ref=e51]:
              - generic [ref=e52]: "# Go"
              - generic [ref=e53]: "# Vue"
        - link "No.02 生活随笔 · 3天前 写博客的初衷 开始写博客，是为了记录学习过程中的思考和总结。" [ref=e54] [cursor=pointer]:
          - /url: /articles/why-blog
          - generic [ref=e55]:
            - text: No.02
            - generic [ref=e56]:
              - generic [ref=e57]: 生活随笔
              - generic [ref=e58]: · 3天前
            - heading "写博客的初衷" [level=3] [ref=e59]
            - paragraph [ref=e60]: 开始写博客，是为了记录学习过程中的思考和总结。
        - 'link "No.03 技术笔记 · 3天前 MySQL 索引优化实践 数据库性能优化是后端开发中的重要环节。 # MySQL" [ref=e61] [cursor=pointer]':
          - /url: /articles/mysql-index
          - generic [ref=e62]:
            - text: No.03
            - generic [ref=e63]:
              - generic [ref=e64]: 技术笔记
              - generic [ref=e65]: · 3天前
            - heading "MySQL 索引优化实践" [level=3] [ref=e66]
            - paragraph [ref=e67]: 数据库性能优化是后端开发中的重要环节。
            - generic [ref=e69]: "# MySQL"
        - 'link "No.04 技术笔记 · 3天前 使用 Go + Vue 搭建前后台分离博客系统 这是一篇关于使用 Go 和 Vue 3 搭建前后台分离博客系统的文章。 # Go # Vue" [ref=e70] [cursor=pointer]':
          - /url: /articles/go-vue-blog
          - generic [ref=e71]:
            - text: No.04
            - generic [ref=e72]:
              - generic [ref=e73]: 技术笔记
              - generic [ref=e74]: · 3天前
            - heading "使用 Go + Vue 搭建前后台分离博客系统" [level=3] [ref=e75]
            - paragraph [ref=e76]: 这是一篇关于使用 Go 和 Vue 3 搭建前后台分离博客系统的文章。
            - generic [ref=e77]:
              - generic [ref=e78]: "# Go"
              - generic [ref=e79]: "# Vue"
    - generic [ref=e80]:
      - text: "\""
      - generic [ref=e81]:
        - paragraph [ref=e82]: — About The Author
        - heading "Jun" [level=2] [ref=e83]
        - paragraph [ref=e84]: 全栈开发者 / 写作者
      - paragraph [ref=e85]: 热爱技术与设计，用代码构建有温度的产品。
  - contentinfo [ref=e86]:
    - generic [ref=e87]:
      - generic [ref=e88]:
        - generic [ref=e89]: J
        - generic [ref=e90]: JunBlog
      - generic [ref=e91]: © 2026 Jun. 用简洁的设计，承载持续表达的内容。
      - generic [ref=e92]:
        - link "文章" [ref=e93] [cursor=pointer]:
          - /url: /articles
        - link "关于" [ref=e94] [cursor=pointer]:
          - /url: /about
        - link "登录" [ref=e95] [cursor=pointer]:
          - /url: /login
```

# Test source

```ts
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
  27  |     await likeBtn.click()
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
> 106 |     expect(newComments).toBe(initialComments + 1)
      |                         ^ Error: expect(received).toBe(expected) // Object.is equality
  107 |   })
  108 | })
  109 | 
```