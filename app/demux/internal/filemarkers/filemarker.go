package filemarkers

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	NO_MARKERS  = "no markers found"
	MarkerError = "ERROR"
	BASE        = "OUTBASE"
	CONTENTTYPE = "TYPE"
	EPISODE     = "EPISODE"
	FILM        = "FILM"
	TRL         = "TRL"
	SER         = "SER"
)

type Marker struct {
	OfType map[string]string
}

func New(name string) Marker {
	m := Marker{}
	m.OfType = make(map[string]string)
	parts := strings.Split(name, "--")
	if len(parts) == 1 {
		m.OfType[MarkerError] = fmt.Sprintf("%v in %v", NO_MARKERS, name)
	}

	for i, marker := range parts {
		if i == 0 {
			m.OfType[BASE] = marker

		}
		if strings.Contains(marker, FILM) {
			m.OfType[CONTENTTYPE] = FILM
		}
		if strings.Contains(marker, SER) {
			m.OfType[CONTENTTYPE] = SER
			re := regexp.MustCompile(`(s[0-9]{1,}_[0-9]{1,})`)
			tag := re.FindString(name)
			if tag != "" {
				m.OfType[EPISODE] = tag
			}
		}
		if strings.Contains(marker, TRL) {
			m.OfType[CONTENTTYPE] = TRL
		}

	}
	return m
}
