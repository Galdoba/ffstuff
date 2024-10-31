package bat

import (
	"fmt"
	"os"
)

type batmaker struct {
	path   string
	prompt string
}

type BatOptions func(*batOption)

type batOption struct {
	noEcho bool
}

func NoEcho() BatOptions {
	return func(bo *batOption) {
		bo.noEcho = true
	}
}

func CreateBat(path, prompt string, opts ...BatOptions) error {
	setings := &batOption{}
	for _, change := range opts {
		change(setings)
	}
	text := ""
	if setings.noEcho {
		text += "@ECHO OFF\n"
	}
	text += prompt
	_, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create batch file: %v", err)
	}
	if err := os.WriteFile(path, []byte(text), 0666); err != nil {
		return fmt.Errorf("failed to write batch file: %v", err)
	}
	return nil
}
