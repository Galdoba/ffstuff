package compare

import (
	"sort"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/namedata"
	"github.com/Galdoba/ffstuff/pkg/spreadsheet"
	"github.com/Galdoba/ffstuff/pkg/spreadsheet/tablemanager"
)

/*
исходные данные:
имя файла - стринг
таблица - csv

цель - выяснить является ли файл - амедией.

продседура:
разбиваем тело файла на слова
ищем эти слова в таблице
если совпадений более 1 фильтруем
если совпадение равно 1 - отдаем результат
*/

var TableData []tablemanager.TaskData

func init() {
	sp, err := spreadsheet.New()
	if err != nil {
		panic(err.Error())
	}
	taskList := tablemanager.TaskListFrom(sp)
	TableData = taskList.ALL()
}

func SuggestNameFromTable(path string) string {
	patern := matchWithTable(path)
	//keys := keysSorted(patern)
	maxMatches := 0
	bestMatch := []string{}
	for _, val := range patern {
		if val > maxMatches {
			maxMatches = val
		}
	}
	for k, val := range patern {
		if val != maxMatches {
			continue
		}
		bestMatch = append(bestMatch, k)
	}
	if len(bestMatch) == 1 {
		return bestMatch[0]
	}
	return ""
}

func keysSorted(paternmap map[string]int) []string {
	ks := []string{}
	for l, v := range paternmap {
		if v > 0 {
			ks = append(ks, l)
		}
	}
	sort.Strings(ks)
	return ks
}

func matchWithTable(path string) map[string]int {
	matchTimes := make(map[string]int)
	ef := namedata.EditForm(path)
	words := ef.Words()
	for _, data := range TableData {
		// if data.Agent() != "Амедиа" { //временно
		// 	continue
		// }

		taskNameRaw := data.Name()
		nameTL := namedata.TransliterateForEdit(taskNameRaw)
		words2 := wordsFromPath(nameTL)
		matchTimes[taskNameRaw] = countMatches(words, words2)
	}
	return matchTimes
}

func wordsFromPath(path string) []string {
	short := namedata.RetrieveShortName(path)
	separators := []string{" ", "(", ")", "-", "."}
	//shortLit := strings.Split(short, "")
	//shortLitSt := []string{}
	shortSt := ""
	for _, sep := range separators {
		shortSt = strings.ReplaceAll(short, sep, "_")

	}

	//shortLitStS := strings.Join(shortLitSt, "")
	for i := 0; i < 20; i++ {
		shortSt = strings.ReplaceAll(shortSt, "__", "_")
	}
	words := strings.Split(shortSt, "_")
	return words
}

func countMatches(sl1, sl2 []string) int {
	m := 0
	for _, s1 := range sl1 {
		for _, s2 := range sl2 {

			if s1 != "" && s1 == s2 {
				m++
			}
		}
	}
	return m
}
