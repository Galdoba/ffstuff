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

var ErrBlankLine = errors.New("blank line detected")

const (
	fcmModeDF         = "DROP FRAME"
	fcmModeNDF        = "NON-DROP FRAME"
	standardStatement = "STANDARD"
)

type edlData struct {
	edlSource string //источник самого edl (не обязательно файл)
	fcmMode   string
	decidions []Decidion
}

func (ed *edlData) String() string {
	str := "Source File: " + ed.edlSource + "\n"
	for i, v := range ed.decidions {
		str += fmt.Sprintf("Decidion %d: %v\n", i, v)
	}
	return str
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
	index     string
	phrase    []Statement
	concluded bool
}

func Parse(r io.Reader) (*edlData, error) {
	fmt.Println("Start Parse Reader")
	eData := edlData{}
	//eData, parseError = parseLine()
	scanner := bufio.NewScanner(r)
	activeDesidion := Decidion{"", []Statement{}, false}
	newStatement, parseError := newNote("Test line")
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		newStatement, parseError = parseLine(line)
		if newStatement == nil {
			parseError = fmt.Errorf("no statemend received for Decidion %v", activeDesidion.index)
		}
		switch parseError {
		default:
			return &eData, parseError
		case ErrBlankLine:
			continue
		case nil:
		}
		if newStatement.Type() == STATEMENT_STANDARD {
			data, _ := newStatement.Declare()
			if data[0] != activeDesidion.index {
				//fmt.Println(activeDesidion)
				eData.decidions = append(eData.decidions, activeDesidion)
				activeDesidion = Decidion{data[0], []Statement{}, false}
				//fmt.Println("\n---INITIATE NEW DECIDION HERE---")
			}
		}
		activeDesidion.phrase = append(activeDesidion.phrase, newStatement)
		//fmt.Println(newStatement.Declare())
		//////////////////////////////////
		//АНАЛИЗИРОВАТЬ СТЭЙТМЕНТЫ ЗДЕСЬ//

		//////////////////////////////////
	}
	eData.decidions = append(eData.decidions, activeDesidion)
	if parseError != nil {
		return nil, parseError
	}
	fmt.Println("End Parse Reader")
	return &eData, nil
}

func parseLine(line string) (Statement, error) {
	switch {
	default:
		return newNote(line)
	case isStandard(line):
		return newStandard(line)
	case line == "":
		return nil, ErrBlankLine
	}
}
