package helper

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// DetectModulePath search go.mod in actual o parent directory
// and return module path (the word after "module").
func DetectModulePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		gm := filepath.Join(dir, "go.mod")
		if info, errStat := os.Stat(gm); errStat == nil && !info.IsDir() {
			f, errOpen := os.Open(gm)
			if errOpen != nil {
				return "", errOpen
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if strings.HasPrefix(line, "module ") {
					campos := strings.Fields(line)
					if len(campos) >= 2 {
						return campos[1], nil
					}
				}
			}
			if errScan := scanner.Err(); errScan != nil {
				return "", errScan
			}
		}
		father := filepath.Dir(dir)
		if father == dir {
			break
		}
		dir = father
	}
	return "", nil
}
