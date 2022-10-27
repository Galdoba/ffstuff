package filelist

import (
	"fmt"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/scanner"
)

type FileList struct {
	paths     []string
	exeptions []string
	er        error
}

func New(root string) FileList {
	l, err := scanner.Scan(root, "")
	fl := FileList{paths: l, er: err}
	return fl
}

func (fl *FileList) Print() {
	for _, path := range fl.paths {
		print := true
		for _, ex := range fl.exeptions {
			if !print {
				break
			}
			if strings.Contains(path, ex) {
				print = false
				continue
			}
		}
		if print {
			fmt.Printf("%v\n", path)
		}
	}
}

func (fl *FileList) AddExeption(ex string) {
	fl.exeptions = append(fl.exeptions, ex)
}
