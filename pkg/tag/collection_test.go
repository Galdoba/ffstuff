package tag

import (
	"fmt"
	"testing"
)

var inputs = []string{
	`Zlo--s04e13--SER--HSUB--Zlo_s04e13_PRT240816080000_SER_05005_18.mp4`,
	`Sem_smertnyh_grehov_s02_02__HD.mp4`,
	`Zlo_s04_13_PRT240816080000__HD.mp4`,
	`Zlo_s04_12_PRT240809084458_AUDIORUS51_proxy_r1.ac3`,
	`Zlo_s04_11_PRT240802080000.srt`,
}

func TestParsing(t *testing.T) {
	for _, input := range inputs {
		tags := parse_INFILE(input)
		fmt.Println(input)
		fmt.Println(tags)
	}
	fmt.Println("/////////////////")
	for _, input := range inputs {
		tags := parse_OUTFILE(input)
		fmt.Println(input)
		fmt.Println(tags)
	}
}
