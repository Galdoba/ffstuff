package dataconnection

import (
	"fmt"
	"testing"
)

func TestCommonPrefix(t *testing.T) {
	inp1 := []string{
		"kino1--FILM--name1.mov",
		"trailer1--FILM--name1.mov",
		"kino1--TRL--name1.mov",
	}
	inp2 := []string{
		"kin5o1--FILM--name1.mov",
		"trai1ler1--FILM--name1.mov",
		"kino12--TRL--name1.mov",
	}
	for i := 0; i < len(inp1); i++ {
		fmt.Println(commonPrefix(inp1[i], inp2[i]))
	}

}
