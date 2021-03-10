package namedata

import (
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
	return strings.Join(pathData[0:len(pathData)-1], "\\")
}

func RetrieveShortName(path string) string {
	return shortFileName(path)
}

func RetrieveBase(path string) string {
	fileName := shortFileName(path)
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
