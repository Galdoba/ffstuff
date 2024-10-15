package logman

import (
	"fmt"
	"testing"

	"github.com/Galdoba/ffstuff/pkg/logman/v2/colorizer"
)

func TestLogMan(t *testing.T) {

	Setup(WithAppLogLevelImportance(ImportanceALL), WithColorizer(colorizer.DefaultScheme()))

	SetOutput("file.txt", ALL)

	msg := NewMessage("process %v complete at %v with pamcake", "testing", 2000)
	fmt.Println(msg)
	err := process(msg, logMan.logLevels[INFO])
	if err != nil {
		t.Errorf(err.Error())
	}
	err2 := process(msg, logMan.logLevels[WARN])
	if err2 != nil {
		t.Errorf(err2.Error())
	}
	Debug(NewMessage("testing debug"), "test: 42")
	Printf("this is message with %v of type string", "argument")

	//Fatalf("testing fatal conv func")

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
