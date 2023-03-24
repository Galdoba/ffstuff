package main

import (
	"fmt"

	"gopkg.in/AlecAivazis/survey.v1"
)

func askSelection(message string, options []string) (string, error) {
	choose := ""
	validator := survey.ComposeValidators()
	promptSelect := &survey.Select{
		Message: message,
		Options: append(options, "Отмена"),
	}
	if err := survey.AskOne(promptSelect, &choose, validator); err != nil {
		return choose, err
	}
	if choose == "Отмена" {
		return "", fmt.Errorf("была выбрана `Отмена`")
	}

	return choose, nil
}

func panicIfErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
