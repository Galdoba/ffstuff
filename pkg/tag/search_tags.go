package tag

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	MarkerInput     TagMarker = "--"
	MarkerOutput    TagMarker = "__"
	MarkerBad       TagMarker = "[BAD MARKER]"
	NoMarkerPresent TagMarker = "[NO MARKER]"
)

type TagMarker string

var ErrTagNotFound = errors.New("no tag detected")
var ErrBothMarkersPresent = errors.New("both markers present")

func markerPresent(name string) TagMarker {
	name = filepath.Base(name)
	inputDetected := false
	outputDetected := false
	if strings.Contains(name, string(MarkerInput)) {
		inputDetected = true
	}
	if strings.Contains(name, string(MarkerOutput)) {
		inputDetected = true
	}
	if inputDetected && !outputDetected {
		return MarkerInput
	}
	if !inputDetected && outputDetected {
		return MarkerOutput
	}
	if inputDetected && outputDetected {
		return MarkerBad
	}
	return NoMarkerPresent
}

//SearchContentType - searching episode tag in filepath provided.
//Returns error if TagMarker is not Input or tags not found.
func SearchContentType(name string) (Tag, error) {
	name = filepath.Base(name)
	switch markerPresent(name) {
	default:
		return Tag{}, fmt.Errorf("no Input Marker present")
	case MarkerBad:
		return Tag{}, ErrBothMarkersPresent
	case MarkerInput:
	}
	if strings.HasPrefix(name, string(FILM)) {
		return TypeFILM(), nil
	}
	if strings.HasPrefix(name, string(TRL)) {
		return TypeTRL(), nil
	}
	if strings.HasPrefix(name, string(SER)) {
		return TypeSER(), nil
	}
	return Tag{}, ErrTagNotFound
}

//SearchEpisode - searching episode tag in filepath provided.
//Returns error if tags found != 1
func SearchEpisode(name string) (Tag, error) {
	name = filepath.Base(name)
	re := regexp.MustCompile(`(s[0-9]{1,}e[0-9]{1,})`)
	re.FindAllString(name, -1)
	episodeTags := re.FindAllString(name, -1)
	if len(episodeTags) == 0 {
		return Tag{}, ErrTagNotFound
	}
	if !hasSameElements(episodeTags) {
		return Tag{}, fmt.Errorf("different episode tags detected")
	}
	return Episode(episodeTags[0]), nil
}

func hasSameElements(slice []string) bool {
	for i, s1 := range slice {
		for j, s2 := range slice {
			if j <= i {
				continue
			}
			if s1 != s2 {
				return false
			}
		}
	}
	return true
}

//SearchRevision - searching episode tag in filepath provided.
//Returns error if TagMarker is not Input or tags not found.
func SearchRevision(name string) (Tag, error) {
	name = filepath.Base(name)
	tagNum := 0
	switch markerPresent(name) {
	case MarkerBad:
		return Tag{}, ErrBothMarkersPresent
	case MarkerOutput:
		tagNum = findRevisionNumberOutput(name)
	case MarkerInput, NoMarkerPresent:
		tagNum = findRevisionNumberInput(name)
	}
	if tagNum == 0 {
		return Tag{}, ErrTagNotFound
	}
	return Revision(tagNum), nil
}

func findRevisionNumberInput(name string) int {
	re := regexp.MustCompile(`(_R[0-9]{1,}.)`)
	tag := re.FindString(name)
	tag = strings.TrimPrefix(tag, "_")
	tag = strings.TrimSuffix(tag, ".")
	tag = strings.TrimPrefix(tag, "R")
	tagNum, _ := strconv.Atoi(tag)
	return tagNum
}

func findRevisionNumberOutput(name string) int {
	re := regexp.MustCompile(`(_R[0-9]{1,}_)`)
	tag := re.FindString(name)
	tag = strings.TrimPrefix(tag, "_")
	tag = strings.TrimSuffix(tag, "_")
	tag = strings.TrimPrefix(tag, "R")
	tagNum, _ := strconv.Atoi(tag)
	return tagNum
}
