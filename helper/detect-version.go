package helper

import (
	"fmt"
	"os/exec"
	"strings"
)

func DetectGoVersion() (string, error) {
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		return "", err
	}

	parts := strings.Fields(string(out))
	if len(parts) > 3 {
		return "", fmt.Errorf("unexpected go version output: %q", out)
	}

	return strings.TrimPrefix(parts[2], "go"), nil
}
