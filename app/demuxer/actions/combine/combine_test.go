package actioncombine

import (
	"fmt"
	"testing"
)

func input() [][]string {
	return [][]string{
		{`SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).C.wav SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).L.wav`, `SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).LFE.wav`, `SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).Ls.wav`, `SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).R.wav`, `SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).Rs.wav`},
		{`SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).C.wav`, `SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).L.wav`, `SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).LFE.wav`, `SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).Ls.wav`, `SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).R.wav`, `SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).Rs.wav`},
	}
}

func TestInput(t *testing.T) {
	for i, inputData := range input() {
		fmt.Println("Test", i+1, inputData)
		channelMap, err := autoMap(inputData)
		if err != nil {
			t.Errorf(`autoMap returned error: %v`, err.Error())
			continue
		}
		for _, tag := range []string{"L", "R", "C", "LFE", "Ls", "Rs"} {
			fmt.Printf("TAG: %v = %v\n", tag, channelMap[tag])
		}
	}
}
