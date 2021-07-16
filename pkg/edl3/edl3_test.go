package edl3

import (
	"fmt"
	"testing"

	"github.com/macroblock/imed/pkg/types"
)

func sampleLines() []string {
	return []string{
		"",
		"aaa",
		"вап",
		"TITLE: Filmz",
		"FCM: NON-DROP FRAME",
		"001  AX       V     C        00:00:00:00 01:52:06:00 00:00:00:00 01:52:06:00",
		"* FROM CLIP NAME: The_conjuring_the_devil_made_me_do_it_HD.mp4",
		"002  AX       A2    C        00:04:47:00 01:52:06:00 00:04:47:00 01:52:06:00",
		"* FROM CLIP NAME: The_conjuring_the_devil_made_me_do_it_AUDIORUS51.m4a",
		"004  AX       A     W001 024 00:28:40:22 01:52:06:00 00:28:40:22 01:52:06:00",
		"EFFECTS NAME IS Constant Power",
		"* TO CLIP NAME: The_conjuring_the_devil_made_me_do_it_AUDIOENG51.m4a",
		"005  BL       V     C        00:00:00:00 00:03:15:13 01:52:06:00 01:55:21:13",
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
		// exp := &note{"example", "example", nil}
		// if n == exp {
		// 	t.Errorf("construct %v == %v, must not be", n, exp)
		// }
		if n.format != "NOTE" {
			fmt.Println("line: ", line)
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
			switch err.Error() {
			default:
				t.Errorf("error: '%v' is unexpected", err.Error())
				continue
			case "input line have no data":
				continue
			case "statement is not a Standard":
				continue
			}
		}
		if st.id == "" {
			t.Errorf("error: index '%v' is unexpected, expect 001-999", st.id)
		}
		if st.reel == "" {
			t.Errorf("error: reel '%v' is unexpected, expect BL or AX", st.reel)
		}
		if st.channel == "" {
			t.Errorf("error: channel '%v' is unexpected, expect V, A, A2 or NONE", st.channel)
		}
		if st.editType == "" {
			t.Errorf("error: editType '%v' is unexpected, expect C, D or Wxxx", st.editType)
		}
		if st.editDuration != 0 && st.editType != "C" {
			t.Errorf("error: editDuration 0 is unexpected if editType is 'C'", st.editDuration)
		}
		_, errFI := types.ParseTimecode(st.fileIN)
		if errFI != nil {
			t.Errorf("error: timeStamp FI have unexpected error %v", err.Error())
		}
		fmt.Println("")
	}
}
