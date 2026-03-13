package main

import (
	"net/http"
	"strings"
)

func (s *Server) setupRoutes() {
	// Use a catch-all handler that routes to specific handlers
	s.mux.HandleFunc("/", s.handleRequest)
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	
	// Check if path starts with basePath
	basePath := s.cfg.basePath
	if basePath == "/" {
		basePath = ""
	}
	
	// API routes
	apiPrefix := basePath + "/api/"
	if strings.HasPrefix(path, apiPrefix) {
		apiPath := strings.TrimPrefix(path, apiPrefix)
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
	
	// Static files
	staticPrefix := basePath + "/static/"
	if strings.HasPrefix(path, staticPrefix) {
		s.handleStatic(w, r)
		return
	}
	
	// Index page
	if path == s.cfg.basePath || path == s.cfg.basePath+"/" || (s.cfg.basePath == "/" && path == "/") {
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

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
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
