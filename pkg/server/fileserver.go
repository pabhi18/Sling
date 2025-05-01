package server

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

func CustomFileServer(root http.FileSystem, templatePath string) http.Handler {
	fs := http.FileServer(root)
	tmpl := template.Must(template.ParseFiles(templatePath))

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
