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

	// File APIs
	http.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		files, _ := os.ReadDir(sharedDir)
		var fileList []File
		for _, f := range files {
			fileInfo, _ := f.Info()
			file := File{
				Name:    fileInfo.Name(),
				IsDir:   fileInfo.IsDir(),
				Size:    FormatFileSize(fileInfo.Size()),
				ModTime: FormatModTime(fileInfo.ModTime().Format(time.RFC3339)),
				Path:    path.Join(sharedDir, fileInfo.Name()),
			}
			fileList = append(fileList, file)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fileList)
	})

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open(sharedDir + "/" + r.URL.Query().Get("file"))
		fmt.Println("Downloading file:", f.Name())
		if err != nil {
			http.Error(w, "File not found", 404)
			return
		}
		defer f.Close()
		fmt.Println("Downloading file:", f.Name())
		io.Copy(w, f)
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("file")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid upload"})
			return
		}
		defer file.Close()
		out, err := os.Create(sharedDir + "/" + header.Filename)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Uploaded", "file": header.Filename})
	})

	ip := GetLocalIP()
	url := fmt.Sprintf("http://%s:%d", ip, config.GetConfig().PORT)

	qr.ShowQrCode(url)
	fmt.Println("Server started at", url, "sharing directory:", sharedDir)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.GetConfig().PORT), nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
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
