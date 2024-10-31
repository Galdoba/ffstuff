package clrequest

import (
	"fmt"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/clrequest/bat"
)

type request struct {
	programm string
	args     []string
}

func New(program string, args ...string) *request {
	r := request{}
	r.programm = program
	r.args = args
	return &r
}

func (r *request) Prompt() string {
	text := r.programm
	for _, arg := range r.args {
		switch argType(arg) {
		case normal, quoted, doubleQuoted:
			text += " " + arg
		case hasSpace:
			text += " " + fmt.Sprintf(`"%v"`, arg)
		}
	}
	return text
}

const (
	quoted       = "quoted"
	doubleQuoted = "doubleQuoted"
	hasSpace     = "has space"
	normal       = "normal"
)

func argType(arg string) string {
	if strings.HasPrefix(arg, `"`) && strings.HasSuffix(arg, `"`) {
		return doubleQuoted
	}
	if strings.HasPrefix(arg, "`") && strings.HasSuffix(arg, "`") {
		return quoted
	}
	if strings.Contains(arg, " ") {
		return hasSpace
	}
	return normal
}

func (r *request) Bat(path string, options ...bat.BatOptions) error {
	return bat.CreateBat(path, r.Prompt(), options...)
}
