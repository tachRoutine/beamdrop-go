package beam

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/tachRoutine/beamdrop-go/config"
	"github.com/tachRoutine/beamdrop-go/pkg/qr"
	"github.com/tachRoutine/beamdrop-go/static"
)

type File struct {
	Name    string `json:"name"`
	Size    string `json:"size"`
	IsDir   bool   `json:"isDir"`
	ModTime string `json:"modTime"`
	Path    string `json:"path"`
}

type ServerConfig struct {
	SharedDir string
	Port      int
	Password  string
	Verbose   bool
	NoQR      bool
}

type Server struct {
	config ServerConfig
	stats  *ServerStats
}

type ServerStats struct {
	Downloads int    `json:"downloads"`
	Uploads   int    `json:"uploads"`
	StartTime string `json:"startTime"`
}

func StartServer(cfg ServerConfig) {
	server := &Server{
		config: cfg,
		stats: &ServerStats{
			StartTime: time.Now().Format(time.RFC3339),
		},
	}

	// Use config port or default
	port := cfg.Port
	if port == 0 {
		port = config.GetConfig().PORT
	}

	server.setupRoutes()

	ip := GetLocalIP()
	url := fmt.Sprintf("http://%s:%d", ip, port)

	if !cfg.NoQR {
		qr.ShowQrCode(url)
	}
	fmt.Printf("Server started at %s sharing directory: %s\n", url, cfg.SharedDir)

	if cfg.Password != "" {
		fmt.Println("🔒 Password protection enabled")
	}

	if cfg.Verbose {
		fmt.Println("📊 Verbose logging enabled")
	}

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

func (s *Server) setupRoutes() {
	// Middleware for authentication
	authMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if s.config.Password != "" {
				password := r.Header.Get("X-Password")
				if password != s.config.Password {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(map[string]string{"error": "Authentication required"})
					return
				}
			}
			next(w, r)
		}
	}

	// Static files and frontend
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		if urlPath == "/" {
			urlPath = "/index.html"
		}

		file, err := static.FrontendFiles.Open("frontend" + urlPath)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Not found"})
			return
		}
		defer file.Close()

		ext := strings.ToLower(path.Ext(urlPath))
		if mimeType := mime.TypeByExtension(ext); mimeType != "" {
			w.Header().Set("Content-Type", mimeType)
		} else {
			w.Header().Set("Content-Type", "application/octet-stream")
		}

		io.Copy(w, file)
	})

	// API Routes
	http.HandleFunc("/files", authMiddleware(s.handleFiles))
	http.HandleFunc("/download", authMiddleware(s.handleDownload))
	http.HandleFunc("/upload", authMiddleware(s.handleUpload))
	http.HandleFunc("/delete", authMiddleware(s.handleDelete))
	http.HandleFunc("/preview", authMiddleware(s.handlePreview))
	http.HandleFunc("/stats", s.handleStats)
}

func (s *Server) handleFiles(w http.ResponseWriter, r *http.Request) {
	relativePath := r.URL.Query().Get("path")
	if relativePath == "" {
		relativePath = "."
	}

	fullPath := filepath.Join(s.config.SharedDir, relativePath)

	// Security check: ensure we don't go outside shared directory
	if !strings.HasPrefix(fullPath, s.config.SharedDir) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Access denied"})
		return
	}

	files, err := os.ReadDir(fullPath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Directory not found"})
		return
	}

	var fileList []File
	for _, f := range files {
		fileInfo, _ := f.Info()
		file := File{
			Name:    fileInfo.Name(),
			IsDir:   fileInfo.IsDir(),
			Size:    FormatFileSize(fileInfo.Size()),
			ModTime: FormatModTime(fileInfo.ModTime().Format(time.RFC3339)),
			Path:    filepath.Join(relativePath, fileInfo.Name()),
		}
		fileList = append(fileList, file)
	}

	if s.config.Verbose {
		fmt.Printf("📁 Listed %d items in directory: %s\n", len(fileList), relativePath)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileList)
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	fullPath := filepath.Join(s.config.SharedDir, fileName)

	// Security check
	if !strings.HasPrefix(fullPath, s.config.SharedDir) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	f, err := os.Open(fullPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	stat, _ := f.Stat()
	w.Header().Set("Content-Disposition", "attachment; filename=\""+stat.Name()+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))

	s.stats.Downloads++
	if s.config.Verbose {
		fmt.Printf("⬇️  Downloaded: %s (%s)\n", fileName, FormatFileSize(stat.Size()))
	}

	io.Copy(w, f)
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid upload"})
		return
	}
	defer file.Close()

	uploadPath := r.FormValue("path")
	if uploadPath == "" {
		uploadPath = "."
	}

	fullDir := filepath.Join(s.config.SharedDir, uploadPath)
	fullPath := filepath.Join(fullDir, header.Filename)

	// Security check
	if !strings.HasPrefix(fullPath, s.config.SharedDir) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Access denied"})
		return
	}

	// Create directory if it doesn't exist
	os.MkdirAll(fullDir, 0755)

	out, err := os.Create(fullPath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save file"})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to write file"})
		return
	}

	s.stats.Uploads++
	if s.config.Verbose {
		fmt.Printf("⬆️  Uploaded: %s (%d bytes)\n", header.Filename, header.Size)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Uploaded successfully", "file": header.Filename})
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileName := r.URL.Query().Get("file")
	fullPath := filepath.Join(s.config.SharedDir, fileName)

	// Security check
	if !strings.HasPrefix(fullPath, s.config.SharedDir) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Access denied"})
		return
	}

	err := os.Remove(fullPath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "File not found or cannot be deleted"})
		return
	}

	if s.config.Verbose {
		fmt.Printf("🗑️  Deleted: %s\n", fileName)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "File deleted successfully"})
}

func (s *Server) handlePreview(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	fullPath := filepath.Join(s.config.SharedDir, fileName)

	// Security check
	if !strings.HasPrefix(fullPath, s.config.SharedDir) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	f, err := os.Open(fullPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	// Detect content type
	ext := strings.ToLower(filepath.Ext(fileName))
	var contentType string

	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		contentType = mime.TypeByExtension(ext)
	case ".txt", ".md", ".log", ".json", ".xml", ".csv":
		contentType = "text/plain"
	default:
		contentType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", contentType)
	io.Copy(w, f)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.stats)
}

func GetLocalIP() string {
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "localhost"
}

func FormatFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB", float64(size)/(1024*1024*1024))
	}
}

func FormatModTime(modTime string) string {
	t, err := time.Parse(time.RFC3339, modTime)
	if err != nil {
		return modTime
	}
	return t.Format("2006-01-02 15:04:05")
}
