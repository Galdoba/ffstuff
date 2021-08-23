package edl3

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/macroblock/imed/pkg/types"
)

func sampleLines() []string {
	return []string{
		//INVALID LINES
		"",
		"aaa <<<",
		"aaa",
		"вап",
		"вап dfg вап",
		"TITLE: Filmz",
		"TITLEE: Filmz",
		"FCM: NON-DROP FRAME",
		"FCM: DROP FRAME",
		"FCM: NOsdOP FRAME",
		"FCMM: NON-DROP FRAME",
		"001  AX       V     C        00:00:00:00 01:52:06:00 00:00:00:00 01:52:06:00",
		"01  AX       V     C        00:00:00:00 01:52:06:00 00:00:00:00 01:52:06:00",
		"*  AX       V     C        00:00:00:00 01:52:06:00 00:00:00:00 01:52:06:00",
		"001111  AX       V     C        00:00:00:00 01:52:06:00 00:00:00:00 01:52:06:00",
		"-001  AX       V     C        00:00:00:00 01:52:06:00 00:00:00:00 01:52:06:00",
		"OOI  AX       V     C        00:00:00:00 01:52:06:00 00:00:00:00 01:52:06:00",
		"  AX       V     C        00:00:00:00 01:52:06:00 00:00:00:00 01:52:06:00",
		"001  BY       V     C        00:00:00:00 01:52:06:00 00:00:00:00 01:52:06:00",
		"001              C        00:00:00:00 01:52:06:00 00:00:00:00 01:52:06:00",
		"001  AX       V     C   17     00:00:00:00 01:52:06:00 00:00:00:00 01:52:06:00",
		"001  AX       V     R   17     00:00:00:00 01:52:06:00 00:00:00:00 01:52:06:00",
		"001  AX       A3     R   17     00:00:00:00 01:52:06:00 00:00:00:00 01:52:06:00",
		"001  AX       A3     C        00:00:00:00 01:52:06:00 00:00:00:00 01:52:06:00",
		"001  AX       V     C        00:0L00:00 01:52:06:00 00:00:00:00 01:52:06:00",
		"001  AX       V     C        00:00:00:00 01:52:06:00 00:00:00:00 01:52:0B:00",
		"001  AX       V     C        00:00:00:00 01:527:06:00 00:00:00:00 01:52:06:00",
		"001  AX       V     C        00:00:00:00 01:52:06:00 00:00:00:010 01:52:06:00",
		"001  AX       V     C        00:00:00:00 01:52:06:00 00:00:00:010 01:52:06:00\n",
		"* FROM CLIP NAME: The_conjuring_the_devil_made_me_do_it_HD.mp4",
		"* FM CLIP NAME: The_conjuring_the_devil_made_me_do_it_HD.mp4",
		"* FM CLIP NAME: The_conjuring_the_devil_made_me_do_it_HD.mp4",
		"002  AX       A2    C        00:04:47:00 01:52:06:00 00:04:47:00 01:52:06:00",
		"* FROM CLIP NAME: The_conjuring_the_devil_made_me_do_it_AUDIORUS51.m4a",
		"004  AX       A     W001 024 00:28:40:22 01:52:06:00 00:28:40:22 01:52:06:00",
		"004  AX       A     W211 024 00:28:40:22 01:52:06:00 00:28:40:22 01:52:06:00",
		"004  AX       A     W101 024.2 00:28:40:22 01:52:06:00 00:28:40:22 01:52:06:00",
		"004  AX       A     W101 024.B 00:28:40:22 01:52:06:00 00:28:40:22 01:52:06:00",
		"EFFECTS NAME IS Constant Power",
		"* TO CLIP NAME: The_conjuring_the_devil_made_me_do_it_AUDIOENG51.m4a",
		"005  BL       V     C        00:00:00:00 00:03:15:13 01:52:06:00 01:55:21:13",
		"AUD       4",
		"AUD  3    ",
		"AUD  3    4",
		"AUD       ",
		"AUD       5",
		"AUD  B    6",
		"AUD  B     ",
		//VALID EDL
		"TITLE: sample_test",
		"FCM: NON-DROP FRAME",
		"",
		"001  BL       V     C        00:00:00:00 00:01:52:07 00:00:00:00 00:01:52:07",
		"",
		"002  BL       A     C        00:00:00:00 00:00:00:04 00:00:00:00 00:00:00:00",
		"",
		"003  BL       A     C        00:00:00:00 00:00:00:00 00:01:52:07 00:01:52:07",
		"FCM: NON-DROP FRAME",
		"003  AX       A     W001 004 00:01:52:07 00:03:21:09 00:01:52:07 00:03:21:09",
		"EFFECTS NAME IS Constant Power",
		"* FROM CLIP NAME: BL",
		"* TO CLIP NAME: vyskochka_s02_01_2020__hd_rus51.m4a",
		"",
		"004  BL       V     C        00:00:00:00 00:00:01:05 00:00:00:00 00:00:00:00",
		"",
		"005  BL       V     C        00:00:00:00 00:00:00:00 00:01:52:07 00:01:52:07",
		"FCM: NON-DROP FRAME",
		"005  AX       V     D    030 00:01:52:07 00:03:20:08 00:01:52:07 00:03:20:08",
		"EFFECTS NAME IS CROSS DISSOLVE",
		"* FROM CLIP NAME: BL",
		"* TO CLIP NAME: vyskochka_s02_01_2020__hd.mp4",
		"",
		"006  AX       V     C        00:03:20:08 00:03:20:08 00:03:20:08 00:03:20:08",
		"FCM: NON-DROP FRAME",
		"006  AX       V     D    030 00:59:58:20 01:00:02:20 00:03:20:08 00:03:24:08",
		"EFFECTS NAME IS CROSS DISSOLVE",
		"* FROM CLIP NAME: vyskochka_s02_01_2020__hd.mp4",
		"* TO CLIP NAME: Title 97",
		"",
		"007  AX       A     C        00:03:21:09 00:03:21:09 00:03:21:09 00:03:21:09",
		"FCM: NON-DROP FRAME",
		"007  BL       A     W001 004 00:00:00:00 03:56:38:12 00:03:21:09 04:00:00:00",
		"EFFECTS NAME IS Constant Power",
		"* FROM CLIP NAME: vyskochka_s02_01_2020__hd_rus51.m4a",
		"* TO CLIP NAME: BL",
		"",
		"008  AX       V     C        00:03:21:13 00:03:21:13 00:03:21:13 00:03:21:13",
		"FCM: NON-DROP FRAME",
		"008  AX       V     D    030 01:00:00:00 01:00:02:20 00:03:21:13 00:03:24:08",
		"EFFECTS NAME IS CROSS DISSOLVE",
		"* FROM CLIP NAME: vyskochka_s02_01_2020__hd.mp4",
		"* TO CLIP NAME: Title 97",
		"",
		"009  AX       V     C        01:00:02:20 01:00:02:20 00:03:24:08 00:03:24:08",
		"FCM: NON-DROP FRAME",
		"009  BL       V     D    030 00:00:00:00 03:56:34:12 00:03:24:08 04:00:00:00",
		"EFFECTS NAME IS CROSS DISSOLVE",
		"* FROM CLIP NAME: Title 97",
		"* TO CLIP NAME: BL",
		"",
	}
}

func TestNewNote(t *testing.T) {
	for _, line := range sampleLines() {
		n, err := newNote(line)
		if err != nil {
			switch {
			default:
				t.Errorf("error: '%v' is unexpected", err.Error())
				continue
			case err.Error() == "input line have no data":
				continue
			}
		}
		if n.format != "NOTE" {
			t.Errorf("format == %v, but expect 'NOTE'", n.format)
		}
		if n.content == "" {
			t.Errorf("n.content = '%v', from line = '%v' must not be blank", n.content, line)
		}
		if n.definedDataFields == nil {
			t.Errorf("n.definedDataFields (%v) is not initiated", n.definedDataFields)
		}
		if len(n.definedDataFields) != 0 {
			t.Errorf("n.definedDataFields (%v) contain data. It should not", n.definedDataFields)
		}
	}
}

func TestNewTitle(t *testing.T) {
	for _, line := range sampleLines() {
		n, errN := newNote(line)
		if errN != nil {
			continue
		}
		ttl, err := newTitle(n)
		if err != nil {
			switch err.Error() {
			default:
				t.Errorf("error: '%v' is unexpected", err.Error())
				continue
			case "input line have no data":
				continue
			case "statement is not a title":
				continue
			}
		}
		if len(ttl.title) == 0 {
			t.Errorf("error: title len is '%v', expected > 0", len(ttl.title))
		}
	}
}

func TestNewFCM(t *testing.T) {
	for _, line := range sampleLines() {
		n, errN := newNote(line)
		if errN != nil {
			continue
		}
		fcm, err := newFCM(n)
		if err != nil {
			switch err.Error() {
			default:
				t.Errorf("error: '%v' is unexpected", err.Error())
				continue
			case "input line have no data":
				continue
			case "statement is not a FCM":
				continue
			case "unknown FCM mode":
				continue
			}
		}
		if fcm.mode != "NON-DROP FRAME" && fcm.mode != "DROP FRAME" {
			t.Errorf("error: fcm.mode = %v, 'NON-DROP FRAME/DROP FRAME' is expected", fcm.mode)
		}
	}
}

func TestNewStandard(t *testing.T) {
	for _, line := range sampleLines() {
		n, errN := newNote(line)
		if errN != nil {
			continue
		}
		st, err := newStandard(n)
		if err != nil {
			switch {
			default:
				fmt.Println("feed:", line)
				t.Errorf("UNEXPECTED error: %v", err.Error())
				continue
			case err.Error() == "statement is not a Standard":
				continue
			case strings.Contains(err.Error(), "invalid statement syntax: "):
				continue
			}
		}
		index, errIndex := strconv.Atoi(st.id)
		if errIndex != nil {
			t.Errorf("error: index '%v' not a number", st.id)
		}
		if index < 1 || index > 999 {
			t.Errorf("error: index '%v' is unexpected, expecting '001-999'", st.id)
		}
		if st.reel != "BL" && st.reel != "AX" {
			t.Errorf("error: reel '%v' is unexpected, expecting BL or AX", st.reel)
		}
		if !listContains([]string{"V", "A", "A2", "NONE"}, st.channel) {
			t.Errorf("error: channel '%v' is unexpected, expecting V, A, A2 or NONE", st.channel)
		}
		if !listContains(validWipeCodes(), st.editType) {
			t.Errorf("error: editType '%v' is unexpected, expecting C, D or Wxxx", st.editType)
		}
		if st.editDuration != 0 && st.editType == "C" {
			t.Errorf("error: editDuration %v is unexpected if editType is 'C'", st.editDuration)
		}
		if !isTimecode(st.fileIN) {
			t.Errorf("error: fileIN is not types.Timecode %v", st.fileIN)
		}
		if !isTimecode(st.fileOUT) {
			t.Errorf("error: fileOUT is not types.Timecode %v", st.fileOUT)
		}
		if !isTimecode(st.seqIN) {
			t.Errorf("error: seqIN is not types.Timecode %v", st.seqIN)
		}
		if !isTimecode(st.seqOUT) {
			t.Errorf("error: seqOUT is not types.Timecode %v", st.seqOUT)
		}
	}
}

func isTimecode(ts types.Timecode) bool {
	test, err := types.ParseTimecode("00:00:00:01")
	switch {
	case reflect.TypeOf(ts) == reflect.TypeOf(test):
		return true
	case err != nil:
		return false
	default:
		return false
	}
}

func TestNewSource(t *testing.T) {
	for _, line := range sampleLines() {
		n, errN := newNote(line)
		if errN != nil {
			continue
		}
		src, err := newSource(n)
		if err != nil {
			switch err.Error() {
			default:
				t.Errorf("error: '%v' is unexpected %v, %v", err.Error(), src, line)
				continue
			case "statement is not a Source":
				continue
			}
		}
		if src.sourceA == "" && src.sourceB == "" {
			t.Errorf("sourceA and sourceB can not be unfilled at the same time: line = '%v'", line)
		}
		if src.sourceA != "" && src.sourceB != "" {
			t.Errorf("sourceA and sourceB can not be filled at the same time: line = '%v'", line)
		}
	}
}

func TestNewAud(t *testing.T) {
	for _, line := range sampleLines() {
		n, errN := newNote(line)
		if errN != nil {
			continue
		}
		a, err := newAud(n)
		if err != nil {
			switch {
			default:
				t.Errorf("error: '%v' is unexpected %v, %v", err.Error(), a, line)
				continue
			case err.Error() == "statement is not a Aud":
				continue
			case strings.Contains(err.Error(), "invalid statement "):
				continue
			}
		}
		if a.channel != 3 && a.channel != 4 {
			t.Errorf("error: a.channel = '%v', expecting 3 or 4 from line '%v'", a.channel, line)
		}

	}
}

func TestEffectStatement(t *testing.T) {
	for _, line := range sampleLines() {
		n, errN := newNote(line)
		if errN != nil {
			continue
		}
		es, err := newEffectStatement(n)
		if err != nil {
			switch {
			default:
				t.Errorf("error: is %v, expect nil ('%v' '%v')", err.Error(), es, line)
				continue
			case err.Error() == "statement is not a EffectStatement":
				continue
			case strings.Contains(err.Error(), "invalid statement "):
				continue
			}
		}
	}
}

func TestStatement(t *testing.T) {
	for _, line := range sampleLines() {
		stType, stData, err := Statement(line)
		if err != nil {
			t.Errorf("error : %v", err.Error())
		}
		if stType != "NOTE" {
			switch len(stData) {
			case 0:
				t.Errorf("Statement defined as %v, but no data detected", stType)
			default:
				if stData[0] == "" {
					t.Errorf("Statement defined as %v, but no data contained", stType)
				}
			}

		}
	}
}
