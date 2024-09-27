package grabberargs

import (
	"github.com/Galdoba/ffstuff/app/grabber/internal/validation"
	"github.com/Galdoba/ffstuff/pkg/logman"
)

func ValidateGrabArguments(args ...string) error {
	if len(args) == 0 {
		return logman.Errorf("grab command MUST have at least one argument")
	}
	for _, arg := range args {
		if err := validation.FileValidation(arg); err != nil {
			return logman.Errorf("argument '%v' invalid: %v", arg, err)
		}
	}
	return nil
}
