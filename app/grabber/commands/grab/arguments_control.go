package grab

import (
	"errors"
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/pkg/logman"
)

func ValidateArgs(args ...string) ([]string, []error) {
	sourcePaths := []string{}
	errs := []error{}
	if len(args) == 0 {
		logman.Warn(("no arguments provided: grab command expects arguments"))
	}
	for _, arg := range args {
		file, errOpen := os.OpenFile(arg, os.O_RDONLY, 0666)
		if errOpen != nil {
			if errors.Is(errOpen, os.ErrNotExist) {
				logman.Warn("bad argument: %v is not exist", arg)
				errs = append(errs, fmt.Errorf("bad argument: %v is not exist", arg))
				continue
			}
			errs = append(errs, fmt.Errorf("bad argument: %v: %v", arg, errOpen))
			continue
		}

		fs, errStats := file.Stat()
		if errStats != nil {
			logman.Warn("bad argument: %v: %v", arg, errStats)
			errs = append(errs, fmt.Errorf("bad argument: %v: %v", arg, errStats))
			continue
		}
		if fs.IsDir() {
			logman.Warn("bad argument: %v is directory: grab command don't do directories", arg)
			continue
		}

		sourcePaths = append(sourcePaths, arg)
	}
	if len(sourcePaths) == 0 {
		errs = append(errs, fmt.Errorf("no valid source paths detected"))
	}
	return sourcePaths, errs
}
