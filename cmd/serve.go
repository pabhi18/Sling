/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "github.com/spf13/cobra"
	"os/exec"
	"strings"
	"html/template"
)

var (
	path string
	port string
	ip string
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve a folder over HTTP",

	Run: func(cmd *cobra.Command, args []string) {
		serveFolder(path, port)
	},
}

func serveFolder(folderPath, port string) {
    if _, err := os.Stat(folderPath); os.IsNotExist(err) {
        fmt.Println("Folder doesn't exist:", err)
        os.Exit(1)
    }
    localIp := getIp()

    fmt.Printf("Serving %s at http://localhost:%s and for Local Networks at http://%s:%s\n", folderPath, port, localIp, port)

    tmpl := "templates/index.html"

	handler := CustomFileServer(http.Dir(folderPath), tmpl)

	http.Handle("/", handler)

    err := http.ListenAndServe(":"+port, nil)
    if err != nil {
        fmt.Println("Error starting server: ", err)
        os.Exit(1)
    }
}

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

			// Filtering by search query
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
				linkPath = strings.TrimPrefix(linkPath, "/") // Avoid double slashes
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

func getIp() string {
	cmd := exec.Command("sh", "-c", "ifconfig | grep inet | grep -v 127.0.0.1 | grep -v inet6 | awk '{ print $2 }'")
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Failed to get local IP:", err)
		return ""
	}

	return strings.TrimSpace(string(output))
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&path, "path", "p", ".", "Path to folder to serve")
    serveCmd.Flags().StringVar(&port, "port", "8080", "Port to serve on")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
