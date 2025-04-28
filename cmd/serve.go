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

    // Load the HTML template
    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        fmt.Println("Error loading template:", err)
        os.Exit(1)
    }

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("search")
		var files []string
		var err error
	
		// If there's a search query, filter files
		if query != "" {
			files, err = searchFiles(folderPath, query)
		} else {
			files, err = getFilesInDir(folderPath)
		}
	
		if err != nil {
			http.Error(w, "Error reading files", http.StatusInternalServerError)
			return
		}
	
		data := struct {
			Path  string
			Files []string
			SearchQuery string
		}{
			Path:       folderPath,
			Files:      files,
			SearchQuery: query,
		}
	
		// Render the template
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
	})

    err = http.ListenAndServe(":"+port, nil)
    if err != nil {
        fmt.Println("Error starting server:", err)
        os.Exit(1)
    }
}

func searchFiles(folderPath, query string) ([]string, error) {
    var result []string
    files, err := getFilesInDir(folderPath)
    if err != nil {
        return nil, err
    }
    for _, file := range files {
        if strings.Contains(file, query) {
            result = append(result, file)
        }
    }
    return result, nil
}



func getFilesInDir(folderPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
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
