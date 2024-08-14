package tag

import "fmt"

const (
	KeyName                  = "Name"
	KeyContentType           = "Content_Type"
	KeyEpisode               = "Episode"
	KeySeason                = "Season"
	KeyRevision              = "Revision"
	KeyVideo                 = "Video"
	KeyAudio                 = "Audio"
	KeyComment               = "Comment"
	FILM           TypeValue = "FILM"
	SER            TypeValue = "SER"
	TRL            TypeValue = "TRL"
	UsageBase      UsageType = "Base"
	UsageInput     UsageType = "Input"
	UsageOutput    UsageType = "Output"
	usageNome                = 0
	usageBase                = 1
	usageIn                  = 2
	usageOut                 = 4
	usageError               = 8
)

//UsageType defines how Tag may be used.
type UsageType string

//TypeValue defines what type of media we are working with.
type TypeValue string

//Tag - datapoint to be used in various scripts to mark file's processing state, demuxin data and purpose.
type Tag struct {
	Key   string
	Value string
	usage int
}

var NoTag = Tag{"No tag", "NoTag", 8}

//New - creates custom tag for prefix/postfix.
//For standard tags use other functions.
func New(k, v string, usages ...UsageType) (Tag, error) {
	tg := Tag{
		Key:   formatTagKey(k),
		Value: v,
		usage: calculatedUsage(usages...),
	}
	//TODO: key must be restricted to filepath allowed runes
	switch tg.usage {
	case 1, 2, 4, 6:
	default:
		return tg, fmt.Errorf("bad usagetypes combination provided")
	}
	return tg, nil
}

func formatTagKey(key string) string {
	if isStandardKey(key) {
		return key
	}
	return fmt.Sprintf("custom_key:%v", key)
}

func isStandardKey(key string) bool {
	switch key {
	default:
		return false
	case KeyName, KeyContentType, KeyEpisode, KeySeason, KeyRevision, KeyVideo, KeyAudio, KeyComment:
		return true
	}
}

func calculatedUsage(uTypes ...UsageType) int {
	sum := 0
	for _, utype := range uTypes {
		switch utype {
		case UsageBase:
			sum += usageBase
		case UsageInput:
			sum += usageIn
		case UsageOutput:
			sum += usageOut
		default:
			sum += usageError
		}
	}
	return sum
}

//Name - creates standard tag. Should be used instead of tag.New() if possibe.
//This tag defines base name for demuxing process.
func CreateNameTag(v string) Tag {
	return Tag{Key: KeyName, Value: v, usage: usageBase}
}

//TypeFILM - creates standard tag. Should be used instead of tag.New() if possibe.
//This tag defines content as Film.
var TypeFILMTag = Tag{Key: KeyContentType, Value: string(FILM), usage: usageIn}

//TypeSER - creates standard tag. Should be used instead of tag.New() if possibe.
//This tag defines content as Serial.
var TypeSERTag = Tag{Key: KeyContentType, Value: string(SER), usage: usageIn}

//TypeTRL - creates standard tag. Should be used instead of tag.New() if possibe.
//This tag defines content as Trailer.
var TypeTRLTag = Tag{Key: KeyContentType, Value: string(TRL), usage: usageIn}

//Episode - creates standard tag. Should be used instead of tag.New() if possibe.
//This tag defines season and episode numbers. MUST be used for Serials only.
func CreateEpisodeTag(v string) Tag {
	return Tag{Key: KeyEpisode, Value: v, usage: usageIn + usageOut}
}

//Season - creates standard tag. Should be used instead of tag.New() if possibe.
//This tag defines season number. MUST be used for Trailers only.
func CreateSeasonTag(v string) Tag {
	return Tag{Key: KeySeason, Value: v, usage: usageIn + usageOut}
}

//Revision - creates standard tag. Should be used instead of tag.New() if possibe.
//This tag defines number of revision for input media if is was corrected by provider.
func CreateRevisionTag(r int) Tag {
	return Tag{Key: KeyRevision, Value: fmt.Sprintf("R%d", r), usage: usageIn + usageOut}
}

//Video - creates standard tag. Should be used instead of tag.New() if possibe.
//This tag defines size of video steam.
//On input MUST be used for media containg size other than 1920x1080.
//On output MUST be used for all media containg video stream.
func CreateVideoTag(v string) Tag {
	return Tag{Key: KeyVideo, Value: v, usage: usageIn + usageOut}
}

//Audio - creates standard tag. Should be used instead of tag.New() if possibe.
//This tag defines language and channel_layout for audio stream in media.
//MUST be used for all media consisting from only one audio stream.
func CreateAudioTag(lang, layout string) Tag {
	return Tag{Key: KeyAudio, Value: formatAudioTagValue(lang, layout), usage: usageOut}
}

func formatAudioTagValue(lang, layout string) string {
	return fmt.Sprintf("AUDIO%v%v", lang, layout)
}

//Comment - creates standard tag. Should be used instead of tag.New() if possibe.
//This tag defines program/user comment for output
func CreateCommentTags(comments ...string) Tag {
	return Tag{Key: KeyComment, Value: formatCommentTagValue(comments...), usage: usageIn + usageOut}
}

func formatCommentTagValue(comments ...string) string {
	cmnt := ""
	for _, comment := range comments {
		cmnt += fmt.Sprintf("%v_", comment)
	}
	return cmnt
}
