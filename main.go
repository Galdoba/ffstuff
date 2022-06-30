package main

import (
	"gopkg.in/AlecAivazis/survey.v1"
)

// the questions to ask
var qs = []*survey.Question{
	{
		Name:      "name",
		Prompt:    &survey.Input{Message: "What is your name?"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name: "color",
		Prompt: &survey.Select{
			Message: "Choose a color:",
			Options: []string{"red", "blue", "green"},
			Default: "red",
		},
	},
	{
		Name:   "age",
		Prompt: &survey.Input{Message: "How old are you?"},
	},
}

func main() {
	// the answers will be written to this struct

	// media, _ := probe.NewMedia(`d:\IN\IN_testInput\trailers\CRUELLA_iEST_TLRE_HD_2398_51_20_16x9_185_RUS_D1415623.mov`)
	// mr := media.MediaFileReport()
	// probe.SelectAudio(mr)
}
