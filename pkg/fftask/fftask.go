package fftask

import "fmt"

const (
	TASK_MUX = "Mux"
)

func validOperations() map[string]bool {
	vo := make(map[string]bool)
	for _, val := range []string{
		"Mux",
	} {
		vo[val] = true
	}
	return vo
}

//Task - описывает состояние задания и параметры задания
type Task struct {
	operation       string
	inputArgs       []string
	cmdOutput       string
	completed       bool
	preCheckErrors  []error
	postCheckErrors []error
}

//New - creates new Task
func New(operation string) (*Task, error) {
	tsk := Task{}
	err := fmt.Errorf("error was not addressed")
	validOperations := validOperations()
	if validOperations[operation] {
		err = nil
		tsk.operation = operation
	}
	return &Task{}, err
}
