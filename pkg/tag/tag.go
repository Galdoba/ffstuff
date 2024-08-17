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

//UsageType defines how Tag may be used.
type TagKey string

//TypeValue defines what type of media we are working with.
type TypeValue string

//TypeValue defines what size of video stream we are working with.
type VideoValue string

//TypeValue defines what size of video stream we are working with.
type SrtValue string

//Tag - datapoint to be used in various scripts to mark file's processing state, demuxin data and purpose.
type Tag struct {
	Key   TagKey
	Value string
}

var NoTag = Tag{NoTagKey, "NoTag"}

//New - creates custom tag for prefix/postfix.
//For standard tags use other functions.
func New(k TagKey, v string) Tag {
	return Tag{Key: k, Value: v}
}

// //Name - creates standard tag. Should be used instead of tag.New() if possibe.
// //This tag defines base name for demuxing process.
// func CreateNameTag(v string) Tag {
// 	return Tag{Key: KeyName, Value: v}
// }

// //TypeFILM - creates standard tag. Should be used instead of tag.New() if possibe.
// //This tag defines content as Film.
// var TypeFILMTag = Tag{Key: KeyContentType, Value: string(FILM)}

// //TypeSER - creates standard tag. Should be used instead of tag.New() if possibe.
// //This tag defines content as Serial.
// var TypeSERTag = Tag{Key: KeyContentType, Value: string(SER)}

// //TypeTRL - creates standard tag. Should be used instead of tag.New() if possibe.
// //This tag defines content as Trailer.
// var TypeTRLTag = Tag{Key: KeyContentType, Value: string(TRL)}

// //Episode - creates standard tag. Should be used instead of tag.New() if possibe.
// //This tag defines season and episode numbers. MUST be used for Serials only.
// func CreateEpisodeTag(v string) Tag {
// 	return Tag{Key: KeyEpisode, Value: v}
// }

// //Season - creates standard tag. Should be used instead of tag.New() if possibe.
// //This tag defines season number. MUST be used for Trailers only.
// func CreateSeasonTag(v string) Tag {
// 	return Tag{Key: KeySeason, Value: v}
// }

// //Revision - creates standard tag. Should be used instead of tag.New() if possibe.
// //This tag defines number of revision for input media if is was corrected by provider.
// func CreateRevisionTag(r int) Tag {
// 	return Tag{Key: KeyRevision, Value: fmt.Sprintf("R%d", r)}
// }

// //Video - creates standard tag. Should be used instead of tag.New() if possibe.
// //This tag defines size of video steam.
// //On input MUST be used for media containg size other than 1920x1080.
// //On output MUST be used for all media containg video stream.
// func CreateVideoTag(v string) Tag {
// 	return Tag{Key: KeyVideo, Value: v}
// }

// //Audio - creates standard tag. Should be used instead of tag.New() if possibe.
// //This tag defines language and channel_layout for audio stream in media.
// //MUST be used for all media consisting from only one audio stream.
// func CreateAudioTag(lang, layout string) Tag {
// 	return Tag{Key: KeyAudio, Value: formatAudioTagValue(lang, layout)}
// }

// func formatAudioTagValue(lang, layout string) string {
// 	return fmt.Sprintf("AUDIO%v%v", lang, layout)
// }

// //Comment - creates standard tag. Should be used instead of tag.New() if possibe.
// //This tag defines program/user comment for output
// func CreateCommentTags(comments ...string) Tag {
// 	return Tag{Key: KeyComment, Value: formatCommentTagValue(comments...)}
// }

// func formatCommentTagValue(comments ...string) string {
// 	cmnt := ""
// 	for _, comment := range comments {
// 		cmnt += fmt.Sprintf("%v_", comment)
// 	}
// 	return cmnt
// }
