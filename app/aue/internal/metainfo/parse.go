package metainfo

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	key "github.com/Galdoba/ffstuff/app/aue/internal/define"
)

func Parse(name string) []Meta {
	found := []Meta{}
	found, err := parseComplex(name)
	if err == nil {
		return found
	}
	//TODO: log error
	found, err = parseSimple(name)
	if err == nil {
		return found
	}
	//TODO: log error
	return parseDesperate(name)

}

func parseComplex(name string) ([]Meta, error) {

	re := regexp.MustCompile(`(*_[0-9]{3,}_PRT[0-9]{12,})`)
	complexFeed := re.FindString(name)
	if complexFeed == "" {
		return nil, ErrNotFound
	}
	found, err := separateComplex(complexFeed)
	if err != nil {
		return found, fmt.Errorf("parseComplex failed: %v", err)
	}

	return []Meta{}, nil
}

var ErrNotFound = errors.New("")

func separateComplex(feed string) ([]Meta, error) {
	found := []Meta{}
	data := strings.Split(feed, "_")
	if len(data) < 3 {
		return nil, fmt.Errorf("can't separate complex: bad feed '%v'", feed)
	}
	prtFeed := data[len(data)-1]
	prt := NewMeta(key.META_PRT, prtFeed)
	found = append(found, prt)

	seasEpisFeed := data[len(data)-2]
	v, err := strconv.Atoi(seasEpisFeed)
	if err != nil {
		return nil, fmt.Errorf("can't separate complex: bad season/episode feed '%v'", feed)
	}
	if v < 101 {
		return nil, fmt.Errorf("can't separate complex: bad season/episode feed '%v', expect 101+", seasEpisFeed)
	}
	season, episode := seasonAndEpisode(seasEpisFeed)
	found = append(found, season, episode)

	baseFeed := strings.Join(data[0:len(data)-3], "_")
	base := NewMeta(key.META_Base, baseFeed)
	found = append(found, base)
	return found, nil
}

func seasonAndEpisode(feed string) (Meta, Meta) {
	parts := strings.Split(feed, "")
	seasonNum := strings.Join(parts[0:len(parts)-3], "")
	seasonNum = twoDigitNum(seasonNum)
	episodeNum := strings.Join(parts[len(parts)-2:], "")
	episodeNum = twoDigitNum(episodeNum)
	return NewMeta(key.META_Season, seasonNum), NewMeta(key.META_Episode, episodeNum)
}

func twoDigitNum(n string) string {
	v, _ := strconv.Atoi(n)
	if v < 10 {
		n = "0" + n
	}
	return n
}

func parseSimple(name string) ([]Meta, error) {

	re := regexp.MustCompile(`(*_PRT[0-9]{12,})`)
	complexFeed := re.FindString(name)
	if complexFeed == "" {
		return nil, ErrNotFound
	}
	found, err := separateComplex(complexFeed)
	if err != nil {
		return found, fmt.Errorf("parseComplex failed: %v", err)
	}

	return found, nil
}

func separateSimple(feed string) ([]Meta, error) {
	found := []Meta{}
	data := strings.Split(feed, "_")
	if len(data) < 2 {
		return nil, fmt.Errorf("can't separate simple: bad feed '%v'", feed)
	}
	prtFeed := data[len(data)-1]
	prt := NewMeta(key.META_PRT, prtFeed)
	found = append(found, prt)

	baseFeed := strings.Join(data[0:len(data)-2], "_")
	base := NewMeta(key.META_Base, baseFeed)
	found = append(found, base)
	return found, nil
}

func parseDesperate(name string) []Meta {
	data := strings.Split(name, ".")
	baseFeed := strings.Join(data[0:len(data)-1], "")
	return append([]Meta{}, NewMeta(key.META_Base, baseFeed))
}
