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
	fmt.Println("|", tsk, "|")
	output, err := xml.MarshalIndent(tsk, "", "  ")
	_ = ioutil.WriteFile("testTask.xml", output, 0644)
	fmt.Println(err)
	fmt.Println(tsk)
	fmt.Println(output)
	os.Stdout.Write(output)
	f, _ := ioutil.ReadFile("testTask.xml")
	tsk2 := Task{}
	errr := xml.Unmarshal(f, &tsk2)
	fmt.Println(errr)
	fmt.Println("---------------")
	fmt.Println(tsk2)
	fmt.Println(tsk2.Sender)
	fmt.Println(tsk2.Receiver)
	fmt.Println("---------------")
	fmt.Println(tsk)
	fmt.Println(tsk.Sender)
	fmt.Println(tsk.Receiver)
}
