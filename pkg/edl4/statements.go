package edl4

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/types"
)

type note struct {
	format            string
	content           string
	definedDataFields map[string]string
}

type Stater interface {
	State() (string, []string)
}

func Statement(str string) (string, []string, error) {
	n, err := newNote(str)
	if err != nil {
		return "NONE", []string{"line separator"}, nil
	}
	if title, err := newTitle(n); err == nil {
		tp, dt := title.State()
		return tp, dt, nil
	}
	if fcm, err := newFCM(n); err == nil {
		tp, dt := fcm.State()
		return tp, dt, nil
	}
	if st, err := newStandard(n); err == nil {
		tp, dt := st.State()
		return tp, dt, nil
	}
	if a, err := newAud(n); err == nil {
		tp, dt := a.State()
		return tp, dt, nil
	}
	if src, err := newSource(n); err == nil {
		tp, dt := src.State()
		return tp, dt, nil
	}
	tp, dt := n.State()
	return tp, dt, err
}

func newNote(line string) (*note, error) {
	if strings.TrimSpace(line) == "" {
		return nil, errors.New("input line have no data")
	}
	n := note{}
	n.definedDataFields = make(map[string]string)
	n.content = line
	n.format = "NOTE"
	return &n, nil
}

func (n *note) State() (string, []string) {
	return n.format, []string{n.content}
}

type title struct {
	title string
}

func newTitle(n *note) (*title, error) {
	ttl := strings.TrimPrefix(n.content, "TITLE: ")
	if ttl == n.content {
		return nil, errors.New("statement is not a title")
	}
	t := title{}
	t.title = ttl
	return &t, nil
}

func (t *title) State() (string, []string) {
	return "TITLE", []string{t.title}
}

type fcm struct {
	mode string
}

func newFCM(n *note) (*fcm, error) {
	mode := strings.TrimPrefix(n.content, "FCM: ")
	if mode == n.content {
		return nil, errors.New("statement is not a FCM")
	}
	switch mode {
	default:
		return nil, errors.New("unknown FCM mode")
	case "NON-DROP FRAME", "DROP FRAME":
	}
	fcm := fcm{}
	fcm.mode = mode
	return &fcm, nil
}

func (f *fcm) State() (string, []string) {
	return "FCM", []string{f.mode}
}

type standard struct {
	id           string
	reel         string
	channel      string
	editType     string
	editDuration float64
	fileIN       types.Timecode
	fileOUT      types.Timecode
	seqIN        types.Timecode
	seqOUT       types.Timecode
}

func newStandard(n *note) (*standard, error) {
	ss := standard{}
	flds := strings.Fields(n.content)
	l := len(flds)
	switch {
	case l != 8 && l != 9:
		return nil, fmt.Errorf("statement is not a Standard")
	}
	if !validIndex(flds[0]) {
		return nil, fmt.Errorf("invalid statement syntax: index '%v' is invalid", ss.id)
	}
	ss.id = flds[0]
	ss.reel = flds[1]
	if !listContains([]string{"AX", "BL"}, ss.reel) {
		return nil, fmt.Errorf("invalid statement syntax: reel '%v' invalid", ss.reel)
	}
	ss.channel = flds[2]
	if !listContains([]string{"V", "A", "A2", "NONE"}, ss.channel) {
		return nil, fmt.Errorf("invalid statement syntax: channel '%v' invalid", ss.channel)
	}
	ss.editType = flds[3]
	if !listContains(validWipeCodes(), ss.editType) {
		return nil, fmt.Errorf("invalid statement syntax: editType '%v' invalid", ss.editType)
	}
	if l == 9 {
		if ss.editType == "C" {
			return nil, fmt.Errorf("invalid statement syntax: editType 'C' can't have any duration field")
		}
		dur, err := strconv.ParseFloat(flds[l-5], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid statement syntax: %v", err.Error())
		}
		ss.editDuration = dur
	}
	for i := 4; i >= 1; i-- {
		ts, err := types.ParseTimecode(flds[l-i])
		if err != nil {
			return nil, fmt.Errorf("invalid statement syntax: %v", err.Error())
		}
		switch i {
		case 4:
			ss.fileIN = ts
		case 3:
			ss.fileOUT = ts
		case 2:
			ss.seqIN = ts
		case 1:
			ss.seqOUT = ts
		}
	}
	return &ss, nil
}

func (s *standard) State() (string, []string) {
	return "STANDARD", []string{
		s.id,
		s.reel,
		s.editType,
		s.channel,
		strconv.FormatFloat(s.editDuration, 'f', 1, 64),
		s.fileIN.String(),
		s.fileOUT.String(),
		s.seqIN.String(),
		s.seqOUT.String(),
	}
}

type source struct {
	sourceA string
	sourceB string
}

func newSource(n *note) (*source, error) {
	src := source{}
	err := errors.New("Not imlemented")
	err = nil
	aSrc := strings.TrimPrefix(n.content, "* FROM CLIP NAME: ")
	bSrc := strings.TrimPrefix(n.content, "* TO CLIP NAME: ")
	switch {
	default:
		err = errors.New("statement is not a Source")
	case aSrc != n.content:
		src.sourceA = aSrc
	case bSrc != n.content:
		src.sourceB = bSrc
	}
	return &src, err
}

func (s *source) State() (string, []string) {
	if s.sourceA != "" {
		return "SOURCE A", []string{s.sourceA}
	}
	return "SOURCE B", []string{s.sourceB}
}

///////////
type aud struct {
	channel int
}

func newAud(n *note) (*aud, error) {
	a := aud{}
	flds := strings.Fields(n.content)
	if flds[0] != "AUD" {
		return nil, errors.New("statement is not a Aud")
	}
	if len(flds) < 2 {
		return nil, fmt.Errorf("invalid statement syntax: '%v' have no data on channels", n.content)
	}
	if len(flds) > 2 {
		return nil, fmt.Errorf("invalid statement syntax: to many fields present '%v'", n.content)
	}
	val, err := strconv.Atoi(flds[1])
	if err != nil {
		return nil, fmt.Errorf("invalid statement syntax: can't parse '%v'", flds[1])
	}
	a.channel = val
	if a.channel != 3 && a.channel != 4 {
		return nil, fmt.Errorf("invalid statement argument: field '%v' must be '3' or '4' ", a.channel)
	}
	return &a, nil
}

func (a *aud) State() (string, []string) {
	return "AUD", []string{strconv.Itoa(a.channel)}
}

type effectStatement struct {
	effectName string
}

func newEffectStatement(n *note) (*effectStatement, error) {
	es := effectStatement{}
	err := errors.New("Unimplemented555")
	if !strings.HasPrefix(n.content, "EFFECTS NAME IS ") {
		return nil, errors.New("statement is not a EffectStatement")
	}

	err = nil
	effect := strings.TrimPrefix(n.content, "EFFECTS NAME IS ")
	es.effectName = effect

	return &es, err
}

/////////////////HELPERS//////////////////

func validWipeCodes() []string {
	sl := []string{"C", "D"}
	for i := 0; i < 24; i++ {
		code := strconv.Itoa(i)
		if len(code) < 2 {
			code = "0" + code
		}
		sl = append(sl, "W0"+code)
		sl = append(sl, "W1"+code)
	}
	return sl
}

func validIndex(in string) bool {
	ind, err := strconv.Atoi(in)
	switch {
	case ind < 1:
		return false
	case ind > 999:
		return false
	case len(in) != 3:
		return false
	case err != nil:
		return false
	}
	return true
}

func listContains(list []string, elem string) bool {
	for _, v := range list {
		if elem == v {
			return true
		}
	}
	return false
}
