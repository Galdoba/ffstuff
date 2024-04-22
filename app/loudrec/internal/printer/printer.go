package printer

import (
	"fmt"
	"io"
	"os"
)

type PrintOptions struct {
	silent  bool
	writers []io.Writer
}

func (p *PrintOptions) ToggleSilent() {
	p.silent = !p.silent
}

func Print(opt *PrintOptions, msg string) {
	if opt.silent {
		return
	}
	fmt.Fprint(os.Stderr, msg)
}
