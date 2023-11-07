package main

import (
	"fmt"
	"regexp"
)

const (
	Action_Reverse    = "REVERSE"
	Action_MergeLists = "MERGE_LISTS"
	Action_Regexp     = "REGEXP"
	Action_Clean      = "CLEAN"

/*
IF_CONTAINS_TRUE
IF_CONTAINS_FALSE
IF_EQUAL_TRUE
IF_EQUAL_FALSE
REGEXP
APPEND_TEXT
PREPEND_TEXT
REPLACE_TEXT (a) (b)
*/

)

func reverseAction(data map[string][]string, seq *sequanceData) (map[string][]string, error) {
	input := data[seq.inputList]
	output, err := reverse(input)
	if err != nil {
		return data, err
	}
	data[seq.outputList] = output
	fmt.Println(data)
	return data, nil
}

func mergeListsAction(data map[string][]string, seq *sequanceData) (map[string][]string, error) {
	input := data[seq.inputList]
	argLists := [][]string{}
	for _, arg := range seq.args {
		argList := data[arg]
		argLists = append(argLists, argList)
	}
	output := input
	for _, argList := range argLists {
		for _, arg := range argList {
			output = append(output, arg)
		}
	}
	data[seq.outputList] = output
	return data, nil
}

func regexpAction(data map[string][]string, seq *sequanceData) (map[string][]string, error) {
	input := data[seq.inputList]
	//argLists := [][]string{}
	output := input
	for _, arg := range seq.args {
		fmt.Println("arg:", arg)
		rg, err := regexp.Compile(arg)
		if err != nil {
			return data, fmt.Errorf("regexp %v can't complile: %v", err.Error())
		}
		for i, out := range output {
			fmt.Println(i, out, "|||")
			output[i] = rg.FindString(out)
			fmt.Println(output[i], "===")
		}
	}

	data[seq.outputList] = output
	return data, nil
}

func cleanAction(data map[string][]string, seq *sequanceData) (map[string][]string, error) {
	input := data[seq.inputList]
	//argLists := [][]string{}
	output := []string{}
	for _, inp := range input {
		if inp != "" {
			output = append(output, inp)
		}
	}
	data[seq.outputList] = output
	return data, nil
}
