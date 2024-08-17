package tag

// import (
// 	"errors"
// 	"fmt"
// 	"path/filepath"
// 	"regexp"
// 	"strconv"
// 	"strings"
// )

// const (
// 	MarkerInput     TagMarker = "--"
// 	MarkerOutput    TagMarker = "__"
// 	MarkerBad       TagMarker = "[BAD MARKER]"
// 	NoMarkerPresent TagMarker = "[NO MARKER]"
// 	hardcode                  = "hardcode"
// )

// type TagMarker string

// var ErrTagNotFound = errors.New("no tag detected")
// var ErrBothMarkersPresent = errors.New("both markers present")

// // SearchContentType - searching episode tag in filepath provided.
// // Returns error if TagMarker is not Input or tags not found.
// func SearchContentType(name string) (Tag, error) {
// 	name = filepath.Base(name)
// 	switch markerPresent(name) {
// 	default:
// 		return NoTag, fmt.Errorf("no Input Marker present")
// 	case MarkerBad:
// 		return NoTag, ErrBothMarkersPresent
// 	case MarkerInput:
// 	}
// 	if strings.HasPrefix(name, string(FILM)) {
// 		return TypeFILMTag, nil
// 	}
// 	if strings.HasPrefix(name, string(TRL)) {
// 		return TypeTRLTag, nil
// 	}
// 	if strings.HasPrefix(name, string(SER)) {
// 		return TypeSERTag, nil
// 	}
// 	return NoTag, ErrTagNotFound
// }

// // SearchEpisode - searching episode tag in filepath provided.
// // Returns error if tags found != 1
// func SearchEpisode(name string) (Tag, error) {
// 	name = filepath.Base(name)
// 	re := regexp.MustCompile(`(s[0-9]{1,}e[0-9]{1,})`)
// 	re.FindAllString(name, -1)
// 	episodeTags := re.FindAllString(name, -1)
// 	if len(episodeTags) == 0 {
// 		return NoTag, ErrTagNotFound
// 	}
// 	if !hasSameElements(episodeTags) {
// 		return NoTag, fmt.Errorf("different episode tags detected")
// 	}
// 	return CreateEpisodeTag(episodeTags[0]), nil
// }

// func hasSameElements(slice []string) bool {
// 	for i, s1 := range slice {
// 		for j, s2 := range slice {
// 			if j <= i {
// 				continue
// 			}
// 			if s1 != s2 {
// 				return false
// 			}
// 		}
// 	}
// 	return true
// }

// // SearchRevision - searching episode tag in filepath provided.
// // Returns error if TagMarker is BadMarker or tag not found.
// func SearchRevision(name string) (Tag, error) {
// 	name = filepath.Base(name)
// 	tagNum := 0
// 	switch markerPresent(name) {
// 	case MarkerBad:
// 		return NoTag, ErrBothMarkersPresent
// 	case MarkerOutput:
// 		tagNum = findRevisionNumberOutput(name)
// 	case MarkerInput, NoMarkerPresent:
// 		tagNum = findRevisionNumberInput(name)
// 	}
// 	if tagNum == 0 {
// 		return NoTag, ErrTagNotFound
// 	}
// 	return CreateRevisionTag(tagNum), nil
// }

// func findRevisionNumberInput(name string) int {
// 	re := regexp.MustCompile(`(_R[0-9]{1,}.)`)
// 	tag := re.FindString(name)
// 	tag = strings.TrimPrefix(tag, "_")
// 	tag = strings.TrimSuffix(tag, ".")
// 	tag = strings.TrimPrefix(tag, "R")
// 	tagNum, _ := strconv.Atoi(tag)
// 	return tagNum
// }

// func findRevisionNumberOutput(name string) int {
// 	re := regexp.MustCompile(`(_R[0-9]{1,}_)`)
// 	tag := re.FindString(name)
// 	tag = strings.TrimPrefix(tag, "_")
// 	tag = strings.TrimSuffix(tag, "_")
// 	tag = strings.TrimPrefix(tag, "R")
// 	tagNum, _ := strconv.Atoi(tag)
// 	return tagNum
// }

// // SearchVideo - searching video size tag in filepath provided.
// // Returns error if TagMarker is BadMarker or tag not found.
// func SearchVideo(path string) (Tag, error) {
// 	name := filepath.Base(path)
// 	tagDetected := ""
// 	marker := markerPresent(name)
// 	switch marker {
// 	case MarkerBad:
// 		return NoTag, ErrBothMarkersPresent
// 	case MarkerOutput, MarkerInput, NoMarkerPresent:
// 		possibleTags := listVideoTags(hardcode) //TODO: make rule selection mechanism
// 		tagDetected = findVideoTag(name, marker, possibleTags)
// 	}
// 	if tagDetected == "" {
// 		return NoTag, ErrTagNotFound
// 	}
// 	return CreateVideoTag(tagDetected), nil
// }

// func listVideoTags(rule string) []string {
// 	switch rule {
// 	default:
// 		panic(fmt.Sprintf("rule '%v' for listVideoTags() is not implemented"))
// 		return nil
// 	case hardcode:
// 		return hardcodeVideoTags()
// 	}
// }

// func hardcodeVideoTags() []string {
// 	return []string{"HD", "SD", "4K"}
// }

// func findVideoTag(filename string, marker TagMarker, possibleTags []string) string {
// 	for _, pTag := range possibleTags {
// 		if strings.Contains(filename, string(marker)+pTag) {
// 			return pTag
// 		}
// 	}
// 	return ""
// }

// // SearchAudio - searching audio tag in filepath provided.
// // Returns error if TagMarker is not output or tag not found.
// func SearchAudio(path string) (Tag, error) {
// 	name := filepath.Base(path)
// 	tagDetected := ""
// 	switch markerPresent(name) {
// 	case MarkerBad, MarkerInput, NoMarkerPresent:
// 		return NoTag, fmt.Errorf("audio tag is in output only")
// 	case MarkerOutput:
// 		tagDetected = findAudioTag(name)
// 	}
// 	if tagDetected == "" {
// 		return NoTag, ErrTagNotFound
// 	}
// 	lang, layout := decoupleAudioTagContent(tagDetected)
// 	return CreateAudioTag(lang, layout), nil
// }

// func findAudioTag(filename string) string {
// 	re := regexp.MustCompile(`(_AUDIO*[0-9]{1,})`)
// 	tag := re.FindString(filename)
// 	return tag
// }

// func decoupleAudioTagContent(atag string) (string, string) {
// 	atag = strings.TrimSuffix(atag, "_AUDIO")
// 	lang := filterPrefix(atag)
// 	layout := strings.TrimPrefix(atag, lang)
// 	return lang, layout
// }

// func filterPrefix(str string, possiblePrefixes ...string) string {
// 	for _, prefix := range possiblePrefixes {
// 		if strings.HasPrefix(str, prefix) {
// 			return prefix
// 		}
// 	}
// 	return ""
// }

// // TODO: draw prefixes from language gost library
// func listLanguagAbbreviations() []string {
// 	return []string{
// 		"rus",
// 		"eng",
// 		"chi",
// 		"heb",
// 		"fra",
// 		"spa",
// 		"dan",
// 		"kaz",
// 	}
// }

// func markerPresent(name string) TagMarker {
// 	name = filepath.Base(name)
// 	inputDetected := false
// 	outputDetected := false
// 	if strings.Contains(name, string(MarkerInput)) {
// 		inputDetected = true
// 	}
// 	if strings.Contains(name, string(MarkerOutput)) {
// 		inputDetected = true
// 	}
// 	if inputDetected && !outputDetected {
// 		return MarkerInput
// 	}
// 	if !inputDetected && outputDetected {
// 		return MarkerOutput
// 	}
// 	if inputDetected && outputDetected {
// 		return MarkerBad
// 	}
// 	return NoMarkerPresent
// }
