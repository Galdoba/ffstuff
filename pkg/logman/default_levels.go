package logman

import (
	"io"
	"os"
)

const (
	stdTagFATAL = "[FATAL]"
	stdTagERROR = "[ERROR]"
	stdTagWARN  = "[WARN ]"
	stdTagINFO  = "[INFO ]"
	stdTagDEBUG = "[DEBUG]"
	stdTagTRACE = "[TRACE]"
)

var LogLevelFATAL = &loggingLevel{
	name:       FATAL,
	tag:        stdTagFATAL,
	importance: ImportanceFATAL,
	callerInfo: true,
	osExit:     true,
	writers:    map[string]io.Writer{Stdout: os.Stdout},
	formatFunc: formatTextComplex,
}

var LogLevelERROR = &loggingLevel{
	name:       ERROR,
	tag:        stdTagERROR,
	importance: ImportanceERROR,
	callerInfo: true,
	osExit:     false,
	writers:    map[string]io.Writer{Stdout: os.Stdout},
	formatFunc: formatTextComplex,
}

var LogLevelWARN = &loggingLevel{
	name:       WARN,
	tag:        stdTagWARN,
	importance: ImportanceWARN,
	callerInfo: false,
	osExit:     false,
	writers:    map[string]io.Writer{Stdout: os.Stdout},
	formatFunc: formatTextSimple,
}

var LogLevelINFO = &loggingLevel{
	name:       INFO,
	tag:        stdTagINFO,
	importance: ImportanceINFO,
	callerInfo: false,
	osExit:     false,
	writers:    map[string]io.Writer{Stderr: os.Stderr},
	formatFunc: formatTextSimple,
}

var LogLevelDEBUG = &loggingLevel{
	name:       DEBUG,
	tag:        stdTagDEBUG,
	importance: ImportanceDEBUG,
	callerInfo: false,
	osExit:     false,
	writers:    map[string]io.Writer{Stderr: os.Stderr},
	formatFunc: formatTextComplex,
}

var LogLevelTRACE = &loggingLevel{
	name:       TRACE,
	tag:        stdTagTRACE,
	importance: ImportanceTRACE,
	callerInfo: true,
	osExit:     false,
	writers:    map[string]io.Writer{Stderr: os.Stderr},
	formatFunc: formatTextComplex,
}

func defaultLoggingLevels() map[string]*loggingLevel {
	levels := make(map[string]*loggingLevel)
	levels[FATAL] = LogLevelFATAL
	levels[ERROR] = LogLevelERROR
	levels[WARN] = LogLevelWARN
	levels[INFO] = LogLevelINFO
	levels[DEBUG] = LogLevelDEBUG
	levels[TRACE] = LogLevelTRACE
	return levels
}
