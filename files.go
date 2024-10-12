package main

import (
	"os"
	"path/filepath"
)

func GetComposeFilePaths(root string) ([]string, error) {
	var composeFilePaths []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		matched, err := filepath.Match("docker-compose.y*ml", filepath.Base(path))
		if err != nil {
			return err
		}
		if !matched {
			matched, err = filepath.Match("compose.y*ml", filepath.Base(path))
			if err != nil {
				return err
			}
		}
		if matched {
			composeFilePaths = append(composeFilePaths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return composeFilePaths, nil
}
