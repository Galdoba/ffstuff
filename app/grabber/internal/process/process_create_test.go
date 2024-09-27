package process

import (
	"fmt"
	"testing"

	"github.com/Galdoba/ffstuff/pkg/logman"
)

func TestNewProcess(t *testing.T) {
	logman.Setup()
	pr, err := New(SetMode("MODE_GRAB"))
	fmt.Println(pr)
	fmt.Println(err)
}
