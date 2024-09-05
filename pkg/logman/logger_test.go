package logman

import (
	"fmt"
	"testing"
)

func TestLogMan(t *testing.T) {

	Setup(WithAppLogLevelImportance(ImportanceALL))

	SetOutput("file.txt", ALL)

	msg := NewMessage("process %v complete", "testing").WithFields(NewField("metric", 0.7))
	err := process(msg, logMan.logLevels[INFO])
	if err != nil {
		t.Errorf(err.Error())
	}
	err2 := process(msg, logMan.logLevels[WARN])
	if err2 != nil {
		t.Errorf(err2.Error())
	}
	Debug(NewMessage("testing debug"), map[string]interface{}{"arg1": 5, "arg2": 7})
	Fatalf("testing fatal conv func")

	// bt, _ := msg.MarshalJSON()
	// msg1 := NewMessage("sss")
	// if err := msg1.UnmarshalJSON(bt); err != nil {
	// 	t.Errorf("bad %v", err)
	// }
	// fmt.Println(string(bt))
	// fmt.Println("1", msg)
	// fmt.Println("2", msg1)
	// js, err := formatJSON(msg)
	// fmt.Println(js, err)

}

func formatTestFunc(msg Message) (string, error) {
	out := "this is a TEST:"
	for _, key := range msg.Fields() {
		out += fmt.Sprintf("|field='%v' : value='%v'", key, msg.Value(key))
	}
	return out, nil
}
