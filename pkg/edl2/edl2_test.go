package edl2

import "testing"

func TestParse(t *testing.T) {
	edlData, err := ParseFile("d:\\IN\\IN_2021-07-12\\Filmz.edl")
	if !edlData.scanningConcluded {
		t.Errorf("scanning not concluded")
	}
	if err != nil {
		t.Errorf("Error: %v", err)
	}

}
