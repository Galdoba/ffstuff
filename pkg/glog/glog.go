package glog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Galdoba/ffstuff/fldr"
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
)

const (
	LogLevelALL    = 0
	LogLevelTRACE  = 1
	LogLevelDEBUG  = 2
	LogLevelINFO   = 3
	LogLevelWARN   = 4
	LogLevelERROR  = 5
	LogLevelFATAL  = 6
	LogLevelOFF    = 7
	LogPathDEFAULT = "log_path_default"
)

func Test() {
	fmt.Println("START LOG TEST")

	logger := New(fldr.MuxPath()+"logfile.txt", LogLevelINFO)
	fmt.Println("process =", process)
	fmt.Println("logFile =", logFile)
	fmt.Println("loglimit =", loglimit)
	logger.ALL("All message")
	logger.TRACE("Trace message")
	logger.DEBUG("Debug message")
	logger.INFO("Info message")
	logger.WARN("Warn message")
	logger.ERROR("Error message")
	logger.FATAL("Fatal message")
	fmt.Println("END LOG TEST")
}

var process string
var logFile string
var loglimit int

func init() {
	//loglimit = LogLevelDEBUG
	// processInit, err := os.Executable()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// processInit = filepath.Base(processInit)
	// switch runtime.GOOS {
	// case "windows":
	// 	processInit = strings.TrimSuffix(processInit, ".exe")
	// }
	// process = processInit
}

type entry struct {
	callerProgram    string //программа инициирующая событие
	timeStamp        string //время записи
	eventImportance  int    //степень важности события
	eventDescription string //описание события
}

func newEntry(eventDescription string, logLevel int) entry {
	en := entry{}
	processInit, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	processInit = filepath.Base(processInit)
	switch runtime.GOOS {
	case "windows":
		processInit = strings.TrimSuffix(processInit, ".exe")
	}
	process = processInit
	for len(process) < 15 {
		process += " "
	}
	if len(process) > 15 {
		process = string([]byte(process)[0:15])
	}
	en.callerProgram = process
	en.eventDescription = eventDescription
	en.eventImportance = logLevel
	en.timeStamp = time.Now().Format("2006-01-02 15:04:05.000") //Mon Jan 2 15:04:05 -0700 MST 2006
	return en
}

func (en entry) String() string {
	return en.timeStamp + " | " + en.callerProgram + " | " + importanceStr(en.eventImportance) + " | " + en.eventDescription + "\n"
}

func importanceStr(imp int) string {
	str := ""
	switch imp {
	case 0:
		str = "ALL  "
	case 1:
		str = "TRACE"
	case 2:
		str = "DEBUG"
	case 3:
		str = "INFO "
	case 4:
		str = "WARN "
	case 5:
		str = "ERROR"
	case 6:
		str = "FATAL"
	case 7:
		str = "OFF  "
	}
	return str
}

func (en entry) write(file string) error {
	//importance := importanceStr(en.eventImportance)

	n := 15
	for len(en.callerProgram) < n {
		en.callerProgram = en.callerProgram + " "
	}
	if len(en.callerProgram) > n {
		en.callerProgram = en.callerProgram[0:n]
	}
	//event := en.timeStamp + " | " + en.callerProgram + " | " + importance + " | " + en.eventDescription + "\n"
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.WriteString(en.String()); err != nil {
		return err
	}
	// if en.eventImportance >= loglimit {
	// 	color.Output = ansi.NewAnsiStdout()
	// 	switch en.eventImportance {
	// 	case 0, 1, 3:
	// 		color.White(event)
	// 	case 2:
	// 		color.Blue(event)
	// 	case 4:
	// 		color.Yellow(event)
	// 	default:
	// 		color.Red(event)
	// 	}
	// 	//fmt.Println(event)
	// }
	return nil
}

type logger struct {
	program         string
	importanceLevel int
	shoutLevel      int //степень важности события необходимая для вывода в терминал
	file            string
}

//New - принимает путь в котором будет находиться файл и минимальный уровень сообщений которые будут выводиться на терминал во время логирования
func New(path string, level int) Logger {
	if path == LogPathDEFAULT {
		path = fldr.LogPathDefault()
	}
	pathFolders := strings.Split(path, "\\")
	dir := strings.Join(pathFolders[0:len(pathFolders)-1], "\\")
	os.MkdirAll(dir, os.ModePerm)

	f, err := os.OpenFile(path, os.O_RDONLY, 0600)

	if err != nil {
		//fmt.Println(err)
		fn, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		//fmt.Println("Create log file:", path)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer fn.Close()
		lgr := logger{}
		lgr.file = path
		lgr.importanceLevel = level
		lgr.INFO("Create this log file")
		return &lgr
	}
	defer f.Close()

	lgr := logger{}
	lgr.file = path
	lgr.importanceLevel = level
	lgr.shoutLevel = LogLevelINFO
	return &lgr
}

//Logger -
type Logger interface {
	ALL(string) error
	TRACE(string) error
	DEBUG(string) error
	INFO(string) error
	WARN(string) error
	ERROR(string) error
	FATAL(string) error
	ShoutWhen(int)
	//OFF(string)
}

func (l *logger) ALL(eventDescription string) error {
	en := newEntry(eventDescription, 0)
	if err := en.write(l.file); err != nil {
		return err
	}
	shout(en, l)
	return nil
}

func (l *logger) TRACE(eventDescription string) error {
	en := newEntry(eventDescription, 1)
	if err := en.write(l.file); err != nil {
		return err
	}
	shout(en, l)
	return nil
}

func (l *logger) DEBUG(eventDescription string) error {
	en := newEntry(eventDescription, 2)
	if err := en.write(l.file); err != nil {
		return err
	}
	shout(en, l)
	return nil
}

func (l *logger) INFO(eventDescription string) error {
	en := newEntry(eventDescription, 3)
	if err := en.write(l.file); err != nil {
		return err
	}
	shout(en, l)
	return nil
}

func (l *logger) WARN(eventDescription string) error {
	en := newEntry(eventDescription, 4)
	if err := en.write(l.file); err != nil {
		return err
	}
	shout(en, l)
	return nil
}

func (l *logger) ERROR(eventDescription string) error {
	en := newEntry(eventDescription, 5)
	if err := en.write(l.file); err != nil {
		return err
	}
	shout(en, l)
	return nil
}

func (l *logger) FATAL(eventDescription string) error {
	en := newEntry(eventDescription, 6)
	if err := en.write(l.file); err != nil {
		return err
	}
	shout(en, l)
	return nil
}

func (l *logger) ShoutWhen(lvl int) {
	switch lvl {
	default:
		return
	case 0, 1, 2, 3, 4, 5, 6, 7:
		l.shoutLevel = lvl
	}
}

// func SetLogLimit(lvl int) {
// 	loglimit = lvl
// }

func shout(en entry, l *logger) {
	if en.eventImportance >= l.shoutLevel {
		color.Output = ansi.NewAnsiStdout()
		switch en.eventImportance {
		case 0, 1, 3:
			color.White(en.String())
		case 2:
			color.Blue(en.String())
		case 4:
			color.Yellow(en.String())
		default:
			color.Red(en.String())
		}
	}
}
