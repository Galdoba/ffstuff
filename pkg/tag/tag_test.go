package tag

import "testing"

type testTagStruct struct {
	inputKey  TagKey
	inputVal  string
	outputKey TagKey
	outputVal string
}

func testTagCases() []testTagStruct {
	return []testTagStruct{
		{"key1", "val1", "key1", "val1"},
	}
}

func TestTags(t *testing.T) {
	for _, testCase := range testTagCases() {
		tag := New(testCase.inputKey, testCase.inputVal)
		if tag.Key != testCase.outputKey {
			t.Errorf("key mis match: have %v, expect %v", tag.Key, testCase.outputKey)
		}
		if tag.Value != testCase.outputVal {
			t.Errorf("value not match: have %v, expect %v", tag.Value, testCase.outputVal)
		}
	}
}
