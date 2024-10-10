package validation

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func FileValidation(path string) error {
	if path == "" {
		return fmt.Errorf("no filepath provided")
	}
	switch filepath.IsAbs(path) {
	case true:
		fi, err := os.Stat(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("path not exists")
			}
			return fmt.Errorf("path exists, but '%v'", err)
		}
		if fi.IsDir() {
			return fmt.Errorf("path exists, but is directory")
		}
		f, err := os.OpenFile(path, os.O_RDWR, 0666)
		defer f.Close()
		if err != nil {
			return fmt.Errorf("file can't be opened: %v", err)
		}

		return nil
	default:
		return fmt.Errorf("absolute path expected")
	}
}

func DirectoryValidation(path string) error {
	if path == "" {
		return fmt.Errorf("directory is not set")
	}
	switch filepath.IsAbs(path) {
	case true:
		fi, err := os.Stat(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("path not exists")
			}
			return fmt.Errorf("path exists, but '%v'", err)
		}
		if !fi.IsDir() {
			return fmt.Errorf("path exists, but is not directory")
		}
		return nil
	default:
		return fmt.Errorf("absolute path expected")
	}
}

func Exists(filepath string) (bool, error) {
	if _, err := os.Stat(filepath); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}
