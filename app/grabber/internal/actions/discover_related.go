package actions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func DiscoverRelatedFiles(markerFile string) ([]string, error) {
	dir := filepath.Dir(markerFile)
	ext := filepath.Ext(markerFile)

	baseMarker := filepath.Base(markerFile)
	fileBase := strings.TrimSuffix(baseMarker, ext)
	out := []string{}
	sep := string(filepath.Separator)
	fi, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", dir)
	}
	for _, f := range fi {
		if f.IsDir() {
			continue
		}
		if !strings.HasPrefix(f.Name(), fileBase) {
			continue
		}
		path := dir + sep + f.Name()
		if ext == filepath.Ext(path) {
			continue
		}
		out = append(out, path)
	}
	return out, nil
}
