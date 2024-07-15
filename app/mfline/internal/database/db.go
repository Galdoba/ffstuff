package database

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func (db *DBjson) List() []string {
	dirEntries, err := os.ReadDir(db.dir)
	if err != nil {
		fmt.Printf("DB dir: %v\n", err)
	}
	list := []string{}
	for _, fi := range dirEntries {
		if !strings.HasSuffix(fi.Name(), ".json") {
			continue
		}
		list = append(list, fi.Name())
	}
	return list
}

type SearchResult struct {
	Err    error
	Result string
	Key    string
}

var ErrNotFound = errors.New("entry not found")

func (db *DBjson) FindEntry(source string) SearchResult {
	sourceName := filepath.Base(source)
	dirEntries, err := os.ReadDir(db.dir)
	if err != nil {
		return SearchResult{Err: fmt.Errorf("read dir: %v", err)}
	}
	for _, entry := range dirEntries {
		fullname := entry.Name()
		if !strings.HasSuffix(fullname, ".json") {
			continue
		}
		parts := strings.Split(sourceName, "--")
		body := parts[len(parts)-1]
		if !strings.Contains(fullname, body) {
			continue
		}
		key := strings.TrimSuffix(body, ".json")
		return SearchResult{Result: db.dir + entry.Name(), Key: key}
	}
	return SearchResult{Err: ErrNotFound}
}
