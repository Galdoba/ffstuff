package tag

import (
	"testing"
)

type newCollectionTestCase struct {
	input             string
	colFileName       string
	colCollectionType TagsType
	colTags           map[TagKey]Tag
	expectedError     error
}

func testNewCollectionTestCases() []newCollectionTestCase {
	return []newCollectionTestCase{
		////////////////////////////////////////////////////////////////////////////////
		{"Darbi_i_dzhoan--SER--s01e01--PRT240818173612--Darbi_i_dzhoan_s01e01_PRT240818173612_SER_04896_18.mp4",
			"Darbi_i_dzhoan--SER--s01e01--PRT240818173612--Darbi_i_dzhoan_s01e01_PRT240818173612_SER_04896_18.mp4",
			IN_FILE,
			map[TagKey]Tag{
				BASE:    {BASE, "Darbi_i_dzhoan"},
				SEASON:  {SEASON, "01"},
				EPISODE: {EPISODE, "01"},
				TYPE:    {TYPE, string(SER)},
				PRT:     {PRT, "PRT240818173612"},
				ORIGIN:  {ORIGIN, "Darbi_i_dzhoan_s01e01_PRT240818173612_SER_04896_18.mp4"},
			},
			nil},
		////////////////////////////////////////////////////////////////////////////////
		{"Darbi_i_dzhoan_s01_01_PRT240806093047__HD_proxy.mp4",
			"Darbi_i_dzhoan_s01_01_PRT240806093047__HD_proxy.mp4",
			OUT_FILE,
			map[TagKey]Tag{
				BASE:     {BASE, "Darbi_i_dzhoan"},
				SEASON:   {SEASON, "01"},
				EPISODE:  {EPISODE, "01"},
				PRT:      {PRT, "PRT240806093047"},
				VIDEO:    {TYPE, string(HD)},
				COMMENTS: {COMMENTS, "proxy"},
				EXT:      {EXT, "mp4"},
			},
			nil},
		////////////////////////////////////////////////////////////////////////////////
		{"Darbi_i_dzhoan_s01e01_PRT240818173612_SER_04896_18.mp4",
			"Darbi_i_dzhoan_s01e01_PRT240818173612_SER_04896_18.mp4",
			UNDEFINED,
			map[TagKey]Tag{},
			ErrPromptExpected},
	}
}

func TestNewCollection(t *testing.T) {
	for _, testcase := range testNewCollectionTestCases() {
		// fmt.Println("testcase", i)
		// fmt.Println(testcase.input)
		col, err := NewCollection(testcase.input)
		if err != testcase.expectedError {
			t.Errorf("error received: %v, expected %v", err, testcase.expectedError)
		}
		if col.FileName != testcase.colFileName {
			t.Errorf("filename received: %v, expected %v", col.FileName, testcase.colFileName)
		}
		if col.CollectionType != testcase.colCollectionType {
			t.Errorf("CollectionType received: %v, expected %v", col.CollectionType, testcase.colCollectionType)
		}
		if !equalMapsOfTags(col.TagWithKey, testcase.colTags) {
			t.Errorf("Collection tags received: %v, expected %v", col.TagWithKey, testcase.colTags)
		}

	}
}

func TestAssembledCollection(t *testing.T) {
	for _, testcase := range testNewCollectionTestCases() {
		col, _ := NewCollection(testcase.input)
		if col.CollectionType != UNDEFINED {
			continue
		}
		//tagConstr := constructor.Default()

		col.Add()

	}
}

//func constructTag(c Constructor)

func equalTags(aTag, bTag Tag) bool {
	if aTag.Key != bTag.Key {
		return false
	}
	if aTag.Value != bTag.Value {
		return false
	}
	return true
}

func equalMapsOfTags(aMap, bMap map[TagKey]Tag) bool {
	if len(aMap) != len(bMap) {
		return false
	}
	for aKey, tagInAmap := range aMap {
		tagInBmap := bMap[aKey]
		if !equalTags(tagInAmap, tagInBmap) {
			return false
		}
	}
	return true
}
