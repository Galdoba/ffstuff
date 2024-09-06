package actions

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/pkg/logman"
)

func RemoveProject(directory string) error {
	err := os.RemoveAll(directory)
	if err != nil {
		err = fmt.Errorf("failed to remove directory '%v': %v", directory, err)
		logman.Error(err)
		return err
	}
	logman.Debug(logman.NewMessage("directory %v removed with all files", directory))
	return nil
}
