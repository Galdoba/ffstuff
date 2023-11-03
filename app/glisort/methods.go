package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	Action_Reverse = " Reverse "
)

var method map[string]func([]string) ([]string, error)
var sortAction map[string]func(map[string][]string, string, string, []string, string) (map[string][]string, error)

func init() {
	sortAction = make(map[string]func(map[string][]string, string, string, []string, string) (map[string][]string, error))
	sortAction[Action_Reverse] = reverseAction
}

func reverseAction(data map[string][]string, source string, action string, args []string, out string) (map[string][]string, error) {
	input := data[source]
	output, err := reverse(input)
	if err != nil {
		return data, err
	}
	data[out] = output
	return data, nil
}

func reverse(sl []string) ([]string, error) {
	outSl := []string{}
	for i := len(sl) - 1; i >= 0; i-- {
		outSl = append(outSl, sl[i])
	}
	return outSl, nil
}

type sortMethod struct {
	Name          string              `json:"Name"`
	ListSeparator string              `json:"ListSeparator",omitempty`
	input         map[string][]string //tag:list
	Sequance      []string            `json:"Sequance"`
	actions       []func(map[string][]string, string, string, []string, string) (map[string][]string, error)
	output        []string
}

func NewSortMethod(name string) *sortMethod {
	sm := sortMethod{}
	sm.Name = name
	sm.input = make(map[string][]string)

	return &sm
}

func (sm *sortMethod) Compile(name string) error {
	dir := strings.TrimSuffix(configPath, "glisort.json")
	fls, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	data := []byte{}
	for _, f := range fls {
		if f.IsDir() {
			fmt.Println("skip dir: ", f.Name())
			continue
		}
		if f.Name() != name+".json" {
			fmt.Println("skip file: ", f.Name())
			continue
		}
		fmt.Println("read:", dir+f.Name())
		data, err = os.ReadFile(dir + f.Name())
		if err != nil {
			return err
		}
		if len(data) == 0 {
			return fmt.Errorf("no data recived from %v", f.Name())
		}
	}
	if len(data) == 0 {
		return fmt.Errorf("no data recived from %v", dir)
	}
	smTemp := &sortMethod{}
	err = json.Unmarshal(data, smTemp)
	if err != nil {
		return err
	}
	if sm.Name != smTemp.Name {

		return fmt.Errorf("expecting to compile %v, but have %v", sm.Name, smTemp.Name)
	}
	sm.ListSeparator = smTemp.ListSeparator
	sm.Sequance = smTemp.Sequance
	sm.constructActions()

	return nil
}

func (sm *sortMethod) Execute() error {
	for i, _ := range sm.Sequance {
		sorc := parseSource(sm.Sequance[i])
		actn := parseActions(sm.Sequance[i])
		args := parseArguments(sm.Sequance[i])
		trgt := parseTarget(sm.Sequance[i])

		output, err := sm.actions[i](sm.input, sorc, actn, args, trgt)
		if err != nil {
			return err
		}
		sm.input = output
	}
	return nil
}

/*
(source) action [argument]... {target}

*/

func (sm *sortMethod) constructActions() error {
	for _, action := range sm.Sequance {
		//sorc := parseSource(action)
		actn := parseActions(action)
		//args := parseArguments(action)
		//trgt := parseTarget(action)
		switch actn {
		case Action_Reverse:
			actFunc := sortAction[actn]
			sm.actions = append(sm.actions, actFunc)
		}
		//sm.Sequance = append(sm.Sequance)
		//fmt.Println(sorc, actn, args, trgt, sm.Sequance)
	}
	return nil
}

func parseSource(a string) string {
	if strings.HasPrefix(a, "(") {
		source := strings.TrimPrefix(a, "(")
		sParts := strings.Split(source, ") ")
		if len(sParts) == 1 {
			return ""
		}
		return sParts[0]
	}
	return ""
}

func parseActions(a string) string {
	if strings.Contains(a, Action_Reverse) {
		return Action_Reverse
	}
	return ""
}

func parseArguments(a string) []string {
	args := []string{}
	aParts := strings.Split(a, " [")
	for _, part := range aParts {
		if strings.HasSuffix(part, "] ") {
			args = append(args, strings.TrimSuffix(part, "] "))
		}
	}
	return args
}

func parseTarget(a string) string {
	sl := strings.Split(a, " {")
	if len(sl) < 1 {
		return ""
	}
	trgt := strings.Join(sl[len(sl)-1:], "")
	if !strings.HasSuffix(trgt, "}") {
		return ""
	}
	return strings.TrimSuffix(trgt, "}")
}

/*
	PartWhite (1) (2)==>
	Reverse (1)
*/
