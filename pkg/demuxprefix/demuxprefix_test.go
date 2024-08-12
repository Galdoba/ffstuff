package demuxprefix

import (
	"testing"
)

var testData = []testCase{
	{`\\192.168.31.4\buffer\IN\_DONE\Dom_drakona--s02e02--SER--Dom_drakona_s02e02_PRT240810214506_SER_04736_18.mp4`,
		"Dom_drakona--s02e02--SER--"},
	{`\\192.168.31.4\buffer\IN\_DONE\Dom_drakona_R1--s02e05--SER--Dom_drakona_s02e05_SER_04840_18_R1.mp4`, "Dom_drakona_R1--s02e05--SER--"},
	{`\\192.168.31.4\buffer\IN\_DONE\Velikolepnaya_chetverka--FILM--FabulousFour_422hq_HD_rec709_51_ru_CutVersion_08082024.mov`, "Velikolepnaya_chetverka--FILM--"},
	{`\\192.168.31.4\root\IN\@TRAILERS\_DONE\Voyna_foylaya_s01_TRL\Voyna_foylaya_s01--TRL--voyna_foyla_a_teka.mp4`, "Voyna_foylaya_s01--TRL--"},
}

type testCase struct {
	input          string
	outputExpected string
}

func TestDemuxprefix(t *testing.T) {
	for _, testCase := range testData {
		prefix := DemuxPrefix(testCase.input)
		if prefix != testCase.outputExpected {
			t.Errorf("input %v: output '%v', expected '%v'", testCase.input, prefix, testCase.outputExpected)
		}
	}
}
