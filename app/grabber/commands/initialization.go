package commands

import (
	"fmt"
	"strings"

	"github.com/Galdoba/ffstuff/app/grabber/config"
	"github.com/Galdoba/ffstuff/pkg/logman"
)

func commandInit() error {
	cfgLoaded, err := config.Load()
	if err != nil {
		return err
	}
	errs := config.Validate(cfgLoaded)
	switch len(errs) {
	default:
		for _, err := range errs {
			if err != nil {
				fmt.Println(err)
			}
		}
		return fmt.Errorf("config errors detected")
	case 0:
		cfg = cfgLoaded
		if err := setupLogger(cfg.LOG_LEVEL, cfg.LOG); err != nil {
			return fmt.Errorf("logger initialization failed: %v", err)
		}
		// if err := setupSourceConstructor(); err != nil {
		// 	return fmt.Errorf("source constructor initialization failed: %v", err)
		// }

	}
	return nil
}

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
