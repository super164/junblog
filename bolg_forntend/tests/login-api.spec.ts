import { test, expect } from '@playwright/test';

// API 基础 URL
const API_BASE_URL = 'http://localhost:8080/api/v1';

// 测试数据
const TEST_USER = {
  username: 'testuser',
  password: 'test123456',
};

test.describe('登录 API 测试', () => {
  test.describe('POST /api/v1/auth/login', () => {
    test('正确的用户名和密码应该返回200和token', async ({ request }) => {
      const response = await request.post(`${API_BASE_URL}/auth/login`, {
        data: TEST_USER,
      });

      expect(response.status()).toBe(200);

      const body = await response.json();
      expect(body.code).toBe(200);
      expect(body.data.token).toBeTruthy();
      expect(body.data.refresh_token).toBeTruthy();
      expect(body.data.user).toBeTruthy();
      expect(body.data.user.username).toBe(TEST_USER.username);
    });

    test('错误的密码应该返回401', async ({ request }) => {
      const response = await request.post(`${API_BASE_URL}/auth/login`, {
        data: {
          username: TEST_USER.username,
          password: 'wrongpassword',
        },
      });

      expect(response.status()).toBe(401);

      const body = await response.json();
      expect(body.code).toBe(401);
      expect(body.msg).toContain('登录失败');
    });

    test('不存在的用户名应该返回401', async ({ request }) => {
      const response = await request.post(`${API_BASE_URL}/auth/login`, {
        data: {
          username: 'nonexistentuser',
          password: 'test123456',
        },
      });

      expect(response.status()).toBe(401);

      const body = await response.json();
      expect(body.code).toBe(401);
    });

    test('空用户名应该返回400', async ({ request }) => {
      const response = await request.post(`${API_BASE_URL}/auth/login`, {
        data: {
          password: TEST_USER.password,
        },
      });

      expect(response.status()).toBe(400);

      const body = await response.json();
      expect(body.code).toBe(400);
    });

    test('空密码应该返回400', async ({ request }) => {
      const response = await request.post(`${API_BASE_URL}/auth/login`, {
        data: {
          username: TEST_USER.username,
        },
      });

      expect(response.status()).toBe(400);

      const body = await response.json();
      expect(body.code).toBe(400);
    });

    test('空请求体应该返回400', async ({ request }) => {
      const response = await request.post(`${API_BASE_URL}/auth/login`, {
        data: {},
      });

      expect(response.status()).toBe(400);
    });

    test('无效的JSON格式应该返回400', async ({ request }) => {
      const response = await request.post(`${API_BASE_URL}/auth/login`, {
        data: 'invalid json',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      expect(response.status()).toBe(400);
    });
  });

  test.describe('登录响应格式验证', () => {
    test('成功登录应该返回正确的响应结构', async ({ request }) => {
      const response = await request.post(`${API_BASE_URL}/auth/login`, {
        data: TEST_USER,
      });

      const body = await response.json();

      // 验证响应结构
      expect(body).toHaveProperty('code');
      expect(body).toHaveProperty('data');
      expect(body.data).toHaveProperty('token');
      expect(body.data).toHaveProperty('refresh_token');
      expect(body.data).toHaveProperty('user');

      // 验证用户信息结构
      const user = body.data.user;
      expect(user).toHaveProperty('id');
      expect(user).toHaveProperty('username');
      expect(user).toHaveProperty('email');
      expect(user).toHaveProperty('role');
      expect(user).toHaveProperty('createdAt');
    });

    test('token 应该是有效的 JWT 格式', async ({ request }) => {
      const response = await request.post(`${API_BASE_URL}/auth/login`, {
        data: TEST_USER,
      });

      const body = await response.json();
      const token = body.data.token;

      // JWT 格式：header.payload.signature
      const parts = token.split('.');
      expect(parts.length).toBe(3);

      // 验证每个部分都是 base64 编码
      parts.forEach((part) => {
        expect(() => atob(part)).not.toThrow();
      });
    });
  });

  test.describe('使用 token 访问受保护的 API', () => {
    test('登录后应该能访问用户资料接口', async ({ request }) => {
      // 先登录获取 token
      const loginResponse = await request.post(`${API_BASE_URL}/auth/login`, {
        data: TEST_USER,
      });

      const loginBody = await loginResponse.json();
      const token = loginBody.data.token;

      // 使用 token 访问用户资料
      const profileResponse = await request.get(`${API_BASE_URL}/user/profile`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      expect(profileResponse.status()).toBe(200);

      const profileBody = await profileResponse.json();
      expect(profileBody.code).toBe(200);
      expect(profileBody.data.username).toBe(TEST_USER.username);
    });

    test('没有 token 应该无法访问受保护的 API', async ({ request }) => {
      const response = await request.get(`${API_BASE_URL}/user/profile`);

      expect(response.status()).toBe(401);
    });

    test('无效的 token 应该返回401', async ({ request }) => {
      const response = await request.get(`${API_BASE_URL}/user/profile`, {
        headers: {
          Authorization: 'Bearer invalid-token',
        },
      });

      expect(response.status()).toBe(401);
    });
  });

  test.describe('Refresh Token 测试', () => {
    test('使用 refresh token 应该能获取新的 access token', async ({ request }) => {
      // 先登录获取 refresh token
      const loginResponse = await request.post(`${API_BASE_URL}/auth/login`, {
        data: TEST_USER,
      });

      const loginBody = await loginResponse.json();
      const refreshToken = loginBody.data.refresh_token;

      // 使用 refresh token 获取新的 access token
      const refreshResponse = await request.post(`${API_BASE_URL}/auth/refresh`, {
        data: {
          refresh_token: refreshToken,
        },
      });

      expect(refreshResponse.status()).toBe(200);

      const refreshBody = await refreshResponse.json();
      expect(refreshBody.code).toBe(200);
      expect(refreshBody.data.token).toBeTruthy();
      expect(refreshBody.data.refresh_token).toBeTruthy();
    });

    test('无效的 refresh token 应该返回401', async ({ request }) => {
      const response = await request.post(`${API_BASE_URL}/auth/refresh`, {
        data: {
          refresh_token: 'invalid-refresh-token',
        },
      });

      expect(response.status()).toBe(401);
    });
  });
});
