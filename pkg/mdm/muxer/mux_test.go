package muxer

// import (
// 	"fmt"
// 	"testing"
// )

// func emulatedInstructions() []Task {
// 	instructions := []string{
// 		"",
// 	}
// 	instructions = append(instructions, validInstructions()...)
// 	video := []string{
// 		"file.mp4", "file.mpeg",
// 		"",
// 	}
// 	audio := []string{
// 		"file_rus20.ac3", "file_rus51.ac3",
// 		"file_eng20.ac3", "file_eng51.ac3",
// 		"file_qqq20.ac3", "file_qqq51.ac3",
// 		"", "file.aac", "file.m4a",
// 		"file_rus51.aac",
// 	}
// 	subs := []string{
// 		"file.srt",
// 		"",
// 	}
// 	var taskList []Task
// 	for _, inst := range instructions {
// 		for _, vid := range video {
// 			for _, aud1 := range audio {
// 				for _, aud2 := range audio {
// 					for _, srt := range subs {
// 						taskList = append(taskList, Task{"", inst, vid, aud1, aud2, srt, nil, "", "", ""})
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return taskList
// }

// func TestMuxTask(t *testing.T) {
// 	taskList := emulatedInstructions()
// 	//alltests := len(taskList)
// 	undecided := 0
// 	for i, task := range taskList {
// 		task.Validate()
// 		if task.err != nil {
// 			if task.err.Error() == "Undecided" {
// 				fmt.Printf("Task %v undecided : %v\n", i, task)
// 				t.Errorf("Task %v (%v): error: %v\n", i, task, task.err.Error())
// 				undecided++
// 				continue
// 			}
// 			continue
// 		}
// 		//fmt.Printf("Test %v of %v: task = %v - VALID \n", i+1, alltests, task)
// 	}
// 	//fmt.Printf("Undesided Tests %v of %v\n", undecided, alltests)
// }

// func emulatedTXTInstructions() []string {
// 	return []string{
// 		"somefile.mp4 ar2e6_sr",
// 		"somefile.mp4 ar2e5_sr",
// 		"somefile.mp4 ar2e6sr",
// 	}
// }

// func TestTaskCreation(t *testing.T) {
// 	for i, val := range emulatedTXTInstructions() {
// 		task := NewTask("disk:\\path\\dir", val)
// 		if task.err != nil {
// 			t.Errorf("Task %v (%v): error: %v\n", i, task, task.err.Error())
// 		}
// 	}
// }
