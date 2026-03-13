package main

import (
	"net/http"
	"strings"
)

func (s *Server) setupRoutes() {
	s.mux.HandleFunc(s.cfg.basePath, s.handleIndex)
	s.mux.HandleFunc(s.cfg.basePath+"/", s.handleIndex)
	s.mux.HandleFunc(s.cfg.basePath+"/static/", s.handleStatic)
	
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
	s.mux.HandleFunc(apiPath+"write", s.handleWrite)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	data, err := staticFS.ReadFile("static/index.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	content := strings.Replace(string(data), "{{basePath}}", s.cfg.basePath, -1)
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
