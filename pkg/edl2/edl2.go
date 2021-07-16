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
	edlSource         string //источник самого edl (не обязательно файл)
	data              []Statement
	scanningConcluded bool
	pseudoResults     [][]string
	//resultMap map[int]trackData
}

func (ed *edlData) String() string {
	str := "Source File: " + ed.edlSource + "\n"
	for i, v := range ed.data {
		str += fmt.Sprintf("Statement %d: %v\n", i, v)
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
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

func Parse(r io.Reader) (*edlData, error) {
	fmt.Println("Start Parse Reader")
	eData := edlData{}
	scanner := bufio.NewScanner(r)

	newStatement, parseError := newNote("NULL STATEMENT")
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		newStatement, parseError = parseLine(line)
		switch {
		default:
			return &eData, fmt.Errorf("Unknown or missed error: %v", parseError.Error())
		case parseError == ErrBlankLine:
			continue
		case parseError == nil:
			eData.data = append(eData.data, newStatement)

		}
	}
	eData.scanningConcluded = true

	//////////////////////////////////
	//АНАЛИЗИРОВАТЬ СТЭЙТМЕНТЫ ЗДЕСЬ//
	//fmt.Println("CONCLUDE DATA HERE:")
	//ShowResults(eData)
	//////////////////////////////////

	fmt.Println("End Parse Reader")

	return &eData, parseError
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

////////////////////////////DATA ASSEMBLER////////////////

type Constructror struct {
	title         string
	fcmMode       string
	activeChannel string
	sourcesUsed   map[string]int
	vClips        []clip
}

func ShowResults(edl edlData) error {
	v0 := []string{}
	a1 := []string{}
	a2 := []string{}
	a3 := []string{}
	a4 := []string{}
	unknwn := []string{}
	waitAUD := false
	/////////////////////////////
	for i, st := range edl.data {
		fmt.Printf("Go res %d\n", i+1)
		if st.Type() == "AUD" && waitAUD {
			fmt.Println("1 IF")
			fld, _ := st.Declare()
			switch fld[1] {
			case "3":
				a3 = append(a3, unknwn...)
			case "4":
				a4 = append(a4, unknwn...)
			}
			unknwn = []string{}
			waitAUD = false
		}
		if st.Type() != STATEMENT_STANDARD {
			fmt.Println("2 IF")
			continue
		}

		fields, err := st.Declare()
		if err != nil {
			fmt.Println("3 IF", err)
			return err
		}
		if fields[1] != "AX" {
			fmt.Println("5 IF")
			continue
		}
		if fields[3] != "C" {
			fmt.Println("6 IF")
			continue
		}
		switch fields[2] {
		default:
			fmt.Println("4 IF", err)
		case "V":
			v0 = append(v0, fields[0])
		case "A":
			a1 = append(a1, fields[0])
		case "A2":
			a2 = append(a2, fields[0])
		case "NONE":
			unknwn = append(unknwn, fields[0])
			waitAUD = true
		}
	}
	/////////////////////////
	fmt.Println(v0)
	fmt.Println(a1)
	fmt.Println(a2)
	fmt.Println(a3)
	fmt.Println(a4)
	fmt.Println(unknwn)
	return nil
}

type clip struct {
	//name        string
	mixType     string
	mixDuration float64
	channel     string
	duration    types.Timecode
	srcA        string
	srcAIn      types.Timecode
	srcADur     types.Timecode
	srcB        string
	srcBIn      types.Timecode
	srcBDur     types.Timecode
	nextClip    *clip
}
