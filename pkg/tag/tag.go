package tag

const (
	BASE     TagKey = "Base"
	MARKER   TagKey = "Marker"
	TYPE     TagKey = "Type"
	SEASON   TagKey = "Season"
	EPISODE  TagKey = "Episode"
	PRT      TagKey = "Prt"
	VIDEO    TagKey = "Video"
	AUDIO    TagKey = "Audio"
	SRT      TagKey = "Srt"
	REVISION TagKey = "Revision"
	COMMENTS TagKey = "Comments"
	EXT      TagKey = "Extention"
	ORIGIN   TagKey = "Origin"
	OUTNAME  TagKey = "OUTNAME"
	NoTagKey TagKey = "No Tag"

	FILM TypeValue = "FILM"
	SER  TypeValue = "SER"
	TRL  TypeValue = "TRL"

	UHD VideoValue = "4K"
	HD  VideoValue = "HD"
	SD  VideoValue = "SD"

	SUB     SrtValue = "SUB"
	HARDSUB SrtValue = "HSUB"
)

// UsageType defines how Tag may be used.
type TagKey string

// TypeValue defines what type of media we are working with.
type TypeValue string

// TypeValue defines what size of video stream we are working with.
type VideoValue string

// TypeValue defines what size of video stream we are working with.
type SrtValue string

// Tag - datapoint to be used in various scripts to mark file's processing state, demuxin data and purpose.
type Tag struct {
	Key   TagKey
	Value string
}

var NoTag = Tag{NoTagKey, "NoTag"}

// New - creates custom tag for prefix/postfix.
// For standard tags use other functions.
func New(k TagKey, v string) Tag {
	return Tag{Key: k, Value: v}
}

func (t Tag) String() string {
	if t.Value == "NoTag" {
		return ""
	}
	return t.Value
}
