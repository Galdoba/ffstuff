package logfile

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
	LogLevelALL   = 0
	LogLevelTRACE = 1
	LogLevelDEBUG = 2
	LogLevelINFO  = 3
	LogLevelWARN  = 4
	LogLevelERROR = 5
	LogLevelFATAL = 6
	LogLevelOFF   = 7
)

func Test() {
	fmt.Println("START LOG TEST")
	fmt.Println("process =", process)
	fmt.Println("logFile =", logFile)
	fmt.Println("loglimit =", loglimit)
	logger := New(fldr.MuxPath()+"logfile.txt", 7)
	logger.TRACE("Trace message")
	logger.DEBUG("Debug message")
	logger.INFO("Info message")
	fmt.Println("END LOG TEST")
}

var process string
var logFile string
var loglimit int

func init() {
	loglimit = LogLevelDEBUG
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
	en.callerProgram = process
	en.eventDescription = eventDescription
	en.eventImportance = logLevel
	en.timeStamp = time.Now().Format("2006-01-02 15:04:05.000") //Mon Jan 2 15:04:05 -0700 MST 2006
	return en
}

func (en entry) write(file string) error {
	importance := ""
	switch en.eventImportance {
	case 0:
		importance = "ALL  "
	case 1:
		importance = "TRACE"
	case 2:
		importance = "DEBUG"
	case 3:
		importance = "INFO "
	case 4:
		importance = "WARN "
	case 5:
		importance = "ERROR"
	case 6:
		importance = "FATAL"
	case 7:
		importance = "OFF  "
	}
	n := 12
	for len(en.callerProgram) < n {
		en.callerProgram = en.callerProgram + " "
	}
	if len(en.callerProgram) > n {
		en.callerProgram = en.callerProgram[0:n]
	}
	event := en.timeStamp + " | " + en.callerProgram + " | " + importance + " | " + en.eventDescription + "\n"
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.WriteString(event); err != nil {
		return err
	}
	if en.eventImportance >= loglimit {
		color.Output = ansi.NewAnsiStdout()
		switch en.eventImportance {
		case 0, 1, 3:
			color.White(event)
		case 2:
			color.Blue(event)
		case 4:
			color.Yellow(event)
		default:
			color.Red(event)
		}
		//fmt.Println(event)
	}
	return nil
}

type logger struct {
	program         string
	importanceLevel int
	file            string
}

//New -
func New(path string, level int) Logger {
	//logFile = path
	//os.MkdirAll(confDir, os.ModePerm) - TODO: сделать создание папки если таковой нет
	f, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		fn, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
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
	//OFF(string)
}

func (l *logger) ALL(eventDescription string) error {
	en := newEntry(eventDescription, 0)
	if err := en.write(l.file); err != nil {
		return err
	}
	return nil
}

func (l *logger) TRACE(eventDescription string) error {
	en := newEntry(eventDescription, 1)
	if err := en.write(l.file); err != nil {
		return err
	}
	return nil
}

func (l *logger) DEBUG(eventDescription string) error {
	en := newEntry(eventDescription, 2)
	if err := en.write(l.file); err != nil {
		return err
	}
	return nil
}

func (l *logger) INFO(eventDescription string) error {
	en := newEntry(eventDescription, 3)
	if err := en.write(l.file); err != nil {
		return err
	}
	return nil
}

func (l *logger) WARN(eventDescription string) error {
	en := newEntry(eventDescription, 4)
	if err := en.write(l.file); err != nil {
		return err
	}
	return nil
}

func (l *logger) ERROR(eventDescription string) error {
	en := newEntry(eventDescription, 5)
	if err := en.write(l.file); err != nil {
		return err
	}
	return nil
}

func (l *logger) FATAL(eventDescription string) error {
	en := newEntry(eventDescription, 6)
	if err := en.write(l.file); err != nil {
		return err
	}
	return nil
}
