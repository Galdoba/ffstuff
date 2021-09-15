package edl4

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var ErrBlankLine = errors.New("blank line detected")

func ParseFile(path string) (parcedData, error) {
	fmt.Println("Start Parse File")
	f, err := os.Open(path)
	if err != nil {
		return parcedData{}, err
	}
	defer f.Close()
	return Parse(f)
}

func Parse(r io.Reader) (parcedData, error) {
	//fmt.Println("Start Parse Reader")
	var eData []statementData
	scanner := bufio.NewScanner(r)

	newStatement := statementData{}
	parseError := errors.New("Initial")
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		newStatement, parseError = parseLine(line)
		switch {
		default:

			return parcedData{eData}, fmt.Errorf("Unknown or missed error: %v", parseError.Error())
		case parseError == ErrBlankLine:
			continue
		case parseError == nil:
			eData = append(eData, statementData{newStatement.sType, newStatement.fields})
			//fmt.Printf("Input: '%v'\n", line)
			//	fmt.Printf("Parsed: %v	'%v'\n", newStatement.sType, newStatement.fields)
		}
	}

	//////////////////////////////////
	//АНАЛИЗИРОВАТЬ СТЭЙТМЕНТЫ ЗДЕСЬ//
	//fmt.Println("CONCLUDE DATA HERE:")
	//ShowResults(eData)
	//////////////////////////////////

	//fmt.Println("End Parse Reader")
	pd := parcedData{eData}
	return pd, parseError
}

type statementData struct {
	sType  string
	fields []string
}

func parseLine(line string) (statementData, error) {
	sd := statementData{}
	sType, sData, err := Statement(line)
	sd.sType = sType
	sd.fields = sData
	if err != nil {
		return statementData{}, err
	}
	return sd, nil
}
