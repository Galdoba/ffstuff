package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/namedata"
)

func onScreenBW(width int) (string, error) {
	infoMap, err := infoFields()
	if err != nil {
		return "", err
	}
	scr := ""
	max := len(infoMap)
	for key := 0; key < max; key++ {
		//fmt.Println(infoMap[key].String())
		s := infoMap[key].String()
		s = format(s, width) + "\n"

		scr += s
	}

	return scr, nil
}

func format(s string, width int) string {
	if strings.HasPrefix(s, "*") {
		if len(s) >= width {
			ss := strings.Split(s, "")
			s = strings.Join(ss[:width-3], "") + ".."
		}
		return s
	}
	fl := strings.Fields(s)
	add := ""
	for len(s) != width {

		if len(s) < width {
			add += " "
			s = fl[0] + " " + fl[1] + " " + add + fl[2] + " " + fl[3] + " " + fl[4]
		}
	}
	return s

}

func infoFields() (map[int]*entry, error) {
	infoMap := make(map[int]*entry)
	dataStore, err := os.OpenFile(storagePath+storageFile, os.O_RDONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("can't open %v", storagePath+storageFile)
	}
	lines := []string{}
	scanner := bufio.NewScanner(dataStore)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	dataStore.Close()
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

func checkFileName(name string) string {
	letters := strings.Split(name, "")
	for _, glyph := range letters {
		glyph = strings.ToLower(glyph)
		switch glyph {
		case " ", ")", "(", "'":
			return fmt.Sprintf("Bad Name (contains |%v|)", glyph)
		case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", ".", "_":
			continue
		default:
			return fmt.Sprintf("Need Transliteration")
		}
	}
	return ""
}
