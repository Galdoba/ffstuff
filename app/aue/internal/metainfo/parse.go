package metainfo

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	key "github.com/Galdoba/ffstuff/app/aue/internal/define"
	logger "github.com/Galdoba/ffstuff/pkg/logman"
)

func Parse(name string) []Meta {
	found := []Meta{}
	found, err := parseComplex(name)
	if err == nil {
		return found
	}
	logger.Error(fmt.Errorf("parseComplex Failed:", name))

	found, err = parseSimple(name)
	if err == nil {
		return found
	}
	logger.Error(fmt.Errorf("parseSimple Failed:", name))

	return parseDesperate(name)

}

func parseComplex(name string) ([]Meta, error) {
	re := regexp.MustCompile(`(.{1,}_s[0-9]{2,}e[0-9]{2,}_PRT[0-9]{12,})`)
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

var ErrNotFound = errors.New("")

func separateComplex(feed string) ([]Meta, error) {
	found := []Meta{}
	data := strings.Split(feed, "_")
	if len(data) <= 3 {
		return nil, fmt.Errorf("can't separate complex: bad feed '%v'", feed)
	}
	prtFeed := data[len(data)-1]
	prt := NewMeta(key.META_PRT, prtFeed)
	found = append(found, prt)

	seasEpisFeed := data[len(data)-2]
	season, episode := seasonAndEpisode(seasEpisFeed)
	found = append(found, episode, season)

	baseFeed := strings.Join(data[0:len(data)-2], "_")
	base := NewMeta(key.META_Base, baseFeed)
	found = append(found, base)
	return found, nil
}

func seasonAndEpisode(feed string) (Meta, Meta) {
	feed = strings.TrimPrefix(feed, "s")
	numbers := strings.Split(feed, "e")

	return NewMeta(key.META_Season, numbers[0]), NewMeta(key.META_Episode, numbers[1])
}

func twoDigitNum(n string) string {
	v, _ := strconv.Atoi(n)
	if v < 10 {
		n = "0" + n
	}
	return n
}

func parseSimple(name string) ([]Meta, error) {

	re := regexp.MustCompile(`(.{1,}_PRT[0-9]{12,})`)
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
