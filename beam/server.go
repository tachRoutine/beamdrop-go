package beam

import (
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/tachRoutine/beamdrop-go/static"
)

func StartServer(sharedDir string) string {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		if urlPath == "/" {
			urlPath = "/index.html"
		}

		// Open embedded file
		file, err := static.FrontendFiles.Open("frontend" + urlPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer file.Close()

		// Detecting MIME type
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
		for _, f := range files {
			fmt.Fprintln(w, f.Name())
		}
	})

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open(sharedDir + "/" + r.URL.Query().Get("file"))
		if err != nil {
			http.Error(w, "File not found", 404)
			return
		}
		defer f.Close()
		io.Copy(w, f)
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, header, _ := r.FormFile("file")
		defer file.Close()
		out, _ := os.Create(sharedDir + "/" + header.Filename)
		defer out.Close()
		io.Copy(out, file)
		fmt.Fprintln(w, "Uploaded")
	})

	ip := getLocalIP()
	url := fmt.Sprintf("http://%s:8080", ip)
	fmt.Println("Server started at", url, "sharing directory:", sharedDir)
	http.ListenAndServe(":8080", nil)
	return url
}

func getLocalIP() string {
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
