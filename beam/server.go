package beam

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/tachRoutine/ekiliBeam-go/static"
)

func StartServer(sharedDir string) string {
	staticDir := static.FrontendFiles
	// Serve assets under /assets/ with correct MIME type
	assetsFs := http.FileServer(http.FS(staticDir))
	http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len("/assets/"):] // remove /assets/ prefix
		switch {
		case len(path) > 3 && path[len(path)-3:] == ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case len(path) > 4 && path[len(path)-4:] == ".css":
			w.Header().Set("Content-Type", "text/css")
		}
		http.StripPrefix("/assets/", assetsFs).ServeHTTP(w, r)
	})
	// Optionally, serve index.html at root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/frontend/index.html")
	})

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
	url := fmt.Sprintf("Open http://%s:8080", ip)
	fmt.Println("Server started at", url)
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
