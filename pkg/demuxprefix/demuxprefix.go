package demuxprefix

import (
	"fmt"
	"path/filepath"
	"strings"
)

var pathSep = string(filepath.Separator)
var prefSep = "--"

const (
	keyName      = "Name"
	keyType      = "Type"
	keyEpisode   = "Episode"
	keyRevision  = "Revision"
	tagValueFILM = "FILM"
	tagValueTRL  = "TRL"
	tagValueSER  = "SER"
)

type Prefix struct {
	tags map[string]string
}

func New() *Prefix {
	pr := Prefix{}
	pr.tags = make(map[string]string)
	return &pr
}

func (pr *Prefix) AddTags(tags ...Tag) error {
	tagmap := pr.tags
	for _, tag := range tags {
		if _, ok := tagmap[tag.key]; ok {
			return fmt.Errorf("tag '%v' alredy provided", tag.key)
		}
		tagmap[tag.key] = tag.val
	}
	if tagmap[keyName] == "" {
		return fmt.Errorf("tag Name is not provided")
	}
	if tagmap[keyType] == "" {
		return fmt.Errorf("tag Type is not provided")
	}
	if tagmap[keyType] == tagValueSER && tagmap[keyEpisode] == "" {
		return fmt.Errorf("tag Type is provided, but tag Episode is not")
	}
	pr.tags = tagmap
	return nil
}

func (pr *Prefix) String() string {
	str := ""
	switch pr.tags[keyType] {
	case tagValueFILM, tagValueTRL:
		str = fmt.Sprintf("%v%v%v", pr.tags[keyName], pr.tags[keyType], pr.tags[keyRevision])
	case tagValueSER:
		str = fmt.Sprintf("%v%v%v%v", pr.tags[keyName], pr.tags[keyType], pr.tags[keyEpisode], pr.tags[keyRevision])
	}
	str = strings.ReplaceAll(str, prefSep+prefSep, prefSep)
	return str
}

type Tag struct {
	key string
	val string
}

func (tg Tag) String() string {
	return tg.val + prefSep
}

func TagName(name string) Tag {
	return Tag{
		key: keyName,
		val: name,
	}
}

func TagFILM() Tag {
	return Tag{
		key: keyType,
		val: "FILM",
	}
}

func TagTRL() Tag {
	return Tag{
		key: keyType,
		val: "TRL",
	}
}

func TagSER() Tag {
	return Tag{
		key: keyType,
		val: "SER",
	}
}

func TagEpisode(s, e string) Tag {
	return Tag{
		key: keyEpisode,
		val: fmt.Sprintf("s%ve%v", s, e),
	}
}

func TagRevision(r int) Tag {
	return Tag{
		key: keyEpisode,
		val: fmt.Sprintf("R%v", r),
	}
}

func DemuxPrefix(path string) string {
	//dir := filepath.Dir(path)
	fullname := filepath.Base(path)
	tags := strings.Split(fullname, prefSep)
	tags = tags[0 : len(tags)-1]
	return strings.Join(tags, prefSep) + prefSep
}
