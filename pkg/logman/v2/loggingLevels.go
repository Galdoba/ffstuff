package logman

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	FATAL = "fatal"
	ERROR = "error"
	WARN  = "warn"
	INFO  = "info"
	DEBUG = "debug"
	TRACE = "trace"
	PING  = "ping"
	ALL   = "All_Levels"
)

type loggingLevel struct {
	name         string
	tag          string
	importance   int
	callerInfo   bool
	osExit       bool
	writers      map[string]io.Writer
	formatters   map[string]func(Message) (string, error)
	colorSchemes map[string]uint8
	formatFunc   func(Message) (string, error)
}

func NewLoggingLevel(name string, opts ...LevelOpts) *loggingLevel {
	lo := loggingLevel{}
	lo.name = name
	options := defaultLevel()
	lo.tag = strings.ToUpper(fmt.Sprintf("[%v]", name))
	for _, enrich := range opts {
		enrich(&options)
	}
	lo.callerInfo = options.callerInfo
	lo.osExit = options.osExit
	lo.formatFunc = options.formatFunc
	lo.importance = options.importance
	lo.tag = options.tag
	lo.writers = options.writers
	return &lo
}

/*
NewLoggingLevel("report", 55,
color bool
writers := []string{}
tagColor = int
ignore fields
enforce fields
timestamp short/full
appStartTime time.Time
sinceStart bool
)
*/

func LevelTag(tag string) LevelOpts {
	return func(lvl *lvlOpts) {
		lvl.tag = tag
	}
}

func LevelImportance(imp int) LevelOpts {
	return func(lvl *lvlOpts) {
		lvl.importance = imp
	}
}

func LevelCallerInfo(callerInfo bool) LevelOpts {
	return func(lvl *lvlOpts) {
		lvl.callerInfo = callerInfo
	}
}

func LevelExitWhenDone(osExit bool) LevelOpts {
	return func(lvl *lvlOpts) {
		lvl.osExit = osExit
	}
}

func LevelWriters(writerKeys ...string) LevelOpts {
	return func(lvl *lvlOpts) {
		wrtrs := make(map[string]io.Writer)
		for _, key := range writerKeys {
			switch key {
			case Stderr:
				wrtrs[key] = os.Stderr
			case Stdout:
				wrtrs[key] = os.Stdout
			default:
				f, err := os.OpenFile(key, flags, perm)
				if err != nil {
					panic("TODO: add file check: " + err.Error())
				}
				wrtrs[key] = f
				f.Close()
			}
		}
		lvl.writers = wrtrs
	}
}

// LevelFormatterFunc sets custom formatter function.
func LevelFormatterFunc(formatter func(Message) (string, error)) LevelOpts {
	return func(lvl *lvlOpts) {
		lvl.formatFunc = formatter
	}
}

func (lvl *loggingLevel) setWriter(key string) {
	switch key {
	case Stderr:
		lvl.writers[key] = os.Stderr
	case Stdout:
		lvl.writers[key] = os.Stdout
	default:
		f, err := os.OpenFile(key, flags, perm)
		if err != nil {
			panic("TODO: add file check: " + err.Error())
		}
		lvl.writers[key] = f
		f.Close()
	}
}

func (lvl *loggingLevel) clearWriters() {
	lvl.writers = make(map[string]io.Writer)
}

type lvlOpts struct {
	tag        string
	importance int
	callerInfo bool
	osExit     bool
	writers    map[string]io.Writer
	formatFunc func(Message) (string, error)
}

func defaultLevel() lvlOpts {
	return lvlOpts{
		tag:        "[CUSTOM]",
		importance: ImportanceINFO,
		callerInfo: false,
		writers:    map[string]io.Writer{Stderr: os.Stderr},
		formatFunc: func(Message) (string, error) {
			return "nothing", fmt.Errorf("formatFunc not set")
		},
	}
}

type LevelOpts func(*lvlOpts)
