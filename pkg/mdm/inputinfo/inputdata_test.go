package inputinfo

import (
	"fmt"
	"testing"

	"github.com/Galdoba/utils"
)

func TestInputReading(t *testing.T) {
	examples := gatherInfo()
	for i, input := range examples {
		mf := CollectData(input)
		name := input.parseName()
		if name == "" {
			t.Errorf("expecting name to be foudn in: \n%v", input.String())
		}
		fmt.Println(i, mf)
	}
}

func gatherInfo() []inputdata {
	allData := []inputdata{}
	data := inputdata{}
	for _, line := range utils.LinesFromTXT(`C:\Users\a.akkulov\Desktop\cmdLine.go`) {
		if line == `/*` {
			data = inputdata{}
			continue
		}
		if line == `*/` {
			allData = append(allData, data)

			continue
		}
		data.data = append(data.data, line)
	}
	fmt.Println(len(allData), "examples found")
	return allData
}
