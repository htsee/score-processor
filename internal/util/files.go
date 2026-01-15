package util

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
)

func CheckFileType(file, filetype string) error {
	if path.Ext(file) != "."+filetype {
		return fmt.Errorf("%q is not a %v", file, filetype)
	}

	if _, err := os.Stat(file); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("File %q does not exist", file)
		}
		return fmt.Errorf("Cannot access file %q: %w", file, err)
	}

	return nil
}
