package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type AuthHandler struct {
	password string
	basePath string
	sessions map[string]*Session
	mu       sync.RWMutex
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

	sessionID := randomString(32)
	h.mu.Lock()
	h.sessions[sessionID] = &Session{
		ID:        sessionID,
		ExpiresAt: time.Now().Add(24 * time.Hours),
	}
	h.mu.Unlock()

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

func (h *AuthHandler) GetCSRFToken(r *http.Request) string {
	csrf, _ := r.Cookie("csrf_token")
	if csrf != nil {
		return csrf.Value
	}
	return ""
}

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
