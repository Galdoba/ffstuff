package logger

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	debug    = "[DEBUG]"
	info     = "[INFO ]"
	warn     = "[WARN ]"
	errLv    = "[ERROR]"
	fatal    = "[FATAL]"
	LvlFATAL = 0
	LvlERROR = 10
	LvlWARN  = 20
	LvlINFO  = 30
	LvlDEBUG = 40
	LvlTRACE = 50
)

var flags int = os.O_CREATE | os.O_WRONLY | os.O_APPEND
var perm fs.FileMode = 0666

type logManager struct {
	logPath   string
	debugMode bool
	loglevel  int
	logger    *log.Logger
}

type OptFunc func(*options)

type options struct {
	logpath   string
	debugMode bool
	loglevel  int
}

func defaultOpts() options {
	return options{
		loglevel: LvlINFO,
	}
}

func LogFilepath(path string) OptFunc {
	return func(o *options) {
		o.logpath = path
	}
}

func DebugMode(val bool) OptFunc {
	return func(o *options) {
		o.debugMode = val
	}
}

var logger *logManager

func Setup(opts ...OptFunc) error {
	al := logManager{}
	opt := defaultOpts()
	for _, set := range opts {
		set(&opt)
	}
	al.logPath = opt.logpath
	al.debugMode = opt.debugMode
	al.logger = log.New(os.Stdout, "", 0)
	logger = &al
	return validateLogger()
}

func validateLogger() error {
	for _, err := range []error{
		assertLogFilePath(logger.logPath),
		assertDebugLog(logger.debugMode, logger.logPath),
	} {
		if err != nil {
			return err
		}
	}
	return nil
}

func assertLogFilePath(path string) error {
	if path == "" {
		return fmt.Errorf("log filepath not provided")
	}
	if ch := invalidPathChar(path); ch != "" {
		return fmt.Errorf("log filepath provided contains bad character: '%v'", ch)
	}
	if err := os.Rename(path, path); err != nil {
		switch errors.Is(err, os.ErrNotExist) {
		case false:
			return fmt.Errorf("logfile assertion failed: %v", err)
		case true:
			dir := filepath.Dir(path)
			sep := string(filepath.Separator)
			name := filepath.Base(path)
			if err = os.MkdirAll(dir, 0666); err != nil {
				return fmt.Errorf("logfile directory creation failed: %v", err)
			}
			f, err := os.Create(dir + sep + name)
			defer f.Close()
			if err != nil {
				return fmt.Errorf("logfile creation failed: %v", err)
			}
		}
	}
	return nil
}

func assertDebugLog(debugMode bool, path string) error {
	if !debugMode {
		return nil
	}
	return assertLogFilePath(debugLogPath(path))
}

func invalidPathChar(path string) string {
	path = strings.ReplaceAll(path, `\`, "/")
	layers := strings.Split(path, "/")
	for _, layer := range layers {
		for _, ch := range strings.Split(layer, "") {
			switch ch {
			case "<", ">", "/", "|", "?", "*", `"`, " ":
				return ch
			default:
			}
		}
	}
	return ""
}

func prefix(s string) string {
	out := timestamp()
	switch s {
	case debug:
		out += " " + debug
	case info:
		out += " " + info
	case warn:
		out += " " + warn
	case errLv:
		out += " " + errLv
	case fatal:
		out += " " + fatal
	default:
		out += " [ ??? ]"
	}
	return out
}

//Debug - Print debug message.
//Very verbose for tracking data. Do not return early
//Format: '{DATE} {TIME} {LEVEL} {file} {line} {funcname} [args] >> {MESSAGE}'
func Debug(msg string, args ...interface{}) error {
	if logger == nil {
		return fmt.Errorf("logger was not initiated")
	}
	al := logger
	if !al.debugMode {
		return nil
	}
	f, err := os.OpenFile(al.logPath, flags, perm)
	if err != nil {
		fmt.Println(err)
	}
	fd, err := os.OpenFile(debugLogPath(al.logPath), flags, perm)
	if err != nil {
		fmt.Println(err)
	}

	argStr := argsInfo(args...)
	codeInfo := callerFunctionInfo(true, true, true) + argStr
	prefix := prefix(debug) + codeInfo
	msg = prefix + " >> " + msg
	al.logger.SetOutput(io.MultiWriter(f, fd))
	al.logger.Printf("%v", msg)
	f.Close()
	fd.Close()
	return nil
}

//Info - Print debug message.
//General messages to file and stderr.
//Format: '{DATE} {TIME} {LEVEL} >> {MESSAGE}'
func Info(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	if logger == nil {
		return fmt.Errorf("logger was not initiated")
	}
	al := logger
	if !al.debugMode {
		return nil
	}
	f, err := os.OpenFile(al.logPath, flags, perm)
	if err != nil {
		fmt.Println(err)
		return err
	}
	prefix := prefix(info)
	msg = prefix + " >> " + msg
	msg = strings.ReplaceAll(msg, "\r", "")
	msg = strings.ReplaceAll(msg, "\n", ": ")
	al.logger.SetOutput(io.MultiWriter(f, os.Stderr))
	al.logger.Printf("%v", msg)
	f.Close()
	return nil
}

//Warn - Print debug message.
//Important messages.
//Format: '{DATE} {TIME} {LEVEL} >> {MESSAGE}'
func Warn(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	if logger == nil {
		return fmt.Errorf("logger was not initiated")
	}
	al := logger
	if !al.debugMode {
		return nil
	}
	f, err := os.OpenFile(al.logPath, flags, perm)
	if err != nil {
		fmt.Println(err)
		return err
	}
	prefix := prefix(warn)
	msg = prefix + " >> " + msg
	msg = strings.ReplaceAll(msg, "\r", "")
	msg = strings.ReplaceAll(msg, "\n", ": ")
	al.logger.SetOutput(io.MultiWriter(f, os.Stdout))
	al.logger.Printf("%v", msg)
	f.Close()
	return nil
}

//Error - Print debug message.
//Vey bad messages. Non-Critical
//Format: '{DATE} {TIME} {LEVEL} >> {MESSAGE}'
func Error(err error) error {
	msg := err.Error()
	if logger == nil {
		return fmt.Errorf("logger was not initiated")
	}
	al := logger
	if !al.debugMode {
		return nil
	}
	f, err := os.OpenFile(al.logPath, flags, perm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	prefix := prefix(errLv)
	msg = prefix + " >> " + msg
	msg = strings.ReplaceAll(msg, "\r", "")
	msg = strings.ReplaceAll(msg, "\n", ": ")
	al.logger.SetOutput(io.MultiWriter(f, os.Stdout))
	al.logger.Printf("%v", msg)
	f.Close()
	return nil
}

//Fatal - Print fatal message.
//Vey bad messages. Critical Error. Exit with code 1.
//Format: '{DATE} {TIME} {LEVEL} >> {MESSAGE}'
func Fatal(msg string, args ...interface{}) error {
	if logger == nil {
		fmt.Println("logger was not initiated")
	}
	al := logger
	if !al.debugMode {
		return nil
	}
	f, err := os.OpenFile(al.logPath, flags, perm)
	if err != nil {
		fmt.Println(err)
	}
	fd, err := os.OpenFile(debugLogPath(al.logPath), flags, perm)
	if err != nil {
		fmt.Println(err)
	}

	argStr := argsInfo(args...)
	codeInfo := callerFunctionInfo(true, true, true) + argStr
	prefix := prefix(fatal) + codeInfo
	msg = prefix + " >> " + msg
	al.logger.SetOutput(io.MultiWriter(f, fd, os.Stdout))
	al.logger.Printf("%v", msg)
	f.Close()
	fd.Close()
	fmt.Fprintf(os.Stdout, "%v\n", msg)
	fmt.Fprintf(os.Stdout, "exit code 1")
	os.Exit(1)
	return nil
}

func argsInfo(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	s := "("
	for _, arg := range args {
		switch v := arg.(type) {
		default:
			s += fmt.Sprintf("'%v' %T, ", v, v)
		}

	}
	s = strings.TrimSuffix(s, ", ")
	return s + ")"
}

func debugLogPath(logpath string) string {
	dir := filepath.Dir(logpath)
	name := filepath.Base(logpath)
	return dir + "/debug_" + name
}

func callerFunctionInfo(funcName, fileName, lineNumber bool) string {
	counter, file, line, success := runtime.Caller(2) //back to stack on 2 levels
	if !success {
		return ""
	}

	fName := runtime.FuncForPC(counter).Name()
	info := ""

	if fileName {
		file := strings.Split(file, "github.com/Galdoba/ffstuff/app/")
		info += fmt.Sprintf(" [%v]", file[len(file)-1])
	}
	if lineNumber && line > 0 {
		info += fmt.Sprintf(" line %v:", line)
	}
	if funcName {
		fName = strings.TrimPrefix(fName, "github.com/Galdoba/ffstuff/app/aue/")
		info += fmt.Sprintf(" func: %v", fName)
	}
	info = strings.TrimSuffix(info, " ")

	return info
}

func timestamp() string {
	format := "2006/01/02 15:04:05.999"
	formatLen := len(format)
	stamp := time.Now().Format(format)
	for len(stamp) < formatLen {
		stamp += "0"
	}
	return stamp
}
