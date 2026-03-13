# File Manager (Go Web 文件管理器)

[English](#english) | [中文](#中文)

---

## English

### About

A lightweight, self-hosted web file manager built with Go. Simply run the binary, access via browser, and manage your server files with ease.

### Features

- 🔐 Password-only authentication (no username required)
- 📁 Full file management: browse, upload, download, view, edit
- 🔧 File operations: rename, copy, move, delete
- 🌐 Supports deployment under any URL path (e.g., `/files`, `/a/b/c/manager`)
- 🛡️ Path traversal protection for security
- 💻 Minimal dependencies - single binary deployment

### Quick Start

```bash
# Download or build
go build -o file-manager

# Run
./file-manager -root /path/to/serve -password YOUR_PASSWORD -port 8080

# Or with custom base path (for nginx reverse proxy)
./file-manager -root /var/www -password secret -port 8080 -basepath /files
```

### Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-root` | Root directory to serve (required) | - |
| `-password` | Login password (required) | - |
| `-port` | HTTP listen port | 8080 |
| `-basepath` | URL base path for routing | / |

### Nginx Configuration

```nginx
location /files/ {
    proxy_pass http://127.0.0.1:8080;
}
```

### Build

```bash
git clone https://github.com/keeply-cn/file-manager.git
cd file-manager
go build -o file-manager
```

---

## 中文

### 项目由来

这是一个轻量级的 Go 语言 Web 文件管理器。在开发过程中，我发现需要一个简单的方式来管理服务器上的文件，而现有的方案要么太复杂，要么需要复杂的配置。因此我决定自己动手，用 Go 标准库编写一个简洁的文件管理器。

### 功能特性

- 🔐 纯密码认证（无需用户名）
- 📁 完整文件管理：浏览、上传、下载、查看、编辑
- 🔧 文件操作：重命名、复制、移动、删除
- 🌐 支持任意 URL 路径部署（如 `/files`、`/a/b/c/manager`）
- 🛡️ 路径遍历防护，确保安全
- 💻 极简依赖 - 单二进制部署

### 快速开始

```bash
# 下载或编译
go build -o file-manager

# 运行
./file-manager -root /path/to/serve -password 你的密码 -port 8080

# 或使用自定义路径（用于 nginx 反向代理）
./file-manager -root /var/www -password secret -port 8080 -basepath /files
```

### 命令行参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-root` | 要托管的根目录（必填） | - |
| `-password` | 登录密码（必填） | - |
| `-port` | HTTP 监听端口 | 8080 |
| `-basepath` | URL 基础路径 | / |

### Nginx 配置示例

```nginx
location /files/ {
    proxy_pass http://127.0.0.1:8080;
}
```

### 编译安装

```bash
git clone https://github.com/keeply-cn/file-manager.git
cd file-manager
go build -o file-manager
```

---

### 个人网站

🌐 https://www.keeply.cn

---

### License

MIT License
