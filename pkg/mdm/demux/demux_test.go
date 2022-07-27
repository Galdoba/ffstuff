package demux

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Galdoba/ffstuff/pkg/mdm/format"
	"github.com/Galdoba/ffstuff/pkg/mdm/probe"
)

type demInp struct {
	taskFormat string
	tableData  []string
	paths      []string
}

func demuxerInput() []demInp {
	input := []demInp{
		{format.FilmHD, []string{"", "", "", "", "", "", "", "O", "Film Name", "y", "", "", "", "CONTRAGENT", "22.22.22"}, []string{`d:\IN\IN_testInput\trailers\AllAboutSex_Trailer_Rus_v2_PSH_1_PSH_2_H264.mp4`}},
		{format.FilmSD, []string{"", "", "", "", "", "", "", "O", "Film Name", "y", "", "", "", "CONTRAGENT", "22.22.22"}, []string{`d:\IN\IN_testInput\trailers\AllAboutSex_Trailer_Rus_v2_PSH_1_PSH_2_H264.mp4`}},
		{format.Film4K, []string{"", "", "", "", "", "", "", "O", "Film Name", "y", "", "", "", "CONTRAGENT", "22.22.22"}, []string{`d:\IN\IN_testInput\trailers\AllAboutSex_Trailer_Rus_v2_PSH_1_PSH_2_H264.mp4`}},
		{format.FilmHD, []string{"", "", "", "", "", "", "", "O", "Film Name SD", "y", "", "", "", "CONTRAGENT", "22.22.22"}, []string{`d:\IN\IN_testInput\trailers\AllAboutSex_Trailer_Rus_v2_PSH_1_PSH_2_H264.mp4`}},
		{format.FilmSD, []string{"", "", "", "", "", "", "", "O", "Film Name SD", "y", "", "", "", "CONTRAGENT", "22.22.22"}, []string{`d:\IN\IN_testInput\trailers\AllAboutSex_Trailer_Rus_v2_PSH_1_PSH_2_H264.mp4`}},
		{format.Film4K, []string{"", "", "", "", "", "", "", "O", "Film Name SD", "y", "", "", "", "CONTRAGENT", "22.22.22"}, []string{`d:\IN\IN_testInput\trailers\AllAboutSex_Trailer_Rus_v2_PSH_1_PSH_2_H264.mp4`}},
		{format.FilmHD, []string{"", "", "", "", "", "", "", "O", "Film Name 4K", "y", "", "", "", "CONTRAGENT", "22.22.22"}, []string{`d:\IN\IN_testInput\trailers\AllAboutSex_Trailer_Rus_v2_PSH_1_PSH_2_H264.mp4`}},
		{format.FilmSD, []string{"", "", "", "", "", "", "", "O", "Film Name 4K", "y", "", "", "", "CONTRAGENT", "22.22.22"}, []string{`d:\IN\IN_testInput\trailers\AllAboutSex_Trailer_Rus_v2_PSH_1_PSH_2_H264.mp4`}},
		{format.Film4K, []string{"", "", "", "", "", "", "", "O", "Film Name 4K", "y", "", "", "", "CONTRAGENT", "22.22.22"}, []string{`d:\IN\IN_testInput\trailers\AllAboutSex_Trailer_Rus_v2_PSH_1_PSH_2_H264.mp4`}},
		{format.Film4K, []string{"", "", "", "", "", "", "", "O", "Film Name", "y", "", "", "", "CONTRAGENT", "22.22.22"}, []string{`d:\IN\IN_testInput\trailers\AllAboutSex_Trailer_Rus_v2_PSH_1_PSH_2_H264.mp4`, `d:\IN\IN_testInput\trailers\SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).C.wav`, `d:\IN\IN_testInput\trailers\SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).L.wav`, `d:\IN\IN_testInput\trailers\SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).LFE.wav`, `d:\IN\IN_testInput\trailers\SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).Ls.wav`, `d:\IN\IN_testInput\trailers\SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).R.wav`, `d:\IN\IN_testInput\trailers\SUPERSTAR_TRL_5.1_MIX_86.0dB(Int).Rs.wav`}},
	}
	return input
}

func input() []string {
	return []string{
		`d:\IN\IN_testInput\trailers\TOBOT_S3_RUS_EP02.mxf`,
		`d:\IN\IN_testInput\trailers\Odin_idealnyy_kadr_s01e06_SER_12167.mp4`,
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
}

func Test_AllAsIs(t *testing.T) {
	return
	paths := input()
	for _, path := range paths {
		break
		com, err := AllAsIs(path)
		if err != nil {
			t.Errorf("AllAsIs(%v)\n\treturned error: %v", path, err.Error())
			continue
		}
		fmt.Println("path:", path)
		fmt.Println("command:", com)
		fmt.Println("  ")

	}
}

//ffmpeg -i d:\IN\IN_testInput\trailers\CRUELLA_iEST_TLRE_HD_2398_51_20_16x9_185_RUS_D1415623.mov -map 0:0:0 -c:v copy d:\IN\IN_testInput\trailers\CRUELLA_iEST_TLRE_HD_2398_51_20_16x9_185_RUS_D1415623_RAW_0.mp4  -map 0:0:0 -c:a copy d:\IN\IN_testInput\trailers\CRUELLA_iEST_TLRE_HD_2398_51_20_16x9_185_RUS_D1415623_RAW_0.wav  -map 0:1:0 -c:a copy d:\IN\IN_testInput\trailers\CRUELLA_iEST_TLRE_HD_2398_51_20_16x9_185_RUS_D1415623_RAW_1.wav  -map 0:2:0 -c:a copy d:\IN\IN_testInput\trailers\CRUELLA_iEST_TLRE_HD_2398_51_20_16x9_185_RUS_D1415623_RAW_2.wav  -map 0:3:0 -c:a copy d:\IN\IN_testInput\trailers\CRUELLA_iEST_TLRE_HD_2398_51_20_16x9_185_RUS_D1415623_RAW_3.wav  -map 0:4:0 -c:a copy d:\IN\IN_testInput\trailers\CRUELLA_iEST_TLRE_HD_2398_51_20_16x9_185_RUS_D1415623_RAW_4.wav  -map 0:5:0 -c:a copy d:\IN\IN_testInput\trailers\CRUELLA_iEST_TLRE_HD_2398_51_20_16x9_185_RUS_D1415623_RAW_5.wav  -map 0:6:0 -c:a copy d:\IN\IN_testInput\trailers\CRUELLA_iEST_TLRE_HD_2398_51_20_16x9_185_RUS_D1415623_RAW_6.wav  -map 0:7:0 -c:a copy d:\IN\IN_testInput\trailers\CRUELLA_iEST_TLRE_HD_2398_51_20_16x9_185_RUS_D1415623_RAW_7.wav

func TestMapping(t *testing.T) {
	return
	for n, path := range input() {
		fmt.Printf("test %v: %v\n", n, path)
		str, err := Mapping(path, probe.MediaTypeFilmHD)
		if err != nil {
			fmt.Printf("func returned error: %v\n", err.Error())
		}
		fmt.Printf("func returned []string: \n")
		for _, sl := range str {
			fmt.Println(sl)
		}

	}
}

func TestDemuxer(t *testing.T) {
	for i, demIn := range demuxerInput() {
		//fmt.Printf("Test %v:\n%v\n%v\n%v\n  \n", i+1, demIn.taskFormat, demIn.paths, demIn.taskData)
		fmt.Printf("\n \n Test %v: %v\n", i+1, demIn)
		name := demIn.tableData[8]
		agent := demIn.tableData[13]
		pubDate := demIn.tableData[14]
		d, err := New(demIn.taskFormat, name, agent, pubDate, demIn.paths...)
		if d == nil {
			t.Errorf("Func return no object")
		}
		if err != nil {
			t.Errorf("Func return error: %v", err.Error())
		}
		if len(d.sourcePaths) < 1 {
			t.Errorf("source paths not set (expect: %v)", demIn.paths)
		}
		if d.targetFormat == nil {
			t.Errorf("Target Format not declares (expect: %v)", demIn.taskFormat)
		}
		if d.tableData.Name() == "" {
			t.Errorf("Task data name not set (expect '%v')", name)
		}
		if d.tableData.Agent() == "" {
			t.Errorf("Task data agent not set (expect '%v')", agent)
		}
		if d.tableData.PublicationDate() == "          " {
			t.Errorf("Task data Publication Date not set (expect '%v')", pubDate)
		}

		if strings.Contains(d.TaskDataName(), " SD") {
			sdPreset, _ := format.SetAs(format.FilmSD)
			if d.targetFormat.Dimention.MaxSize() != sdPreset.Dimention.MaxSize() {
				t.Errorf("Task data declares `FILM SD` as target format, have '%v' != '%v'", d.targetFormat, sdPreset.Dimention.MaxSize())
			}
		}
		if strings.Contains(d.TaskDataName(), " 4K") {
			preset4K, _ := format.SetAs(format.Film4K)
			if d.targetFormat != preset4K {
				t.Errorf("Task data declares `FILM 4K` as target format, but have '%v'", demIn.taskFormat)
			}
		}

		argMap := d.MapArguments()
		for _, v := range argMap {
			if v == "" {
				//t.Errorf("map argument '%v' is unset", k)
			}
		}
	}
}
