package edl2

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/macroblock/imed/pkg/types"
)

const (
	TITLE              = "TITLE:"
	FCM                = "FCM:"
	SPLIT              = "SPLIT:"
	GPI                = "GPI"
	MS                 = "M/S"
	SWN                = "SWM"
	M2                 = "M2"
	AUD                = "AUD"
	EVENT              = "*"
	STATEMENT_HEADER   = "HEADER"
	STATEMENT_NOTE     = "NOTE"
	STATEMENT_STANDARD = "STANDARD"
)

type Statement interface {
	Declare() ([]string, error)
	Type() string
}

type header struct {
	id   string
	data string
}

func (h *header) Declare() ([]string, error) {
	return []string{h.id, h.data}, nil
}

func (h *header) Type() string {
	return "HEADER"
}

func isHeader(line string) bool {
	fld := strings.Fields(line)
	if len(fld) < 1 {
		return false
	}
	if fld[0] != TITLE {
		return false
	}
	return true
}

func newHeader(line string) (*header, error) {
	if !isHeader(line) {
		return nil, fmt.Errorf("line IS NOT a header statement:\n%v", line)
	}
	h := header{}
	h.id = TITLE
	h.data = strings.Split(line, TITLE+" ")[1]
	return &h, nil
}

type note struct {
	id      string
	message string
}

func (n *note) Declare() ([]string, error) {
	return []string{n.message}, nil
}

func (n *note) Type() string {
	return "NOTE"
}

func newNote(line string) (Statement, error) {
	fld := strings.Fields(line)
	switch fld[0] {
	default:
		n := note{"note", line}
		return &n, nil
	case TITLE:
		return newHeader(line)
	case FCM:
		return newFCM(line)
	case M2:
		return newM2(line)
	case AUD:
		return newAud(line)
	case EVENT:
		return newEvent(line)
		// case fcmID, splitID, gpiID, mstrSlvID, m2ID, machineDataID:
		// return nil, nil
	}

}

// func isValidNote(line string) bool { - бессмысленно, любая строка является стейтментом. если это не стандартное, то это ноут
// 	fld := strings.Fields(line)
// 	if len(fld) < 1 {
// 		return false
// 	}
// 	switch fld[0] {
// 	case fcmID, splitID, gpiID, mstrSlvID, m2ID, machineDataID:
// 		return
// 	}
// 	if fld[0] != fcmID {
// 		return false
// 	}
// 	return true
// }

type fcm struct {
	id   string
	data string
}

func isFCM(line string) bool {
	fld := strings.Fields(line)
	if len(fld) < 1 {
		return false
	}
	if fld[0] != FCM {
		return false
	}
	return true
}

func (f *fcm) Declare() ([]string, error) {
	return []string{f.id, f.data}, nil
}

func (f *fcm) Type() string {
	return "FCM"
}

func newFCM(line string) (*fcm, error) {
	if !isFCM(line) {
		return nil, fmt.Errorf("line IS NOT a FCM statement:\n%v", line)
	}
	h := fcm{}
	h.id = FCM
	h.data = strings.Split(line, FCM+" ")[1]
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
			ss.sequanceTime.out = timeCode
		case 1:
			ss.sequanceTime.in = timeCode
		case 2:
			ss.fileTime.out = timeCode
		case 3:
			ss.fileTime.in = timeCode
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

func (std *standard) Declare() ([]string, error) {
	dc := []string{std.index, std.reel, std.channels}
	dc = append(dc, std.editDesidion.editStatement)
	dc = append(dc, strconv.FormatFloat(std.editDesidion.transitionDiration, 'f', 1, 64))
	dc = append(dc, strconv.FormatFloat(std.fileTime.in.InSeconds(), 'f', 1, 64))
	dc = append(dc, strconv.FormatFloat(std.fileTime.out.InSeconds(), 'f', 1, 64))
	dc = append(dc, strconv.FormatFloat(std.fileTime.lenght.InSeconds(), 'f', 1, 64))
	dc = append(dc, strconv.FormatFloat(std.sequanceTime.in.InSeconds(), 'f', 1, 64))
	dc = append(dc, strconv.FormatFloat(std.sequanceTime.out.InSeconds(), 'f', 1, 64))
	dc = append(dc, strconv.FormatFloat(std.sequanceTime.lenght.InSeconds(), 'f', 1, 64))
	return dc, nil
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
	if !isEvent(line) {
		return nil, fmt.Errorf("statement is not Event:\n%v", line)
	}
	ev := event{}
	data := strings.TrimPrefix(line, "* ")
	ev.data = data
	return &ev, nil
}

func (ev *event) Declare() ([]string, error) {
	return []string{ev.data}, nil
}

func (ev *event) Type() string {
	return "EVENT"
}

type motionMemory struct {
	id                  string
	sourceID            string
	reelSpeed           string //или TRIG1 или VARIABLE
	reelReference       string //ждем data если reelSpeed содержит VARIABLE
	timingRelationships string
	timecodeValue       types.Timecode //абсолютный относительный или кол-во фреймов
	disabilityTrigger   bool
}

func newM2(line string) (*motionMemory, error) {
	mm := motionMemory{}
	//TODO: низко приоритетно. закончить когда будет готов сам парсер
	return &mm, nil
}

func isM2(line string) bool {
	fld := strings.Fields(line)
	if len(fld) < 1 {
		return false
	}
	if fld[0] != "M2" {
		return false
	}
	return true
}

func (mm *motionMemory) Declare() ([]string, error) {
	//TODO: низко приоритетно. закончить когда будет готов сам парсер
	return []string{"[TEMPLATE OF MOTION MEMORY DECLARATION]"}, nil
}

func (mm *motionMemory) Type() string {
	//TODO: низко приоритетно. закончить когда будет готов сам парсер
	return "M2"
}

type aud struct {
	id            string
	chanIndicator string
}

func newAud(line string) (*aud, error) {
	if !isAud(line) {
		return nil, fmt.Errorf("statement is not Aud:\n%v", line)
	}
	a := aud{}
	fields := strings.Fields(line)
	a.id = fields[0]
	switch fields[1] {
	default:
		return nil, fmt.Errorf("unknown channel indecator '%v'", fields[1])
	case "3", "4":
		a.chanIndicator = fields[1]
	}
	return &a, nil
}

func isAud(line string) bool {
	fld := strings.Fields(line)
	if len(fld) < 1 || len(fld) > 2 {
		return false
	}
	if fld[0] != "AUD" {
		return false
	}
	return true
}

func (a *aud) Declare() ([]string, error) {
	return []string{a.id, a.chanIndicator}, nil
}

func (a *aud) Type() string {
	return "AUD"
}

///////////////
