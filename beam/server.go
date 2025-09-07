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
	"strings"
	"time"

	"github.com/tachRoutine/beamdrop-go/config"
	"github.com/tachRoutine/beamdrop-go/pkg/logger"
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

func StartServer(sharedDir string) {
	logger.Info("Initializing HTTP handlers")
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		if urlPath == "/" {
			urlPath = "/index.html"
		}

		logger.Debug("Serving static file: %s", urlPath)
		file, err := static.FrontendFiles.Open("frontend" + urlPath)
		if err != nil {
			logger.Warn("Static file not found: %s", urlPath)
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

	// File APIs
	http.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Listing files from directory: %s", sharedDir)
		files, err := os.ReadDir(sharedDir)
		if err != nil {
			logger.Error("Failed to read directory %s: %v", sharedDir, err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to read directory"})
			return
		}
		
		var fileList []File
		for _, f := range files {
			fileInfo, err := f.Info()
			if err != nil {
				logger.Warn("Failed to get file info for %s: %v", f.Name(), err)
				continue
			}
			file := File{
				Name:    fileInfo.Name(),
				IsDir:   fileInfo.IsDir(),
				Size:    FormatFileSize(fileInfo.Size()),
				ModTime: FormatModTime(fileInfo.ModTime().Format(time.RFC3339)),
				Path:    path.Join(sharedDir, fileInfo.Name()),
			}
			fileList = append(fileList, file)
		}
		logger.Debug("Found %d files/directories", len(fileList))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fileList)
	})

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("file")
		filePath := sharedDir + "/" + filename
		
		logger.Info("Download request for file: %s", filename)
		f, err := os.Open(filePath)
		if err != nil {
			logger.Error("Failed to open file %s: %v", filePath, err)
			http.Error(w, "File not found", 404)
			return
		}
		defer f.Close()
		
		logger.Info("Serving download for file: %s", filename)
		io.Copy(w, f)
		logger.Info("Download completed for file: %s", filename)
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Upload request received")
		file, header, err := r.FormFile("file")
		if err != nil {
			logger.Error("Invalid upload request: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid upload"})
			return
		}
		defer file.Close()
		
		filePath := sharedDir + "/" + header.Filename
		logger.Info("Uploading file: %s (size: %d bytes)", header.Filename, header.Size)
		
		out, err := os.Create(filePath)
		if err != nil {
			logger.Error("Failed to create file %s: %v", filePath, err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save file"})
			return
		}
		defer out.Close()
		
		_, err = io.Copy(out, file)
		if err != nil {
			logger.Error("Failed to write file %s: %v", filePath, err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to write file"})
			return
		}
		
		logger.Info("File uploaded successfully: %s", header.Filename)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Uploaded", "file": header.Filename})
	})

	ip := GetLocalIP()
	url := fmt.Sprintf("http://%s:%d", ip, config.GetConfig().PORT)

	logger.Info("Generating QR code for URL: %s", url)
	qr.ShowQrCode(url)
	logger.Info("Server started at %s sharing directory: %s", url, sharedDir)
	
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.GetConfig().PORT), nil)
	if err != nil {
		logger.Fatal("Server error: %v", err)
	}
}

func GetLocalIP() string {
	logger.Debug("Detecting local IP address")
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logger.Warn("Failed to get network interfaces: %v", err)
		return "localhost"
	}
	
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				logger.Debug("Found local IP: %s", ipnet.IP.String())
				return ipnet.IP.String()
			}
		}
	}
	
	logger.Warn("No local IP found, using localhost")
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
