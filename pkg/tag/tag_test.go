package tag

import (
	"fmt"
	"strings"
	"testing"
)

type testTagCase struct {
	tag           Tag
	expectedError error
}

func customTag(tg Tag, err error) Tag {
	return tg
}

func customCase(tg Tag, err error) testTagCase {
	return testTagCase{tg, err}
}

var testTagCases = []testTagCase{
	//{Tag{Key: "", Value: "", usage: 0}, fmt.Errorf("invalid key '%v'", "")},
	{Tag{Key: "Name", Value: "some name", usage: 1}, nil},
	{CreateNameTag("aaa"), nil},
	{TypeFILMTag, nil},
	{TypeSERTag, nil},
	{TypeTRLTag, nil},
	{CreateEpisodeTag("s01e02"), nil},
	{CreateEpisodeTag("lala"), nil},
	{CreateSeasonTag("s01"), nil},
	{CreateSeasonTag("la"), nil},
	{CreateRevisionTag(1), nil},
	{CreateRevisionTag(0), nil},
	{CreateRevisionTag(-1), nil},
	{CreateVideoTag("HD"), nil},
	{CreateVideoTag("4K"), nil},
	{CreateVideoTag("SD"), nil},
	{CreateVideoTag(""), nil},
	{CreateVideoTag("lala"), nil},
	{CreateAudioTag("rus", "51"), nil},
	{CreateAudioTag("rus", "20"), nil},
	{CreateAudioTag("eng", "stereo"), nil},
	{CreateCommentTags("comment text 1", "comment 2"), nil},
	customCase(customTag(New("translation", "Goblin", UsageOutput)), nil),
	//customCase(customTag(New("translationasdkjf-3409*a/sdf/er", "SDfkhakdcasf/bd&Goblin", UsageOutput, UsageInput, UsageBase)), fmt.Errorf("bad usagetypes combination provided")),
	//TODO: переписсать тест
}

func TestTag(t *testing.T) {
	for testCaseNum, testdata := range testTagCases {
		if !hasStandardKey(testdata.tag) {
			if !strings.HasPrefix(testdata.tag.Key, "custom_key:") {
				t.Errorf("test case %v: %v\nerror: %v", testCaseNum+1, testdata.tag, fmt.Errorf("invalid key '%v'", testdata.tag.Key))
			}
			if testdata.expectedError != nil {
				if testdata.expectedError.Error() != fmt.Errorf("invalid key '%v'", testdata.tag.Key).Error() {
					t.Errorf("test case %v: %v\nerror: %v", testCaseNum+1, testdata.tag, testdata.expectedError)
				}
			}
		}
		switch hasValidUsage(testdata.tag) {
		case true:
		case false:
			if testdata.expectedError == nil {
				t.Errorf("test case %v: %v\nerror: %v", testCaseNum+1, testdata.tag, fmt.Errorf("bad usagetypes combination provided"))
			}
		}

	}
}

func hasStandardKey(tg Tag) bool {
	switch tg.Key {
	default:
		return false
	case KeyName, KeyContentType, KeyEpisode, KeySeason, KeyRevision, KeyVideo, KeyAudio, KeyComment:
		return true
	}
}

func hasValidUsage(tg Tag) bool {
	switch tg.Key {
	default:
		switch tg.usage {
		case 1, 2, 4, 6:
			return true
		}
	case KeyName:
		switch tg.usage {
		case 1:
			return true
		}
	case KeyContentType:
		switch tg.usage {
		case 2:
			return true
		}
	case KeyAudio:
		switch tg.usage {
		case 4:
			return true
		}
	case KeyEpisode, KeySeason, KeyRevision, KeyVideo, KeyComment:
		switch tg.usage {
		case 6:
			return true
		}
	}
	return false
}