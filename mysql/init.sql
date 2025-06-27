-- 创建数据库
CREATE DATABASE IF NOT EXISTS xx_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE xx_admin;

-- 创建角色表
CREATE TABLE IF NOT EXISTS roles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE,
    nickname VARCHAR(50),
    avatar VARCHAR(255),
    status INT DEFAULT 1,
    role_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

-- 创建菜单表
CREATE TABLE IF NOT EXISTS menus (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    path VARCHAR(100),
    component VARCHAR(100),
    icon VARCHAR(50),
    sort INT DEFAULT 0,
    parent_id INT NULL,
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (parent_id) REFERENCES menus(id)
);

-- 插入初始数据
INSERT INTO roles (name, description) VALUES 
('admin', '系统管理员'),
('user', '普通用户')
ON DUPLICATE KEY UPDATE description = VALUES(description);

-- 插入管理员用户 (密码: admin123)
INSERT INTO users (username, password, email, nickname, role_id) VALUES 
('admin', 'e99a18c428cb38d5f260853678922e03', 'admin@example.com', '系统管理员', 1)
ON DUPLICATE KEY UPDATE email = VALUES(email), nickname = VALUES(nickname);

-- 插入菜单数据
INSERT INTO menus (name, path, component, icon, sort, parent_id) VALUES 
('首页', '/', 'Home', 'House', 1, NULL),
('用户管理', '/user', 'User', 'User', 2, NULL),
('角色管理', '/role', 'Role', 'Setting', 3, NULL),
('菜单管理', '/menu', 'Menu', 'Menu', 4, NULL),
('数据表格', '/table', 'Table', 'List', 5, NULL)
ON DUPLICATE KEY UPDATE 
    path = VALUES(path), 
    component = VALUES(component), 
    icon = VALUES(icon), 
    sort = VALUES(sort); 