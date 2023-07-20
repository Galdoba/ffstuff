package task

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"testing"
)

func TestTask(t *testing.T) {
	usr, _ := user.Current()
	userName := usr.Name
	tsk, err := NewTask(
		Input(KEY_deadline, "2023-06-18 13:34"),
		Input(KEY_receiver, "ffautoBatch"),
		Input(KEY_title, "Write Todosher"),
		Input(KEY_descr, "ffmpeg -i test.mp4"),
		Input(KEY_sender, userName),
	)
	output, err := xml.MarshalIndent(tsk, "", "  ")
	_ = ioutil.WriteFile("testTask.xml", output, 0644)
	fmt.Println(err)
	fmt.Println(tsk)
	fmt.Println(output)
	os.Stdout.Write(output)

}
