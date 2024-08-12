package tag

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	MarkerInput     TagMarker = "--"
	MarkerOutput    TagMarker = "__"
	MarkerBad       TagMarker = "[BAD MARKER]"
	NoMarkerPresent TagMarker = "[NO MARKER]"
)

type TagMarker string

func markerPresent(name string) TagMarker {
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

//SearchEpisode - searching episode tag in filepath provided.
//Returns error if tags found != 1
func SearchEpisode(path string) (Tag, error) {
	name := filepath.Base(path)
	re := regexp.MustCompile(`(s[0-9]{1,}e[0-9]{1,})`)
	re.FindAllString(name, -1)
	episodeTags := re.FindAllString(name, -1)
	if len(episodeTags) == 0 {
		return Tag{}, fmt.Errorf("no tag detected")
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
