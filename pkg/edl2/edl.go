package edl2

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/macroblock/imed/pkg/types"
)

const (
	titleID           = "TITLE:"
	fcmID             = "FCM:"
	fcmModeDF         = "DROP FRAME"
	fcmModeNDF        = "NON-DROP FRAME"
	standardStatement = "STANDARD"
)

type edlData struct {
	edlSource string //источник самого edl (не обязательно файл)
	decidions []Decidion
}

type timeSegment struct {
	in     types.Timecode
	out    types.Timecode
	lenght types.Timecode
}

type editType struct {
	editStatement      string  //тип перехода с предыдущего сегмента.
	transitionDiration float64 //длинна перехода во фреймах. если editStatement == "C", то transitionDiration равно нулю
}

//////////////////////

func ParseFile(path string) (*edlData, error) {
	fmt.Println("Start Parse File")

	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer f.Close()

	return Parse(f)
}

type Decidion struct {
	index    string
	phrase   []Statement
	endAlert bool
}

func Parse(r io.Reader) (*edlData, error) {
	fmt.Println("Start Parse Reader")
	eData := edlData{}
	//eData, parseError = parseLine()
	scanner := bufio.NewScanner(r)
	parseError := errors.New("Initial")
	parseError = nil
	i := 0

	activeDesidion := Decidion{"", []Statement{}, false}
	for scanner.Scan() {
		// parseLine(state, &eData, scanner.Text()) (state, err)

		i++
		line := strings.TrimSpace(scanner.Text())
		newStatement, err := parseLine(line)
		switch { //DEBUG
		default:
			if newStatement.Type() == standardStatement {
				data, _ := newStatement.Declare()
				if data[0] != activeDesidion.index {
					fmt.Println("\n---INITIATE NEW DECIDION HERE---")
					activeDesidion.index = data[0]
				}

			}
			tp := newStatement.Type()
			for len(tp) < 10 {
				tp += " "
			}
			fmt.Printf("%v | %v\n", tp, line)

		case err != nil:
			if err.Error() == "BLANK LINE" {
				continue
			}
			fmt.Printf("UNKNOWN    | %v \n", line)
			parseError = err
		}
		//////////////////////////////////
		//АНАЛИЗИРОВАТЬ СТЭЙТМЕНТЫ ЗДЕСЬ//

		//////////////////////////////////
		if newStatement == nil {
			continue
		}
		fmt.Println(newStatement)
		parseError = err
		//eData.statment = append(eData.statment, newStatement)
	}
	if parseError != nil {
		return nil, parseError
	}
	fmt.Println("End Parse Reader")
	return &eData, nil
}

func parseLine(line string) (Statement, error) {
	switch {
	default:
		return nil, fmt.Errorf("line is UNKNOWN type statement:\n%v", line)
	case isHeader(line):
		return newHeader(line)
	case isFCM(line):
		return newFCM(line)
	case isM2(line):
		return newM2(line)
	case isStandard(line):
		return newStandard(line)
	case isEvent(line):
		return newEvent(line)
	case line == "":
		return nil, fmt.Errorf("BLANK LINE")
	case strings.Contains(line, "EFFECTS NAME IS"):
		line = "* " + line
		return newEvent(line)
	}

}

// type blankLine struct {
// 	err       error
// 	blConfirm bool
// }

// func (bl *blankLine) Declare() ([]string, error) {
// 	return []string{}, nil
// }

// func (bl *blankLine) Type() string {
// 	return "BLANK LINE"
// }
