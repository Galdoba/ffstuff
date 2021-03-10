package scan

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

//ScanRoot - walks through all paths under the root and returns evaluation results
func ScanRoot(root string) []string {
	var result []string
	if err := filepath.Walk(root, evaluate); err != nil {
		//logger.ERROR(err.Error())
		result = append(result, err.Error())
	}

	return result
}

//evaluate - возвращает ошибку с путем до найденного ready файла
func evaluate(path string, f os.FileInfo, err error) error {
	if f.IsDir() {
		fmt.Println("Scanning: ", path)
		return nil
	}
	if filepath.Ext(path) != "ready" {
		return nil
	}
	return errors.New(path)
}
