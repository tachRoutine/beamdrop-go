package beam

import (
    "fmt"
    "net"
    "net/http"
    "os"
    "io"
)

func StartServer() {
	fs := http.FileServer(http.Dir("./static"))
    http.Handle("/", fs)

    http.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
        files, _ := os.ReadDir("./shared")
        for _, f := range files {
            fmt.Fprintln(w, f.Name())
        }
    })

    http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
        f, err := os.Open("./shared/" + r.URL.Query().Get("file"))
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
        out, _ := os.Create("./shared/" + header.Filename)
        defer out.Close()
        io.Copy(out, file)
        fmt.Fprintln(w, "Uploaded")
    })

    ip := getLocalIP()
    fmt.Println("Open http://" + ip + ":8080")
    http.ListenAndServe(":8080", nil)
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