package muxer

import (
	"fmt"
	"testing"
)

func instructionConstructor() []string {
	prefix := "a"
	l1 := []string{"r", "e", "qqq"}
	c1 := []string{"2", "6"}
	l2 := []string{"e", "qqq"}
	c2 := []string{"2", "6"}
	postfix := []string{"", "_sr"}
	var instructions []string
	for _, a := range l1 {
		for _, b := range c1 {
			for _, c := range l2 {
				for _, d := range c2 {
					for _, e := range postfix {
						instructions = append(instructions, prefix+a+b+c+d+e)
					}
				}
			}
		}
	}
	return instructions
}

func emulatedTasks() []Task {
	instructions := []string{
		"",
	}
	instructions = append(instructions, instructionConstructor()...)
	video := []string{
		"file.mp4", "file.mpeg",
		"",
	}
	audio := []string{
		"file_rus20.ac3", "file_rus51.ac3",
		"file_eng20.ac3", "file_eng51.ac3",
		"file_qqq20.ac3", "file_qqq51.ac3",
		"", "file.aac", "file.m4a",
	}
	subs := []string{
		"file.srt",
		"",
	}
	var taskList []Task
	for _, inst := range instructions {
		for _, vid := range video {
			for _, aud1 := range audio {
				for _, aud2 := range audio {
					for _, srt := range subs {
						taskList = append(taskList, Task{inst, vid, aud1, aud2, srt, nil})
					}
				}
			}
		}
	}
	return taskList
}

func TestMuxTask(t *testing.T) {
	taskList := emulatedTasks()
	alltests := len(taskList)
	valid := 0
	for i, task := range taskList {
		//fmt.Printf("Test %v: task = %v \n", i+1, task)
		task.Validate()
		if task.err != nil {
			valid++
			continue
		}
		fmt.Printf("Test %v of %v: task = %v \n", i+1, alltests, task)
		//fmt.Printf("Task invalid. reason: %v\n", task.err.Error())
	}
	fmt.Printf("Undesided Tests %v of %v\n", alltests-valid, alltests)
}
