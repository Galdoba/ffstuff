package edl3

import (
	"errors"
	"strings"
)

type note struct {
	format            string
	content           string
	definedDataFields map[string]string
}

type Stater interface {
	State() (string, []string)
}

func newNote(line string) (*note, error) {
	if strings.TrimSpace(line) == "" {
		return nil, errors.New("input line have no data")
	}
	n := note{}
	n.definedDataFields = make(map[string]string)
	n.format = "NOTE"
	n.content = line
	return &n, nil
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

type standard struct {
	id           string
	reel         string
	channel      string
	editType     string
	editDuration float64
	fileIN       string
	fileOUT      string
	seqIN        string
	seqOUT       string
}

func newStandard(n *note) (*standard, error) {
	ss := standard{}

	return &ss, nil
}
