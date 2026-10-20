package web

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed all:web/dist
var webFiles embed.FS

func RegisterStaticFiles(mux *http.ServeMux) {
	distFS, err := fs.Sub(webFiles, "web/dist")
	if err != nil {
		panic("failed to embed web/dist: " + err.Error())
	}

	fileServer := http.FileServer(http.FS(distFS))
	index, err := fs.ReadFile(distFS, "index.html")
	if err != nil {
		panic("failed to read embedded index.html: " + err.Error())
	}

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cleanPath := strings.TrimPrefix(path.Clean("/"+r.URL.Path), "/")
		if cleanPath == "" || cleanPath == "." {
			cleanPath = "index.html"
		}
		if f, err := distFS.Open(cleanPath); err == nil {
			_ = f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}

		// Preserve API 404s rather than returning the SPA shell for misspelled
		// API URLs. All other extensionless paths are client-side routes.
		if r.URL.Path == "/info" || r.URL.Path == "/history" || strings.HasPrefix(r.URL.Path, "/history/") || path.Ext(r.URL.Path) != "" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
			return
		}
		_, _ = w.Write(index)
	}))
}
