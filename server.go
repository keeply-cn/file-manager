package main

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	"file-manager/handlers"
)

//go:embed static
var staticFS embed.FS

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
