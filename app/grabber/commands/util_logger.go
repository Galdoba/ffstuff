package commands

import (
	"fmt"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/logman"
)

func setupLogger(level, logpath string) error {
	switch level {
	case "":
		return fmt.Errorf("config field 'Minimum Log Level' contains no data")
	default:
		return fmt.Errorf("config field 'Minimum Log Level' contains invalid data: '%v'\n"+
			"expecting: TRACE, DEBUG, INFO, WARN, ERROR или FATAL", level)
	case strings.ToUpper(logman.TRACE):
		if err := logman.Setup(logman.WithAppLogLevelImportance(logman.ImportanceTRACE)); err != nil {
			return fmt.Errorf("logger setup failed: %v", err)
		}
	case strings.ToUpper(logman.DEBUG):
		if err := logman.Setup(logman.WithAppLogLevelImportance(logman.ImportanceDEBUG)); err != nil {
			return fmt.Errorf("logger setup failed: %v", err)
		}
	case strings.ToUpper(logman.INFO):
		if err := logman.Setup(logman.WithAppLogLevelImportance(logman.ImportanceINFO)); err != nil {
			return fmt.Errorf("logger setup failed: %v", err)
		}
	case strings.ToUpper(logman.WARN):
		if err := logman.Setup(logman.WithAppLogLevelImportance(logman.ImportanceWARN)); err != nil {
			return fmt.Errorf("logger setup failed: %v", err)
		}
	case strings.ToUpper(logman.ERROR):
		if err := logman.Setup(logman.WithAppLogLevelImportance(logman.ImportanceERROR)); err != nil {
			return fmt.Errorf("logger setup failed: %v", err)
		}
	case strings.ToUpper(logman.FATAL):
		if err := logman.Setup(logman.WithAppLogLevelImportance(logman.ImportanceFATAL)); err != nil {
			return fmt.Errorf("logger setup failed: %v", err)
		}
	}
	logman.ClearOutput(logman.ALL)
	logman.SetOutput(logpath, logman.ALL)
	logman.SetOutput(logman.Stderr, logman.ALL)
	return nil
}
