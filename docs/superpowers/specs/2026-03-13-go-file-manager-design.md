# Go Web 文件管理器设计文档

**日期**: 2026-03-13

## 1. 项目概述

一个基于 Go 的轻量级 Web 文件管理器，通过浏览器访问，支持基本文件操作（上传/下载/查看/编辑/重命名/复制/移动/删除）。部署时支持任意 base path（如 `/files` 或 `/a/b/c/manager`）。

## 2. 功能需求

### 2.1 认证
- 仅密码验证，无需用户名
- 密码通过命令行参数 `-password` 传入
- 登录成功后使用 Session/Cookie 保持登录状态

### 2.2 文件操作
| 操作 | 描述 |
|------|------|
| 浏览 | 列出目录下的文件和文件夹，支持导航 |
| 上传 | 支持单文件/多文件上传，拖拽上传 |
| 下载 | 文件下载，文件夹打包下载（zip） |
| 查看 | 文本文件内容预览，图片预览 |
| 编辑 | 纯文本编辑器（textarea），保存修改 |
| 重命名 | 文件/文件夹重命名 |
| 复制 | 复制文件/文件夹到目标目录 |
| 移动 | 移动文件/文件夹到目标目录 |
| 删除 | 删除文件/文件夹（支持批量） |
| 新建 | 创建新文件夹 |

### 2.3 配置
- 根目录：通过命令行 `-root` 指定（必选）
- 端口：通过命令行 `-port` 指定（默认 8080）
- Base Path：通过命令行 `-basepath` 指定（默认 `/`）
- 密码：通过命令行 `-password` 指定（必选）

## 3. 技术架构

### 3.1 技术栈
- **后端**: Go 1.21+, 标准库 `net/http`
- **前端**: 原生 HTML/CSS/JS（无框架依赖）
- **静态资源**: Go 1.16+ `embed` 嵌入

### 3.2 架构模式
前后端分离架构，Go 服务嵌入静态前端文件：
```
/                   -> index.html（入口）
{basePath}/api/*    -> REST API
{basePath}/static/* -> 静态资源（js/css/images）
```

### 3.3 核心模块

```
main.go             # 入口，命令行参数解析
server.go           # HTTP 服务器启动
router.go           # 路由定义
handlers/
  auth.go           # 登录/登出处理
  files.go          # 文件操作 API
static/             # 前端静态资源（embed）
```

## 4. API 设计

### 4.1 认证 API
| 方法 | 路径 | 描述 |
|------|------|------|
| POST | {basePath}/api/login | 登录 |
| POST | {basePath}/api/logout | 登出 |
| GET  | {basePath}/api/check | 检查登录状态 |

### 4.2 文件 API
| 方法 | 路径 | 描述 |
|------|------|------|
| GET  | {basePath}/api/list | 列出目录内容 |
| GET  | {basePath}/api/read | 读取文件内容 |
| POST | {basePath}/api/upload | 上传文件 |
| GET  | {basePath}/api/download | 下载文件 |
| POST | {basePath}/api/create | 创建文件夹 |
| POST | {basePath}/api/rename | 重命名 |
| POST | {basePath}/api/copy | 复制 |
| POST | {basePath}/api/move | 移动 |
| POST | {basePath}/api/delete | 删除 |

### 4.3 响应格式
```json
{
  "code": 0,
  "msg": "success",
  "data": {...}
}
```

## 5. 前端设计

### 5.1 页面结构
- **登录页**: 密码输入框 + 登录按钮
- **文件管理器**: 
  - 顶部：当前路径面包屑
  - 侧边栏/工具栏：操作按钮
  - 主区域：文件列表（表格/网格）
  - 弹窗：编辑、确认对话框

### 5.2 关键特性
- 所有资源路径使用动态 base path：`${basePath}/static/xxx.js`
- 前端通过 `window.location.pathname` 解析 base path（去除路由部分）
- 响应式布局
- 拖拽上传
- 文件类型图标
- 操作确认对话框

## 6. Base Path 支持

**核心要求**：服务可部署在任意 nginx 二级/三级目录

### 6.1 实现方式
1. 启动时接收 `-basepath` 参数（如 `/manager` 或 `/a/b/c/files`）
2. 所有 API 和静态资源路径都拼接 base path
3. 前端通过 JS 获取当前 base path

### 6.2 路径示例
假设 `-basepath=/files`：
```
/files/              -> 首页
/files/api/login     -> 登录 API
/files/static/app.js -> 前端 JS
/files/api/list?path=/docs
```

## 7. 安全性考虑

### 7.1 路径遍历防护
禁止访问根目录之外的文件，实现逻辑：
```go
func safePath(rootDir, userPath string) (string, error) {
    absRoot, _ := filepath.Abs(rootDir)
    absPath, _ := filepath.Abs(filepath.Join(rootDir, userPath))
    if !strings.HasPrefix(absPath, absRoot) {
        return "", errors.New("access denied")
    }
    return absPath, nil
}
```

### 7.2 Session 认证
- 使用内存存储 Session（map[token]User）
- 登录成功后设置 HttpOnly Cookie（名称：`session_id`）
- Session 过期时间：24小时
- Token：UUID 随机字符串

### 7.3 文件限制
- 文件大小限制：100MB（可通过参数配置）
- 禁止上传可执行文件（.exe, .sh, .bat 等）

### 7.4 CSRF 防护
- 使用 SameSite Cookie
- 登录成功后生成 CSRF Token（UUID）存入 Cookie
- 前端从 Cookie 读取 Token，添加到请求头 `X-CSRF-Token`
- 后端验证 Token 与 Cookie 是否匹配

## 8. 部署

### 8.1 构建
```bash
go build -o file-manager
```

### 8.2 启动
```bash
./file-manager -root /var/www -password secret -port 8080 -basepath /files
```

### 8.3 Nginx 配置示例
```nginx
location /files/ {
    proxy_pass http://127.0.0.1:8080;
}
```

## 9. 验收标准

- [ ] 密码登录功能正常
- [ ] 可浏览指定根目录下的文件
- [ ] 可上传/下载文件
- [ ] 可查看文本文件内容
- [ ] 可编辑并保存文本文件
- [ ] 可重命名/复制/移动/删除
- [ ] 部署在二级/三级目录时所有功能正常
- [ ] 路径遍历防护有效
