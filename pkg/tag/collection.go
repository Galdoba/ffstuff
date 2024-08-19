package tag

import (
	"errors"
	"fmt"
	"path/filepath"
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
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		case "_", ".", "-":
		}
	}
	return true
}

var ErrConflictingMarkers = errors.New("parsing forbiden: conflicting markers detected")

func NewCollection(path string) (*collection, error) {
	name := filepath.Base(path)
	if !nameIsLatinOnly(name) {
		return nil, fmt.Errorf("filename contains spaces or non-latin symbols '%v'", name)
	}
	tcol := collection{}
	tcol.FileName = name
	tcol.TagWithKey = make(map[TagKey]Tag)

	tagTypes := tagsTypeExpected(name)
	switch tagTypes {
	default:
		return nil, fmt.Errorf("unidentified tag type marker present '%v'", tagTypes)
	case CONFLICTING:
		return nil, ErrConflictingMarkers
	case IN_FILE, OUT_FILE, UNDEFINED:
		tcol.CollectionType = tagTypes
	}

	parsedTags, err := parse(&tcol)
	if err == nil {
		_, err = tcol.Add(parsedTags...)
	}
	return &tcol, err
}

var ErrPromptExpected = errors.New("parsing forbiden: expecting prompt to fill collection")

func parse(tc *collection) ([]Tag, error) {
	parsedTags := []Tag{}
	switch tc.CollectionType {
	case IN_FILE:
		parsedTags = append(parsedTags, parse_INFILE(tc.FileName)...)
	case OUT_FILE:
		parsedTags = append(parsedTags, parse_OUTFILE(tc.FileName)...)
	case UNDEFINED:
		return parsedTags, ErrPromptExpected
	default:
		return parsedTags, fmt.Errorf("can't parse filename with '%v' collection type", tc.CollectionType)
	}
	return parsedTags, nil
}

func parse_INFILE(name string) []Tag {
	parsedTags := []Tag{}
	parts := strings.Split(name, IN_SEP)
	for i, part := range parts {
		switch i {
		case 0:
			parsedTags = append(parsedTags, New(BASE, part))
		case len(parts) - 1:
			parsedTags = append(parsedTags, New(ORIGIN, part))
		default:
			for _, tag := range parseNextINFILETag(part) {
				parsedTags = append(parsedTags, tag)
			}
		}
	}
	return parsedTags

}

func parseNextINFILETag(part string) []Tag {
	tags := []Tag{}
	for _, tag := range []Tag{
		parse_in_typeTag(part),
		parse_in_seasonTag(part),
		parse_in_episodeTags(part),
		parse_in_prtTag(part),
		parse_in_videoTag(part),
		parse_in_srtTag(part),
		parse_in_revisionTag(part),
	} {
		if tag.Key != NoTagKey {
			tags = append(tags, tag)
		}
	}
	return tags
}

func parse_OUTFILE(name string) []Tag {
	parsedTags := []Tag{}
	parsedTags = append(parsedTags,
		parse_out_SeasonTag(name),
		parse_out_EpisodeTag(name),
		parse_out_prtTag(name),
		parse_out_VideoTag(name),
		parse_out_AudioTag(name),
		parse_out_SrtTag(name),
		parse_out_RevisionTag(name),
		parse_out_Outname(name),
		parse_out_Extention(name),
	)
	return parsedTags
}

func parse_UNKNOWN(name string) []Tag {
	parsedTags := []Tag{}
	parsedTags = append(parsedTags,
		parse_out_SrtTag(name),
		parse_out_RevisionTag(name),
		parse_out_Extention(name),
	)
	return parsedTags
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

type TagCollection interface {
	Add(...Tag) (int, error)
	Review() (string, TagsType, []Tag)
	//ProjectOutname() string
}

// Add - adds tags that's key are not exists to a collection.
// Special tag 'NoTag' - will be ignored.
// Return number of tags added and error, if any tags were rejected.
func (tc *collection) Add(tags ...Tag) (int, error) {
	errText := ""
	added := 0
	rejected := 0
	for _, tag := range tags {
		if tag.Key == NoTagKey {
			continue
		}
		if _, ok := tc.TagWithKey[tag.Key]; ok {
			rejected++
			errText += fmt.Sprintf("%v, ", tag)
			continue
		}
		tc.TagWithKey[tag.Key] = tag
		added++
	}
	if rejected > 0 {
		errText = strings.TrimSuffix(errText, ", ")
		return added, fmt.Errorf("tags rejected: %v", errText)
	}

	return added, nil
}

// Review - returns collection's filename, collectionType and tags (including ones with missing type) in orderly list
func (tc *collection) Review() (string, TagsType, []Tag) {
	keys := []TagKey{BASE, TYPE, SEASON, EPISODE, PRT, VIDEO, AUDIO, SRT, REVISION, COMMENTS, EXT, ORIGIN, OUTNAME}
	tags := []Tag{}
	for _, key := range keys {
		if val, ok := tc.TagWithKey[key]; ok {
			tags = append(tags, val)
			continue
		}
		//tags = append(tags, New(key, NoTag.Value))
	}
	return tc.FileName, tc.CollectionType, tags
}

// ProjectOutfileName - create name for inputfile for demuxing process
// Return empty string if colection if not type 'UNDEFINED'
func ProjectInfileName(c TagCollection, newTags ...Tag) string {
	file, colType, oldTags := c.Review()
	if colType != UNDEFINED {
		return ""
	}
	newCol, _ := NewCollection(file)
	newCol.Add(oldTags...)
	newCol.Add(newTags...)
	//	IN_FILE = {BASE}{TYPE}[SEASON][EPISODE][PRT][VIDEO][SRT][REVISION]{MARKER}{ORIGIN}
	iName := ""
	if val, ok := newCol.TagWithKey[BASE]; ok {
		iName += val.String()
	}
	if val, ok := newCol.TagWithKey[TYPE]; ok {
		iName += IN_SEP + val.String()
	}
	if val, ok := newCol.TagWithKey[SEASON]; ok {
		iName += IN_SEP + "s" + val.String()
	}
	if val, ok := newCol.TagWithKey[EPISODE]; ok {
		iName += "e" + val.String()
	}
	if val, ok := newCol.TagWithKey[PRT]; ok {
		iName += IN_SEP + val.String()
	}
	if val, ok := newCol.TagWithKey[VIDEO]; ok {
		iName += IN_SEP + val.String()
	}
	if val, ok := newCol.TagWithKey[AUDIO]; ok {
		iName += IN_SEP + val.String()
	}
	if val, ok := newCol.TagWithKey[SRT]; ok {
		iName += IN_SEP + val.String()
	}
	if val, ok := newCol.TagWithKey[REVISION]; ok {
		iName += IN_SEP + val.String()
	}
	if val, ok := newCol.TagWithKey[COMMENTS]; ok {
		iName += IN_SEP + val.String()
	}
	iName += IN_MARKER
	if val, ok := newCol.TagWithKey[ORIGIN]; ok {
		iName += val.String()
	}
	//iName = strings.ReplaceAll(iName, "___", "__")
	return iName
}

// ProjectOutfileName - create name for outputfile for demuxing process
// Return empty string if colection if not type 'IN_FILE'
func ProjectOutfileName(c TagCollection, newTags ...Tag) string {
	file, colType, oldTags := c.Review()
	if colType != IN_FILE {
		return ""
	}
	newCol, _ := NewCollection(file)
	newCol.Add(oldTags...)
	newCol.Add(newTags...)
	oName := ""
	if val, ok := newCol.TagWithKey[BASE]; ok {
		oName += val.String()
	}
	if val, ok := newCol.TagWithKey[SEASON]; ok {
		oName += OUT_SEP + "s" + val.String()
	}
	if val, ok := newCol.TagWithKey[EPISODE]; ok {
		oName += OUT_SEP + val.String()
	}
	if val, ok := newCol.TagWithKey[PRT]; ok {
		oName += OUT_SEP + val.String()
	}
	oName += OUT_MARKER
	if val, ok := newCol.TagWithKey[VIDEO]; ok {
		oName += OUT_SEP + val.String()
	}
	if val, ok := newCol.TagWithKey[AUDIO]; ok {
		oName += OUT_SEP + val.String()
	}
	if val, ok := newCol.TagWithKey[SRT]; ok {
		oName += OUT_SEP + val.String()
	}
	if val, ok := newCol.TagWithKey[REVISION]; ok {
		oName += OUT_SEP + val.String()
	}
	if val, ok := newCol.TagWithKey[COMMENTS]; ok {
		oName += OUT_SEP + val.String()
	}
	if val, ok := newCol.TagWithKey[EXT]; ok {
		oName += "." + val.String()
	}
	oName = strings.ReplaceAll(oName, "___", "__")
	return oName
}

//другая библиотека
// func (tc *collection) FillWithPrompt() error {
// 	if tc.CollectionType != UNDEFINED {
// 		return fmt.Errorf("FillWithPrompt method must be used with unidentified collections only")
// 	}
// 	baseSuggestions := suggestBase(tc.FileName)
// 	base, err := operator.Select("select base", baseSuggestions...)
// 	if err != nil {
// 		return err
// 	}
// 	tc.Add(New(BASE, base))
// 	return nil
// }
//
// func suggestBase(name string) []string {
// 	parts := strings.Split(name, "_")
// 	suggestions := []string{}
// 	if len(parts) == 1 {
// 		return []string{name, "NONE"}
// 	}
// 	for i := range parts {
// 		if i == 0 {
// 			continue
// 		}
// 		suggestions = append(suggestions, strings.Join(parts[0:i], "_"))
// 	}
// 	return suggestions
// }
