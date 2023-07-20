package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/namedata"
	"github.com/Galdoba/ffstuff/pkg/spreadsheet"
	"github.com/Galdoba/ffstuff/pkg/spreadsheet/tablemanager"
	"github.com/Galdoba/ffstuff/pkg/translit"
	"gopkg.in/AlecAivazis/survey.v1"
)

// the questions to ask
var simpleQs = []*survey.Question{
	{
		Name: "name",
		Prompt: &survey.Select{
			Message: "Выбери название из таблицы:",
			Options: []string{},
		},
		Validate: survey.Required,
	},
}

func execute(f func() error) {
	if err := f(); err != nil {
		panic(err.Error())
	}
}

/*
todo:
аргументы - файлы
спрашивает к какой записи они относятся
добавляет соответствующий тег к их названию и транслитерирует остальное

*/

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Не полученно аргументов")
		askUser("Нажми ENTER для завершения", []string{})
		return
	}

	edits := preapareEditForms(args)
	names := drawOptionsFromTable()
	if len(names) == 0 {
		fmt.Println("В таблице нет доступных названий")
		askUser("Нажми ENTER для завершения", []string{})
		return
	}
	answer := askUser("Выбери название из таблицы:", names)
	errors := addPrefixToFiles(answer, edits)
	for _, err := range errors {
		fmt.Println("error: ", err.Error())
	}
	if len(errors) != 0 {
		askUser("Нажми ENTER для завершения", []string{})
	}

}

func preapareEditForms(args []string) []*namedata.EditNameForm {
	edits := []*namedata.EditNameForm{}
	for _, arg := range args {
		thisEdit := namedata.EditForm(arg)
		edits = append(edits, thisEdit)

	}
	return edits
}

func drawOptionsFromTable() []string {
	sheet, err := spreadsheet.New()
	if err != nil {
		fmt.Println(err.Error())
	}
	if err := sheet.Update(); err != nil {
		fmt.Println("Can't update spreadsheet info")
		panic(err.Error())
	}

	fullList := tablemanager.TaskListFrom(sheet)
	readyList := fullList.ReadyForDemux()
	names := []string{}
	for _, entry := range readyList {
		names = append(names, entry.Name())
	}
	return names
}

func askUser(message string, options []string) string {
	question := append([]*survey.Question{}, &survey.Question{Name: "name",
		Prompt: &survey.Select{
			Message: "Выбери название из таблицы:",
			Options: options,
		},
		Validate: survey.Required,
	})

	answers := struct {
		Name string
	}{}

	// ask the question
	err := survey.Ask(question, &answers)
	if err != nil {
		fmt.Println("не могу получить ответ")
		panic(err.Error())
	}
	return answers.Name
}

func addPrefixToFiles(prefix string, edits []*namedata.EditNameForm) []error {
	errors := []error{}
	prefix = translit.CleanName(prefix)
	translName := translit.Transliterate(prefix)
	fmt.Printf("'%s' ==> %s.\n", prefix, translName)
	for _, edit := range edits {
		switch strings.HasPrefix(edit.ShortName(), translName) {
		case false:
			if err := edit.AddPrefix(translName + "--"); err != nil {
				errors = append(errors, err)
			}
		case true:
		}
	}
	return errors
}
