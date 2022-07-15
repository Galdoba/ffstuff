package handle

import (
	"fmt"
	"os"

	"gopkg.in/AlecAivazis/survey.v1"
)

func SelectionSingle(message string, options ...string) string {
	selected, err := askSelection(survey.ComposeValidators(), message, options)
	Error(err)
	return selected
}

func askSelection(validator survey.Validator, message string, options []string) (string, error) {
	choose := ""
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

func askInput(val survey.Validator, message string) (string, error) {
	result := ""
	promptInput := &survey.Input{
		Message: message,
	}
	if err := survey.AskOne(promptInput, &result, val); err != nil {
		return result, err
	}
	return result, nil
}

func Error(err error, notFatal ...error) {
	if err != nil {
		for _, errMet := range notFatal {
			if err.Error() == errMet.Error() {
				fmt.Printf("Error: %v", err)
				return
			}
		}
		fmt.Printf("Fatal Error: %v", err)
		os.Exit(1)
	}
}
