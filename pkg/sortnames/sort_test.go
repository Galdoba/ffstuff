package sortnames

import (
	"fmt"
	"strings"
	"testing"
)

func TestBumpToTopIndex(t *testing.T) {
	slInt := []int{0, 1, 2, 3, 4}
	ind := -5
	newSl := BumpToTopIndex(slInt, ind)
	fmt.Println(slInt)
	fmt.Println(newSl)
	ind2 := 3
	newSl2 := BumpIndexUpByOne(slInt, ind2)
	fmt.Println(slInt)
	fmt.Println(newSl2)
}

func TestNoDuplicates(t *testing.T) {
	test1 := []string{"aa", "aa", "bb", "cc", "bb", "dd", "dd"}
	res1 := OmitDuplicates(test1)
	if strings.Join(res1, "") != "aabbccdd" {
		t.Errorf("test 1 fail")
	}
	//printSliceLn(emulatedNames())
}

func printSliceLn(sl []string) {
	for _, s := range sl {
		fmt.Println(s)
	}
}

func emulatedNames() []string {
	var emulatedNames []string
	base := []string{"serial_AAA", "serial_BBB", "film_AA", "film_BB"}
	sTag := []string{"_s01", "_s02", ""}
	epTag := []string{"03", "06", ""}
	revTag := []string{"", "_R1", "_R2"}
	vidTag := []string{"_sd", "_hd", "_4k"}
	audTag := []string{"_AUDIORUS20", "_AUDIORUS51", "_AUDIOENG20", "_AUDIOENG51", ""}
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

/*
[PROGRESS][SPEED____][SIZE____][FILENAME____________________]
01234567890123456789012345678901234567890123456789012345678901234567890123456789

name1_s01_02_hd.mp4            Complete  124.7 Gb
name1_s01_03_hd.mp4              56 %    158.3 Gb
name1_s01_03_hd.mp4              23 %    258.3 Gb



*/
