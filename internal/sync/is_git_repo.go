package sync

import (
	"os"
	"path/filepath"
)

func isGitRepo(path string) (bool, error) {
	stat, err := os.Stat(filepath.Join(path, ".git"))
	if err != nil {
		return false, err
	}
	if !stat.IsDir() {
		return false, nil
	}
	return true, nil
}
