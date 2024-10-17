package logman

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Galdoba/ffstuff/pkg/logman/v2/colorizer"
)

type formatterExpanded struct {
	//formatFunc        func(Message) string
	fieldFormaFuncMap map[string]func(Message) (string, error)
	requestedFields   []string
	ignoredFields     []string
}

func NewFE(requestedFields []string) *formatterExpanded {
	fe := formatterExpanded{}
	fe.fieldFormaFuncMap = make(map[string]func(Message) (string, error))
	fe.requestedFields = requestedFields
	return &fe
}

func basicFormatter(field string, val interface{}) (string, error) {
	format := "%v=%v"
	return fmt.Sprintf(format, field, val), nil
}

func (fe *formatterExpanded) Format(msg Message) string {
	output := ""
	for _, field := range fe.requestedFields {
		if mustIgnore(fe.ignoredFields, field) {
			continue
		}
		val := msg.Value(field)
		formatted := ""
		err := errors.New("not formatted")
		fn := fe.fieldFormaFuncMap[field]
		switch fn {
		case nil:
			formatted, err = basicFormatter(field, val)
			if err != nil {
				return output + formatted
			}
		default:
			formatted, err = fn(msg)
			if err != nil {
				return output + formatted + "!<> " + err.Error()
			}
		}
		output += formatted + " "

	}
	return output
}

func mustIgnore(fields []string, field string) bool {
	for _, f := range fields {
		if f == field {
			return true
		}
	}
	return false
}

func (fe *formatterExpanded) AddFormatterFunc(field string, fn func(Message) (string, error)) {
	fe.fieldFormaFuncMap[field] = fn
}

func stdFormatFunc_time(msg Message) (string, error) {
	tm, err := validateTimeArg(msg.Value("time"))
	if err != nil {
		return "", err
	}
	return tm.Format("[06/01/02 15:04:05.999]"), nil
}

func stdFormatFunc_timeSince(args ...any) (string, error) {
	tm, err := validateTimeArg(args)
	if err != nil {
		return "", err
	}
	duration := time.Since(tm)
	return fmt.Sprintf("[%v]", float64(duration.Milliseconds())/1000), nil
}

func validateTimeArg(args ...any) (time.Time, error) {
	if len(args) != 1 {
		return time.Time{}, fmt.Errorf("stdTimeFormat function expect 1 argument (have %v)", len(args))
	}
	val := args[0]
	str := fmt.Sprintf("%v", val)
	str = strings.TrimPrefix(str, "[")
	str = strings.TrimSuffix(str, "]")
	tm, err := time.Parse(time.RFC3339Nano, fmt.Sprintf("%v", str))
	if err != nil {
		return time.Time{}, err
	}
	return tm, nil
}

func stdMessageFormat(msg Message) (string, error) {
	colors := logMan.colorizer
	inputs := msg.InputArgs()
	var format string
	var args []interface{}
	for i := -1; i < len(inputs)-1; i++ {
		switch i {
		case -1:
			format = fmt.Sprintf("%v", inputs[-1])
		default:
			args = append(args, inputs[i])
		}
	}
	switch colors {
	case nil:
		return fmt.Sprintf(format, args...), nil
	default:
		var coloredArgs []string
		for _, arg := range args {
			colored := colors.ColorizeByType(arg)
			coloredArgs = append(coloredArgs, colored)
		}
		text := combineColored(format, coloredArgs...)
		return fmt.Sprintf(text), nil
	}
}

func stdLevelFormat(msg Message) (string, error) {
	colors := logMan.colorizer
	level := msg.Value(keyLevel)
	if level == nil {
		return "", errNoField(keyLevel)
	}
	tag := fmt.Sprintf("[%v]", level)
	switch colors {
	case nil:
		return tag, nil
	default:
		keyFg := colorizer.NewKey(colorizer.FG_KEY, fmt.Sprintf("%v", level))
		keyBg := colorizer.NewKey(colorizer.BG_KEY, fmt.Sprintf("%v", level))
		return fmt.Sprintf("%v", colors.ColorizeByKeys(tag, keyFg, keyBg)), nil
	}
}

func errNoField(field string) error {
	return fmt.Errorf("no field with key '%v'", field)
}

func stdCallerLongFormat(msg Message) (string, error) {
	colors := logMan.colorizer
	file := msg.Value(keyFile)
	line := msg.Value(keyLine)
	funk := msg.Value(keyFunc)
	if file == nil {
		return "", errNoField(keyFile)
	}
	if line == nil {
		return "", errNoField(keyLine)
	}
	sep := string(filepath.Separator)
	pathParts := strings.Split(fmt.Sprint(file), fmt.Sprintf("%vsrc%v", sep, sep))
	switch len(pathParts) {
	case 2:
		file = pathParts[1]
	}
	text := fmt.Sprintf("\n  [caller=%v:%v]", file, line)
	if funk != nil {
		text += fmt.Sprintf(" [func=%v]", funk)
	}

	switch colors {
	case nil:
		return fmt.Sprintf("%v", text), nil
	default:
		keyFg := colorizer.NewKey(colorizer.FG_KEY, "caller")
		keyBg := colorizer.NewKey(colorizer.BG_KEY, "caller")
		return fmt.Sprintf("%v", colors.ColorizeByKeys(text, keyFg, keyBg)), nil
	}
}
