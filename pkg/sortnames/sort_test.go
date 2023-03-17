package sortnames

import (
	"fmt"
	"strings"
	"testing"
)

func TestNoDuplicates(t *testing.T) {
	test1 := []string{"aa", "bb", "cc", "bb", "dd"}
	res1 := NoDuplicates(test1)
	if strings.Join(res1, "") != "aabbccdd" {
		t.Errorf("test 1 fail")
	}
	printSliceLn(emulatedNames())
}

//name = base + sTag + epTag + revTag + vtag + atag + censtag  + extn

func printSliceLn(sl []string) {
	for _, s := range sl {
		fmt.Println(s)
	}
}

func emulatedNames() []string {
	var emulatedNames []string
	base := []string{"serial_AAA", "serial_BBB", "film_AA", "film_BB"}
	sTag := []string{"_s01", "_s02", ""}
	epTag := []string{"01", "02", ""}
	revTag := []string{"", "_R1", "_R2"}
	vidTag := []string{"_sd", "_hd", "_4k"}
	audTag := []string{"_AUDIORUS20", "_AUDIORUSRUS51", "_AUDIOENG20", "_AUDIOENG51", ""}
	extn := []string{".txt", ".m4a", ".srt", ".mp4", ".ready"}
	total := len(base) * len(sTag) * len(epTag) * len(revTag) * len(vidTag) * len(audTag) * len(extn)
	i := 0
	for _, ba := range base {
		for _, sT := range sTag {
			for _, ep := range epTag {
				for _, re := range revTag {
					for _, vi := range vidTag {
						for _, au := range audTag {
							for _, ex := range extn {
								emulatedNames = append(emulatedNames, ba+sT+ep+re+vi+au+ex)
								i++
								fmt.Printf("emulating names %v/%v\r", i, total)
							}
						}
					}
				}
			}
		}
	}
	fmt.Println("")
	return emulatedNames
}
