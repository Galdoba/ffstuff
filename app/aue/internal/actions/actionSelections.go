package actions

import (
	"fmt"
	"os"

	"github.com/Galdoba/ffstuff/pkg/logman"
)

const (
	ActionSkip          = "Skip"
	ActionSetup         = "Setup Source Files"
	ActionRemoveProject = "Remove Project"
)

func SelectSourceAction(sourceDir string) (string, error) {
	fi, err := os.ReadDir(sourceDir)
	if err != nil {
		return ActionSkip, fmt.Errorf("failed to read project directory: %v", err)
	}
	filesFound := []string{}
	for _, f := range fi {
		filesFound = append(filesFound, f.Name())
	}
	if len(filesFound) == 0 {
		logman.Warn("%v is empty", sourceDir)
		return ActionRemoveProject, nil
	}
	if len(filesFound) == 1 && filesFound[0] == "metadata.json" {
		logman.Info("%v is subject to remove", sourceDir)
		return ActionRemoveProject, nil
	}
	for _, file := range filesFound {
		if file == "lock" {
			logman.Debug(logman.NewMessage("%v is locked", sourceDir))
			return ActionSkip, nil
		}
	}
	return ActionSetup, nil
}
