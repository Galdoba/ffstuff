package logman

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type message struct {
	fields map[string]interface{}
}

func NewMessage(format string, args ...interface{}) *message {
	m := message{}
	m.fields = make(map[string]interface{})
	timeCreated := time.Now()
	m.fields[keyMessage] = fmt.Sprintf(format, args...)
	if logMan.colorizer != nil {
		coloredArgs := []string{}
		for _, arg := range args {
			coloredArgs = append(coloredArgs, fmt.Sprintf("%v", logMan.colorizer.Colorize(arg)))
			fmt.Println(logMan.colorizer.Colorize(arg))
		}

		m.fields[keyMessageColor] = combineColored(format, coloredArgs...)
		fmt.Println(m.fields[keyMessageColor])
	}
	m.fields[keyTime] = timeCreated.Format(time.RFC3339Nano)
	return &m
}

func combineColored(format string, args ...string) string {
	fmtParts := strings.Split(format, `%v`)
	combined := ""
	for i, part := range fmtParts {
		combined += part
		if i < len(args) {
			combined += args[i]
		}
	}
	return combined
}

// Message is an interface to a struct to set/get/list data fields.
type Message interface {
	Fields() []string
	Value(string) interface{}
	SetField(string, interface{})
}

// Value return variable that serves as state if field.
func (m *message) Value(key string) interface{} {
	if val, ok := m.fields[key]; ok {
		return val
	}
	return nil
}

// Fields return sorted list of keys for contained fields.
func (m *message) Fields() []string {
	keys := []string{}
	for k := range m.fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// SetField - sets/override fields value.
func (m *message) SetField(key string, value interface{}) {
	m.fields[key] = value
}

// WithFields sets multiple fields to a message.
func (m *message) WithFields(flds ...messageField) *message {
	for _, fld := range flds {
		m.fields[fld.key] = fld.value
	}
	return m
}

// WithArgs overrided special type of fields with keys "arg №", where № is number in order of appearence.
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

func (m *message) clearLevel() {
	delete(m.fields, keyLevel)
}

type messageField struct {
	key   string
	value interface{}
}

// NewField created messageField. Used as argument in message.WithFields().
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

// // MarshalJSON return JSON encoding of message and
// func (m *message) MarshalJSON() ([]byte, error) {
// 	j, err := json.Marshal(struct {
// 		Fields map[string]interface{}
// 	}{
// 		Fields: m.fields,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return j, nil
// }

// MarshalJSON return JSON encoding of Message.
func MarshalJSON(m Message) ([]byte, error) {
	fields := make(map[string]interface{})
	keys := m.Fields()
	for _, k := range keys {
		fields[k] = m.Value(k)
	}
	bt, err := json.Marshal(struct {
		Fields map[string]interface{}
	}{
		Fields: fields,
	})
	if err != nil {
		return nil, err
	}
	return bt, nil
}

func (m *message) unmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	// Alias for message (temp struct)
	var realMessage struct {
		Fields map[string]interface{} `json:"Fields"`
	}
	// Unmarshal the json into the realMessage.
	if err := json.Unmarshal(data, &realMessage); err != nil {
		return err
	}
	// Set the fields to the new struct,
	m.fields = realMessage.Fields
	return nil
}
