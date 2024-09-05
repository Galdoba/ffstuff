package logman

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

var logMan *logManager
var flags int = os.O_CREATE | os.O_WRONLY | os.O_APPEND
var perm fs.FileMode = 0666

const (
	ImportanceNONE  = 100
	ImportanceFATAL = 99
	ImportanceERROR = 80
	ImportanceWARN  = 70
	ImportanceINFO  = 50
	ImportanceDEBUG = 30
	ImportanceTRACE = 10
	ImportanceALL   = 0

	//fieldKeys
	keyTime    = "time"
	keyLevel   = "level"
	keyMessage = "message"
	keyFile    = "file"
	keyLine    = "line"
	keyFunc    = "callerFuncName"

	Stdout = "StdOut"
	Stderr = "StdErr"
)

type logManager struct {
	appMinimumLoglevel int
	logLevels          map[string]*loggingLevel
	logger             *log.Logger
}

// Setup sets logging levels for logMan. If no
func Setup(opts ...LogmanOptions) error {
	al := logManager{}
	al.appMinimumLoglevel = ImportanceALL
	al.logLevels = defaultLoggingLevels()
	opt := defaultOpts()
	for _, set := range opts {
		set(&opt)
	}
	for _, lvl := range opt.logLevels {
		al.logLevels[lvl.name] = lvl
	}

	al.appMinimumLoglevel = opt.appMinimumLoglevel
	al.logger = log.New(os.Stdout, "", 0)
	logMan = &al
	return nil
}

// SetOutput ADD new and REPLACE existing writers for loggingLevels.
// writerKey cases:
// Stdout - sets os.Stdout as writer (default for LvlFATAL, LvlERROR and LvlWARN )
// Stderr - sets os.Stderr as writer (default for LvlINFO , LvlDEBUG and LvlTRACE)
// If writerKey is a filepath to new or existing file, that file will be used as writer.
// TODO: If writerKey is a filepath to directory each message will be written to own file for this directory.
func SetOutput(writerKey string, levels ...string) {
	for _, levelName := range levels {
		for lvlName, lvlOpts := range logMan.logLevels {
			if levelName != lvlName && levelName != ALL {
				continue
			}
			lvlOpts.setWriter(writerKey)
		}
	}
}

// ClearOutput is used to REMOVE ALL writers for given level.
// Used removing default writers.
func ClearOutput(levels ...string) {
	for _, levelName := range levels {
		for lvlName, lvlOpts := range logMan.logLevels {
			if levelName != lvlName && levelName != ALL {
				continue
			}
			lvlOpts.clearWriters()
		}
	}
}

// LogmanOptions - settings for logMan object.
type LogmanOptions func(*options)

type options struct {
	appMinimumLoglevel int
	logLevels          map[string]*loggingLevel
}

func defaultOpts() options {
	return options{
		appMinimumLoglevel: ImportanceALL,
		logLevels:          defaultLoggingLevels(),
	}

}

// WithLogLevels sets loglevels to logman with slice of NewLogLevel functions.
// Used to create custom logLevels.
// Caution: It overrides default levels if new loglevel has standard key ("fatal", "error", "warn", "info", "debug", "trace").
func WithLogLevels(lvls ...*loggingLevel) LogmanOptions {
	return func(o *options) {
		o.logLevels = make(map[string]*loggingLevel)
		for _, lvl := range lvls {
			o.logLevels[lvl.name] = lvl
		}
	}
}

// WithAppLogLevelImportance sets minimum message importance level logMan will process.
// If input is below ImportanceNone importance will be set to ImportanceNone.
// If input is above ImportanceALL importance will be set to ImportanceALL.
func WithAppLogLevelImportance(importance int) LogmanOptions {
	return func(o *options) {
		if importance < ImportanceNONE {
			importance = ImportanceNONE
		}
		if importance > ImportanceALL {
			importance = ImportanceALL
		}
		o.appMinimumLoglevel = importance
	}
}

func isInBounds(n, min, max int) bool {
	if n < min || n > max {
		return false
	}
	return true
}

// ProcessMessage is a general call for processing message.
// Must be used if custom log levels are used.
func ProcessMessage(msg Message, levels ...string) error {
	loggingLevels := []*loggingLevel{}
	for _, level := range levels {
		loggingLevels = append(loggingLevels, logMan.logLevels[level])
	}
	return process(msg, loggingLevels...)
}

// This is main func for processing messages on levels provided.
// It return processing error of nil if processing successful.
func process(msg Message, lvls ...*loggingLevel) error {
	errorStack := []error{}
	fatalCalled := false
	for _, lvl := range lvls {
		if lvl == nil {
			errorStack = append(errorStack, fmt.Errorf("logginglevel provided was not set"))
			continue
		}
		if logMan.appMinimumLoglevel >= lvl.importance {
			continue
		}
		if !isPresent(lvl) {
			errorStack = append(errorStack, fmt.Errorf("level %v was not set properly", lvl.name))
			continue
		}

		for _, present := range logMan.logLevels {
			if lvl.name != present.name {
				continue
			}
			if lvl.formatFunc == nil {
				errorStack = append(errorStack, fmt.Errorf("level %v have no formatFunc", lvl.name))
				continue
			}

			msg.SetField(keyLevel, lvl.tag)
			if lvl.callerInfo {
				file, line, fn := callerFunctionInfo(3)
				if msg.Value(keyFile) == nil {
					msg.SetField(keyFile, file)
				}
				if msg.Value(keyLine) == nil {
					msg.SetField(keyLine, line)
				}
				if msg.Value(keyFunc) == nil {
					msg.SetField(keyFunc, fn)
				}
			}
			text, err := lvl.formatFunc(msg)
			if err != nil {
				errorStack = append(errorStack, fmt.Errorf("formatting message failed: '%v' level: %v", lvl.name, err))
				continue
			}
			if err = lvl.write(text); err != nil {
				errorStack = append(errorStack, fmt.Errorf("writting message failed: %v", err))
				continue
			}
			if lvl.osExit {
				fatalCalled = true
			}
		}
	}
	if err := joinErrors("processing message failed", errorStack...); err != nil {
		return err
	}
	if fatalCalled {
		os.Exit(1)
	}
	return nil
}

func isPresent(lvl *loggingLevel) bool {
	for _, present := range logMan.logLevels {
		if lvl.name == present.name && lvl.tag == present.tag {
			return true
		}
	}
	return false
}

func (lvl *loggingLevel) write(text string) error {
	text = strings.TrimSuffix(text, "\n") + "\n"
	errorStack := []error{}
	for key, writer := range lvl.writers {
		switch {
		case key == Stderr, key == Stdout:
			textColorized := colorizeTags(text)
			bt := []byte(textColorized)
			if _, err := writer.Write(bt); err != nil {
				errorStack = append(errorStack, fmt.Errorf("'%v' level: writer %v failed: %v", lvl.name, key, err))
			}
		default:
			f, err := os.OpenFile(key, flags, perm)
			if err != nil {
				errorStack = append(errorStack, fmt.Errorf("'%v' level: open file failed: %v", lvl.name, err))
				continue
			}
			if _, err := f.WriteString(text); err != nil {
				errorStack = append(errorStack, fmt.Errorf("'%v' level: write to file failed: %v", lvl.name, err))
				continue
			}
			f.Close()
		}

	}
	if err := joinErrors("writing message failed", errorStack...); err != nil {
		return err
	}
	return nil
}

func colorizeTags(text string) string {
	//TODO: temporary
	//need to make formatter struct which will colorize text
	text = strings.ReplaceAll(text, stdTagFATAL, color.RedString(stdTagFATAL))
	text = strings.ReplaceAll(text, stdTagERROR, color.HiRedString(stdTagERROR))
	text = strings.ReplaceAll(text, stdTagWARN, color.YellowString(stdTagWARN))
	text = strings.ReplaceAll(text, stdTagDEBUG, color.CyanString(stdTagDEBUG))
	text = strings.ReplaceAll(text, stdTagTRACE, color.HiBlackString(stdTagTRACE))
	return text
}

func joinErrors(message string, errs ...error) error {
	if len(errs) == 0 {
		return nil
	}
	errCombined := fmt.Errorf("%v:", message)
	for _, err := range errs {
		errCombined = fmt.Errorf("%v\n%v", errCombined, err)
	}
	return errCombined
}

////////////////////////////////

func callerFunctionInfo(n int) (string, int, string) {
	counter, file, line, success := runtime.Caller(n) //back to stack on n levels
	if !success {
		return "", 0, ""
	}
	funcName := runtime.FuncForPC(counter).Name()
	return file, line, funcName
}

func timestampFormat(tm time.Time) string {
	format := "2006/01/02 15:04:05.999"
	formatLen := len(format)
	stamp := tm.Format(format)
	for len(stamp) < formatLen {
		stamp += "0"
	}
	return stamp
}
