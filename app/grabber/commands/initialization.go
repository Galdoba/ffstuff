package commands

import (
	"fmt"

	"github.com/Galdoba/ffstuff/app/grabber/config"
	"github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/Galdoba/ffstuff/pkg/logman/colorizer"
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
		if err := setupLogger(cfg.CONSOLE_LOG_LEVEL, cfg.FILE_LOG_LEVEL, cfg.LOG); err != nil {
			return fmt.Errorf("logger initialization failed: %v", err)
		}
		// if err := setupSourceConstructor(); err != nil {
		// 	return fmt.Errorf("source constructor initialization failed: %v", err)
		// }

	}
	return nil
}

// func setupLogger(levelConsole, levelFile, logpath string) error {
// 	levelConsole = strings.ToUpper(levelConsole)
// 	levelFile = strings.ToUpper(levelFile)
// 	keepToConsole := []string{}
// 	switch levelConsole {
// 	case "":
// 		return fmt.Errorf("config field 'Minimum Log Level: Terminal' contains no data")
// 	default:
// 		return fmt.Errorf("config field 'Minimum Log Level: Terminal' contains invalid data: '%v' "+
// 			"(expecting: TRACE, DEBUG, INFO, WARN, ERROR or FATAL)", levelConsole)
// 	case strings.ToUpper(logman.TRACE):
// 		keepToConsole = []string{logman.TRACE, logman.DEBUG, logman.INFO, logman.WARN, logman.ERROR, logman.FATAL}
// 	case strings.ToUpper(logman.DEBUG):
// 		keepToConsole = []string{logman.DEBUG, logman.INFO, logman.WARN, logman.ERROR, logman.FATAL}
// 	case strings.ToUpper(logman.INFO):
// 		keepToConsole = []string{logman.INFO, logman.WARN, logman.ERROR, logman.FATAL}
// 	case strings.ToUpper(logman.WARN):
// 		keepToConsole = []string{logman.WARN, logman.ERROR, logman.FATAL}
// 	case strings.ToUpper(logman.ERROR):
// 		keepToConsole = []string{logman.ERROR, logman.FATAL}
// 	case strings.ToUpper(logman.FATAL):
// 		keepToConsole = []string{logman.FATAL}
// 	}
// 	switch levelFile {
// 	case "":
// 		return fmt.Errorf("config field 'Minimum Log Level: File' contains no data")
// 	default:
// 		return fmt.Errorf("config field 'Minimum Log Level: File' contains invalid data: '%v' "+
// 			"(expecting: TRACE, DEBUG, INFO, WARN, ERROR or FATAL)", levelConsole)
// 	case strings.ToUpper(logman.TRACE):
// 		if err := logman.Setup(logman.WithAppLogLevelImportance(logman.ImportanceTRACE)); err != nil {
// 			return fmt.Errorf("logger setup failed: %v", err)
// 		}
// 	case strings.ToUpper(logman.DEBUG):
// 		if err := logman.Setup(logman.WithAppLogLevelImportance(logman.ImportanceDEBUG)); err != nil {
// 			return fmt.Errorf("logger setup failed: %v", err)
// 		}
// 	case strings.ToUpper(logman.INFO):
// 		if err := logman.Setup(logman.WithAppLogLevelImportance(logman.ImportanceINFO)); err != nil {
// 			return fmt.Errorf("logger setup failed: %v", err)
// 		}
// 	case strings.ToUpper(logman.WARN):
// 		if err := logman.Setup(logman.WithAppLogLevelImportance(logman.ImportanceWARN)); err != nil {
// 			return fmt.Errorf("logger setup failed: %v", err)
// 		}
// 	case strings.ToUpper(logman.ERROR):
// 		if err := logman.Setup(logman.WithAppLogLevelImportance(logman.ImportanceERROR)); err != nil {
// 			return fmt.Errorf("logger setup failed: %v", err)
// 		}
// 	case strings.ToUpper(logman.FATAL):
// 		if err := logman.Setup(logman.WithAppLogLevelImportance(logman.ImportanceFATAL)); err != nil {
// 			return fmt.Errorf("logger setup failed: %v", err)
// 		}
// 	}
// 	logman.ClearOutput(logman.ALL)
// 	logman.SetOutput(logpath, logman.ALL)
// 	keepOutput := []string{}
// 	keep := false
// 	for _, level := range keepToConsole {
// 		if strings.ToUpper(level) == levelConsole {
// 			keep = true
// 		}
// 		if keep {
// 			keepOutput = append(keepOutput, level)
// 		}
// 	}
// 	if len(keepOutput) == 0 {
// 		keepOutput = []string{logman.FATAL}
// 	}
// 	logman.SetOutput(logman.Stderr, keepOutput...)
// 	return nil
// }

func setupLogger(levelConsole, levelFile, logpath string) error {
	withColor := logman.WithColor(colorizer.DefaultScheme())
	noColor := logman.WithColor(nil)
	withAllFields := logman.WithRequestedFields(logman.Request_Full)
	withErrorFields := logman.WithRequestedFields(logman.Request_Medium)
	withLogFields := logman.WithRequestedFields(logman.Request_ShortTime)
	withReportFields := logman.WithRequestedFields(logman.Request_ShortTime)
	return logman.Setup(
		logman.WithAppLogLevelImportance(logman.ImportanceALL),
		logman.WithLogLevels(
			logman.NewLoggingLevel(logman.FATAL,
				logman.WithWriter(logman.Stdout, logman.NewFormatter(withAllFields, withColor)),
				logman.WithWriter(logpath, logman.NewFormatter(withAllFields, noColor)),
			),
			logman.NewLoggingLevel(logman.ERROR,
				logman.WithWriter(logman.Stdout, logman.NewFormatter(withReportFields, withColor)),
				logman.WithWriter(logpath, logman.NewFormatter(withErrorFields, noColor)),
			),
			logman.NewLoggingLevel(logman.WARN,
				logman.WithWriter(logman.Stderr, logman.NewFormatter(withReportFields, withColor)),
				logman.WithWriter(logpath, logman.NewFormatter(withLogFields, noColor)),
			),
			logman.NewLoggingLevel(logman.INFO,
				logman.WithWriter(logman.Stderr, logman.NewFormatter(withReportFields, withColor)),
				logman.WithWriter(logpath, logman.NewFormatter(withLogFields, noColor)),
			),
			logman.NewLoggingLevel(logman.DEBUG,
				logman.WithWriter(logman.Stderr, logman.NewFormatter(withReportFields, withColor)),
				logman.WithWriter(logpath, logman.NewFormatter(withLogFields, noColor)),
			),
			logman.NewLoggingLevel(logman.TRACE,
				logman.WithWriter(logman.Stderr, logman.NewFormatter(withReportFields, withColor)),
				logman.WithWriter(logpath, logman.NewFormatter(withAllFields, noColor)),
			),
		),
		logman.WithGlobalColorizer(colorizer.DefaultScheme()),
		logman.WithGlobalWriterFormatter(logpath,
			logman.NewFormatter(logman.WithRequestedFields(logman.Request_ShortTime)),
		),
	)
}
