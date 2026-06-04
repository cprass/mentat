package core

import (
	"fmt"
	"os"
)

// Check whether the given directory exists.
// If the given directory doesn't exist, create it (using filemode 600)
func EnsureDir(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(dir, 0600)
			if err != nil {
				return fmt.Errorf("unable to create directory: %w", err)
			}
			return nil
		}
		return fmt.Errorf("can't read directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("given path is not a directory: %s", dir)
	}

	return nil
}
