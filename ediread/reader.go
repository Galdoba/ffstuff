package ediread

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

type ediReader interface {
	ReadEDI(string) (string, error)
	AddToBatch(string, []string) error
}

type edlData struct {
	entry []string
}

//Entry - возвращает базовую информацию о клипах после прочтения EDL
func (edl *edlData) Entry() []string {
	return edl.entry
}

func NewEdlData(path string) (edlData, error) {
	eDt := edlData{}
	file, err := os.Open(path)
	if err != nil {
		if strings.Contains(err.Error(), "The system cannot find the file specified") {
			return eDt, errors.New("EDL-file not Found")
		}
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	text := ""
	for scanner.Scan() {
		text += scanner.Text() + " "
		if scanner.Text() == "" {
			eDt.entry = append(eDt.entry, text)
			text = ""
		}
		//fmt.Println(scanner.Text())
	}
	eDt.entry = append(eDt.entry, text)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	eDt = cleanEDLdata(eDt)
	return eDt, nil
}

func folder(path string) string {
	flr := strings.Split(path, "\\")
	flr = flr[0 : len(flr)-1]
	return strings.Join(flr, "\\") + "\\"
}

func cleanEDLdata(e edlData) edlData {
	e2 := edlData{}
	for _, val := range e.entry {
		if strings.Contains(val, "FROM CLIP") {
			e2.entry = append(e2.entry, val)
		}
	}
	return e2
}
