package util

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func CheckFileType(file, filetype string) error {
	if path.Ext(file) != "."+filetype {
		return fmt.Errorf("file %q is not a %v", file, filetype)
	}

	if _, err := os.Stat(file); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("file %q does not exist", file)
		}
		return fmt.Errorf("cannot access file %q: %w", file, err)
	}

	return nil
}

func FileList(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var fileList []string
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		fileList = append(fileList, path)
	}
	return fileList, nil
}
