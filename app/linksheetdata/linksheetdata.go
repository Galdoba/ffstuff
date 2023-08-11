package main

import (
	"fmt"
	"io/fs"
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
		Prompt: &survey.Select{
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
	for _, arg := range args {
		fmt.Println(arg)
	}

	fmt.Println("================================================================================")
	optType := askUser("Что это за контент?:", []string{"Фильм", "Трейлер", "Сериал", "[ВЫХОД]"})
	switch optType {
	case "Фильм":
		optType = namedata.CONTENT_TYPE_FILM
	case "Трейлер":
		optType = namedata.CONTENT_TYPE_TRL
	case "Сериал":
		optType = namedata.CONTENT_TYPE_SER
	case "[ВЫХОД]":
		os.Exit(2)
	}

	edits := preapareEditForms(args)
	names := drawOptionsFromTable(optType)
	if len(names) == 0 {
		fmt.Println("В таблице нет доступных названий")
		askUser("", []string{"Нажми ENTER для завершения"})
		os.Exit(1)
	}
	if optType == namedata.CONTENT_TYPE_TRL {
		names = excludeNames(names)
		if len(names) == 0 {
			fmt.Println("В таблице нет доступных названий")
			askUser("", []string{"Нажми ENTER для завершения"})
			os.Exit(1)
		}
	}

	answer := askUser("Выбери название из таблицы:", names)
	prefix := translit.CleanName(answer)
	translName := translit.Transliterate(prefix)
	fmt.Printf("'%s' ==> %s.\n", prefix, translName)
	errors := []error{}
	switch optType {
	case namedata.CONTENT_TYPE_SER:
		numbs := catchNumbersFromTableName(answer)
		seas := askUser("Какой это сезон? ", numbs)
		if strings.Contains(translName, seas+"_sezon") {
			translName = strings.Split(translName, seas+"_sezon")[0] + "s" + seas + "_"
			//			askUser("Верно? ", []string{translName, "No"})
		}
		///////////////
		chosenNames := []string{}
		for _, ed := range edits {
			nm := strings.TrimSuffix(ed.ShortName(), "."+namedata.RetrieveExtention(ed.Source()))
			chosenNames = append(chosenNames, nm)
		}
		allNumbrs := catchNumbersFromNames(chosenNames)
		numOpts := []string{}
		for _, nums := range allNumbrs {
			numOpts = append(numOpts, strings.Join(nums, " "))
		}
		seas2 := askUser("Какие это серии? ", numOpts)

		eps := strings.Fields(seas2)
		for j, ep := range eps {
			translNameSer := translName + ep
			editOne := []*namedata.EditNameForm{edits[j]}
			errors = append(errors, addPrefixToFiles(translNameSer, optType, editOne)...)
		}
		if strings.Contains(translName, seas+"_sezon") {
			translName = strings.Split(translName, seas+"_sezon")[0] + "s" + seas + "_xx"
			//			askUser("Верно? ", []string{translName, "No"})
		}
	default:
		errors = append(errors, addPrefixToFiles(translName, optType, edits)...)
	}

	for _, err := range errors {
		fmt.Println("error: ", err.Error())
	}

	if len(errors) != 0 {
		askUser("", []string{"Нажми ENTER для завершения"})
	}
	newNames := []string{}
	for _, edit := range edits {
		newNames = append(newNames, edit.Source())
	}
	namedata.NormalizeSoundNames(newNames)

}

func preapareEditForms(args []string) []*namedata.EditNameForm {
	edits := []*namedata.EditNameForm{}
	for _, arg := range args {
		thisEdit := namedata.EditForm(arg)
		edits = append(edits, thisEdit)

	}
	return edits
}

func drawOptionsFromTable(optType string) []string {
	sheet, err := spreadsheet.New()
	if err != nil {
		fmt.Println(err.Error())
	}
	if err := sheet.Update(); err != nil {
		fmt.Println("Can't update spreadsheet info")
		panic(err.Error())
	}

	fullList := tablemanager.TaskListFrom(sheet)
	names := []string{}
	switch optType {
	default:
		return nil
		//"Фильм", "Трейлер", "Сериал"
	case namedata.CONTENT_TYPE_FILM:
		readyList := fullList.ChooseFilm()
		for _, entry := range readyList {
			names = append(names, entry.Name())
		}
	case namedata.CONTENT_TYPE_TRL:
		readyList := fullList.ChooseTrailer()
		for _, entry := range readyList {
			names = append(names, entry.Name())
		}
	case namedata.CONTENT_TYPE_SER:
		readyList := fullList.ChooseSeason()
		for _, entry := range readyList {
			names = append(names, entry.Name())
		}
	}

	return names
}

func askUser(message string, options []string) string {
	question := append([]*survey.Question{}, &survey.Question{Name: "name",
		Prompt: &survey.Select{
			Message:  message,
			Options:  options,
			PageSize: 15,
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

func addPrefixToFiles(translName, contentType string, edits []*namedata.EditNameForm) []error {
	errors := []error{}

	for _, edit := range edits {
		switch strings.HasPrefix(edit.ShortName(), translName) {
		case false:

			if err := edit.AddPrefix(translName + "--" + contentType + "--"); err != nil {
				errors = append(errors, err)
			}
		case true:
		}
	}
	return errors
}

func catchNumbersFromTableName(fullName string) []string {
	//Perevozchik_01_sezon_12_seriy
	//Martimertv_01_sezon_05_seriya
	buf := ""
	found := []string{}
	for _, gl := range strings.Split(fullName, "") {
		switch gl {
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			buf += gl
			fmt.Printf("found %v: buffer = '%v'\n", gl, buf)
		default:
			if buf != "" {
				fmt.Printf("found %v: buffer = '%v' writing...\n", gl, buf)
				found = append(found, buf)
				buf = ""
			}
		}
	}
	return found
}

func catchNumbersFromNames(names []string) [][]string {
	numbers := [][]string{}
	for _, name := range names {
		numbrs := catchNumbersFromTableName(name)
		for i, n := range numbrs {
			if len(numbers) < i+1 {
				numbers = append(numbers, []string{})
			}
			numbers[i] = append(numbers[i], n)
		}
	}
	return numbers
}

func excludeNames(names []string) []string {
	leftoverNames := []string{}
	entries, _ := os.ReadDir(`\\192.168.31.4\root\EDIT\@trailers_temp\`)
	for _, n := range names {
		if hasReadyTrailer(n, entries) {
			continue
		}
		leftoverNames = append(leftoverNames, n)
	}
	return leftoverNames
}

func hasReadyTrailer(name string, entries []fs.DirEntry) bool {
	trlNM := translit.Transliterate(name)

	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".mp4") {
			continue
		}
		namestart := strings.Split(e.Name(), "_TRL")[0]
		if strings.HasPrefix(trlNM, namestart) {
			return true
		}
	}
	return false
}
