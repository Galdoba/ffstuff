package commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Galdoba/ffstuff/app/grabber/config"
	"github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/Galdoba/ffstuff/pkg/logman/colorizer"
	"github.com/Galdoba/ffstuff/pkg/stdpath"
	"github.com/urfave/cli/v2"
)

func commandInit(c *cli.Context) error {
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
		if err := setupLogger(c, cfg); err != nil {
			return fmt.Errorf("logger initialization failed: %v", err)
		}
		// if err := setupSourceConstructor(); err != nil {
		// 	return fmt.Errorf("source constructor initialization failed: %v", err)
		// }

	}
	return nil
}

var stdLogLevels = []string{logman.TRACE, logman.DEBUG, logman.INFO, logman.WARN, logman.ERROR, logman.FATAL}

func setupLogger(c *cli.Context, cfg *config.Configuration) error {
	consoleLevels, fileLevels := filterLevels(cfg.CONSOLE_LOG_LEVEL, cfg.FILE_LOG_LEVEL)
	jsonLogDir := stdpath.LogDir() + "structured" + string(filepath.Separator)
	//consoleFormatter := logman.NewFormatter(logman.WithRequestedFields(logman.Request_ShortReport), logman.WithColor(colorizer.DefaultScheme()))
	// fileReportFormatter := logman.NewFormatter(logman.WithRequestedFields(logman.Request_ShortTime), logman.WithColor(nil))
	// fileErrorFormatter := logman.NewFormatter(logman.WithRequestedFields(logman.Request_Full), logman.WithColor(nil))
	// if err := logman.Setup(
	// 	logman.WithAppLogLevelImportance(logman.ImportanceALL),
	// 	logman.WithGlobalColorizer(colorizer.DefaultScheme()),
	// 	logman.WithAppName(c.App.Name),
	// 	logman.WithJSON(stdpath.LogDir()),
	// ); err != nil {
	// 	return err
	// }
	// if err := logman.ResetWriters(stdLogLevels...); err != nil {
	// 	return err
	// }

	// for _, level := range consoleLevels {
	// 	if err := logman.SetLevelWriterFormatter(level, logman.Stderr, consoleFormatter); err != nil {
	// 		return fmt.Errorf("failed to set formatter to level's '%v' writer (%v): %v", level, logman.Stderr, consoleFormatter)
	// 	}
	// }
	// for _, level := range fileLevels {
	// 	switch level {
	// 	case logman.ERROR, logman.FATAL:
	// 		if err := logman.SetLevelWriterFormatter(level, cfg.LOG, fileErrorFormatter); err != nil {
	// 			return fmt.Errorf("failed to set formatter to level's '%v' writer (%v): %v", level, cfg.LOG, fileErrorFormatter)
	// 		}
	// 	default:
	// 		if err := logman.SetLevelWriterFormatter(level, cfg.LOG, fileReportFormatter); err != nil {
	// 			return fmt.Errorf("failed to set formatter to level's '%v' writer (%v): %v", level, cfg.LOG, fileReportFormatter)
	// 		}
	// 	}
	// }
	logman.Setup(
		logman.WithAppName(c.App.Name),
		logman.WithAppLogLevelImportance(logman.ImportanceALL),
		logman.WithGlobalColorizer(colorizer.DefaultScheme()),
		logman.WithGlobalWriterFormatter(stdpath.LogFile(), logman.NewFormatter(logman.WithRequestedFields(logman.Request_ShortTime))),
		logman.WithJSON(jsonLogDir),
	)
	for _, cLevel := range stdLogLevels {
		if notListed(consoleLevels, cLevel) {
			if err := logman.RemovetWriter(cLevel, logman.Stderr); err != nil {
				return err
			}
		}
	}
	for _, fLevel := range stdLogLevels {
		if notListed(fileLevels, fLevel) {
			fmt.Println("delete flevel", fLevel)
			if err := logman.RemovetWriter(fLevel, stdpath.LogFile()); err != nil {
				return err
			}
			if err := logman.RemovetWriter(fLevel, jsonLogDir); err != nil {
				return err
			}
		}
	}

	return nil
}

func filterLevels(console, file string) ([]string, []string) {
	console = strings.ToLower(console)
	file = strings.ToLower(file)
	cMet := false
	fMet := false
	consoleLevels := []string{}
	fileLevels := []string{}
	for _, val := range stdLogLevels {
		if val == console {
			cMet = true
		}
		if val == file {
			fMet = true
		}
		if cMet {
			consoleLevels = append(consoleLevels, val)
		}
		if fMet {
			fileLevels = append(fileLevels, val)
		}
	}
	return consoleLevels, fileLevels
}

func notListed(levels []string, lvl string) bool {
	for _, val := range levels {
		if val == lvl {
			return false
		}
	}
	return true
}
