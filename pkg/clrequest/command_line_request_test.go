package clrequest

import (
	"fmt"
	"testing"

	"github.com/Galdoba/ffstuff/pkg/clrequest/bat"
)

func Test_request_Bat(t *testing.T) {
	req := New("echo", "i", "am", "echo")
	err := req.Bat(`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\_template\echoer.bat`, bat.NoEcho())
	fmt.Println(err)
}
