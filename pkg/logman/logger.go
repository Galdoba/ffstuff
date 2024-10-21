package logman

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Galdoba/ffstuff/pkg/logman/colorizer"
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
	ImportancePING  = 1
	ImportanceALL   = 0

	//fieldKeys
	keyTime        = "time"
	keySince       = "since"
	keyLevel       = "level"
	keyMessage     = "message"
	keyFile        = "file"
	keyLine        = "line"
	keyFunc        = "callerFuncName"
	keyCaller      = "caller"
	keyCallerShort = "caller_short"
	keyCallerLong  = "caller_long"

	Stdout = "StdOut"
	Stderr = "StdErr"
)

type logManager struct {
	appMinimumLoglevel int
	logLevels          map[string]*loggingLevel
	longCallerNames    bool
	logger             *log.Logger
	colorizer          Colorizer
	startTime          time.Time
	//activeWriter       string
}

// Colorizer - uses Color Schema to make console output colored depending on fariable type
type Colorizer interface {
	ColorizeByType(interface{}) string
	ColorizeByKeys(interface{}, ...colorizer.ColorKey) string
}

// Setup sets logMan options. Place it at start of the program.
func Setup(opts ...LogmanOptions) error {
	al := logManager{}
	al.startTime = time.Now()
	al.appMinimumLoglevel = ImportanceALL
	al.logLevels = make(map[string]*loggingLevel)
	//al.logLevels = defaultLoggingLevels()
	opt := defaultOpts()
	for _, set := range opts {
		set(&opt)
	}
	for _, lvl := range opt.logLevels {
		al.logLevels[lvl.name] = lvl
	}
	al.appMinimumLoglevel = opt.appMinimumLoglevel
	al.longCallerNames = opt.longCallerNames
	al.colorizer = opt.colorizer
	//add colors to all console writers.
	if al.colorizer != nil {
		for _, lvl := range opt.logLevels {
			for wrtr, formatter := range lvl.writerFormatterMap {
				switch wrtr {
				case Stdout, Stderr:
					if !formatter.customColorizer {
						formatter.colorizer = al.colorizer
					}
				default:

				}
			}

		}
	}
	//add global writers and formatters to all levels
	for i, writerKey := range opt.globalWriterKeys {
		for _, lvl := range al.logLevels {
			if _, ok := lvl.writerFormatterMap[writerKey]; ok {
				continue
			}
			if al.logLevels[lvl.name].writerFormatterMap == nil {
				al.logLevels[lvl.name].writerFormatterMap = make(map[string]*formatterExpanded)
			}
			al.logLevels[lvl.name].writerFormatterMap[writerKey] = opt.globalFormatters[i]
		}
	}
	logMan = &al
	return nil
}

// LogmanOptions - settings for logMan object.
type LogmanOptions func(*options)

type options struct {
	appMinimumLoglevel int
	longCallerNames    bool
	logLevels          map[string]*loggingLevel
	colorizer          Colorizer
	globalWriterKeys   []string
	globalFormatters   []*formatterExpanded
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
		//o.logLevels = make(map[string]*loggingLevel)
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
		if importance > ImportanceNONE {
			importance = ImportanceNONE
		}
		if importance < ImportanceALL {
			importance = ImportanceALL
		}
		o.appMinimumLoglevel = importance
	}
}

// WithGlobalColorizer - sets global color scheme for logman
func WithGlobalColorizer(colorizer Colorizer) LogmanOptions {
	return func(o *options) {
		o.colorizer = colorizer
	}
}

// WithGlobalWriterFormatter - Add writer to all level.
// Useful to setup logfile.
func WithGlobalWriterFormatter(writer string, formatter *formatterExpanded) LogmanOptions {
	return func(o *options) {
		o.globalWriterKeys = append(o.globalWriterKeys, writer)
		o.globalFormatters = append(o.globalFormatters, formatter)
	}
}

// func WithLongCallerNames(val bool) LogmanOptions {
// 	return func(o *options) {
// 		o.longCallerNames = val
// 	}
// }

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
// If Message is nil function will return with no error.
func process(msg Message, lvls ...*loggingLevel) error {
	errorStack := []error{}
	fatalCalled := false
	if msg == nil {
		return nil
	}
	for _, lvl := range lvls {
		if lvl == nil {
			errorStack = append(errorStack, fmt.Errorf("logginglevel provided was not set"))
			continue
		}
		if lvl.importance < logMan.appMinimumLoglevel {
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

			if err := lvl.write(msg); err != nil {
				errorStack = append(errorStack, fmt.Errorf("writting message failed: %v", err))
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

func (lvl *loggingLevel) write(message Message) error {
	errorStack := []error{}
	var writer io.Writer
	for writerKey, formatter := range lvl.writerFormatterMap {
		switch writerKey {
		case Stderr:
			writer = os.Stderr
		case Stdout:
			writer = os.Stdout
		default:
			switch writerInfo(writerKey) {
			case "file":
				wr, err := os.OpenFile(writerKey, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
				switch err {
				case nil:
					writer = wr
				default:
					errorStack = append(errorStack, fmt.Errorf("failed to open writer '%v'", writerKey))
					continue
				}
			case "dir":
				sep := string(filepath.Separator)
				dirPath := strings.TrimSuffix(writerKey, sep) + sep
				msgTime, err := time.Parse(time.RFC3339Nano, fmt.Sprintf("%v", message.Value("time")))
				if err != nil {
					msgTime = time.Now()
				}
				msgFile := fmt.Sprintf("%v%v_%v.lmm", dirPath, msgTime.UnixNano(), lvl.name)
				wr, err := os.OpenFile(msgFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
				switch err {
				case nil:
					writer = wr
				default:
					errorStack = append(errorStack, fmt.Errorf("failed to open writer '%v'", writerKey))
					continue
				}
			}
		}
		text := formatter.Format(message, true)
		text = strings.TrimSuffix(text, "\n") + "\n"
		bt := []byte(text)
		_, err := writer.Write(bt)
		if err != nil {
			errorStack = append(errorStack, err)
		}
	}
	if err := joinErrors("writing message failed", errorStack...); err != nil {
		return err
	}
	return nil
}

func writerInfo(path string) string {
	f, err := os.Stat(path)
	if err != nil {
		return "bad"
	}
	if f.IsDir() {
		return "dir"
	}
	if f.Mode().IsRegular() {
		fl, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			fl.Close()
			return "bad"
		}
		fl.Close()
		return "file"
	}
	return "bad"
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
	// if !logMan.longCallerNames {
	// 	file = filepath.Base(file)
	// }
	if !success {
		return "", 0, ""
	}
	funcName := runtime.FuncForPC(counter).Name()
	return file, line, funcName
}
