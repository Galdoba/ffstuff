package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	LogLevelALL = iota
	LogLevelTRACE
	LogLevelDEBUG
	LogLevelINFO
	LogLevelWARN
	LogLevelERROR
	LogLevelFATAL
	LogLevelOFF
)

var process string
var logFile string
var loglimit int

func init() {
	loglimit = LogLevelOFF
	process, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	process = filepath.Base(process)
}

type entry struct {
	callerProgram    string //программа инициирующая событие
	timeStamp        string //время записи
	eventImportance  int    //степень важности события
	eventDescription string //описание события
}

func newEntry(eventDescription string, logLevel int) entry {
	en := entry{}
	en.callerProgram = process
	en.eventDescription = eventDescription
	en.eventImportance = logLevel
	en.timeStamp = time.Now().Format("2021-Mar-1-17:27:00.000")
	return en
}

func (en entry) write(file string) error {
	importance := ""
	switch en.eventImportance {
	case 0:
		importance = "ALL  :"
	case 1:
		importance = "TRACE:"
	case 2:
		importance = "DEBUG:"
	case 3:
		importance = "INFO :"
	case 4:
		importance = "WARN :"
	case 5:
		importance = "ERROR:"
	case 6:
		importance = "FATAL:"
	case 7:
		importance = "OFF  :"
	}
	event := importance + " " + en.timeStamp + " " + en.callerProgram + " | " + en.eventDescription + "\n"
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.WriteString(event); err != nil {
		return err
	}
	return nil
}

type logger struct {
	program string
}

func NewLog(path string, level int) {
	logFile = path
}

type Logger interface {
	ALL(string)
	TRACE(string)
	DEBUG(string)
	INFO(string)
	WARN(string)
	ERROR(string)
	FATAL(string)
	OFF(string)
}

func (l *logger) ALL(eventDescription string) error {
	en := newEntry(eventDescription, 0)
	if err := en.write(logFile); err != nil {
		return err
	}
	return nil
}
