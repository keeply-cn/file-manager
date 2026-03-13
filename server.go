package main

import (
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	cfg          *Config
	mux          *http.ServeMux
	server       *http.Server
	authHandler  interface{}
	fileHandler  interface{}
}

func newServer(cfg *Config) *Server {
	s := &Server{
		cfg: cfg,
		mux: http.NewServeMux(),
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
