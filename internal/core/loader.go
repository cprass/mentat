package core

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func LoadFiles(root string, ext string) ([]string, error) {
	var files []string

	if root == "" {
		return nil, fmt.Errorf("root must be a path")
	}

	if ext == "" {
		return nil, fmt.Errorf("ext must be set")
	}

	e := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories - WalkDir handles them automatically
		if d.IsDir() {
			return nil
		}

		if ext == strings.ToLower(filepath.Ext(path)) {
			files = append(files, path)
		}

		return nil
	})
	if e != nil {
		return nil, e
	}

	return files, nil
}
