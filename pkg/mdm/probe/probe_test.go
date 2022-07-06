package probe

import (
	"fmt"
	"strings"
	"testing"
)

func TestProbe(t *testing.T) {
	paths := []string{
		`d:\IN\IN_testInput\trailers\AllAboutSex_Trailer_Rus_v2_PSH_1_PSH_2_H264.mp4`,
		`d:\IN\IN_testInput\trailers\Belosnezhka_skazka_dlya_vzroslyh_TRL.mp4`,
		`d:\IN\IN_testInput\trailers\ChungkingExpress_stereo_Trailer.mp4`,
		`d:\IN\IN_testInput\trailers\CIApe_TRL_HD24_RU51.mov`,
		`d:\IN\IN_testInput\trailers\CRUELLA_iEST_TLRE_HD_2398_51_20_16x9_185_RUS_D1415623.mov`,
		`d:\IN\IN_testInput\trailers\Dixieland_1080p_185_16x9LB_2398_rus_preview_V1_PRO.mov`,
		`d:\IN\IN_testInput\trailers\EIFFEL_TRL_2.0.mp4`,
		`d:\IN\IN_testInput\trailers\eng20.png`,
		`d:\IN\IN_testInput\trailers\House_of_D_HD23_cutout_RU51.mov`,
		`d:\IN\IN_testInput\trailers\Journal64_TRL_HD24_RU51.mov`,
		`d:\IN\IN_testInput\trailers\Kogda_ona_prihodit_trailer_24fps_444_20_2022_2.mov`,
		`d:\IN\IN_testInput\trailers\Kogda_ona_prihodit_trailer_24fps_444_20_2022_2.wav`,
		`d:\IN\IN_testInput\trailers\RunGun_TrailerB_RunGun_1920x1080_RU.mp4`,
		`d:\IN\IN_testInput\trailers\rus20.png`,
		`d:\IN\IN_testInput\trailers\rus51.png`,
		`d:\IN\IN_testInput\trailers\SOMNAMBULII.Trl.HD.20.rus.clear.VO.mp4`,
		`d:\IN\IN_testInput\trailers\Spasatelnaya_missiya_krausov_AUDIOENG20.m4a`,
		`d:\IN\IN_testInput\trailers\Spasatelnaya_missiya_krausov_AUDIORUS20.m4a`,
		`d:\IN\IN_testInput\trailers\Spasatelnaya_missiya_krausov_AUDIORUS51.m4a`,
		`d:\IN\IN_testInput\trailers\SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).C.wav`,
		`d:\IN\IN_testInput\trailers\SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).L.wav`,
		`d:\IN\IN_testInput\trailers\SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).LFE.wav`,
		`d:\IN\IN_testInput\trailers\SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).Ls.wav`,
		`d:\IN\IN_testInput\trailers\SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).R.wav`,
		`d:\IN\IN_testInput\trailers\SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).Rs.wav`,
		`d:\IN\IN_testInput\trailers\SUPERSTAR_TRL_v4_24fps_SCOPE.mov`,
		`d:\IN\IN_testInput\trailers\TheIllusionofControl_TRL_1080p_RU-XX_20_24_h264-30mbit.mp4`,
		`d:\IN\IN_testInput\trailers\wave3.png`,
		`d:\IN\IN_testInput\trailers\waveformTest.bat	`,
		`\\192.168.31.4\edit\_exchange\#PETR\s05e01_Rostelecom_FLASH_YR05_18_19_NORA_16x9_STEREO_5_1_2_0_LTRT_EPISODE_E2291774_RUSSIAN_ENGLISH_10750107.mpg`,
	}
	for _, path := range paths {
		fmt.Println("  ")
		mo, err := MediaFileReport(path, mediaTypeTrailerHD)
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
				t.Errorf("Video stream dimentions contains error: %v", vs.dimentions)
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
	}
}

func TestIssues(t *testing.T) {
	dim := []dimentions{
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
	target := dimentions{1920, 1080}
	for _, d := range dim {
		issue := dimentionIssue(d, target)
		switch issue {
		default:
			t.Errorf("dimention: %vx%v / %vx%v == %v", d.width, d.height, target.width, target.height, issue)
		case "", "Need Downscale", "Dimention to small for target":
		}

	}

}