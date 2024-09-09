package reportfile

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Report struct {
	Fields map[string]string `json:"Stats"`
}

type Field struct {
	Key   string
	Value string
}

func NewField(key string, info string) Field {
	return Field{key, info}
}

func New(flds ...Field) *Report {
	rep := Report{}
	rep.Fields = make(map[string]string)
	rep.Fields["Created At"] = time.Now().String()
	for _, fld := range flds {
		rep.Fields[fld.Key] = fld.Value
	}
	return &rep
}

func (rep *Report) AddFields(flds ...Field) error {
	overwriten := []string{}
	for _, fld := range flds {
		if _, ok := rep.Fields[fld.Key]; ok {
			overwriten = append(overwriten, fld.Key)
		}
		rep.Fields[fld.Key] = fld.Value
	}
	switch len(overwriten) {
	default:
		return fmt.Errorf("fields %v were overwritten", strings.Join(overwriten, ", "))
	case 1:
		return fmt.Errorf("field %v was overwritten", overwriten[0])
	case 0:
		return nil
	}
}

func (rep *Report) Marshal() ([]byte, error) {
	return json.MarshalIndent(rep, "", "  ")
}

func (rep *Report) Unmarshal(bt []byte) error {
	return json.Unmarshal(bt, rep)
}

func (rep *Report) String() string {
	bt, err := rep.Marshal()
	if err != nil {
		return "unmarshaling err"
	}
	return string(bt)
}

func (rep *Report) CreateFile(path string) error {
	f, err := os.Create(path)
	f.Truncate(0)
	if err != nil {
		return fmt.Errorf("create report file failed: %v", err)
	}
	_, err = f.WriteString(rep.String())
	if err != nil {
		return fmt.Errorf("write report file failed: %v", err)
	}
	f.Close()
	return nil
}

func ReadFile(path string) (*Report, error) {
	bt, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read report file failed: %v", err)
	}
	rep := &Report{}
	err = rep.Unmarshal(bt)
	if err != nil {
		return nil, fmt.Errorf("unmarshal report file failed: %v", err)
	}
	return rep, nil
}

func (rep *Report) Find(key string) Field {
	if val, ok := rep.Fields[key]; ok {
		return NewField(key, val)
	}
	return Field{}
}

func (rep *Report) FindAll(partOfKey string) []Field {
	found := []Field{}
	for k, v := range rep.Fields {
		if strings.Contains(k, partOfKey) {
			found = append(found, NewField(k, v))
		}
	}
	return found
}
