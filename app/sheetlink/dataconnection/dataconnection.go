package dataconnection

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/translit"
)

type dataConnection struct {
	filenames   []string
	projectType int
	editName    string
	destpath    string
	season      int
	episode     int
	agent       string
	releseDate  string
	comment     string
}

func New(filenames ...string) *dataConnection {
	dc := dataConnection{}
	dc.filenames = filenames
	return &dc
}

func (dc *dataConnection) Link(tabledata [][]string) error {
	//	fmt.Println(dc)
	matchValue := make(map[int][]string)

	//destPath := ""
	best := 0
	key := 0

	for l, data := range tabledata {
		if data[7] != "O" {

			r := regexp.MustCompile(`\d{2}\_\d{2}\_\d{4}`)
			found := r.FindString(translit.Transliterate(data[8]))
			if found != "" {
				dc.destpath = found
			}

		}
		//fmt.Println(data[8])
		transl := translit.Transliterate(data[8])
		//fmt.Println(transl)

		common := commonPrefix(dc.filenames[0], transl)
		//	fmt.Println(commonPrefix(data[8], transl))
		//fmt.Println(transl)
		if common != "" {
			fmt.Println(dc)
			fmt.Println("")
			fmt.Println(data[8])
			fmt.Println(transl)
			fmt.Println(common, "<-- common")
			matchValue[len(common)] = append(matchValue[len(common)], dc.destpath+"||"+common)
			if len(common) > best {
				best = len(common)
				key = l
			}
			fmt.Println("=====")
		}

	}
	fmt.Println("best match:")
	fmt.Println(matchValue[best])
	fmt.Println(tabledata[key])
	switch len(matchValue[best]) {
	case 1:
		fmt.Println("ONLY")
		dc.comment = tabledata[key][0]
		dc.editName = matchValue[best][0] + tag(dc.filenames[0])
		dc.agent = tabledata[key][13]
		fmt.Println(dc)
		return nil
	}

	fmt.Println(dc)
	return fmt.Errorf("not implemented")
}

func tag(s string) string {
	if strings.Contains(s, "--FILM--") {
		return "--FILM--"
	}
	if strings.Contains(s, "--SER--") {
		return "--SER--"
	}
	if strings.Contains(s, "--TRL--") {
		return "--TRL--"
	}
	return ""
}

func commonPrefix(s1, s2 string) string {
	sl1 := strings.Split(s1, "")
	sl2 := strings.Split(s2, "")
	slP := []string{}
	switch len(sl1) > len(sl2) {
	case true:
		for i, s := range sl1 {
			if i < len(sl2) && s == sl2[i] {
				slP = append(slP, s)
				continue
			}
			break
		}
	case false:
		for i, s := range sl2 {
			if s == sl1[i] {
				slP = append(slP, s)
				continue
			}
			break
		}
	}
	return strings.Join(slP, "")
}
