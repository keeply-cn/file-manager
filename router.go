package main

import (
	"net/http"
	"strings"
)

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/", s.handleRequest)
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	
	// basePath 用于 API 和静态资源路径匹配（去除外部路径前缀）
	// 但路由始终注册在 "/" 上，由端口转发处理外部路径
	externalBasePath := s.cfg.basePath
	if externalBasePath == "/" {
		externalBasePath = ""
	}
	
	// 去除外部 basepath 前缀，得到内部路径
	internalPath := path
	if externalBasePath != "" && strings.HasPrefix(path, externalBasePath) {
		internalPath = strings.TrimPrefix(path, externalBasePath)
		if !strings.HasPrefix(internalPath, "/") {
			internalPath = "/" + internalPath
		}
	}
	
	// API routes - 匹配内部路径 /api/*
	if strings.HasPrefix(internalPath, "/api/") {
		apiPath := strings.TrimPrefix(internalPath, "/api/")
		switch apiPath {
		case "login":
			s.handleLogin(w, r)
		case "logout":
			s.handleLogout(w, r)
		case "check":
			s.handleCheck(w, r)
		case "list":
			s.handleList(w, r)
		case "read":
			s.handleRead(w, r)
		case "upload":
			s.handleUpload(w, r)
		case "download":
			s.handleDownload(w, r)
		case "create":
			s.handleCreate(w, r)
		case "rename":
			s.handleRename(w, r)
		case "copy":
			s.handleCopy(w, r)
		case "move":
			s.handleMove(w, r)
		case "delete":
			s.handleDelete(w, r)
		case "write":
			s.handleWrite(w, r)
		default:
			http.NotFound(w, r)
		}
		return
	}
	
	// Static files - 匹配内部路径 /static/*
	if strings.HasPrefix(internalPath, "/static/") {
		s.handleStatic(w, r, internalPath)
		return
	}
	
	// Index page - 匹配内部路径 / 或空
	if internalPath == "/" || internalPath == "" {
		s.handleIndex(w, r)
		return
	}
	
	http.NotFound(w, r)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	data, err := staticFS.ReadFile("static/index.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	basePath := s.cfg.basePath
	if basePath == "/" {
		basePath = ""
	}
	content := strings.Replace(string(data), "{{basePath}}", basePath, -1)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(content))
}

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request, internalPath string) {
	// internalPath is like /static/xxx, extract xxx
	path := strings.TrimPrefix(internalPath, "/static/")
	if path == "" {
		path = "index.html"
	}
	
	data, err := staticFS.ReadFile("static/" + path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	
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
