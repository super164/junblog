-- JunBlog 数据库初始化脚本
CREATE DATABASE IF NOT EXISTS junblog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE junblog;

-- 创建用户表（基础结构，其他表由 GORM AutoMigrate 创建）
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码',
    email VARCHAR(100) UNIQUE COMMENT '邮箱',
    phone VARCHAR(20) COMMENT '手机号',
    role VARCHAR(20) DEFAULT 'user' COMMENT '角色',
    avatar VARCHAR(255) COMMENT '头像',
    status BOOLEAN DEFAULT TRUE COMMENT '状态',
    github_id VARCHAR(50) UNIQUE COMMENT 'GitHub用户ID',
    github_login VARCHAR(100) COMMENT 'GitHub用户名',
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_github_id (github_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建默认管理员账号（密码: admin123，生产环境请修改）
INSERT INTO users (created_at, updated_at, username, password, email, role, status)
VALUES (NOW(), NOW(), 'admin', '$2a$10$N.ZOn9G6w3Fz4nFHRXn5GOe9Th2jKZqK7TAKpXv4pG1wFkBmvUYCi', 'admin@junblog.com', 'admin', true)
ON DUPLICATE KEY UPDATE updated_at = NOW();
