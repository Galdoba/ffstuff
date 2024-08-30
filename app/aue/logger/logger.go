package logger

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var flags int = os.O_CREATE | os.O_WRONLY | os.O_APPEND
var perm fs.FileMode = 0666

type aueLogManager struct {
	logPath   string
	debugMode bool
	perm      fs.FileMode
	logger    *log.Logger
}

type OptFunc func(*options)

type options struct {
	logpath   string
	debugMode bool
}

func defaultOpts() options {
	return options{
		logpath:   "aue.log",
		debugMode: true,
	}
}

func LogFilepath(path string) OptFunc {
	return func(o *options) {
		o.logpath = path
	}
}

func DebugMode() OptFunc {
	return func(o *options) {
		o.debugMode = true
	}
}

var LOG *logrus.Logger

var logg *aueLogManager

func New(opts ...OptFunc) *aueLogManager {
	al := aueLogManager{}
	opt := defaultOpts()
	for _, set := range opts {
		set(&opt)
	}
	al.logPath = opt.logpath
	al.debugMode = opt.debugMode
	al.logger = log.New(os.Stdout, "", 0)
	return &al
}

func debugPrefix() string {
	return "DEBUG"
}

func Func1() {
	fmt.Println("f1")
	Func2()
}

type some struct {
}

type Some interface {
	Func3(string)
}

func Func2() {
	fmt.Println("f2")
	sm := &some{}
	SM := Some(sm)
	SM.Func3("++")
}
func (sm *some) Func3(s string) {
	fmt.Println("f3", s)
	logg.Debug("this is dbg", 3, "___", sm)
	Func4()

}
func Func4() {
	fmt.Println("f4")
	logg.Debug("this is dbg 22")
}

func (al *aueLogManager) Debug(msg string, args ...interface{}) {
	if !al.debugMode {
		return
	}
	f, err := os.OpenFile(al.logPath, flags, perm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fd, err := os.OpenFile(debugLogPath(al.logPath), flags, perm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	prefix := debugPrefix() + " "
	argStr := argString(args...)
	msg = prefix + callerFunctionInfo(true, true, true) + " " + argStr + " " + msg
	al.logger.SetOutput(io.MultiWriter(f, fd))
	al.logger.Printf("%v", msg)
	//f.Close()
	//fd.Close()
}

func argString(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	s := "args: "
	for _, arg := range args {
		switch v := arg.(type) {
		default:
			s += fmt.Sprintf("{%T : '%v'}", v, v)
		}

	}
	return s
}

func debugLogPath(logpath string) string {
	dir := filepath.Dir(logpath)
	name := filepath.Base(logpath)
	return dir + "/debug_" + name
}

func callerFunctionInfo(funcName, fileName, lineNumber bool) string {
	counter, file, line, success := runtime.Caller(2)
	if !success {
		return ""
	}

	fName := runtime.FuncForPC(counter).Name()
	info := ""
	if funcName {
		fName = strings.TrimPrefix(fName, "github.com/Galdoba/ffstuff/app/aue/")
		info += fName + " "
	}
	if fileName {
		file := strings.Split(file, "github.com/Galdoba/ffstuff/app/")
		info += file[len(file)-1] + " "
	}
	if lineNumber && line > 0 {
		info += fmt.Sprintf("line %v", line) + " "
	}
	info = strings.TrimSuffix(info, " ")

	return info
}

/*
LOG EXAMPLE:
2024/08/30 16:43:11.540 [DEBUG]: func name: file.go line #666 args(agr1, arg2) message
2024/08/30 16:43:11.540 [INFO ]: operation complete
2024/08/30 16:43:11.540 [WARN ]: failed to do X
2024/08/30 16:43:11.540 [ERROR]: failed to do Y
2024/08/30 16:43:11.540 [FATAL]: failed to do Z, exiting...

2024/08/29 17:58:41.318 [DEBUG]: logger.(*some).Func3 aue/logger/logger.go line 92 args: {int : '3'}{string : '___'}{*logger.some : '&{}'} this is dbg


2024/08/30 16:43:11.540 [ERROR]: failed to do Y
2024/08/29 17:58:41.318 [DEBUG]: file:[aue/logger/logger.go] line 92: logger.(*some).Func3(3 int, ___ string, &{} *logger.some)
    debug message: this is dbg
*/
