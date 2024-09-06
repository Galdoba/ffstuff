package logman

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
)

type formatter struct {
	writerKey string
	color     bool
}

func formatTextComplex(msg Message) (string, error) {
	fldKeys := []string{}
	for _, key := range msg.Fields() {
		switch key {
		case keyTime, keyLevel, keyFile, keyLine, keyFunc, keyMessage:
		default:
			fldKeys = append(fldKeys, key)
		}
	}
	s := ""
	tm, err := time.Parse(time.RFC3339Nano, fmt.Sprintf("%v", msg.Value(keyTime)))
	if err != nil {
		return "", fmt.Errorf("time parsing failed: %v", err)
	}
	s += timestampFormat(tm)

	for i, val := range []interface{}{
		msg.Value(keyLevel),
		msg.Value(keyMessage),
		msg.Value(keyFile),
		msg.Value(keyLine),
		msg.Value(keyFunc),
	} {
		if val != nil {
			switch i {
			default:
				s += " " + fmt.Sprintf("%v", val)
			case 0:
				s += " " + fmt.Sprintf("%v:", val)
			case 1:
				s += " " + fmt.Sprintf("%v", val)
			case 2:
				s += " " + fmt.Sprintf("\n  caller: file={%v}", val)
			case 3:
				s += " " + fmt.Sprintf("line={%v}", val)
			case 4:
				s += " " + fmt.Sprintf("func={%v}", val)

			}
		}
	}

	if len(fldKeys) != 0 {
		s += "\n  fields: "
		for _, key := range fldKeys {
			val := msg.Value(key)
			s += fmt.Sprintf("{'%v'='%v' %T}; ", key, val, val)
		}
		s = strings.TrimSuffix(s, "; ")
	}

	return s, nil
}

func formatTextSimple(msg Message) (string, error) {
	s := ""
	tm, err := time.Parse(time.RFC3339Nano, fmt.Sprintf("%v", msg.Value(keyTime)))
	if err != nil {
		return "", fmt.Errorf("time parsing failed: %v", err)
	}
	s += timestampFormat(tm)

	for i, val := range []interface{}{
		msg.Value(keyLevel),
		msg.Value(keyMessage),
	} {
		if val != nil {
			switch i {
			default:
				s += " " + fmt.Sprintf("%v:", val)
			case 1:
				s += " " + fmt.Sprintf(`%v`, val)
			}
		}
	}
	return s, nil
}

func formatJSON(msg Message) (string, error) {
	s := `{"Fields":`
	keys := msg.Fields()
	switch len(keys) {
	default:
		s += "{"
		for _, key := range keys {
			switch val := msg.Value(key).(type) {
			case string:
				fmt.Println("string", val)
				s += fmt.Sprintf(`"%v":"%v",`, key, val)
			default:
				s += fmt.Sprintf(`"%v":%v,`, key, val)
			}

		}
		s = strings.TrimSuffix(s, ",") + "}"

	case 0:
		s += `null`
	}
	s += `}`
	return s, nil
}

func formatPing(msg Message) (string, error) {
	fldKeys := []string{}
	for _, key := range msg.Fields() {
		switch key {
		case keyTime, keyLevel, keyFile, keyLine, keyFunc, keyMessage:
		default:
			fldKeys = append(fldKeys, key)
		}
	}
	s := ""
	for i, val := range []interface{}{
		msg.Value(keyFile),
		msg.Value(keyLine),
		msg.Value(keyFunc),
	} {
		if val != nil {
			switch i {
			default:
				//s += " " + fmt.Sprintf("%v:", val)
			case 0:
				file := filepath.Base(fmt.Sprintf(`%v`, val))
				s += file + " "
			case 1:
				s += fmt.Sprintf(`%v`, val) + " "
			case 2:
				fun := filepath.Base(fmt.Sprintf(`%v`, val))
				s += fun + " "
			}
		}
	}
	for _, k := range fldKeys {
		s += fmt.Sprintf("%v", msg.Value(k)) + " "
	}
	s = strings.TrimSuffix(s, " ")
	s = color.HiBlackString(s)
	return s, nil
}
