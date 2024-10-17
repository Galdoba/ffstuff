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
	//msg.SetField(keyLevel, "report")
	fmt.Println(msg)
	fmt.Println("----")
	fe := NewFE([]string{"time", keyLevel, keyMessage})
	fe.AddFormatterFunc("time", stdFormatFunc_time)
	fe.AddFormatterFunc("message", stdMessageFormat)
	formated := fe.Format(msg)
	fmt.Println("----")
	fmt.Println(formated)
	fmt.Println("----")
	err := process(msg, logMan.logLevels[WARN])
	if err != nil {
		t.Errorf(err.Error())
	}

}

func formatTestFunc(msg Message) (string, error) {
	out := "this is a TEST:"
	for _, key := range msg.Fields() {
		out += fmt.Sprintf("|field='%v' : value='%v'", key, msg.Value(key))
	}
	return out, nil
}
