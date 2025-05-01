package server

import (
	"os"
	"fmt"
	"io/ioutil"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

func CustomFileServer(root http.FileSystem) http.Handler {
	fs := http.FileServer(root)

	tmplURL := "https://raw.githubusercontent.com/pabhi18/Sling/main/templates/index.html"
	tmplResp, err := http.Get(tmplURL)
	if err != nil {
		fmt.Printf("❌ Error fetching template: %v\n", err)
		os.Exit(1)
	}
	defer tmplResp.Body.Close()

	tmplContent, err := ioutil.ReadAll(tmplResp.Body)
	if err != nil {
		fmt.Printf("❌ Error reading template body: %v\n", err)
		os.Exit(1)
	}

	tmpl, err := template.New("index").Parse(string(tmplContent))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqPath := filepath.Clean(r.URL.Path)

		f, err := root.Open(reqPath)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer f.Close()

		info, err := f.Stat()
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		if info.IsDir() {
			entries, err := f.Readdir(-1)
			if err != nil {
				http.Error(w, "Error reading directory", http.StatusInternalServerError)
				return
			}

			searchQuery := r.URL.Query().Get("search")
			files := []struct {
				Name string
				URL  string
			}{}

			for _, entry := range entries {
				name := entry.Name()
				if searchQuery != "" && !strings.Contains(strings.ToLower(name), strings.ToLower(searchQuery)) {
					continue
				}

				linkPath := filepath.Join(reqPath, name)
				linkPath = strings.TrimPrefix(linkPath, "/")
				files = append(files, struct {
					Name string
					URL  string
				}{
					Name: name,
					URL:  "/" + strings.ReplaceAll(linkPath, "\\", "/"),
				})
			}

			data := struct {
				Path        string
				Files       []struct {
					Name string
					URL  string
				}
				SearchQuery string
			}{
				Path:        r.URL.Path,
				Files:       files,
				SearchQuery: searchQuery,
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := tmpl.Execute(w, data); err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
			}
			return
		}

		fs.ServeHTTP(w, r)
	})
}
