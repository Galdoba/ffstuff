package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/namedata"
	"github.com/fatih/color"
)

func onScreenBW(width int) (string, error) {
	infoMap, err := infoFields()
	if err != nil {
		return "", err
	}
	dirsSorted := dirsSort()

	scr := ""
	//max := len(infoMap)
	dir := ""
	for _, key := range keysByDirOrder(infoMap, dirsSorted) {
		if dir != infoMap[key].Dir() { //print new Dir as a header
			dir = infoMap[key].Dir()
			scr += dir + "\n"
		}

		s := format(infoMap[key], width) + "\n"

		scr += s
	}

	return scr, nil
}

func keysByDirOrder(infoMap map[int]*entry, dirsSorted []string) []int {
	resultSl := []int{}
	cntnt := []string{"TRL", "FILM", "SER", ""}
	for _, dir := range dirsSorted {
		for _, cnt := range cntnt {
			dirSl := []int{}
			for i := 0; i < len(infoMap); i++ {
				ed := namedata.EditForm(infoMap[i].file)
				if ed.ContentType() != cnt {
					continue
				}
				if dir == namedata.RetrieveDirectory(infoMap[i].file) {
					dirSl = append(dirSl, i)
				}
			}
			resultSl = append(resultSl, dirSl...)
		}
	}
	return resultSl
}

func dirsSort() []string {
	rtKeys := []string{}
	for k, _ := range Conf.Roots {
		rtKeys = append(rtKeys, k)
	}
	sort.Strings(rtKeys)
	rtVals := []string{}
	for _, key := range rtKeys {
		for _, val := range Conf.Roots[key] {
			rtVals = append(rtVals, val)
		}
	}
	return rtVals
}

func format(en *entry, width int) string {
	s := en.String()
	fmt.Println("fmt", en.file)
	if strings.HasPrefix(s, "*") {
		ss := strings.Split(s, "")
		fields := strings.Split(s, "WARNING:")
		w := width
		if len(fields) > 1 {
			w = width / 3 * 2
			ss = strings.Split(fields[0], "")
		}
		if len(ss) >= w {
			s = strings.Join(ss[:(w)-3], "") + ".."
		}
		return s
	}
	fl := strings.Fields(s)
	add := " "
	size := en.data["fSize"]
	for len(size) < 9 {
		size = " " + size
	}
	mPRF := colormPRF(en.data["mProfile"])

	//cntnt := en.data["mTag"]
	//maxNameLen := width - len(size) - len(mPRF) - 2
	s1 := strings.Join(fl[0:2], " ")
	s = s1 + add + mPRF + " " + size
	for len(s) != width {
		fmt.Println(".", s)
		switch len(s) >= width {
		case false:
			add += " "
			s = s1 + add + mPRF + " " + size
			fmt.Println(":", len(s), s)
		case true:
			s1 = trimEnd(s1)
			s = s1 + add + mPRF + " " + size
			fmt.Println("|", len(s), s)
		}

	}
	return s
}

func colormPRF(str string) string {
	switch str {
	default:
		return str
	case "1100-0", "1110-0":
		return color.HiGreenString(str)
	}
}

func trimEnd(str string) string {
	if len(str) < 3 {
		return ""
	}
	ltrs := strings.Split(str, "")
	ltrs = ltrs[0 : len(str)-3]
	ltrs = append(ltrs, ".", ".")
	return strings.Join(ltrs, "")
}

func infoFields() (map[int]*entry, error) {
	infoMap := make(map[int]*entry)
	lines, err := readLinesFromStorage()
	if err != nil {
		return nil, fmt.Errorf("readLinesFromStorage(): %v")
	}
	gathered := 0
	errs := []string{}
	for _, line := range lines {
		info, err := newEntry(line)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		if info == nil {
			continue
		}
		infoMap[gathered] = info
		gathered++
	}
	if len(errs) > 0 {
		errText := fmt.Sprintf("%v error(s) met:\n")
		for _, errtext := range errs {
			errText += errtext + "\n"
		}
		return infoMap, fmt.Errorf(errText)
	}
	return infoMap, nil
}

type entry struct {
	file string
	data map[string]string
}

func newEntry(line string) (*entry, error) {
	dl := entry{}
	dl.data = make(map[string]string)
	data := strings.Split(line, "  ")
	path := data[0]
	exists, err := exists(path)
	if !exists {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("file not exists: %v", path)
	}
	stats, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("os.Stat(path): %v", err.Error())
	}
	if stats.IsDir() {
		return nil, nil
	}
	dl.file = path
	fields := strings.Split(data[1], "|")
	for _, pair := range fields {
		if pair == "" {
			continue
		}
		info := strings.Split(pair, ":")
		dl.data[info[0]] = info[1]
	}
	return &dl, nil
}

func (dl *entry) keys() []string {
	ks := []string{}
	for k := range dl.data {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

func exists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil

	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil

	} else {
		return false, fmt.Errorf("file of schrodinger: file may or may not exist. err:=%v", err.Error())
	}
}

func (dl *entry) String() string {
	str := ""
	shortName := namedata.RetrieveShortName(dl.file)
	edit := namedata.EditForm(dl.file)
	editName := edit.EditName()
	switch editName {
	case "":
		str += "*"
		str += shortName
		warning := checkFileName(shortName)
		if warning != "" {
			str += " " + warning
		}

	default:
		str += dl.data["mTag"]
		if dl.data["mTag"] == "SER" {
			str += " "
		}
		str += " " + editName
		str += " " + dl.data["mProfile"]
		str += " " + dl.data["fSize"]
	}
	return str
}

func (dl *entry) Dir() string {
	return namedata.RetrieveDirectory(dl.file)
}

func checkFileName(name string) string {
	letters := strings.Split(name, "")
	for _, glyph := range letters {
		glyph = strings.ToLower(glyph)
		switch glyph {
		case " ", ")", "(", "'":
			return fmt.Sprintf("WARNING: Bad Name (contains |%v|)", glyph)
		case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
			"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
			"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", ".", "_", "-":
			continue
		default:
			return fmt.Sprintf("WARNING: Need Transliteration")
		}
	}
	return ""
}

func readLinesFromStorage() ([]string, error) {
	dataStore, err := os.OpenFile(storagePath+storageFile, os.O_RDONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("can't open %v", storagePath+storageFile)
	}
	defer dataStore.Close()
	lines := []string{}
	scanner := bufio.NewScanner(dataStore)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}
