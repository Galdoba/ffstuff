package configbuilder

import (
	"fmt"
	"os"
	"testing"
)

func TestSourcePath(t *testing.T) {
	cb := &configBuilder{}
	err := cb.SetSourceDir(`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\loudrec\config\`)
	if err != nil {
		t.Errorf("met error: %v", err.Error())
	}
	fmt.Println(cb.sourcePath)
	fmt.Println(cb.app)
	fmt.Println(cb.pathToConfig)

	cb.AddField(NewField("ShowLogs", "bool", "Show Logs of file Generation", true))
	cb.AddField(NewField("HumanOnly", "bool", "Character is human only", true))
	cb.AddField(NewField("OpenCareers", "map[string]bool", "Careers Paths Allowed", true))

	/*
		VideoPaths map[string]string `yaml:"Video Paths    ,omitempty"`
		AudioPaths map[string]string `yaml:"Audio Paths    ,omitempty"`
		Subtitle   string            `yaml:"Subtitle Path  ,omitempty"`
		Storage    string            `yaml:"Storage Dir    ,omitempty"`
	*/

	for _, fields := range cb.fields {
		fmt.Println(fields)
	}
	text := cb.GenerateSource()
	f, _ := os.OpenFile(cb.sourcePath, os.O_CREATE|os.O_WRONLY, 0777)
	defer f.Close()
	f.WriteString(text)

}
