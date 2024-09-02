package v2

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var logMan *logManager
var flags int = os.O_CREATE | os.O_WRONLY | os.O_APPEND
var perm fs.FileMode = 0666

const (
	LvlFATAL = 0
	LvlERROR = 10
	LvlWARN  = 20
	LvlINFO  = 30
	LvlDEBUG = 40
	LvlTRACE = 50

	//fieldKeys
	keyTime    = "time"
	keyLevel   = "level"
	keyMessage = "message"
	keyFile    = "file"
	keyLine    = "line"
	keyFunc    = "func"
)

type logManager struct {
	logPaths    map[int]string
	appLoglevel int
	logLevels   map[int]*levelOptions
	logger      *log.Logger
}

func Setup(opts ...OptFunc) error {
	al := logManager{}
	al.appLoglevel = -1
	opt := defaultOpts()
	for _, set := range opts {
		set(&opt)
	}
	al.logLevels = opt.logLevels
	al.logPaths = opt.logpaths
	al.logger = log.New(os.Stdout, "", 0)
	logMan = &al
	return nil
}

type OptFunc func(*options)

type options struct {
	logpaths  map[int]string
	logLevels map[int]*levelOptions
}

func defaultOpts() options {
	opts := options{}
	opts.logLevels = make(map[int]*levelOptions)
	opts.logpaths = make(map[int]string)
	opts.logLevels[LvlDEBUG] = &levelOptions{
		levelTag:       "[DEBUG]",
		callerFilename: true,
		callerFileline: true,
		callerFuncName: true,
	}
	opts.logLevels[LvlINFO] = &levelOptions{
		levelTag: "[INFO ]",
		writers:  map[string]io.Writer{"stderr": os.Stderr},
	}
	return options{
		//loglevel: LvlINFO,
	}
}

func WithLogPathAll(path string) OptFunc {
	return func(o *options) {
		for lvl := range o.logpaths {
			o.logpaths[lvl] = path
		}
	}
}

func WithLogPathLevel(lvl int, path string) OptFunc {
	return func(o *options) {
		o.logpaths[lvl] = path
	}
}

func WithLogLevel(lvl int, level levelOptions) OptFunc {
	return func(o *options) {
		o.logLevels[lvl] = &level
	}
}

type levelOptions struct {
	levelTag       string
	callerFilename bool
	callerFileline bool
	callerFuncName bool
	//textFormat     string
	//writerFormat   string
	writers map[string]io.Writer
	//colorsFG map[string]string
	//colorBG  string
	execFunc func(*message) error
}

func NewLogLevel() levelOptions {
	lo := levelOptions{}
	return lo
}

func (ll *levelOptions) formatAsText(msg *message) string {
	output := ""
	timestamp := fmt.Sprintf("%v ", timestampFormat(msg.timeCreated))
	output += fmt.Sprintf("%v ", timestamp)
	output += fmt.Sprintf("%v ", ll.levelTag)
	output += fmt.Sprintf("%v ", msg.fields[keyMessage])
	if ll.callerFilename {
		if val, ok := msg.fields[keyFile]; ok {
			output += fmt.Sprintf("%v ", val)
		}
	}
	if ll.callerFilename {
		if val, ok := msg.fields[keyLine]; ok {
			output += fmt.Sprintf("%v ", val)
		}
	}
	if ll.callerFilename {
		if val, ok := msg.fields[keyFunc]; ok {
			output += fmt.Sprintf("%v ", val)
		}
	}
	if len(msg.args()) == 0 {
		output = strings.TrimSuffix(output, " ")
	}
	output += "{"
	for _, arg := range msg.args() {
		output += fmt.Sprintf("%v:%v; ", arg.key, arg.value)
	}
	output = strings.TrimSuffix(output, "; ") + "}"

	return output
}

func (ll *levelOptions) formatAsJson(msg *message) string {
	mandFlds, argFlds := msg.sortFields()
	text := "{"
	for _, fld := range mandFlds {
		text += fmt.Sprintf(`"%v" : "%v",`, fld.key, fld.value)
	}
	for _, fld := range argFlds {
		text += fmt.Sprintf(`"%v" : "%v",`, fld.key, fld.value)
	}
	text = fmt.Sprintf(text, ",")
	return text + "}"
}

/*
example:
log.Message("fail process: %v", process).Info()
log.Message("fail process: %v", process).Warn()
log.Message("fail process: %v", process).Error()
//convinience
log.Warn(format string, args ...interface{})
log.Error(err)
*/

//logger.Process(logger.NewMessage("event %v complete", eventName), logger.LevelINFO, logger.LevelSECURITY)
//logger.Debug("event %v complete", eventName)

func (lm *logManager) Process(msg *message, lvls ...int) error {
	errProcessing := []error{}
	for _, lvl := range lvls {
		loggingLevel := lm.logLevels[lvl]
		if loggingLevel == nil {
			errProcessing = append(errProcessing, fmt.Errorf("message processing failed: level %v does not exist", lvl))
		}

		if loggingLevel.callerFilename || loggingLevel.callerFileline || loggingLevel.callerFuncName {
			fileName, line, funcName := callerFunctionInfo()
			if loggingLevel.callerFilename {
				msg = msg.WithFields(NewField("file", fileName))
			}
			if loggingLevel.callerFilename {
				msg = msg.WithFields(NewField("line", line))
			}
			if loggingLevel.callerFilename {
				msg = msg.WithFields(NewField("funcName", funcName))
			}
		}

		for key, _ := range loggingLevel.writers {
			switch key {
			case "stderr":
				text := loggingLevel.formatAsText(msg)
				fmt.Fprintf(os.Stderr, text)
			case "stdout":
				text := loggingLevel.formatAsText(msg)
				fmt.Fprintf(os.Stdout, text)
			default:
				f, _ := os.OpenFile(key, flags, perm)
				text := loggingLevel.formatAsJson(msg)
				f.WriteString(text)
			}

		}

	}
	return assembleErrors(errProcessing...)
}

func assembleErrors(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}
	errOut := fmt.Errorf("combined error:")
	for _, err := range errs {
		errOut = fmt.Errorf("%v\n%v", errOut, err)
	}
	return errOut
}

func (msg *message) sortFields() ([]messageField, []messageField) {
	mandatoryFlds := []messageField{}
	for _, key := range fieldKeysMandatory() {
		if val, ok := msg.fields[key]; ok {
			mandatoryFlds = append(mandatoryFlds, messageField{key, val})
		}
	}
	argFlds := msg.args()
	return mandatoryFlds, argFlds
}

type message struct {
	fields      map[string]interface{}
	timeCreated time.Time
}

func Message(format string, args ...interface{}) *message {
	m := message{}
	m.fields = make(map[string]interface{})
	m.timeCreated = time.Now()
	m.fields[keyMessage] = fmt.Sprintf(format, args...)
	m.fields[keyTime] = fmt.Sprintf(timestampFormat(m.timeCreated))
	return &m
}

func (m *message) WithFields(flds ...messageField) *message {
	for _, fld := range flds {
		m.fields[fld.key] = fld.value
	}
	return m
}

func (m *message) WithArgs(args ...interface{}) *message {
	fldLen := len(m.fields)
	for i := 0; i < fldLen; i++ {
		key := fmt.Sprintf("arg %v", i)
		if _, ok := m.fields[fmt.Sprintf("arg %v", i)]; ok {
			delete(m.fields, key)
		}
	}
	for i, arg := range args {
		newKey := fmt.Sprintf("arg %v", i)
		m.fields[newKey] = fmt.Sprintf("%v %T", arg, arg)
	}
	return m
}

func (m *message) args() []messageField {
	fldLen := len(m.fields)
	argFields := []messageField{}
	for i := 0; i < fldLen; i++ {
		key := fmt.Sprintf("arg %v", i)
		if val, ok := m.fields[fmt.Sprintf("arg %v", i)]; ok {
			argFields = append(argFields, messageField{key, val})
		}
	}
	return argFields
}

type messageField struct {
	key   string
	value interface{}
}

func NewField(key string, value interface{}) messageField {
	return messageField{key, value}
}

func isArgField(f messageField) bool {
	if !strings.HasPrefix(f.key, "arg ") {
		return false
	}
	num := strings.TrimPrefix(f.key, "arg ")
	if _, err := strconv.Atoi(num); err != nil {
		return false
	}
	return true
}

func fieldKeysMandatory() []string {
	return []string{
		keyTime,
		keyLevel,
		keyMessage,
		keyFile,
		keyLine,
		keyFunc,
	}
}

////////////////////////////////

func callerFunctionInfo() (string, int, string) {
	counter, file, line, success := runtime.Caller(2) //back to stack on 2 levels
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
