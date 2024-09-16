package actions

import (
	"fmt"
	"os"
	"time"

	"github.com/Galdoba/ffstuff/pkg/logman"
)

func ScanDir(dir string) ([]string, []string, error) {
	files := []string{}
	dirs := []string{}
	for i := 1; i <= 10; i++ {
		files = []string{}
		dirs = []string{}
		fi, err := os.ReadDir(dir)
		if err != nil {
			logman.Warn(fmt.Sprintf("read directory '%v' (atempt %v): failed: %v", dir, i, err.Error()))
			time.Sleep(time.Second * 6)
			switch i {
			default:
				continue
			case 10:
				return files, dirs, fmt.Errorf("failed to scan directory %v: max retry reached", dir)
			}
		}
		for _, f := range fi {
			switch f.IsDir() {
			case true:
				dirs = append(dirs, f.Name())
			case false:
				files = append(files, f.Name())
			}
		}
		break
	}
	return files, dirs, nil
}
