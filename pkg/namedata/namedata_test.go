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
	base := []string{"serial_name_s01_01", "serial_name_s01_02", "serial_name_s01_03"}
	vidTag := []string{"sd", "hd", "4k"}
	audTag := []string{"rus20", "rus51", "eng20", "eng51", "qqq20", "qqq51"}
	eburTag := []string{"-ebur128", "-ebur128-stereo", ""}
	resolutions := []string{".ac3", ".aac", ".m4a", ".txt", ".mp4", ".ready"}
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

func emulateEditAmediaNames() []string {
	var emulatedNames []string
	base := []string{"serial_name_s01_01", "serial_name_s01_02", "serial_name_s01_03"}
	vidTag := []string{"_SD", "_HD", "_4K", ""}
	audTag := []string{"_AUDIORUS20", "_AUDIOENG51", ""}
	proxyTag := []string{"_proxy", ""}
	resolutions := []string{".m4a", ".mp4", ".srt", ".ready"}
	for _, name := range base {
		for _, vt := range vidTag {
			for _, at := range audTag {
				for _, pt := range proxyTag {
					for _, rt := range resolutions {
						switch rt {
						case ".ready", ".srt":
							if vt+at+pt != "" {
								continue
							}
						case ".mp4":
							if at != "" {
								continue
							}
							if vt == "" {
								continue
							}
						case ".m4a":
							if at == "" {
								continue
							}
							if vt != "" {
								continue
							}
						}

						emulatedNames = append(emulatedNames, name+vt+at+pt+rt)
					}
				}
			}
		}
	}
	return emulatedNames
}

func TestSort(t *testing.T) {
	return
	list := emulateEditAmediaNames()
	for i := range list {
		fmt.Println(i, list[i])
	}
	sorted, err := SortFileNames(list, EDIT)
	if err != nil {
		t.Errorf("SortFileNames() вернула ошибку: %v", err.Error())
		panic(0)
	}
	fmt.Println("sorted:")
	for i := range sorted {
		fmt.Println(sorted[i])
	}
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
		//"Край земли. 01 сезон. 01 серия (Edge of The Earth)",
		//`Мы владеем этим городом. 01 сезон. 01 серия`,
		//`Под угрозой (Endangered)`,
		//`По волчьим законам. 06 сезон. 05 серия`,
		//`Обратная сторона красоты. 01 сезон. 04 серия`,
		//`Лестница. 01 сезон. 01 серия`,
		//`Лестница. 01 сезон. 02 серия`,
		//`Лестница. 01 сезон. 03 серия`,
		//`Лестница. 01 сезон. 04 серия`,
		//`Край земли. 01 сезон. 01 серия (Edge of The Earth)`,
		//`Малыш. 01 сезон. 01 серия`,
		//`Малыш. 01 сезон. 02 серия`,
		//`Малыш. 01 сезон. 03 серия`,
		//`The Baby. 01 сезон. 04 серия`,
		//`The Staircase. 01 сезон. 01 серия`,
		//`The Staircase. 01 сезон. 02 серия`,
		//`The Staircase. 01 сезон. `,
		//`Род мужской (Men) 4K`,
		//"Jumong_Prince_Of_the_Legend_13.mov",
		//"Jumong_Prince_Of_the_Legend_14.mov",
		//"Jumong_Prince_Of_the_Legend_15.mov",
		//"Jumong_Prince_Of_the_Legend_16.mov",
		//"Jumong_Prince_Of_the_Legend_17.mov",
		//"Spy_Myeong_Wol_03.mov",
		//"Spy_Myeong_Wol_04.mov",
		//"Spy_Myeong_Wol_05.mov",
		//"Spy_Myeong_Wol_06.mov",
		//"Spy_Myeong_Wol_07.mov",
		//"Spy_Myeong_Wol_08.mov",
		//"Spy_Myeong_Wol_09.mov",
		//"Spy_Myeong_Wol_10.mov",
		//"Spy_Myeong_Wol_11.mov",
		//"Spy_Myeong_Wol_12.mov",
		//"Spy_Myeong_Wol_13.mov",
		//"Spy_Myeong_Wol_14.mov",
		//"Spy_Myeong_Wol_15.mov",
		//"Spy_Myeong_Wol_16.mov",
		//"Spy_Myeong_Wol_17.mov",
		//"Spy_Myeong_Wol_18.mov",
		//`Chumon_s01_01_AUDIORUS51.m4a`,
		//`Chumon_s01_01_AUDIORUS51_proxy.ac3`,
		//`Chumon_s01_01_HD.mp4`,
		//`Chumon_s01_01_HD_proxy.mp4`,
		//`Chumon_s01_02_AUDIORUS51.m4a`,
		//`Chumon_s01_02_AUDIORUS51_proxy.ac3`,
		//`Chumon_s01_02_HD.mp4`,
		//`Chumon_s01_02_HD_proxy.mp4`,
		//"VRfighter_feature_ProRes422HQ_1080p24_RU20_RU51.mov",
		`Bortprovodnica_s02e01_SER_12143.mp4`,
		`Bortprovodnica_s02e01_SER_12143.RUS.srt`,
		`Bortprovodnica_s02e02_SER_12144.mp4`,
		`Bortprovodnica_s02e02_SER_12144.RUS.srt`,
		`Bortprovodnica_s02e03_SER_12145.mp4`,
		`Bortprovodnica_s02e03_SER_12145.RUS.srt`,
		`Bortprovodnica_s02e04_SER_12146.mp4`,
		`Bortprovodnica_s02e04_SER_12146.RUS.srt`,
		`Bortprovodnica_s02e05_SER_12147.mp4`,
		`Bortprovodnica_s02e05_SER_12147.RUS.srt`,
		`Bortprovodnica_s02e06_SER_12148.mp4`,
		`Bortprovodnica_s02e06_SER_12148.RUS.srt`,
		`Bortprovodnica_s02e07_SER_12149.mp4`,
		`Bortprovodnica_s02e07_SER_12149.RUS.srt`,
		`Bortprovodnica_s02e08_SER_12150.mp4`,
		`Bortprovodnica_s02e08_SER_12150.RUS.srt`,
		//`Zalozhniki_s01_01.srt`,
		//`Zalozhniki_s01_01.ready`,
		//`Zalozhniki_s01_01_AUDIOENG20.m4a`,
		//`Zalozhniki_s01_01_AUDIOENG20_proxy.ac3`,
		//`Zalozhniki_s01_01_AUDIORUS51.m4a`,
		//`Zalozhniki_s01_01_AUDIORUS51_proxy.ac3`,
		//		`\\nas\ROOT\EDIT\_amedia\Zalozhniki_s01\Zalozhniki_s01_01_AUDIOENG20.m4a`,
		//		`\\nas\ROOT\EDIT\_amedia\Zalozhniki_s01\Zalozhniki_s01_01.ready`,
	}
}

func Test_SearchNameMask(t *testing.T) {
	return
	fmt.Println("start Test_SearchNameMask")
	mask, err := SearchMask(testNames())
	for _, name := range testNames() {
		fmt.Printf("%v\n", name)
	}
	fmt.Println("")
	fmt.Println(mask.matchPattern)
	fmt.Println(mask.typePattern)
	fmt.Println(err)
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
	return
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

func TestNamesWithMonos(t *testing.T) {
	names := []string{
		`\\192.168.31.4\buffer\IN\_AWAIT\FLOOD_SOUND_C.wav`,
		`\\192.168.31.4\buffer\IN\_AWAIT\FLOOD_SOUND_L.wav`,
		`\\192.168.31.4\buffer\IN\_AWAIT\FLOOD_SOUND_Lfe.wav`,
		`\\192.168.31.4\buffer\IN\_AWAIT\FLOOD_SOUND_Ls.wav`,
		`\\192.168.31.4\buffer\IN\_AWAIT\FLOOD_SOUND_R.wav`,
		`\\192.168.31.4\buffer\IN\_AWAIT\FLOOD_SOUND_Rs.wav`,
		`\\192.168.31.4\buffer\IN\_AWAIT\hd_2019_detstvo_sheldona_s03_02__ar2_vkurazhbambey_xmAsIDJ2MwW_film.mp4`,
	}
	fmt.Println(have6mono(names))
}
