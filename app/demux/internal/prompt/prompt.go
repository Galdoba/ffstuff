package prompt

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

func MultiSelect(question string, args ...string) ([]string, error) {
	selected := []string{}
	options := huh.NewOptions(args...)
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(question).
				Options(options...,
				).
				Value(&selected),
		),
	)
	err := form.Run()
	if err != nil {
		return selected, fmt.Errorf("selection failed: %v")
	}
	return selected, nil
}
