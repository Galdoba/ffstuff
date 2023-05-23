package sortnames

import (
	"fmt"
	"strings"
	"testing"
)

func TestBumpToTopIndex(t *testing.T) {
	return
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

func TestSeacrhFileNameBases(t *testing.T) {
	nm := names()
	//words := searchWords(nm)
	//fmt.Println(words)
	// for _, n := range nm {
	// 	//fmt.Println(namedata.EditForm(n), namedata.EditForm(n).HasTags("AUDIO", "proxy"))
	// }
	for _, base := range GrabberOrder(nm) {
		fmt.Println(base)
	}
}

/*
[PROGRESS][SPEED____][SIZE____][FILENAME____________________]
01234567890123456789012345678901234567890123456789012345678901234567890123456789

name1_s01_02_hd.mp4            Complete  124.7 Gb
name1_s01_03_hd.mp4              56 %    158.3 Gb
name1_s01_03_hd.mp4              23 %    258.3 Gb



*/
