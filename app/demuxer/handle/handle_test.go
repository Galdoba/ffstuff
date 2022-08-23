package handle

import (
	"fmt"
	"testing"
)

func TestPathConversion(t *testing.T) {
	pathO := `ffmpeg -r 25 -i \\nas\buffer\IN\_IN_PROGRESS\z__Film_dlya_testa_MOV_0142.mp4 -filter_complex "[0:a:0]aresample=48000,atempo=25/(25/1)[aud0]; [0:a:1]aresample=48000,atempo=25/(25/1)[aud1]" -map 0:v:0 -c:v libx264 -preset medium -crf 16 -pix_fmt yuv420p -g 0 -map_metadata -1 -map_chapters -1 \\nas\ROOT\EDIT\_amedia\film_dlya_testa_HD.mp4 -map [aud0] -c:a -vn -acodec alac -compression_level 0 -map_metadata -1 -map_chapters -1 \\nas\ROOT\EDIT\_amedia\Film_dlya_testa_AUDIORUS51.m4a -map [aud1] -c:a -vn -acodec alac -compression_level 0 -map_metadata -1 -map_chapters -1 \\nas\ROOT\EDIT\_amedia\Film_dlya_testa_AUDIOENG20.m4a`
	fmt.Println(pathO)
	fmt.Println(" ")
	pathAbs := ConvertToLinux(pathO)
	fmt.Println(pathAbs)
}
