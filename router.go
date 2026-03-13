package main

import (
	"net/http"
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
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != s.cfg.basePath {
		http.NotFound(w, r)
		return
	}
	s.handleStatic(w, r)
}

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现静态文件服务
}
