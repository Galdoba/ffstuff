package fftask

import (
	"fmt"
	"testing"
)

func TestTask(t *testing.T) {
	testOperations := validOperations()
	for operation, validity := range testOperations {
		fmt.Printf("START testing operation '%v':\n", operation)
		tsk, err := New(operation)
		if err != nil {
			t.Errorf("unexpected error: %v", err.Error())
		}
		fmt.Println(validOperations(), operation)
		if validity == false {
			t.Errorf("operation %v is not valid or not implemented", tsk.operation)
		}
		fmt.Printf("END   testing operation '%v':\n", operation)
	}
}
