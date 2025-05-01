package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetIp() string {
	cmd := exec.Command("sh", "-c", "ifconfig | grep inet | grep -v 127.0.0.1 | grep -v inet6 | awk '{ print $2 }'")
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("⚠️  Failed to get local IP address: %v\n", err)
		return ""
	}

	return strings.TrimSpace(string(output))
}