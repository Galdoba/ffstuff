package info

import (
	"fmt"
	"testing"
)

func TestLoudnormData(t *testing.T) {
	data := LoudnormData("d:\\IN\\IN_2022-02-17\\proxy\\The_Magicians_s03e04_AUDIORUS20_loudnorm_report.txt")
	fmt.Println(data, "----")
	for i, d := range data {
		fmt.Println(i, d)
	}
}
