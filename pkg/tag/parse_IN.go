package tag

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func parse_in_typeTag(s string) Tag {
	str := strings.ToUpper(s)
	switch str {
	case string(UHD):
		return New(TYPE, string(FILM))
	case string(HD):
		return New(TYPE, string(SER))
	case string(SD):
		return New(TYPE, string(TRL))
	default:
		return NoTag
	}
}

func parse_in_seasonTag(s string) Tag {
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

func parse_in_episodeTags(s string) Tag {
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

func parse_in_prtTag(s string) Tag {
	re := regexp.MustCompile(`(PRT[0-9]{6,})`) //expect 12 numbers, but 6 is reserved for 'YYMMDD' format
	pTagValue := re.FindString(s)
	if pTagValue == "" {
		return NoTag
	}
	return New(PRT, pTagValue)
}

func parse_in_videoTag(s string) Tag {
	str := strings.ToUpper(s)
	switch str {
	case string(UHD):
		return New(VIDEO, string(UHD))
	case string(HD):
		return New(VIDEO, string(HD))
	case string(SD):
		return New(VIDEO, string(SD))
	default:
		return NoTag
	}
}

func parse_in_srtTag(s string) Tag {
	str := strings.ToUpper(s)
	switch str {
	case string(SUB):
		return New(SRT, string(SUB))
	case string(HARDSUB):
		return New(SRT, string(HARDSUB))
	default:
		return NoTag
	}
}

func parse_in_revisionTag(s string) Tag {
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
