package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/grabber"
)

const (
	DIR_DONE    = `//192.168.31.4/buffer/IN/_DONE/`
	DIR_ARCHIVE = `//192.168.31.4/root/IN/_AMEDIA/_DONE/`
)

func listAmediaFiles() ([]string, error) {
	fi, err := os.ReadDir(DIR_DONE)
	if err != nil {
		return nil, fmt.Errorf("failed read directory: %v", err)
	}
	list := []string{}
	for _, f := range fi {
		if isAmediaName(f.Name()) {
			list = append(list, f.Name())
		}
	}
	debug(fmt.Sprintf("%v files", len(list)))
	return list, nil
}

func isAmediaName(name string) bool {
	re := regexp.MustCompile(`_PRT(\d){12}`)
	found := re.FindString(name)
	debug(found)
	if found == "" {
		return false
	}
	return true
}

func mapAmediaSeasons(list []string) map[string]string {
	fileMap := make(map[string]string)
	for _, file := range list {
		words := strings.Split(file, "_")
		dir := ""
	wordLoop:
		for _, word := range words {

			if strings.HasPrefix(word, "PRT") {
				break wordLoop
			}
			for i := 1; i < 20; i++ {
				if strings.HasPrefix(word, fmt.Sprintf("s0%ve", i)) {
					dir += fmt.Sprintf("s0%v", i)
					break wordLoop
				}
			}
			dir += word + "_"
		}
		dir = strings.TrimSuffix(dir, "_")
		fileMap[file] = dir + "/"

	}
	for k, v := range fileMap {
		debug(k, v)
	}
	return fileMap
}

func transferFiles(list []string, fileMap map[string]string) error {
	errs := []error{}
	for _, file := range list {
		oldPath := DIR_DONE + file
		newDir := DIR_ARCHIVE + fileMap[file]
		errs = append(errs, os.MkdirAll(newDir, 0666))
		newPath := DIR_ARCHIVE + fileMap[file] + file
		fmt.Println("transfer:")
		fmt.Println(oldPath)
		fmt.Println(newPath)
		err := grabber.CopyFile(oldPath, newDir)
		switch err {
		case nil:
			os.Remove(oldPath)
		default:
			errs = append(errs, err)
		}

	}

	return combinedError(errs...)
}

func debug(a ...string) {
	// fmt.Println("DEBUG:", a)
}

func combinedError(errs ...error) error {
	trueErrs := []error{}
	for _, err := range errs {
		if err != nil {
			trueErrs = append(trueErrs, err)
		}
	}
	if len(trueErrs) == 0 {
		return nil
	}
	eNum := 0
	errText := "combined errors:\n"
	for _, err := range trueErrs {
		if err != nil {
			errText += err.Error() + "\n"
			eNum++
		}

	}
	errText += fmt.Sprintf("total %v error(s)", eNum)
	return fmt.Errorf(errText)
}

func main() {
	list, err := listAmediaFiles()
	if err != nil {
		err = fmt.Errorf("failed list amedia files: %v", err)
		fmt.Println(err)
		os.Exit(1)
	}
	fileMap := mapAmediaSeasons(list)
	if err := transferFiles(list, fileMap); err != nil {
		fmt.Println(err)
	}

}

/*
1
прочитать файлы
выделить амедию (через PRT)
2
создать карту файл => сезон
3
создать папки с сезонами
4
отправить файлы по папкам с сезонами
*/
