package silence

import (
	"fmt"
	"testing"
)

func TestDetect(t *testing.T) {
	for n, path := range testPaths() {
		debugMsg(fmt.Sprintf("start test %v (path: %v)\n", n, path))
		si, err := Detect(path)
		if err != nil {
			t.Errorf("Detect returned error:\nfile = '%v'\nerror = '%v'\n", si, err)
		}
		if si == nil {
			t.Errorf("'Silence' object not returned")
		}
		debugMsg(fmt.Sprintf("end test"))
	}
}

func testPaths() []string {
	return []string{
		// "d:\\MUX\\tests\\Shang-Chi_and_the_Legend_of_the_Ten_Rings_HD.mp4",
		// "d:\\MUX\\tests\\s05e01_Rostelecom_FLASH_YR05_18_19_NORA_16x9_STEREO_5_1_2_0_LTRT_EPISODE_E2291774_RUSSIAN_ENGLISH_10750107.mpg",
		// "d:\\MUX\\tests\\strela_s03_03_2014__hd_ar2.mp4",
		"d:\\MUX\\tests\\Ryad_19_TRL_AUDIORUS51.m4a",
		// "d:\\MUX\\tests\\Shang-Chi_and_the_Legend_of_the_Ten_Rings_AUDIORUS51.m4a",
		// "d:\\MUX\\tests\\s05e01_Rostelecom_FLASH_YR05_18_19_NORA_16x9_STEREO_5_1_2_0_LTRT_EPISODE_E2291774_RUSSIAN_ENGLISH_10750107.ac3",
		// "d:\\MUX\\tests\\screenshot_bl1.bmp",
		// "d:\\MUX\\tests\\screenshot_bl2.bmp",
		// "d:\\MUX\\tests\\screenshot_bl3.bmp",
		// "d:\\MUX\\tests\\log.txt",
		// "d:\\MUX\\tests\\output2.png",
		// "d:\\MUX\\tests\\output1.png",
		// "d:\\MUX\\tests\\waveform.bat",
		// "d:\\MUX\\tests\\mauris.bat",
		// "d:\\MUX\\tests\\s05e01_Rostelecom_FLASH_YR05_18_19_NORA_16x9_STEREO_5_1_2_0_LTRT_EPISODE_E2291774_RUSSIAN_ENGLISH_10750107.m4a",
	}
}
