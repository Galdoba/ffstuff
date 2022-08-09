package probe

import (
	"fmt"
	"strings"
	"testing"
)

func input() []string {
	return []string{
		`\\nas\buffer\IN\_REJECTED\Artek_bolshoe_puteshestvie_TRL_AUDIORUS20.m4a`,
		`\\nas\buffer\IN\_REJECTED\План_побега_2_Escape_Plan_2_Hades_2018_2.mkv`,
		`\\nas\buffer\IN\_REJECTED\Corrective_Measures_HDSDR25f_RUS20LR_RUS51LRCLfeLsRs.mov`,
		`\\nas\buffer\IN\_REJECTED\LOVE_HATE_AND_SECURITY_MOVIE_RUS.mov`,
		`\\nas\buffer\IN\_REJECTED\LOVE_HATE_AND_SECURITY_MOVIE_RUS_R1.mov`,
		`\\nas\buffer\IN\_REJECTED\TURN_BACK_MOVIES_RUS.mov`,
		`\\nas\buffer\IN\_REJECTED\Artek_bolshoe_puteshestvie_TRL_HD.mp4`,
		`\\nas\buffer\IN\_REJECTED\Esli_by_steny_mogli_govorit_2_MOV_00079.mp4`,
		`\\nas\buffer\IN\_REJECTED\out.mp4`,
		`\\nas\buffer\IN\_REJECTED\Transit_Trailer.mp4`,
		`\\nas\buffer\IN\_REJECTED\Красный_штат.mp4`,
		`\\nas\buffer\IN\_REJECTED\Крутой_поворот.mp4`,
		`\\nas\buffer\IN\_REJECTED\Плохое_поведение_Behaving-Badly_en.mp4`,
		`\\nas\buffer\IN\_REJECTED\Подстава_TRL.mp4`,
		`\\nas\buffer\IN\_REJECTED\Приключения_мышонка_TRL.mp4`,
		`\\nas\buffer\IN\_REJECTED\Прощай_моя_королева_Les_adieux_a_la_reine_or.mp4`,
		`\\nas\buffer\IN\_REJECTED\Шеф_Comme_un_chef_ru.mp4`,
		`\\nas\buffer\IN\_REJECTED\Powder_Blue_XDCAM_HD422_1080i_ENG_2.0.mxf`,
		`\\nas\buffer\IN\_REJECTED\Zoe_XDCAM_HD422_1080i_ENG_2.0.mxf`,
		`\\nas\buffer\IN\_REJECTED\TURN_BACK_MOVIES_RUS.mov_wavespic300.png`,
		`\\nas\buffer\IN\_REJECTED\Esli_by_steny_mogli_govorit_2_MOV_00079.srt`,
	}
}

func TestProbe(t *testing.T) {
	for _, path := range input() {
		fmt.Println("  ")
		fmt.Println(path)
		mo, err := NewReport(path)
		fmt.Println(mo.String())
		if err != nil {
			fmt.Println(err.Error())
			t.Errorf("Media(path) returned error: %v", err.Error())
		}
		if mo == nil {
			t.Errorf("Media(path) returned no object")
			continue
		}
		for _, vs := range mo.vData {
			if strings.Contains(vs.fps, "unknown") {
				t.Errorf("Video stream fps contains 'unknown': %v", vs)
			}
			if vs.dimentions.height == 0 || vs.dimentions.width == 0 {
				t.Errorf("Video stream Dimentions contains error: %vx%v", vs.dimentions.width, vs.dimentions.height)
			}
		}
		for _, as := range mo.aData {
			if strings.Contains(as.chanLayout, "unknown") {

				t.Errorf("Audio stream chanLayout contains 'unknown': %v/%v", as.chanLayout, mo.filename)
			}
			if as.chanNum <= 0 {
				t.Errorf("Audio stream chanelNum expected to be atleast 1 (hase %v)", as.chanNum)
			}
			if as.sampleRate <= 0 {
				t.Errorf("Audio stream sampleRate expected to be atleast 1 (hase %v)", as.sampleRate)
			}
		}
		fmt.Println("End")
	}
}

func TestIssues(t *testing.T) {
	return
	dim := []Dimentions{
		{1920, 1080},
		{1940, 1080},
		{1900, 1080},
		{1920, 1060},
		{1940, 1060},
		{1900, 1060},
		{1920, 1090},
		{1940, 1090},
		{1900, 1090},
	}
	target := Dimentions{1920, 1080}
	for _, d := range dim {
		issue := dimentionIssue(d, target)
		switch {
		default:
			t.Errorf("dimention: %vx%v / %vx%v == %v", d.width, d.height, target.width, target.height, issue)

		}

	}

}

func TestInterlaceDetect(t *testing.T) {
	return
	for i, path := range input() {
		fmt.Printf("   \n")
		fmt.Printf("test %v: %v", i, path)
		confirmed, err := InterlaceByIdet(path)
		if err != nil {
			t.Errorf("func returned error: %v", err)
		}
		fmt.Println("confirmed", confirmed)

	}
}
