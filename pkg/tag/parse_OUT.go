package tag

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func parse_out_SeasonTag(s string) Tag {
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

func parse_out_EpisodeTag(s string) Tag {
	sTag := parse_out_SeasonTag(s)
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

func parse_out_prtTag(s string) Tag {
	if !strings.Contains(s, OUT_SEP) {
		return NoTag
	}
	re := regexp.MustCompile(`(PRT[0-9]{6,})`) //expect 12 numbers, but 6 is reserved for 'YYMMDD' format
	pTagValue := re.FindString(s)
	if pTagValue == "" {
		return NoTag
	}
	return New(PRT, pTagValue)
}

func parse_out_VideoTag(s string) Tag {
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

func parse_out_AudioTag(s string) Tag {
	re := regexp.MustCompile(`(AUDIO[A-Z]{3,}[0-9]{2,})`)
	found := re.FindString(s)
	if found == "" {
		return NoTag
	}
	return New(AUDIO, found)
}

func parse_out_SrtTag(s string) Tag {
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

func parse_out_RevisionTag(s string) Tag {
	s = strings.ToUpper(s)
	re := regexp.MustCompile(fmt.Sprintf(`(%vR[0-9]{1,})`, OUT_SEP))
	found := re.FindString(s)
	if found == "" {
		return NoTag
	}
	value := strings.TrimPrefix(found, OUT_SEP+"R")
	return New(REVISION, value)
}

func parse_out_Outname(name string) Tag {
	outname := strings.Split(name, OUT_MARKER)[0]
	parts := strings.Split(outname, ".")
	outname = strings.Join(parts[0:len(parts)-1], ".")
	return New(OUTNAME, outname)
}

func parse_out_Extention(name string) Tag {
	parts := strings.Split(name, ".")
	ext := parts[len(parts)-1]
	return New(EXT, ext)
}
