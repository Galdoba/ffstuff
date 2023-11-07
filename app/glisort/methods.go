package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var method map[string]func([]string) ([]string, error)
var sortAction map[string]func(map[string][]string, *sequanceData) (map[string][]string, error)
var definedActions = []string{
	Action_Reverse,
	Action_MergeLists,
	Action_Regexp,
	Action_Clean,
}

func init() {
	sortAction = make(map[string]func(map[string][]string, *sequanceData) (map[string][]string, error))
	sortAction[Action_Reverse] = reverseAction
	sortAction[Action_MergeLists] = mergeListsAction
	sortAction[Action_Regexp] = regexpAction
	sortAction[Action_Clean] = cleanAction
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
	actions       []func(map[string][]string, *sequanceData) (map[string][]string, error)
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
		seqData, err := parseSequance(sm.Sequance[i])
		if err != nil {
			return err
		}
		output, err := sm.actions[i](sm.input, seqData)
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
		actn := parseActions(action)
		for _, action := range definedActions {
			if action == actn {
				actFunc := sortAction[actn]
				sm.actions = append(sm.actions, actFunc)
				continue
			}
		}
	}
	return nil
}

// func parseSource(a string) string {
// 	if strings.HasPrefix(a, "(") {
// 		source := strings.TrimPrefix(a, "(")
// 		sParts := strings.Split(source, ") ")
// 		if len(sParts) == 1 {
// 			return ""
// 		}
// 		return sParts[0]
// 	}
// 	return ""
// }

func parseActions(a string) string {
	for _, action := range definedActions {
		if strings.Contains(a, action) {
			return action
		}
	}
	panic(fmt.Sprintf("action undefined: -%v-", a))
	return ""
}

// func parseArguments(a string) []string {
// 	argLine := strings.ReplaceAll(a, " [", "$$$$$$")
// 	argLine = strings.ReplaceAll(argLine, "]", "$$$$$$")
// 	args := strings.Split(argLine, "$$$")
// 	out := []string{}
// 	for i, ar := range args {
// 		switch i {
// 		case 0, len(args) - 1:
// 			continue
// 		default:
// 		}
// 		if ar != "" {
// 			out = append(out, ar)
// 		}
// 	}
// 	return out
// }

// func parseTarget(a string) string {
// 	sl := strings.Split(a, " {")
// 	if len(sl) < 1 {
// 		return ""
// 	}
// 	trgt := strings.Join(sl[len(sl)-1:], "")
// 	if !strings.HasSuffix(trgt, "}") {
// 		return ""
// 	}
// 	return strings.TrimSuffix(trgt, "}")
// }

//1 ==> REGEXP aaa "bbb ccc" ddd ==> paths
//1 ==> REVERSE ==> 2
type sequanceData struct {
	inputList  string
	action     string
	args       []string
	outputList string
}

func parseSequance(seq string) (*sequanceData, error) {
	sq := sequanceData{}
	words := strings.Split(seq, " ==>")
	seqParts := []string{}
	for _, w := range words {
		if w != "" {
			seqParts = append(seqParts, w)
		}
	}
	if len(seqParts) != 3 {
		return nil, fmt.Errorf("can't parse sequance [%v]:\n expecting 3 part structure with ' ==>' delimeters, have %v", seq, len(seqParts))
	}
	sq.inputList = strings.TrimSpace(words[0])
	sq.outputList = strings.TrimSpace(words[2])
	sq.action = parseActions(seqParts[1])
	argsStr := strings.Replace(seqParts[1], sq.action, "", 1)
	argsStr = strings.TrimSpace(argsStr)
	switch sq.action {
	case Action_Regexp:
		fmt.Printf("argStr |%v|\n", argsStr)
	}

	args := strings.Split(argsStr, " ")
	clearedArgs := []string{}
	currentArg := ""
	closed := true
	for _, arg := range args {
		if strings.HasPrefix(arg, `"`) || strings.HasPrefix(arg, "`") {
			closed = false
		}
		if strings.HasSuffix(arg, `"`) || strings.HasSuffix(arg, "`") {
			closed = true
		}
		switch currentArg {
		case "":
			currentArg = arg
		default:
			currentArg += " " + arg
		}
		if closed {
			clearedArgs = append(clearedArgs, currentArg)
			currentArg = ""
		}
	}
	sq.args = clearedArgs
	return &sq, nil
}
