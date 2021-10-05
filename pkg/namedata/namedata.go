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
	base, vid, aud, ebur := ungroupeName(name)
	if aud == "" {
		return newName, fmt.Errorf("audio tag can't be detected '%v'", name)
	}
	if ebur == "" {
		return name, nil
	}
	//fmt.Println("UNGROUPE:", base, vid, aud, ebur)
	switch {
	default:
		newName = ""
		return newName, fmt.Errorf("New name undecided for '%v'", name)

	case (vid == "hd" || vid == "4k") && strings.Contains(aud, "51") && strings.Contains(ebur, "-stereo"):
		vid = "sd"
		aud = strings.TrimSuffix(aud, "51") + "20"
		newName = base + "__" + vid + "_" + aud + ".ac3"
	case (vid == "hd" || vid == "4k") && strings.Contains(aud, "51") && !strings.Contains(ebur, "-stereo"):
		newName = base + "__" + vid + "_" + aud + ".ac3"
	case vid == "sd" && strings.Contains(aud, "51"):
		aud = strings.TrimSuffix(aud, "51") + "20"
		newName = base + "__" + vid + "_" + aud + ".ac3"
	case strings.Contains(aud, "20"):
		newName = base + "__" + vid + "_" + aud + ".ac3"
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
	if len(strings.Split(data[1], "_")) != 2 {
		return fmt.Errorf("invalid name [%v] - can not define audio and/or video tags", name)
	}
	/////HD20

	return nil
}

func ungroupeName(name string) (base, video, audio, ebur string) {
	data := strings.Split(name, "__")
	base = data[0]
	data2 := strings.Split(data[1], "_")
	video = data2[0]
	if strings.Contains(data2[1], "-ebur128-stereo.ac3") {
		audio = strings.TrimSuffix(data2[1], "-ebur128-stereo.ac3")
		ebur = "-ebur128-stereo"
		return
	}
	if strings.Contains(data2[1], "-ebur128.ac3") {
		audio = strings.TrimSuffix(data2[1], "-ebur128.ac3")
		ebur = "-ebur128"
		return
	}
	if strings.Contains(data2[1], "51.ac3") || strings.Contains(data2[1], "20.ac3") {
		audio = data2[1]
	}
	return
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
