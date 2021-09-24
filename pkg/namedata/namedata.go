package namedata

import (
	"fmt"
	"strings"
)

/*
from Name:
-basename
-extention
-tags


*/

//RetrieveAll -
func RetrieveAll(path string) (string, string, []string) {
	fileName := shortFileName(path)
	tags, ext := splitName(fileName)
	base, tags2 := nameBase(tags)
	return base, ext, tags2
}

func RetrieveDirectory(path string) string {
	pathData := strings.Split(path, "\\")
	return strings.Join(pathData[0:len(pathData)-1], "\\") + "\\"
}

func RetrieveShortName(path string) string {
	return shortFileName(path)
}

func RetrieveBase(path string) string {
	fileName := shortFileName(path)
	if strings.Contains(fileName, "__") {
		return strings.Split(fileName, "__")[0]
	}
	tags, _ := splitName(fileName)
	base, _ := nameBase(tags)
	return base
}

func RetrieveExtention(path string) string {
	fileName := shortFileName(path)
	_, ext := splitName(fileName)
	return ext
}

func RetrieveTags(path string) []string {
	fileName := shortFileName(path)
	tags, _ := splitName(fileName)
	_, tags2 := nameBase(tags)
	return tags2
}

func RetrieveDrive(path string) string {
	data := strings.Split(path, "\\")
	return data[0]
}

func shortFileName(path string) string {
	data := strings.Split(path, "\\")
	fileName := path
	if len(data) > 1 {
		fileName = data[len(data)-1]
	}
	return fileName
}

func splitName(fileName string) ([]string, string) {
	data := strings.Split(fileName, "_")
	tags := []string{}
	ext := ""
	for index, val := range data {
		if index == len(data)-1 {
			p := strings.Split(val, ".")
			ext = p[len(p)-1]
			tags = append(tags, p[0])
			continue
		}
		tags = append(tags, val)
	}
	return tags, ext
}

func nameBase(tags []string) (string, []string) {
	base := ""
	tags2 := []string{}
	for i, val := range tags {
		for _, tg := range KnownTags() {
			if tg == val {
				tags2 = append(tags2, tg)
			}
			if tg != val {
				continue
			}
		}
		if i != len(tags)-1 {
			base += val + "_"
		}
	}
	base = strings.TrimSuffix(base, "_")
	return base, tags2
}

func KnownTags() []string {
	return []string{
		"HD",
		"SD",
		"43",
		"AUDIOENG20",
		"AUDIORUS20",
		"AUDIOENG51",
		"AUDIORUS51",
		"TRL",
		"Proxy",
		"ar2",
		"ar6",
		"ar2e2",
		"ar2e6",
		"ar6e2",
		"ar6e6",
		"rus51",
		"rus20",
		"eng51",
		"eng20",
	}
}

func TrimLoudnormPrefix(name string) (string, error) {
	newName := ""
	if err := validateOldname(name); err != nil {
		return newName, err
	}
	return newName, nil
}

func validateOldname(name string) error {
	if strings.TrimSuffix(name, ".ac3") == name {
		return fmt.Errorf("invalid name [%v] - is not ac3 file", name)
	}
	data := strings.Split(name, "__")
	if len(data) != 2 {
		return fmt.Errorf("invalid name [%v] - does not contain '__'", name)
	}
	/////HD20

	return nil
}

/*
Дата:  212-1099
Место сбора: Pax Rulin (A402231-E)
Миссия: Вызволить VIP находящегося на территории частной коммерческой организации на планете Magen.
Ожидаемый состав отряда:
  -специалист по стрелковому вооружению
  -боевой инженер
  -полевой разведчик
  -медик
  -пилот антигравитационной техники
Оплата: 20  000 Cr каждому участнику перед началом миссии. 50 000 Сr каждому при доставке VIP'а к заказчику.
Комментарии: Наёмный отряд будет снабжен экипировкой, транспортом и актуальными разведданными.  Дополнительные вопросы от Исполнителей не приветствуются.
*/
