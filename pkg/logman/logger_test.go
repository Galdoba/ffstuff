package logman

import (
	"fmt"
	"testing"

	"github.com/Galdoba/ffstuff/pkg/logman/colorizer"
)

func TestLogMan(t *testing.T) {
	//	logMan.appName = "scribe"
	//jsonFormatter := NewFormatter(WithCustomFunc("json", stdJSON), WithRequestedFields([]string{"json"}))

	Setup(WithAppLogLevelImportance(ImportanceALL),
		WithGlobalColorizer(colorizer.DefaultScheme()),
		// WithLogLevels(
		// 	NewLoggingLevel(WARN,
		// 		WithWriter(Stderr, NewFormatter(WithRequestedFields(Request_ShortSince), WithColor(colorizer.DefaultScheme()))),
		// 		WithWriter(`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\pkg\logman\v2\file.txt`, NewFormatter(WithRequestedFields(Request_ShortTime))),
		// 	),
		// 	NewLoggingLevel(INFO,
		// 		WithWriter(Stderr, NewFormatter(WithRequestedFields(Request_ShortSince), WithColor(colorizer.DefaultScheme()))),
		// 		WithWriter(`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\pkg\logman\v2\file.txt`, NewFormatter(WithRequestedFields(Request_ShortSince))),
		// 	),
		// 	NewLoggingLevel(DEBUG,
		// 		WithWriter(Stderr, NewFormatter(WithRequestedFields(Request_ShortReport), WithColor(colorizer.DefaultScheme()))),
		// 		WithWriter(`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\pkg\logman\v2\file.txt`, NewFormatter(WithRequestedFields(Request_Medium))),
		// 	),
		// ),
		//WithGlobalWriterFormatter(Stderr, NewFormatter(WithRequestedFields(Request_ShortSince))),
		//WithGlobalWriterFormatter(Stdout, NewFormatter(WithRequestedFields(Request_ShortReport))),
		WithGlobalWriterFormatter(`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\pkg\logman\v2\file.txt`, NewFormatter(WithRequestedFields(Request_ShortTime))),
		//WithGlobalWriterFormatter(`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\pkg\logman\v2\`, jsonFormatter),
		WithJSON(`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\pkg\logman\v2\`),
		WithAppName("scribe_test"),
	)
	//SetOutput("file.txt", ALL)
	msg := NewMessage("this is test messsage 2 with arg %v, which is float64 and '%v': is string", 3.14, "some string")
	fmt.Println("-------------------")
	fmt.Println("-------------------")
	ProcessMessage(msg, ERROR)
	fmt.Println("----W")
	ProcessMessage(msg, WARN)
	fmt.Println("----I")
	ProcessMessage(msg, INFO)
	fmt.Println("----D")
	ProcessMessage(msg, DEBUG)
	fmt.Println("----")
	ProcessMessage(msg, TRACE)
	fmt.Println("----")
	ProcessMessage(msg, FATAL)
	fmt.Println("-------------------")

}
