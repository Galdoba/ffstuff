package mxprjct

import (
	"fmt"
	"strings"
)

type muxProject struct {
	inputV         string
	inputA         []string
	inputS         string
	expectedResult string
}

func Create(inputPaths []string) (*muxProject, error) {
	mp := muxProject{}
	if err := mp.fillInput(inputPaths); err != nil {
		return &mp, err
	}
	mp.sortAudio()
	fmt.Println(mp, "------")
	mp.projectResult()
	fmt.Println(mp.expectedResult, "------")
	return &mp, nil
}

func (mp *muxProject) fillInput(inputPaths []string) error {
	for _, path := range inputPaths {
		pathUP := strings.ToUpper(path)
		decodedData := decodeWorkFileName(path)
		if decodedData.err != nil {
			//return fmt.Errorf("input file not formated: %v", decodedData.err.Error())
		}
		switch {
		default:
			return fmt.Errorf("unknown input type [%v]", path)
		case strings.Contains(pathUP, ".MP4"): //главный параметр наличие видеопотока
			mp.inputV = path
		case strings.Contains(pathUP, ".M4A"): //сортируем каналы звука следующим этапом
			mp.inputA = append(mp.inputA, path)
		case strings.Contains(pathUP, ".SRT"): //
			mp.inputS = path
		}
	}
	return nil
}

func (mp *muxProject) sortAudio() {
	newOrder := []string{}
	newOrder = append(newOrder, pickTag(mp.inputA, "rus"))
	newOrder = append(newOrder, pickTag(mp.inputA, "eng"))
	newOrder = append(newOrder, pickTag(mp.inputA, "qqq"))
	mp.inputA = newOrder
}

func pickTag(allTags []string, tag string) string {
	for _, aud := range allTags {
		if strings.Contains(aud, "_"+tag+"51.") || strings.Contains(aud, "_"+tag+"20.") {
			return aud
		}
	}
	return ""
}

func (mp *muxProject) projectResult() error {
	result := ""
	dVid := decodeWorkFileName(mp.inputV)
	fmt.Println(dVid)
	result = dVid.base + "_a"
	for _, aud := range mp.inputA {
		dAud := decodeWorkFileName(aud)
		if dVid.base != dAud.base {
			return fmt.Errorf("not the same base")
		}
		for k, v := range dAud.tags {
			if v {
				result += shortenTag(k)
			}
		}
	}
	if mp.inputS != "" {
		result += "_sr"
	}
	result += ".mp4"
	mp.expectedResult = result
	return nil
}

func shortenTag(t string) string {
	switch t {
	default:
		return ""
	case "_rus20":
		return "r2"
	case "_rus51":
		return "r6"
	case "_eng20":
		return "r2"
	case "_eng51":
		return "r6"
	case "_qqq20":
		return "r2"
	case "_qqq51":
		return "r6"
	}
}

type workFileData struct {
	base      string
	directory string
	tags      map[string]bool
	err       error
}

func decodeWorkFileName(path string) workFileData {
	wfd := workFileData{}
	wfd.tags = make(map[string]bool)
	segm := strings.Split(path, "\\")
	for i, s := range segm {
		if i < len(segm)-1 {
			wfd.directory = wfd.directory + s + "\\"
		}
	}
	bd := strings.Split(segm[len(segm)-1], ".")
	for i, b := range bd {
		if i < len(bd)-1 {
			wfd.base = wfd.base + b + "."
		}
	}
	wfd.base = strings.TrimSuffix(wfd.base, ".")
	switch {
	case strings.Contains(wfd.base, "__hd"):
		wfd.tags["HD"] = true
	case strings.Contains(wfd.base, "__sd"):
		wfd.tags["SD"] = true
	case strings.Contains(wfd.base, "_smk"):
		wfd.tags["_smk"] = true
	}
	for _, tag := range soundTags() {
		if strings.Contains(path, tag) {
			wfd.tags[tag] = true
		}
	}
	if len(wfd.tags) == 0 {
		wfd.err = fmt.Errorf("no tags detected")
	}
	return wfd
}

func soundTags() []string {
	tags := []string{}
	for _, lang := range []string{"_rus", "_eng", "_qqq"} {
		for _, chn := range []string{"20", "51"} {
			tags = append(tags, lang+chn)
		}
	}
	return tags
}
