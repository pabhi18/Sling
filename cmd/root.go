package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "1.0.0"

const logo = `

▒█▀▀▀█ ▒█░░░ ▀█▀ ▒█▄░▒█ ▒█▀▀█ 
░▀▀▀▄▄ ▒█░░░ ▒█░ ▒█▒█▒█ ▒█░▄▄ 
▒█▄▄▄█ ▒█▄▄█ ▄█▄ ▒█░░▀█ ▒█▄▄█
`

var rootCmd = &cobra.Command{
	Use:   "sling",
	Short: "Sling - A simple CLI to host any local folder over HTTP",
	Long: fmt.Sprintf(`%s
Sling is a lightweight CLI tool that allows you to serve any folder 
on your system as a local web server using HTTP.

Examples:
  sling serve -p . --port 3000     # Serve current folder on port 3000
  sling serve --path ./site        # Serve 'site' folder on default port 8080

Use 'sling serve --help' to see available options.
`, logo),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(logo)
		fmt.Println("Welcome to Sling - serve any folder via HTTP")
		fmt.Println("Use 'sling serve --help' to get started.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}