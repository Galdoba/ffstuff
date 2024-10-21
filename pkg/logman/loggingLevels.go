package logman

import (
	"io"
	"os"
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
	name               string
	tag                string
	importance         int
	callerInfo         bool
	osExit             bool
	colorSchemes       map[string]uint8
	formatFunc         func(Message) (string, error)
	FMTE               *formatterExpanded
	writerFormatterMap map[string]*formatterExpanded
}

func NewLoggingLevel(name string, opts ...LevelOpts) *loggingLevel {
	lo := loggingLevel{}
	lo.name = name
	options := defaultLevel()
	lo.tag = name
	for _, enrich := range opts {
		enrich(&options)
	}
	lo.callerInfo = options.callerInfo
	lo.osExit = options.osExit
	lo.formatFunc = options.formatFunc
	lo.importance = options.importance
	lo.writerFormatterMap = options.writerFormatterMap
	return &lo
}

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

type lvlOpts struct {
	tag                string
	importance         int
	callerInfo         bool
	osExit             bool
	writers            map[string]io.Writer
	formatFunc         func(Message) (string, error)
	writerFormatterMap map[string]*formatterExpanded
}

func defaultLevel() lvlOpts {
	return lvlOpts{
		tag:                "[CUSTOM]",
		importance:         ImportanceINFO,
		callerInfo:         false,
		writerFormatterMap: make(map[string]*formatterExpanded),
	}
}

type LevelOpts func(*lvlOpts)

func WithWriter(writerKey string, expandedFormatter *formatterExpanded) LevelOpts {
	return func(lo *lvlOpts) {
		lo.writerFormatterMap[writerKey] = expandedFormatter
	}
}
