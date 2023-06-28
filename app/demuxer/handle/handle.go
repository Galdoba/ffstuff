package handle

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/spreadsheet"

	"github.com/Galdoba/ffstuff/pkg/spreadsheet/tablemanager"

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
		Message:  message,
		Options:  append(options, "Отмена"),
		PageSize: 30,
	}
	if err := survey.AskOne(promptSelect, &choose, validator); err != nil {
		return choose, err
	}
	if choose == "Отмена" {
		return "", fmt.Errorf("\nбыла выбрана `Отмена`\n")
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

func SelectFromTable(listType string) []tablemanager.TaskData {
	sp, _ := spreadsheet.New()
	tlist := tablemanager.TaskListFrom(sp)
	//for _, err := range tlist.ParseErrors() {
	//		fmt.Println(err.Error())
	//}
	switch listType {
	default:
		fmt.Printf("SelectFromTable(listType string): listType=%v (unknown listType)", listType)
		return nil
	case "Фильм":
		return tlist.ReadyForDemux()
	case "Трейлер":
		return tlist.ChooseTrailer()

	}
	return nil

}

func TaskListFull() *tablemanager.TaskList {
	sp, _ := spreadsheet.New()
	return tablemanager.TaskListFrom(sp)
}

func ConvertToLinux(command string) string {
	args := strings.Fields(command)
	convertedCommand := ""
	for _, arg := range args {
		switch {
		default:
		case strings.Contains(arg, `\\nas\ROOT\EDIT\`):
			arg = strings.ReplaceAll(arg, `\\nas\ROOT\EDIT\`, `/mnt/pemaltynov/ROOT/EDIT/`)
		case strings.Contains(arg, `\\nas\ROOT\IN\`):
			arg = strings.ReplaceAll(arg, `\\nas\ROOT\IN\`, `/mnt/pemaltynov/ROOT/IN/`)
			//\\nas\buffer\IN\_IN_PROGRESS\
		case strings.Contains(arg, `\\nas\buffer\IN\_IN_PROGRESS\`):
			arg = strings.ReplaceAll(arg, `\\nas\buffer\IN\_IN_PROGRESS\`, `/home/pemaltynov/IN/_IN_PROGRESS/`)
		case strings.Contains(arg, `\\nas\buffer\IN\`):
			arg = strings.ReplaceAll(arg, `\\nas\buffer\IN\`, `/home/pemaltynov/IN/`)
			//\\nas\buffer\IN\_IN_PROGRESS\
		case strings.Contains(arg, `\\nas\buffer\IN\_DONE\`):
			arg = strings.ReplaceAll(arg, `\\nas\buffer\IN\_DONE\`, `/home/pemaltynov/IN/_DONE/`)
		case strings.Contains(arg, `cls`):
			arg = strings.ReplaceAll(arg, `cls`, `clear`)
		case strings.Contains(arg, `mkdir`):
			arg = strings.ReplaceAll(arg, `mkdir`, `mkdir -p`)
		case strings.Contains(arg, `move`):
			arg = strings.ReplaceAll(arg, `move`, `mv`)

			// echo 1>
		}
		arg = strings.ReplaceAll(arg, `\`, `/`)
		convertedCommand += arg + " "
	}
	convertedCommand = strings.ReplaceAll(convertedCommand, ` echo 1>`, ` touch `)
	convertedCommand = strings.ReplaceAll(convertedCommand, ` echo 1>`, ` touch `)
	return convertedCommand
}

func ArchiveDelay(file, destination string) error {
	f, err := os.OpenFile(`\\nas\buffer\IN\bats\sendToArchive.bat`, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	doneFolder := `\\nas\buffer\IN\_DONE\`
	text := fmt.Sprintf("mkdir %v\nmove %v%v %v\n", destination, doneFolder, file, destination)
	if _, err = f.WriteString(text); err != nil {

		return err
	}
	return nil
}
