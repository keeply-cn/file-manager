# Go Web 文件管理器实现计划

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 开发一个可通过浏览器访问的Go文件管理器，支持登录认证和基本文件操作

**Architecture:** 前后端分离架构，Go服务嵌入静态前端文件，所有路径支持动态base path

**Tech Stack:** Go 1.21+ 标准库 net/http, html/template, embed

---

## 文件结构

```
.
├── main.go                 # 入口，命令行参数解析，启动服务器
├── server.go                # HTTP 服务器配置
├── router.go                # 路由定义（静态资源/API）
├── handlers/
│   ├── auth.go              # 登录/登出/状态检查
│   └── files.go             # 文件操作API
├── static/
│   ├── index.html           # 前端入口页面
│   ├── css/
│   │   └── style.css        # 样式
│   └── js/
│       └── app.js           # 前端逻辑
└── utils/
    └── path.go              # 路径处理工具函数
```

---

## Chunk 1: 项目骨架与基础配置

### Task 1: 初始化 Go 模块

- [ ] **Step 1: 初始化 Go 模块**

Run: `go mod init file-manager`

- [ ] **Step 2: 创建目录结构**

```bash
mkdir -p handlers static/css static/js utils
```

- [ ] **Step 3: Commit**

```bash
git init && git add . && git commit -m "chore: init project"
```

---

### Task 2: main.go - 命令行参数解析

**Files:**
- Create: `main.go`

- [ ] **Step 1: 编写 main.go**

```go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Config struct {
	root     string
	password string
	port     int
	basePath string
}

func main() {
	cfg := &Config{}
	flag.StringVar(&cfg.root, "root", "", "Root directory to serve (required)")
	flag.StringVar(&cfg.password, "password", "", "Password for authentication (required)")
	flag.IntVar(&cfg.port, "port", 8080, "Port to listen on")
	flag.StringVar(&cfg.basePath, "basepath", "/", "Base path for routing")
	flag.Parse()

	if cfg.root == "" || cfg.password == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := os.MkdirAll(cfg.root, 0755); err != nil {
		log.Fatalf("Failed to create root directory: %v", err)
	}

	log.Printf("Starting server on :%d with base path %s", cfg.port, cfg.basePath)
	log.Printf("Serving files from: %s", cfg.root)

	if err := run(cfg); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
```

- [ ] **Step 2: 创建 run 函数占位**

```go
// main.go 添加
func run(cfg *Config) error {
	return nil
}
```

- [ ] **Step 3: Commit**

```bash
git add main.go && git commit -m "feat: add main.go with CLI args"
```

---

## Chunk 2: 核心服务与路由

### Task 3: server.go - HTTP 服务器

**Files:**
- Create: `server.go`

- [ ] **Step 1: 编写 server.go**

```go
package main

import (
	"net/http"
	"time"
)

type Server struct {
	cfg    *Config
	mux    *http.ServeMux
	server *http.Server
}

func newServer(cfg *Config) *Server {
	s := &Server{
		cfg: cfg,
		mux: http.NewServeMux(),
	}

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      s.mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return s
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}
```

- [ ] **Step 2: 更新 main.go 引用 run**

```go
func run(cfg *Config) error {
	srv := newServer(cfg)
	return srv.Start()
}
```

- [ ] **Step 3: Commit**

```bash
git add server.go main.go && git commit -m "feat: add HTTP server structure"
```

---

### Task 4: router.go - 路由定义

**Files:**
- Create: `router.go`

- [ ] **Step 1: 编写 router.go**

```go
package main

import (
	"net/http"
)

func (s *Server) setupRoutes() {
	// 静态资源 - 需要处理 basePath
	s.mux.HandleFunc(s.cfg.basePath, s.handleIndex)
	s.mux.HandleFunc(s.cfg.basePath+"/", s.handleIndex)
	s.mux.HandleFunc(s.cfg.basePath+"/static/", s.handleStatic)
	
	// API 路由
	apiPath := s.cfg.basePath + "api/"
	s.mux.HandleFunc(apiPath+"login", s.handleLogin)
	s.mux.HandleFunc(apiPath+"logout", s.handleLogout)
	s.mux.HandleFunc(apiPath+"check", s.handleCheck)
	s.mux.HandleFunc(apiPath+"list", s.handleList)
	s.mux.HandleFunc(apiPath+"read", s.handleRead)
	s.mux.HandleFunc(apiPath+"upload", s.handleUpload)
	s.mux.HandleFunc(apiPath+"download", s.handleDownload)
	s.mux.HandleFunc(apiPath+"create", s.handleCreate)
	s.mux.HandleFunc(apiPath+"rename", s.handleRename)
	s.mux.HandleFunc(apiPath+"copy", s.handleCopy)
	s.mux.HandleFunc(apiPath+"move", s.handleMove)
	s.mux.HandleFunc(apiPath+"delete", s.handleDelete)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != s.cfg.basePath {
		http.NotFound(w, r)
		return
	}
	// 转发到静态文件处理
	s.handleStatic(w, r)
}
```

- [ ] **Step 2: 添加 handleStatic 存根**

```go
func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现静态文件服务
}
```

- [ ] **Step 3: 更新 server.go 调用路由**

```go
func newServer(cfg *Config) *Server {
	s := &Server{
		cfg: cfg,
		mux: http.NewServeMux(),
	}
	s.setupRoutes()  // 添加这行
	// ...
}
```

- [ ] **Step 4: Commit**

```bash
git add router.go server.go && git commit -m "feat: add router setup"
```

---

## Chunk 3: 认证模块

### Task 5: utils/path.go - 路径安全工具

**Files:**
- Create: `utils/path.go`

- [ ] **Step 1: 编写路径安全检查**

```go
package utils

import (
	"errors"
	"path/filepath"
	"strings"
)

var ErrAccessDenied = errors.New("access denied")

func SafePath(rootDir, userPath string) (string, error) {
	// 防止空路径攻击
	if userPath == "" {
		userPath = "/"
	}

	// 拼接路径
	absRoot, _ := filepath.Abs(rootDir)
	absPath, _ := filepath.Abs(filepath.Join(rootDir, userPath))

	// 路径遍历防护
	if !strings.HasPrefix(absPath, absRoot) {
		return "", ErrAccessDenied
	}

	return absPath, nil
}

func GetBaseName(path string) string {
	return filepath.Base(path)
}

func GetDir(path string) string {
	return filepath.Dir(path)
}
```

- [ ] **Step 2: Commit**

```bash
git add utils/path.go && git commit -m "feat: add path safety utilities"
```

---

### Task 6: handlers/auth.go - 认证处理

**Files:**
- Create: `handlers/auth.go`

- [ ] **Step 1: 编写认证模块**

```go
package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type AuthHandler struct {
	password   string
	basePath   string
	sessions    map[string]*Session
	mu          sync.RWMutex
}

type Session struct {
	ID        string
	ExpiresAt time.Time
}

type LoginRequest struct {
	Password string `json:"password"`
}

type ApiResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func NewAuthHandler(password, basePath string) *AuthHandler {
	return &AuthHandler{
		password: password,
		basePath: basePath,
		sessions: make(map[string]*Session),
	}
}

// 简单的随机字符串生成
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Password != h.password {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: "Invalid password"})
		return
	}

	// 生成 session
	sessionID := randomString(32)
	h.mu.Lock()
	h.sessions[sessionID] = &Session{
		ID:        sessionID,
		ExpiresAt: time.Now().Add(24 * time.Hours),
	}
	h.mu.Unlock()

	// 设置 Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     h.basePath,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{Code: 0, Msg: "success"})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID, _ := r.Cookie("session_id")
	if sessionID != nil {
		h.mu.Lock()
		delete(h.sessions, sessionID.Value)
		h.mu.Unlock()
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   h.basePath,
		MaxAge: -1,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{Code: 0, Msg: "success"})
}

func (h *AuthHandler) Check(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: "not logged in"})
		return
	}

	h.mu.RLock()
	session, exists := h.sessions[sessionID.Value]
	h.mu.RUnlock()

	if !exists || time.Now().After(session.ExpiresAt) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: "session expired"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{Code: 0, Msg: "ok"})
}

func (h *AuthHandler) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: "not logged in"})
			return
		}

		h.mu.RLock()
		session, exists := h.sessions[sessionID.Value]
		h.mu.RUnlock()

		if !exists || time.Now().After(session.ExpiresAt) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: "session expired"})
			return
		}

		next(w, r)
	}
}

// GetCSRFToken 获取 CSRF Token
func (h *AuthHandler) GetCSRFToken(r *http.Request) string {
	csrf, _ := r.Cookie("csrf_token")
	if csrf != nil {
		return csrf.Value
	}
	return ""
}

// SetCSRFToken 设置 CSRF Token
func (h *AuthHandler) SetCSRFToken(w http.ResponseWriter) {
	csrf := randomString(32)
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrf,
		Path:     h.basePath,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}
```

- [ ] **Step 2: Commit**

```bash
git add handlers/auth.go && git commit -m "feat: add auth handlers"
```

---

## Chunk 4: 文件操作模块

### Task 7: handlers/files.go - 文件操作API

**Files:**
- Create: `handlers/files.go`

- [ ] **Step 1: 编写文件操作 handler**

```go
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"file-manager/utils"
)

type FileHandler struct {
	rootDir string
}

type FileInfo struct {
	Name    string      `json:"name"`
	Path    string      `json:"path"`
	IsDir   bool        `json:"isDir"`
	Size    int64       `json:"size"`
	ModTime interface{} `json:"modTime"`
}

type ListRequest struct {
	Path string `json:"path"`
}

func NewFileHandler(rootDir string) *FileHandler {
	return &FileHandler{rootDir: rootDir}
}

func (h *FileHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		path = "/"
	}

	safePath, err := utils.SafePath(h.rootDir, path)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: "Access denied"})
		return
	}

	entries, err := os.ReadDir(safePath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: err.Error()})
		return
	}

	var files []FileInfo
	for _, entry := range entries {
		info, _ := entry.Info()
		files = append(files, FileInfo{
			Name:    entry.Name(),
			Path:    filepath.Join(path, entry.Name()),
			IsDir:   entry.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime().Format("2006-01-02 15:04:05"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{Code: 0, Msg: "success", Data: files})
}

func (h *FileHandler) Read(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Query().Get("path")
	safePath, err := utils.SafePath(h.rootDir, path)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: "Access denied"})
		return
	}

	data, err := os.ReadFile(safePath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{Code: 0, Msg: "success", Data: string(data)})
}

func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.FormValue("path")
	if path == "" {
		path = "/"
	}

	safePath, err := utils.SafePath(h.rootDir, path)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: "Access denied"})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: err.Error()})
		return
	}
	defer file.Close()

	dstPath := filepath.Join(safePath, header.Filename)
	out, err := os.Create(dstPath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: err.Error()})
		return
	}
	defer out.Close()

	io.Copy(out, file)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{Code: 0, Msg: "success"})
}

func (h *FileHandler) Download(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	safePath, err := utils.SafePath(h.rootDir, path)
	if err != nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(path))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, safePath)
}

type CreateRequest struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

func (h *FileHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateRequest
	json.NewDecoder(r.Body).Decode(&req)

	safePath, err := utils.SafePath(h.rootDir, req.Path)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: "Access denied"})
		return
	}

	dirPath := filepath.Join(safePath, req.Name)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{Code: 0, Msg: "success"})
}

type RenameRequest struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

func (h *FileHandler) Rename(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RenameRequest
	json.NewDecoder(r.Body).Decode(&req)

	safePath, err := utils.SafePath(h.rootDir, req.Path)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: "Access denied"})
		return
	}

	newPath := filepath.Join(filepath.Dir(safePath), req.Name)
	if err := os.Rename(safePath, newPath); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{Code: 0, Msg: "success"})
}

type CopyMoveRequest struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

func (h *FileHandler) Copy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CopyMoveRequest
	json.NewDecoder(r.Body).Decode(&req)

	srcPath, err := utils.SafePath(h.rootDir, req.Src)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: "Access denied"})
		return
	}

	dstPath, err := utils.SafePath(h.rootDir, req.Dst)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: "Access denied"})
		return
	}

	// 简单实现：如果是目录则递归复制
	if err := copyFile(srcPath, dstPath); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{Code: 0, Msg: "success"})
}

func (h *FileHandler) Move(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CopyMoveRequest
	json.NewDecoder(r.Body).Decode(&req)

	srcPath, err := utils.SafePath(h.rootDir, req.Src)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: "Access denied"})
		return
	}

	dstPath, err := utils.SafePath(h.rootDir, req.Dst)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: "Access denied"})
		return
	}

	if err := os.Rename(srcPath, dstPath); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{Code: 0, Msg: "success"})
}

type DeleteRequest struct {
	Paths []string `json:"paths"`
}

func (h *FileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DeleteRequest
	json.NewDecoder(r.Body).Decode(&req)

	for _, p := range req.Paths {
		safePath, err := utils.SafePath(h.rootDir, p)
		if err != nil {
			continue
		}

		info, _ := os.Stat(safePath)
		if info.IsDir() {
			os.RemoveAll(safePath)
		} else {
			os.Remove(safePath)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{Code: 0, Msg: "success"})
}

// Write 实现文件写入（编辑保存）
func (h *FileHandler) Write(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Path string `json:"path"`
		Content string `json:"content"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	safePath, err := utils.SafePath(h.rootDir, req.Path)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: "Access denied"})
		return
	}

	if err := os.WriteFile(safePath, []byte(req.Content), 0644); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{Code: 1, Msg: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{Code: 0, Msg: "success"})
}

// 辅助函数：复制文件
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	info, err := srcFile.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		return copyDir(src, dst)
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		rel, _ := filepath.Rel(src, path)
		dstPath := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath)
	})
}
```

- [ ] **Step 2: 更新 router.go 添加 API 路由**

```go
// handlers 添加到 server
type Server struct {
	cfg          *Config
	mux          *http.ServeMux
	server       *http.Server
	authHandler  *handlers.AuthHandler
	fileHandler  *handlers.FileHandler
}

func newServer(cfg *Config) *Server {
	s := &Server{
		cfg:          cfg,
		mux:          http.NewServeMux(),
		authHandler:  handlers.NewAuthHandler(cfg.password, cfg.basePath),
		fileHandler:  handlers.NewFileHandler(cfg.root),
	}
	s.setupRoutes()
	// ...
}
```

- [ ] **Step 3: 更新 router.go 添加 handler 函数**

```go
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	s.authHandler.Login(w, r)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	s.authHandler.Logout(w, r)
}

func (s *Server) handleCheck(w http.ResponseWriter, r *http.Request) {
	s.authHandler.Check(w, r)
}

func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	s.authHandler.RequireAuth(s.fileHandler.List)(w, r)
}

func (s *Server) handleRead(w http.ResponseWriter, r *http.Request) {
	s.authHandler.RequireAuth(s.fileHandler.Read)(w, r)
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	s.authHandler.RequireAuth(s.fileHandler.Upload)(w, r)
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	s.authHandler.RequireAuth(s.fileHandler.Download)(w, r)
}

func (s *Server) handleCreate(w http.ResponseWriter, r *http.Request) {
	s.authHandler.RequireAuth(s.fileHandler.Create)(w, r)
}

func (s *Server) handleRename(w http.ResponseWriter, r *http.Request) {
	s.authHandler.RequireAuth(s.fileHandler.Rename)(w, r)
}

func (s *Server) handleCopy(w http.ResponseWriter, r *http.Request) {
	s.authHandler.RequireAuth(s.fileHandler.Copy)(w, r)
}

func (s *Server) handleMove(w http.ResponseWriter, r *http.Request) {
	s.authHandler.RequireAuth(s.fileHandler.Move)(w, r)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	s.authHandler.RequireAuth(s.fileHandler.Delete)(w, r)
}

func (s *Server) handleWrite(w http.ResponseWriter, r *http.Request) {
	s.authHandler.RequireAuth(s.fileHandler.Write)(w, r)
}
```

- [ ] **Step 4: 更新 router.go 添加 write 路由**

```go
s.mux.HandleFunc(apiPath+"write", s.handleWrite)
```

- [ ] **Step 5: Commit**

```bash
git add handlers/files.go router.go server.go && git commit -m "feat: add file handlers and API routes"
```

---

## Chunk 5: 静态前端资源

### Task 8: 嵌入静态文件 - index.html

**Files:**
- Create: `static/index.html`

- [ ] **Step 1: 编写 index.html**

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>文件管理器</title>
    <link rel="stylesheet" href="{{basePath}}/static/css/style.css">
</head>
<body>
    <div id="app">
        <!-- 登录页 -->
        <div id="login-page" class="page">
            <div class="login-box">
                <h1>文件管理器</h1>
                <form id="login-form">
                    <input type="password" id="password" placeholder="请输入密码" required>
                    <button type="submit">登录</button>
                </form>
            </div>
        </div>

        <!-- 文件管理页 -->
        <div id="file-page" class="page hidden">
            <div class="toolbar">
                <button id="btn-upload">上传</button>
                <button id="btn-new-folder">新建文件夹</button>
                <button id="btn-refresh">刷新</button>
                <button id="btn-logout">登出</button>
            </div>
            
            <div class="breadcrumb" id="breadcrumb">
                <span class="path-item" data-path="/">根目录</span>
            </div>

            <div class="file-list" id="file-list">
                <table>
                    <thead>
                        <tr>
                            <th><input type="checkbox" id="check-all"></th>
                            <th>名称</th>
                            <th>大小</th>
                            <th>修改时间</th>
                            <th>操作</th>
                        </tr>
                    </thead>
                    <tbody id="file-tbody"></tbody>
                </table>
            </div>

            <div class="drop-zone" id="drop-zone">
                拖拽文件到这里上传
            </div>
        </div>

        <!-- 编辑弹窗 -->
        <div id="edit-modal" class="modal hidden">
            <div class="modal-content">
                <h3>编辑文件</h3>
                <textarea id="editor"></textarea>
                <div class="modal-actions">
                    <button id="btn-save-edit">保存</button>
                    <button id="btn-cancel-edit">取消</button>
                </div>
            </div>
        </div>

        <!-- 确认弹窗 -->
        <div id="confirm-modal" class="modal hidden">
            <div class="modal-content">
                <p id="confirm-message"></p>
                <div class="modal-actions">
                    <button id="btn-confirm">确定</button>
                    <button id="btn-cancel">取消</button>
                </div>
            </div>
        </div>

        <!-- 新建文件夹弹窗 -->
        <div id="new-folder-modal" class="modal hidden">
            <div class="modal-content">
                <h3>新建文件夹</h3>
                <input type="text" id="new-folder-name" placeholder="文件夹名称">
                <div class="modal-actions">
                    <button id="btn-create-folder">创建</button>
                    <button id="btn-cancel-folder">取消</button>
                </div>
            </div>
        </div>

        <!-- 重命名弹窗 -->
        <div id="rename-modal" class="modal hidden">
            <div class="modal-content">
                <h3>重命名</h3>
                <input type="text" id="rename-input" placeholder="新名称">
                <div class="modal-actions">
                    <button id="btn-rename">确定</button>
                    <button id="btn-cancel-rename">取消</button>
                </div>
            </div>
        </div>
    </div>

    <input type="file" id="file-input" multiple style="display: none;">
    
    <script src="{{basePath}}/static/js/app.js"></script>
</body>
</html>
```

- [ ] **Step 2: Commit**

```bash
git add static/index.html && git commit -m "feat: add index.html"
```

---

### Task 9: 样式文件 - style.css

**Files:**
- Create: `static/css/style.css`

- [ ] **Step 1: 编写 style.css**

```css
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    background: #f5f5f5;
}

.hidden {
    display: none !important;
}

.page {
    min-height: 100vh;
}

/* 登录页 */
#login-page {
    display: flex;
    align-items: center;
    justify-content: center;
}

.login-box {
    background: white;
    padding: 40px;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
    width: 320px;
}

.login-box h1 {
    text-align: center;
    margin-bottom: 24px;
    color: #333;
}

.login-box input {
    width: 100%;
    padding: 12px;
    margin-bottom: 16px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 14px;
}

.login-box button {
    width: 100%;
    padding: 12px;
    background: #007bff;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 16px;
    cursor: pointer;
}

.login-box button:hover {
    background: #0056b3;
}

/* 工具栏 */
.toolbar {
    background: white;
    padding: 12px 20px;
    border-bottom: 1px solid #eee;
    display: flex;
    gap: 8px;
}

.toolbar button {
    padding: 8px 16px;
    background: white;
    border: 1px solid #ddd;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
}

.toolbar button:hover {
    background: #f5f5f5;
}

/* 面包屑 */
.breadcrumb {
    background: white;
    padding: 12px 20px;
    border-bottom: 1px solid #eee;
}

.path-item {
    color: #007bff;
    cursor: pointer;
}

.path-item:hover {
    text-decoration: underline;
}

/* 文件列表 */
.file-list {
    padding: 20px;
}

table {
    width: 100%;
    background: white;
    border-collapse: collapse;
    border-radius: 4px;
    overflow: hidden;
}

th, td {
    padding: 12px;
    text-align: left;
    border-bottom: 1px solid #eee;
}

th {
    background: #f8f9fa;
    font-weight: 600;
}

tr:hover {
    background: #f8f9fa;
}

.file-name {
    cursor: pointer;
    color: #007bff;
}

.file-name:hover {
    text-decoration: underline;
}

.file-icon {
    margin-right: 8px;
}

.actions button {
    padding: 4px 8px;
    margin-right: 4px;
    border: 1px solid #ddd;
    background: white;
    border-radius: 4px;
    cursor: pointer;
    font-size: 12px;
}

.actions button:hover {
    background: #f5f5f5;
}

/* 拖拽上传 */
.drop-zone {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 123, 255, 0.1);
    border: 3px dashed #007bff;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;
    color: #007bff;
    z-index: 1000;
    display: none;
}

.drop-zone.active {
    display: flex;
}

/* 弹窗 */
.modal {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0,0,0,0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2000;
}

.modal-content {
    background: white;
    padding: 24px;
    border-radius: 8px;
    width: 90%;
    max-width: 600px;
    max-height: 80vh;
    display: flex;
    flex-direction: column;
}

.modal-content h3 {
    margin-bottom: 16px;
}

.modal-content textarea {
    flex: 1;
    min-height: 300px;
    padding: 12px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-family: monospace;
    resize: none;
}

.modal-content input[type="text"] {
    width: 100%;
    padding: 12px;
    margin-bottom: 16px;
    border: 1px solid #ddd;
    border-radius: 4px;
}

.modal-actions {
    margin-top: 16px;
    display: flex;
    gap: 8px;
    justify-content: flex-end;
}

.modal-actions button {
    padding: 8px 20px;
    border: 1px solid #ddd;
    border-radius: 4px;
    cursor: pointer;
}

.modal-actions button.primary {
    background: #007bff;
    color: white;
    border: none;
}

.modal-actions button.primary:hover {
    background: #0056b3;
}
```

- [ ] **Step 2: Commit**

```bash
git add static/css/style.css && git commit -m "feat: add styles"
```

---

### Task 10: 前端逻辑 - app.js

**Files:**
- Create: `static/js/app.js`

- [ ] **Step 1: 编写 app.js**

```javascript
(function() {
    // 获取 basePath
    const getBasePath = function() {
        const path = window.location.pathname;
        const parts = path.split('/');
        // 移除最后的空字符串和文件名
        parts.pop();
        // 移除最后的空字符串（因为 split 产生）
        if (parts[parts.length - 1] === '') parts.pop();
        return parts.join('/') || '/';
    };

    const basePath = getBasePath();
    const api = function(endpoint) {
        return basePath + endpoint;
    };

    let currentPath = '/';
    let selectedFiles = [];

    // DOM 元素
    const loginPage = document.getElementById('login-page');
    const filePage = document.getElementById('file-page');
    const loginForm = document.getElementById('login-form');
    const fileList = document.getElementById('file-tbody');
    const breadcrumb = document.getElementById('breadcrumb');
    const fileInput = document.getElementById('file-input');
    const dropZone = document.getElementById('drop-zone');

    // 初始化
    function init() {
        checkAuth();
        bindEvents();
    }

    // 检查登录状态
    function checkAuth() {
        fetch(api('/api/check'))
            .then(r => r.json())
            .then(data => {
                if (data.code === 0) {
                    showFilePage();
                    loadFiles(currentPath);
                } else {
                    showLoginPage();
                }
            })
            .catch(() => showLoginPage());
    }

    function showLoginPage() {
        loginPage.classList.remove('hidden');
        filePage.classList.add('hidden');
    }

    function showFilePage() {
        loginPage.classList.add('hidden');
        filePage.classList.remove('hidden');
    }

    // 事件绑定
    function bindEvents() {
        // 登录
        loginForm.addEventListener('submit', function(e) {
            e.preventDefault();
            const password = document.getElementById('password').value;
            
            fetch(api('/api/login'), {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({password: password})
            })
            .then(r => r.json())
            .then(data => {
                if (data.code === 0) {
                    showFilePage();
                    loadFiles(currentPath);
                } else {
                    alert(data.msg || '登录失败');
                }
            });
        });

        // 工具栏按钮
        document.getElementById('btn-upload').addEventListener('click', () => fileInput.click());
        document.getElementById('btn-new-folder').addEventListener('click', showNewFolderModal);
        document.getElementById('btn-refresh').addEventListener('click', () => loadFiles(currentPath));
        document.getElementById('btn-logout').addEventListener('click', logout);

        // 文件选择
        fileInput.addEventListener('change', handleUpload);

        // 拖拽上传
        document.addEventListener('dragover', e => {
            e.preventDefault();
            dropZone.classList.add('active');
        });
        document.addEventListener('dragleave', e => {
            if (e.target === dropZone) {
                dropZone.classList.remove('active');
            }
        });
        document.addEventListener('drop', e => {
            e.preventDefault();
            dropZone.classList.remove('active');
            if (e.dataTransfer.files.length > 0) {
                uploadFiles(e.dataTransfer.files);
            }
        });

        // 全选
        document.getElementById('check-all').addEventListener('change', function() {
            const checks = document.querySelectorAll('.file-check');
            checks.forEach(c => c.checked = this.checked);
            updateSelectedFiles();
        });
    }

    // 加载文件列表
    function loadFiles(path) {
        currentPath = path;
        fetch(api('/api/list?path=' + encodeURIComponent(path)))
            .then(r => r.json())
            .then(data => {
                if (data.code === 0) {
                    renderFileList(data.data);
                    renderBreadcrumb(path);
                }
            });
    }

    // 渲染文件列表
    function renderFileList(files) {
        fileList.innerHTML = '';
        
        files.forEach(file => {
            const tr = document.createElement('tr');
            tr.innerHTML = `
                <td><input type="checkbox" class="file-check" data-path="${file.path}"></td>
                <td>
                    <span class="file-icon">${file.isDir ? '📁' : '📄'}</span>
                    <span class="file-name" data-path="${file.path}" data-isdir="${file.isDir}">${file.name}</span>
                </td>
                <td>${file.isDir ? '-' : formatSize(file.size)}</td>
                <td>${file.modTime}</td>
                <td class="actions">
                    ${!file.isDir ? '<button class="btn-view">查看</button><button class="btn-edit">编辑</button>' : ''}
                    <button class="btn-download">下载</button>
                    <button class="btn-rename">重命名</button>
                    <button class="btn-copy">复制</button>
                    <button class="btn-move">移动</button>
                    <button class="btn-delete">删除</button>
                </td>
            `;
            
            // 点击文件名
            const nameEl = tr.querySelector('.file-name');
            if (file.isDir) {
                nameEl.addEventListener('click', () => loadFiles(file.path));
            } else {
                nameEl.addEventListener('click', () => viewFile(file.path));
            }

            // 操作按钮
            const actions = tr.querySelector('.actions');
            if (!file.isDir) {
                actions.querySelector('.btn-view')?.addEventListener('click', () => viewFile(file.path));
                actions.querySelector('.btn-edit')?.addEventListener('click', () => editFile(file.path));
                actions.querySelector('.btn-download')?.addEventListener('click', () => downloadFile(file.path));
            }
            actions.querySelector('.btn-rename').addEventListener('click', () => showRenameModal(file.path, file.name));
            actions.querySelector('.btn-copy').addEventListener('click', () => showCopyModal(file.path));
            actions.querySelector('.btn-move')?.addEventListener('click', () => showMoveModal(file.path));
            actions.querySelector('.btn-delete').addEventListener('click', () => deleteFile(file.path));

            fileList.appendChild(tr);
        });
    }

    // 渲染面包屑
    function renderBreadcrumb(path) {
        breadcrumb.innerHTML = '';
        const parts = path.split('/').filter(p => p);
        
        let accPath = '';
        const rootSpan = document.createElement('span');
        rootSpan.className = 'path-item';
        rootSpan.textContent = '根目录';
        rootSpan.addEventListener('click', () => loadFiles('/'));
        breadcrumb.appendChild(rootSpan);

        parts.forEach(part => {
            accPath += '/' + part;
            const sep = document.createElement('span');
            sep.textContent = ' > ';
            breadcrumb.appendChild(sep);

            const span = document.createElement('span');
            span.className = 'path-item';
            span.textContent = part;
            span.addEventListener('click', () => loadFiles(accPath));
            breadcrumb.appendChild(span);
        });
    }

    // 格式化文件大小
    function formatSize(bytes) {
        if (bytes < 1024) return bytes + ' B';
        if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
        return (bytes / 1024 / 1024).toFixed(1) + ' MB';
    }

    // 查看文件
    function viewFile(path) {
        fetch(api('/api/read?path=' + encodeURIComponent(path)))
            .then(r => r.json())
            .then(data => {
                if (data.code === 0) {
                    const modal = document.getElementById('edit-modal');
                    document.getElementById('editor').value = data.data;
                    document.getElementById('editor').readOnly = true;
                    document.getElementById('btn-save-edit').classList.add('hidden');
                    modal.classList.remove('hidden');
                }
            });
    }

    // 编辑文件
    function editFile(path) {
        let fileContent = '';
        fetch(api('/api/read?path=' + encodeURIComponent(path)))
            .then(r => r.json())
            .then(data => {
                if (data.code === 0) {
                    fileContent = data.data;
                    const modal = document.getElementById('edit-modal');
                    document.getElementById('editor').value = fileContent;
                    document.getElementById('editor').readOnly = false;
                    document.getElementById('btn-save-edit').classList.remove('hidden');
                    document.getElementById('btn-save-edit').onclick = () => saveFile(path);
                    modal.classList.remove('hidden');
                }
            });
    }

    // 保存文件
    function saveFile(path) {
        const content = document.getElementById('editor').value;
        fetch(api('/api/write'), {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({path: path, content: content})
        })
        .then(r => r.json())
        .then(data => {
            if (data.code === 0) {
                document.getElementById('edit-modal').classList.add('hidden');
                loadFiles(currentPath);
            } else {
                alert(data.msg);
            }
        });
    }

    // 下载文件
    function downloadFile(path) {
        window.location.href = api('/api/download?path=' + encodeURIComponent(path));
    }

    // 上传文件
    function handleUpload(e) {
        uploadFiles(e.target.files);
    }

    function uploadFiles(files) {
        const formData = new FormData();
        formData.append('path', currentPath);
        
        for (let i = 0; i < files.length; i++) {
            formData.append('file', files[i]);
        }

        fetch(api('/api/upload'), {
            method: 'POST',
            body: formData
        })
        .then(r => r.json())
        .then(data => {
            if (data.code === 0) {
                loadFiles(currentPath);
            } else {
                alert(data.msg);
            }
        });
    }

    // 新建文件夹
    function showNewFolderModal() {
        document.getElementById('new-folder-modal').classList.remove('hidden');
    }

    document.getElementById('btn-create-folder').addEventListener('click', function() {
        const name = document.getElementById('new-folder-name').value;
        if (!name) return;

        fetch(api('/api/create'), {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({path: currentPath, name: name})
        })
        .then(r => r.json())
        .then(data => {
            if (data.code === 0) {
                document.getElementById('new-folder-modal').classList.add('hidden');
                document.getElementById('new-folder-name').value = '';
                loadFiles(currentPath);
            } else {
                alert(data.msg);
            }
        });
    });

    document.getElementById('btn-cancel-folder').addEventListener('click', () => {
        document.getElementById('new-folder-modal').classList.add('hidden');
    });

    // 重命名
    let renamePath = '';
    function showRenameModal(path, name) {
        renamePath = path;
        document.getElementById('rename-input').value = name;
        document.getElementById('rename-modal').classList.remove('hidden');
    }

    document.getElementById('btn-rename').addEventListener('click', function() {
        const name = document.getElementById('rename-input').value;
        if (!name) return;

        fetch(api('/api/rename'), {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({path: renamePath, name: name})
        })
        .then(r => r.json())
        .then(data => {
            if (data.code === 0) {
                document.getElementById('rename-modal').classList.add('hidden');
                loadFiles(currentPath);
            } else {
                alert(data.msg);
            }
        });
    });

    document.getElementById('btn-cancel-rename').addEventListener('click', () => {
        document.getElementById('rename-modal').classList.add('hidden');
    });

    // 删除
    function deleteFile(path) {
        if (!confirm('确定要删除吗？')) return;

        fetch(api('/api/delete'), {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({paths: [path]})
        })
        .then(r => r.json())
        .then(data => {
            if (data.code === 0) {
                loadFiles(currentPath);
            } else {
                alert(data.msg);
            }
        });
    }

    // 复制/移动 (简化版)
    function showCopyModal(src) {
        const dst = prompt('请输入目标路径（如 /folder/name）:');
        if (!dst) return;

        fetch(api('/api/copy'), {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({src: src, dst: dst})
        })
        .then(r => r.json())
        .then(data => {
            if (data.code === 0) {
                loadFiles(currentPath);
            } else {
                alert(data.msg);
            }
        });
    }

    function showMoveModal(src) {
        const dst = prompt('请输入目标路径（如 /folder/name）:');
        if (!dst) return;

        fetch(api('/api/move'), {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({src: src, dst: dst})
        })
        .then(r => r.json())
        .then(data => {
            if (data.code === 0) {
                loadFiles(currentPath);
            } else {
                alert(data.msg);
            }
        });
    }

    // 取消编辑
    document.getElementById('btn-cancel-edit').addEventListener('click', () => {
        document.getElementById('edit-modal').classList.add('hidden');
    });

    // 登出
    function logout() {
        fetch(api('/api/logout'), {method: 'POST'})
            .then(() => {
                showLoginPage();
            });
    }

    // 更新选中的文件
    function updateSelectedFiles() {
        selectedFiles = Array.from(document.querySelectorAll('.file-check:checked')).map(c => c.dataset.path);
    }

    init();
})();
```

- [ ] **Step 2: Commit**

```bash
git add static/js/app.js && git commit -m "feat: add frontend JavaScript"
```

---

### Task 11: 使用 embed 嵌入静态文件

**Files:**
- Modify: `server.go`

- [ ] **Step 1: 创建 embed 存根**

```go
//go:embed static
var staticFS embed.FS
```

- [ ] **Step 2: 更新 handleStatic**

```go
func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
    // 去掉 basePath 和 /static/ 前缀
    path := strings.TrimPrefix(r.URL.Path, s.cfg.basePath)
    path = strings.TrimPrefix(path, "/static/")
    if path == "" {
        path = "index.html"
    }
    
    data, err := staticFS.ReadFile("static/" + path)
    if err != nil {
        http.NotFound(w, r)
        return
    }
    
    // 根据文件类型设置 Content-Type
    var contentType string
    switch {
    case strings.HasSuffix(path, ".html"):
        contentType = "text/html"
    case strings.HasSuffix(path, ".css"):
        contentType = "text/css"
    case strings.HasSuffix(path, ".js"):
        contentType = "application/javascript"
    case strings.HasSuffix(path, ".png"):
        contentType = "image/png"
    case strings.HasSuffix(path, ".jpg"), strings.HasSuffix(path, ".jpeg"):
        contentType = "image/jpeg"
    default:
        contentType = "text/plain"
    }
    
    w.Header().Set("Content-Type", contentType)
    w.Write(data)
}
```

- [ ] **Step 3: 更新 handleIndex**

```go
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
    data, err := staticFS.ReadFile("static/index.html")
    if err != nil {
        http.NotFound(w, r)
        return
    }
    // 替换占位符为实际的 basePath
    content := strings.Replace(string(data), "{{basePath}}", s.cfg.basePath, -1)
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte(content))
}
```

- [ ] **Step 4: Commit**

```bash
git add server.go && git commit -m "feat: embed static files"
```

---

## Chunk 6: 测试与构建

### Task 12: 构建与测试

- [ ] **Step 1: 构建项目**

Run: `go build -o file-manager`

- [ ] **Step 2: 测试运行**

```bash
./file-manager -root /tmp/test -password test123 -port 8080 -basepath /
```

- [ ] **Step 3: 测试各项功能**
- 登录功能
- 文件浏览
- 文件上传/下载
- 文件编辑
- 重命名/复制/移动/删除

- [ ] **Step 4: 测试 basepath**

```bash
./file-manager -root /tmp/test -password test123 -port 8081 -basepath /files
```

- [ ] **Step 5: 最终提交**

```bash
git add . && git commit -m "feat: complete file manager with all features"
```

---

## 验收检查清单

- [ ] 密码登录功能正常
- [ ] 可浏览指定根目录下的文件
- [ ] 可上传/下载文件
- [ ] 可查看文本文件内容
- [ ] 可编辑并保存文本文件
- [ ] 可重命名/复制/移动/删除
- [ ] 部署在 basepath 时所有功能正常
- [ ] 路径遍历防护有效
