package edl2

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/macroblock/imed/pkg/types"
)

var currentEDIndex string

const (
	titleID    = "TITLE:"
	fcmID      = "FCM:"
	fcmModeDF  = "DROP FRAME"
	fcmModeNDF = "NON-DROP FRAME"
)

type edlData struct {
	edlSource string //источник самого edl (не обязательно файл)
	statment  []Statement
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

type Statement interface {
	Parse() ([]string, error)
	Type() string
}

type header struct {
	id   string
	data string
}

func (h *header) Parse() ([]string, error) {
	return []string{h.id, h.data}, nil
}

func (h *header) Type() string {
	return "HEADER"
}

func isHeader(line string) bool {
	title := strings.TrimPrefix(line, titleID)
	if line == title {
		return false
	}
	return true
}

func newHeader(line string) (*header, error) {
	if !isHeader(line) {
		return nil, fmt.Errorf("line IS NOT a header statement:\n%v", line)
	}
	h := header{}
	h.id = titleID
	h.data = strings.Split(line, titleID+" ")[1]
	return &h, nil
}

type fcm struct {
	id   string
	data string
}

func isFCM(line string) bool {
	fcm := strings.TrimPrefix(line, fcmID)
	if line == fcm {
		return false
	}
	return true
}

func (f *fcm) Parse() ([]string, error) {
	return []string{f.id, f.data}, nil
}

func (f *fcm) Type() string {
	return "FCM"
}

func newFCM(line string) (*header, error) {
	if !isFCM(line) {
		return nil, fmt.Errorf("line IS NOT a FCM statement:\n%v", line)
	}
	h := header{}
	h.id = fcmID
	h.data = strings.Split(line, fcmID+" ")[1]
	return &h, nil
}

type standard struct {
	index        string
	reel         string
	channels     string
	editDesidion editType
	fileTime     timeSegment
	sequanceTime timeSegment
}

func isStandard(line string) bool {
	fields := strings.Fields(line)
	if len(fields) != 8 && len(fields) != 9 {
		return false
	}
	if _, err := strconv.Atoi(fields[0]); err != nil { //если не число, то не std
		return false
	}
	if len(fields[1]) > 4 || len(fields[2]) > 4 || len(fields[3]) > 4 {
		return false
	}
	for _, v := range last4IndexesOf(fields) { //последние 4 поля
		if _, err := types.ParseTimecode(fields[v]); err != nil {
			return false
		}
	}
	return true
}

func newStandard(line string) (*standard, error) {
	if !isStandard(line) {
		return nil, fmt.Errorf("statement is not Standard:\n%v", line)
	}
	ss := standard{}
	fields := strings.Fields(line)
	ss.index = fields[0]
	if ss.index != currentEDIndex {
		currentEDIndex = ss.index

		fmt.Println(" ")
		fmt.Println("//////////NEW DECIDION//////////")
	}
	ss.reel = fields[1]
	ss.channels = fields[2]
	for i, v := range last4IndexesOf(fields) {
		timeCode, err := types.ParseTimecode(fields[v])
		if err != nil {
			return nil, fmt.Errorf("timeCode expecting %v | %v", fields[v], line)
		}
		switch i {
		default:
			return nil, fmt.Errorf("newStandard(): unexpected index %v", i)
		case 0:
			ss.fileTime.in = timeCode
		case 1:
			ss.fileTime.out = timeCode
		case 2:
			ss.sequanceTime.in = timeCode
		case 3:
			ss.sequanceTime.out = timeCode
		}
	}
	ss.fileTime.lenght = ss.fileTime.out - ss.fileTime.in
	ss.sequanceTime.lenght = ss.sequanceTime.out - ss.sequanceTime.in
	ss.editDesidion.editStatement = fields[3]
	if len(fields) == 9 {
		trDur, err := strconv.ParseFloat(fields[4], 64)
		if err != nil {
			return nil, fmt.Errorf("variable %v is not a float64", fields[4])
		}
		ss.editDesidion.transitionDiration = trDur
	}
	return &ss, nil
}

func (std *standard) Parse() ([]string, error) {
	return []string{"//data from standard"}, nil
}

func (std *standard) Type() string {
	return "STANDARD"
}

func last4IndexesOf(sl []string) []int {
	l := len(sl)
	res := []int{}
	for i := 0; i < 4; i++ {
		if l-1-i >= 0 {
			res = append(res, l-1-i)
		}
	}
	return res
}

type event struct {
	data string
}

func isEvent(line string) bool {
	event := strings.TrimPrefix(line, "* ")
	if event == line {
		return false
	}
	return true
}

func newEvent(line string) (*event, error) {
	ev := event{}
	data := strings.TrimPrefix(line, "* ")
	ev.data = data
	return &ev, nil
}

func (ev *event) Parse() ([]string, error) {
	return []string{ev.data}, nil
}

func (ev *event) Type() string {
	return "EVENT"
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

func Parse(r io.Reader) (*edlData, error) {
	fmt.Println("Start Parse Reader")
	eData := edlData{}
	//eData, parseError = parseLine()
	scanner := bufio.NewScanner(r)
	parseError := errors.New("Initial")
	parseError = nil
	i := 0
	for scanner.Scan() {
		// parseLine(state, &eData, scanner.Text()) (state, err)

		i++
		line := strings.TrimSpace(scanner.Text())
		newStatement, err := parseLine(line)
		//fmt.Println("statment:", statement)
		parseError = err
		eData.statment = append(eData.statment, newStatement)
		//fmt.Println("err:", parseError)
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
		fmt.Println("found UNDEFINED:", line)
		return nil, fmt.Errorf("line is UNKNOWN type statement:\n%v", line)
	case isHeader(line):
		fmt.Println("found HEADER   :", line)
		return newHeader(line)
	case isFCM(line):
		fmt.Println("found FCM      :", line)
		return newFCM(line)
	case isStandard(line):
		fmt.Println("found STANDARD :", line)
		return newStandard(line)
	case isEvent(line):
		fmt.Println("found EVENT    :", line)
		return newEvent(line)
	case line == "":
		return nil, fmt.Errorf("line is BLANK:\n%v", line)
	}
}
