package edl3

import "testing"

func TestCanCollectInfo(t *testing.T) {
	mix := &mix{}
	stm1 := statementData{"STANDARD", []string{"001", "BL", "C", "0.0", "00:00:00.000", "00:00:00.000", "00:00:00.000", "00:00:00.000"}}
	stm2 := statementData{"STANDARD", []string{"001", "AX", "D", "25.0", "00:01:52.280", "00:01:54.280", "00:00:00.000", "00:00:02.000"}}
	stm3 := statementData{"SOURCE A", []string{"BL"}}
	stm4 := statementData{"SOURCE B", []string{"[filePathExample.mp4]"}}
	info := []statementData{stm1, stm2, stm3, stm4}
	mix.CollectInfo(info)
	if mix.sourceA == "" {
		t.Errorf("Source A == '%v', expect %v", mix.sourceA, stm3.fields[0])
	}
	if mix.sourceB == "" {
		t.Errorf("Source B == '%v', expect %v", mix.sourceA, stm4.fields[0])
	}
	if mix.inPointA.String() == stm1.fields[4] {
		t.Errorf("Source inPoint A == '%v', expect %v", mix.inPointA.String(), stm1.fields[4])
	}
}
