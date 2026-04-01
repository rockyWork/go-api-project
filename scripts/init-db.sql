-- 数据库初始化脚本
-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS go_api_project CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE go_api_project;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(32) NOT NULL,
    email VARCHAR(128),
    phone VARCHAR(16),
    password_hash VARCHAR(255) NOT NULL,
    avatar VARCHAR(255),
    status TINYINT DEFAULT 1 COMMENT '1:正常, 2:禁用',
    role TINYINT DEFAULT 1 COMMENT '1:用户, 2:管理员',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    UNIQUE KEY uk_username (username),
    UNIQUE KEY uk_email (email),
    INDEX idx_status (status),
    INDEX idx_role (role),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- Refresh Token 表（可选，用于Token黑名单或持久化存储）
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_token_hash (token_hash),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Refresh Token 表';

-- 插入默认管理员账号（密码: Admin123456）
-- 注意：生产环境应删除或修改默认账号
INSERT INTO users (username, email, password_hash, status, role, created_at, updated_at) 
VALUES (
    'admin',
    'admin@example.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMy.Mqrqhc6WOUlxRaP3USvYFNb/AN9QAQO',
    1,
    2,
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE updated_at = NOW();

-- 注意：上面的密码哈希是 bcrypt 加密后的 "Admin123456"
-- 可以使用以下命令生成新密码：
-- go run -e 'package main; import ("fmt"; "golang.org/x/crypto/bcrypt"); func main() { b, _ := bcrypt.GenerateFromPassword([]byte("your-password"), bcrypt.DefaultCost); fmt.Println(string(b)) }'
