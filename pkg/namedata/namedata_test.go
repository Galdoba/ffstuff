package namedata

import (
	"fmt"
	"strings"
	"testing"
)

type testData struct {
	input          string
	expectedOutput string
}

func emulatedNames0() []testData {
	return []testData{
		{"some_name_xx21__hd_eng20-ebur128.ac3", "some_name_xx21__hd_eng20.ac3"},
		{"some_name_xx21__hd_eng51-ebur128.ac3", "some_name_xx21__hd_eng51.ac3"},
		{"some_name_xx21__hd_eng51-ebur128-stereo.ac3", "some_name_xx21__sd_eng20.ac3"},
		{"some_name_xx21__hd_eng20-ebur128-stereo.ac3", "some_name_xx21__sd_eng20.ac3"},
		{"some_name_xx21__hd_eng20_ebur128_stereo.ac3", "some_name_xx21__hd_eng20.ac3"},
		{"some_name_xx21_hd_eng20-ebur128.ac3", "some_name_xx21__hd_eng20.ac3"},
	}
}

func emulatedNames() []string {
	var emulatedNames []string
	base := []string{"film_name_0000", "serial_name_s00_00_0000"}
	vidTag := []string{"sd", "hd", "4k"}
	audTag := []string{"rus20", "rus51", "eng20", "eng51", "qqq20", "qqq51"}
	eburTag := []string{"-ebur128", "-ebur128-stereo", ""}
	resolutions := []string{".ac3", ".aac", ".m4a", ".txt", ".mp4"}
	for _, name := range base {
		for _, vt := range vidTag {
			for _, at := range audTag {
				for _, et := range eburTag {
					for _, rt := range resolutions {
						emulatedNames = append(emulatedNames, name+"__"+vt+"_"+at+et+rt)
					}
				}
			}
		}
	}
	emulatedNames = append(emulatedNames, "film_name_0000__sd.mp4")
	return emulatedNames
}

func TestEburTrimmer(t *testing.T) {
	return
	undecided := 0
	for i, oldName := range emulatedNames() {
		newName, err := TrimLoudnormPrefix(oldName)

		if err != nil {
			if strings.Contains(err.Error(), "invalid name [") {
				continue
			}
			t.Errorf("Test %v:	Error: %v", i, err.Error())
			undecided++
		}
		if newName != oldName {
			t.Errorf("oldname > newname (%v > %v)", oldName, newName)
		}

	}
	fmt.Println("Undecided:", undecided)
}

func testNames() []string {
	return []string{
		// "Banshi_predystoriya_s01e07_HD.mp4",
		// "Banshi_predystoriya_s01e08_AUDIOENG20.m4a",
		// "Banshi_predystoriya_s01e08_AUDIORUS51.m4a",
		// "Banshi_predystoriya_s01e08_HD.mp4",
		// "Banshi_predystoriya_s01e09_AUDIOENG20.m4a",
		// "Banshi_predystoriya_s01e09_AUDIORUS51.m4a",
		// "Banshi_predystoriya_s01e09_HD.mp4",
		// "Davay_znakomitsya_kino_MOV_00350_AUDIOENG20.m4a",
		// "Davay_znakomitsya_kino_MOV_00350_AUDIORUS51.m4a",
		// "Davay_znakomitsya_kino_MOV_00350_HD.mp4",
		// "Dryu_maykl_MOV_00855_AUDIOENG20.m4a",
		// "Dryu_maykl_MOV_00855_AUDIORUS20.m4a",
		// "Dryu_maykl_MOV_00855_HD.mp4",
		// "Dryu_maykl_MOV_00855_SUB.srt",
		// "Dve_zvezdy_kerri_fisher_i_debbi_reynolds_MOV_00554_AUDIOENG20.m4a",
		// "Dve_zvezdy_kerri_fisher_i_debbi_reynolds_MOV_00554_AUDIORUS51.m4a",
		// "Dve_zvezdy_kerri_fisher_i_debbi_reynolds_MOV_00554_HD.mp4",
		// "Titane_AUDIORUS51.m4a",
		// "Titane_HD.mp4",
		// "Titane_SD.mp4",
		"Край земли. 01 сезон. 01 серия (Edge of The Earth)",
		`Мы владеем этим городом. 01 сезон. 01 серия`,
		`Под угрозой (Endangered)`,
		`По волчьим законам. 06 сезон. 05 серия`,
		`Обратная сторона красоты. 01 сезон. 04 серия`,
		`Лестница. 01 сезон. 01 серия`,
		`Лестница. 01 сезон. 02 серия`,
		`Лестница. 01 сезон. 03 серия`,
		`Лестница. 01 сезон. 04 серия`,
		`Край земли. 01 сезон. 01 серия (Edge of The Earth)`,
		`Малыш. 01 сезон. 01 серия`,
		`Малыш. 01 сезон. 02 серия`,
		`Малыш. 01 сезон. 03 серия`,
		`The Baby. 01 сезон. 04 серия`,
		`The Staircase. 01 сезон. 01 серия`,
		`The Staircase. 01 сезон. 02 серия`,
		`The Staircase. 01 сезон. `,
	}
}

func TestNameSplitting(t *testing.T) {
	return
	for _, val := range testNames() {
		fmt.Println(val)
		nf := ParseName(val)
		fmt.Println(nf)
		fmt.Println(nf.ReconstructName())
	}
	rnMap, err := RenamerMap()
	fmt.Println(err)
	for k, v := range rnMap {
		fmt.Println(k, "*", v)
	}
}

func TestTransliterate(t *testing.T) {
	for _, name := range testNames() {
		res := TransliterateForEdit(name)
		byRune := strings.Split(res, "")
		fmt.Println("Name  :", name)
		fmt.Println("Result:", res)
		for _, rn := range byRune {

			switch rn {
			default:

				t.Errorf("must not have letter '%v'", rn)
			case "A", "a", "B", "b", "C", "c", "D", "d", "E", "e", "F", "f", "G", "g", "H", "h", "I", "i", "J", "j", "K", "k", "L", "l", "M", "m", "N", "n", "O", "o", "P", "p", "Q", "q", "R", "r", "S", "s", "T", "t", "U", "u", "V", "v", "W", "w", "X", "x", "Y", "y", "Z", "z", "_", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			}
		}
	}
}
