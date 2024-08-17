package tag

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	IN_MARKER  = "--"
	IN_SEP     = "--"
	OUT_MARKER = "__"
	OUT_SEP    = "_"
)

type collection struct {
	FileName       string
	CollectionType TagsType
	TagWithKey     map[TagKey]Tag
}

type TagsType string

var IN_FILE TagsType = "IN"
var OUT_FILE TagsType = "OUT"
var CONFLICTING TagsType = "Conflicting"
var UNDEFINED TagsType = "Undefined_File"

/*
обычный файл --> IN_FILE --> OUT_File

*/

func nameIsLatinOnly(value string) bool {
	value = strings.ToLower(value)
	for _, v := range strings.Split(value, "") {
		switch v {
		default:
			return false
		case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z":
		case "_", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		}
	}
	return true
}

func NewCollection(path string) (*collection, error) {
	name := filepath.Base(path)
	if !nameIsLatinOnly(name) {
		return nil, fmt.Errorf("filename contains spaces or non-latin symbols '%v'", name)
	}
	tcol := collection{}
	tcol.FileName = name
	tagTypes := tagsTypeExpected(name)
	switch tagTypes {
	default:
		return nil, fmt.Errorf("unidentified tag type marker present '%v'", tagTypes)
	case CONFLICTING:
		return nil, fmt.Errorf("both tag markers are present")
	case IN_FILE, OUT_FILE, UNDEFINED:
		tcol.CollectionType = tagTypes
	}
	tcol.TagWithKey = make(map[TagKey]Tag)
	return &tcol, nil
}

func (tc *collection) parse() ([]Tag, error) {
	parsedTags := []Tag{}
	switch tc.CollectionType {
	case IN_FILE:
		parsedTags = append(parsedTags, parse_INFILE(tc.FileName)...)
	case OUT_FILE:
		panic("parse OUT_FILE not implemented")
	default:
		return parsedTags, fmt.Errorf("can't parse filename with '%v' collection type", tc.CollectionType)
	}
	return parsedTags, nil
}

func parse_OUTFILE(name string) []Tag {
	parsedTags := []Tag{}

	parsedTags = append(parsedTags,
		outSeasonTag(name),
		outEpisodeTag(name),
		prtTag(name),
		outVideoTag(name), outAudioTag(name), outSrtTag(name), outRevisionTag(name))

	return parsedTags
}

func parse_INFILE(name string) []Tag {
	parsedTags := []Tag{}
	parts := strings.Split(name, IN_SEP)
	for i, part := range parts {
		if i == 0 {
			parsedTags = append(parsedTags, New(BASE, part))
			continue
		}
		if i == len(parts)-1 {
			continue
		}
		for _, tag := range []Tag{
			typeTag(part),
			seasonTag(part),
			episodeTags(part),
			prtTag(part),
			videoTag(part),
			srtTag(part),
			revisionTag(part),
		} {
			if tag.Key != NoTagKey {
				parsedTags = append(parsedTags, tag)
				continue
			}
		}
	}
	return parsedTags

}

func typeTag(val string) Tag {
	if val == string(FILM) {
		return New(TYPE, string(FILM))
	}
	if val == string(SER) {
		return New(TYPE, string(SER))
	}
	if val == string(TRL) {
		return New(TYPE, string(TRL))
	}
	return NoTag
}

func tagsTypeExpected(base string) TagsType {
	if strings.Contains(base, IN_MARKER) && strings.Contains(base, OUT_MARKER) {
		return CONFLICTING
	}
	if strings.Contains(base, IN_MARKER) {
		return IN_FILE
	}
	if strings.Contains(base, OUT_MARKER) {
		return OUT_FILE
	}
	return UNDEFINED
}

// func infileSeasonEpisodeTags(s string) (Tag, Tag) {
// 	re := regexp.MustCompile(`(s[0-9]{1,}e[0-9]{1,})`)
// 	tagsStr := re.FindString(s)
// 	if tagsStr == "" {
// 		return NoTag, NoTag
// 	}
// 	parts := strings.Split(tagsStr, "e")
// 	parts[0] = strings.TrimPrefix(parts[0], "s")
// 	sTag := NoTag
// 	if _, err := strconv.Atoi(parts[0]); err == nil {
// 		sTag = New(SEASON, parts[0])
// 	}
// 	eTag := NoTag
// 	if _, err := strconv.Atoi(parts[1]); err == nil {
// 		eTag = New(EPISODE, parts[1])
// 	}
// 	return sTag, eTag
// }

func seasonTag(s string) Tag {
	re := regexp.MustCompile(`(s[0-9]{1,}e[0-9]{1,})`)
	tagsStr := re.FindString(s)
	if tagsStr == "" {
		return NoTag
	}
	parts := strings.Split(tagsStr, "e")
	parts[0] = strings.TrimPrefix(parts[0], "s")
	sTag := NoTag
	if _, err := strconv.Atoi(parts[0]); err == nil {
		sTag = New(SEASON, parts[0])
	}

	return sTag
}

func episodeTags(s string) Tag {
	re := regexp.MustCompile(`(s[0-9]{1,}e[0-9]{1,})`)
	tagsStr := re.FindString(s)
	if tagsStr == "" {
		return NoTag
	}
	parts := strings.Split(tagsStr, "e")
	eTag := NoTag
	if _, err := strconv.Atoi(parts[1]); err == nil {
		eTag = New(EPISODE, parts[1])
	}

	return eTag
}

func prtTag(s string) Tag {
	re := regexp.MustCompile(`(PRT[0-9]{6,})`)
	pTagValue := re.FindString(s)
	if pTagValue == "" {
		return NoTag
	}
	return New(PRT, pTagValue)
}

func videoTag(s string) Tag {
	str := strings.ToUpper(s)
	if str == string(UHD) {
		return New(VIDEO, string(UHD))
	}
	if str == string(HD) {
		return New(VIDEO, string(HD))
	}
	if str == string(SD) {
		return New(VIDEO, string(SD))
	}
	return NoTag
}

func outSeasonTag(s string) Tag {
	re := regexp.MustCompile(`(s[0-9]{1,}_[0-9]{1,})`)
	tagsStr := re.FindString(s)
	if tagsStr == "" {
		return NoTag
	}
	parts := strings.Split(tagsStr, "_")
	parts[0] = strings.TrimPrefix(parts[0], "s")
	sTag := NoTag
	if _, err := strconv.Atoi(parts[0]); err == nil {
		sTag = New(SEASON, parts[0])
	}
	return sTag
}

func outEpisodeTag(s string) Tag {
	sTag := outSeasonTag(s)
	if sTag == NoTag {
		return NoTag
	}
	re := regexp.MustCompile(`(s[0-9]{1,}_[0-9]{1,})`)
	tagsStr := re.FindString(s)
	if tagsStr == "" {
		return NoTag
	}
	parts := strings.Split(tagsStr, "_")
	parts[0] = strings.TrimPrefix(parts[0], "s")
	eTag := NoTag
	if _, err := strconv.Atoi(parts[1]); err == nil {
		eTag = New(EPISODE, parts[1])
	}
	return eTag
}

func outVideoTag(s string) Tag {
	str := strings.ToUpper(s)
	if strings.Contains(str, OUT_MARKER+string(UHD)) {
		return New(VIDEO, string(UHD))
	}
	if strings.Contains(str, OUT_MARKER+string(HD)) {
		return New(VIDEO, string(HD))
	}
	if strings.Contains(str, OUT_MARKER+string(SD)) {
		return New(VIDEO, string(SD))
	}
	return NoTag
}

func outAudioTag(s string) Tag {
	re := regexp.MustCompile(`(AUDIO[A-Z]{3,}[0-9]{2,})`)
	found := re.FindString(s)
	if found == "" {
		return NoTag
	}
	return New(AUDIO, found)
}

func srtTag(s string) Tag {
	str := strings.ToUpper(s)
	if str == string(SUB) {
		return New(SRT, string(SUB))
	}
	if str == string(HARDSUB) {
		return New(SRT, string(HARDSUB))
	}
	return NoTag
}

func outSrtTag(s string) Tag {
	if strings.Contains(s, OUT_SEP+string(SUB)) {
		return New(SRT, string(SUB))
	}
	if strings.Contains(s, OUT_SEP+string(HARDSUB)) {
		return New(SRT, string(HARDSUB))
	}
	if strings.HasSuffix(s, ".srt") {
		return New(SRT, string(SUB))
	}
	return NoTag
}

func revisionTag(s string) Tag {
	str := strings.ToUpper(s)
	if !strings.HasPrefix(str, "R") {
		return NoTag
	}
	str = strings.TrimPrefix(str, "R")
	if _, err := strconv.Atoi(str); err == nil {
		return New(REVISION, fmt.Sprintf("R%v", str))
	}
	return NoTag
}

func outRevisionTag(s string) Tag {
	s = strings.ToUpper(s)
	re := regexp.MustCompile(fmt.Sprintf(`(%vR[0-9]{1,})`, OUT_SEP))
	found := re.FindString(s)
	if found == "" {
		return NoTag
	}
	value := strings.TrimPrefix(found, OUT_SEP+"R")
	return New(REVISION, value)
}

//AddUnique - adds tags that's key are not exists to a collection.
//Special tag 'NoTag' - will be ignored.
//Error is reserved for custom implementations.
func (tc *collection) AddUnique(tags ...Tag) error {
	for _, tag := range tags {
		if tag.Key == NoTagKey {
			continue
		}
		if _, ok := tc.TagWithKey[tag.Key]; ok {
			continue
		}
		tc.TagWithKey[tag.Key] = tag
	}
	return nil
}

type Collection interface {
	AddTags(...Tag) error
	VerifyTags() error
	Base() string
	InputTag(string) string
	OutptTag(string) string
}
