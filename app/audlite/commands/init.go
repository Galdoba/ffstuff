package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Galdoba/ffstuff/app/audlite/config"
	"github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/Galdoba/ffstuff/pkg/logman/colorizer"
	"github.com/urfave/cli/v2"
)

var cfg *config.Configuration

func commandInit(c *cli.Context) error {
	cfgLoaded, err := config.Load()
	if err != nil {
		switch errors.Is(err, config.ErrNoConfig) {
		case true:
			cfgLoaded = config.NewConfig(c.App.Version)
			for _, log := range cfgLoaded.LOGS {
				os.MkdirAll(filepath.Dir(log), 0666)
				os.Create(log)
			}
		case false:
			return err
		}
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

		if err := setupLogger(cfg); err != nil {
			return fmt.Errorf("logger initialization failed: %v", err)
		}
		if err := config.Save(cfg); err != nil {

			return fmt.Errorf("initialization failed: %v", err)
		}

	}
	return nil
}

var stdLogLevels = []string{logman.TRACE, logman.DEBUG, logman.INFO, logman.WARN, logman.ERROR, logman.FATAL}

func setupLogger(cfg *config.Configuration) error {
	fmt.Println("start logger setup")
	consoleLevels, fileLevels := filterLevels(cfg.CONSOLE_LOG_LEVEL, cfg.FILE_LOG_LEVEL)
	consoleFormatter := logman.NewFormatter(logman.WithRequestedFields(logman.Request_ShortReport), logman.WithColor(colorizer.DefaultScheme()))
	fileReportFormatter := logman.NewFormatter(logman.WithRequestedFields(logman.Request_ShortTime), logman.WithColor(nil))
	fileErrorFormatter := logman.NewFormatter(logman.WithRequestedFields(logman.Request_Full), logman.WithColor(nil))
	if err := logman.Setup(
		logman.WithAppLogLevelImportance(logman.ImportanceALL),
		logman.WithGlobalColorizer(colorizer.DefaultScheme()),
	); err != nil {
		return err
	}
	if err := logman.ResetWriters(stdLogLevels...); err != nil {
		return err
	}

	for _, level := range consoleLevels {
		if err := logman.SetLevelWriterFormatter(level, logman.Stderr, consoleFormatter); err != nil {
			return fmt.Errorf("failed to set formatter to level's '%v' writer (%v): %v", level, logman.Stderr, consoleFormatter)
		}
	}
	for _, log := range cfg.LOGS {
		for _, level := range fileLevels {
			switch level {
			case logman.ERROR, logman.FATAL:
				if err := logman.SetLevelWriterFormatter(level, log, fileErrorFormatter); err != nil {
					return fmt.Errorf("failed to set formatter to level's '%v' writer (%v): %v", level, log, fileErrorFormatter)
				}
			default:
				if err := logman.SetLevelWriterFormatter(level, log, fileReportFormatter); err != nil {
					return fmt.Errorf("failed to set formatter to level's '%v' writer (%v): %v", level, log, fileReportFormatter)
				}
			}
		}
	}

	fmt.Println("end logger setup")
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
