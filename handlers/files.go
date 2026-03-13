package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

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

func (h *FileHandler) Write(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Path    string `json:"path"`
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
