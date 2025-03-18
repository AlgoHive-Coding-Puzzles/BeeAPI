package services

import (
	"os"
	"path/filepath"
)

// GetDirSize calculates the total size of a directory in bytes
func GetDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// RemoveAll removes a directory and all its contents
func RemoveAll(path string) error {
	return os.RemoveAll(path)
}