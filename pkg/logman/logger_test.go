package logman

import (
	"fmt"
	"testing"

	"github.com/Galdoba/ffstuff/pkg/logman/colorizer"
)

func TestLogMan(t *testing.T) {

	Setup(WithAppLogLevelImportance(ImportanceALL),
		WithGlobalColorizer(nil),

		WithLogLevels(
			NewLoggingLevel(INFO,
				WithWriter(Stderr, NewFormatter(WithRequestedFields(Request_ShortSince), WithColor(colorizer.DefaultScheme()))),
				WithWriter(`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\pkg\logman\v2\file.txt`, NewFormatter(WithRequestedFields(Request_Full))),
			),
			NewLoggingLevel(DEBUG,
				WithWriter(Stderr, NewFormatter(WithRequestedFields(Request_ShortReport), WithColor(colorizer.DefaultScheme()))),
				WithWriter(`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\pkg\logman\v2\file.txt`, NewFormatter(WithRequestedFields(Request_Medium))),
			),
		),
		WithGlobalWriterFormatter(`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\pkg\logman\v2\file.txt`, NewFormatter(WithRequestedFields(Request_ShortSince))),
	)
	//SetOutput("file.txt", ALL)
	msg := NewMessage("this is test messsage 2 with arg %v, which is float64 and '%v': is string", 3.14, "some string")
	fmt.Println("-------------------")
	fmt.Println("-------------------")
	ProcessMessage(msg, ERROR)
	fmt.Println("----")
	ProcessMessage(msg, WARN)
	fmt.Println("----")
	ProcessMessage(msg, INFO)
	fmt.Println("----")
	ProcessMessage(msg, DEBUG)
	fmt.Println("----")
	ProcessMessage(msg, TRACE)
	fmt.Println("----")
	ProcessMessage(msg, FATAL)
	fmt.Println("-------------------")

}
