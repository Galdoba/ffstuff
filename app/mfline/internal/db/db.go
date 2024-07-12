package db

import (
	"fmt"
	"os"
)

var DB_PATH string

type DBjson struct {
	dir string
}

func New(dir string) (*DBjson, error) {
	if dir == "" {
		return nil, fmt.Errorf("directory must be provided")
	}

	fi, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("Db dir stat: %v", err)
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("path provided is not dir: %v", dir)
	}
	dBase := DBjson{}
	dBase.dir = dir
	return &dBase, nil
}

func (db *DBjson) Dir() string {
	return db.dir
}

func (db *DBjson) Validate() error {
	if db.dir == "" {
		return fmt.Errorf("Db not set")
	}
	return nil
}

func validateDir(dir string) error {
	fi, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("dir stat: %v", err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("not a dir: %v", dir)
	}
	return nil
}

func validateFile(file string) error {
	fi, err := os.Stat(file)
	if err != nil {
		return fmt.Errorf("file stat: %v", err)
	}
	if fi.IsDir() {
		return fmt.Errorf("not a file: %v", file)
	}
	return nil
}
