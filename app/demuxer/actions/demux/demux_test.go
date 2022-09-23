package actiondemux

import (
	"fmt"
	"testing"

	"github.com/Galdoba/ffstuff/pkg/mdm/probe"
)

func TestSoundSelection(t *testing.T) {
	return
	rep, _ := probe.StreamsData(`\\nas\ROOT\IN\_UMS\_DONE\Haker\HACKED_Prores422HQ_8ch_Master_012722.mov`, `\\nas\ROOT\IN\_UMS\_DONE\Haker\Hacked_[rus].ac3`, `\\nas\ROOT\IN\_UMS\_DONE\Proklyatie_mayya\CurseOfTheMayansThe_Feat_1080p2398_178_EN51_EN20.mov`, `\\nas\ROOT\IN\_UMS\_DONE\Proklyatie_mayya\The.Curse.of.the.Mayans_[rus].ac3`, `\\nas\ROOT\IN\_MEGO_DISTRIBUSHN\_DONE\Vse_putem\Все_путем.mkv`)
	for i, str := range rep {
		fmt.Printf("%v	%v\n", i, str.PrintStreamData())
	}
	vid, aud := probe.SeparateByTypes(rep)
	for _, v := range vid {
		fmt.Println(v.String())
	}
	fmt.Println(vid[0].String())
	fmt.Println(aud[0].String())

}

func TestInterlaceSearch(t *testing.T) {
	for _, path := range []string{
		`Rusalki_ostrova_mako_s01e01_130-255335_.mov`,
	} {
		found, err := probe.InterlaceByIdet(path)
		switch found {
		case true:
			fmt.Println("FOUND")
		case false:
			fmt.Println("NOT FOUND")
		}
		fmt.Println(err)
	}
}

func selectorEmulation(options []probe.AudioData) []probe.AudioData {
	selected := []probe.AudioData{}
	for i := 0; i < len(options); i++ {
		switch i {
		default:
		case 0, 1, 2, 3, 4, 5:
			selected = append(selected, options[i])
		}
	}
	return selected
}

func merge6monoTo51(str []probe.AudioData) (probe.AudioData, error) {
	aud := probe.AudioData{}
	if len(str) != 6 {
		return aud, fmt.Errorf("ожидаем 6 стримов (получили %v)", len(str))
	}
	aud.SetChanLayout("5.1")

	return aud, nil
}
