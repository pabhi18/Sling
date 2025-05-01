package cmd

import (
	"fmt"
	"os"
	"net/http"
	"github.com/spf13/cobra"
	"github.com/pabhi18/sling/pkg/server"
	"github.com/pabhi18/sling/pkg/utils"
)

var (
	path string
	port string
	ip   string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve a folder over HTTP",

	Run: func(cmd *cobra.Command, args []string) {
		serveFolder(path, port)
	},
}

func serveFolder(folderPath, port string) {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		fmt.Printf("❌ Folder does not exist: %v\n", err)
		os.Exit(1)
	}

	localIp := utils.GetIp()

	fmt.Println("✅ Server started successfully!")
	fmt.Println("🌐 Access URLs:")
	fmt.Printf("   ➤ Localhost:      http://localhost:%s\n", port)
	fmt.Printf("   ➤ Local Network:  http://%s:%s\n", localIp, port)
	fmt.Printf("📁 Serving directory: %s\n", folderPath)

	tmpl := "templates/index.html"
	handler := server.CustomFileServer(http.Dir(folderPath), tmpl)
	http.Handle("/", handler)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("🚨 Error starting server: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&path, "path", "p", ".", "Path to folder to serve")
	serveCmd.Flags().StringVar(&port, "port", "8080", "Port to serve on")
}
