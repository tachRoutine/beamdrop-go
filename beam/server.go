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
	fs := http.FileServer(http.FS(staticDir))
    http.Handle("/", fs)

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