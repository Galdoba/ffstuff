package atempo

import "fmt"

func calibrate(freq int, dur float64) string {
	sttr := ""
	sttr = fmt.Sprintf("aresample=48000,atempo=%v", "asd")
	return sttr
}

//ffmpeg -i audio100Sec.aac -af aresample=48000,atempo=25/24 out2524.aac (1.041666666666667)  00:01:40.00 100.0 00:01:36:01 96.040
//ffmpeg -i audio100Sec.aac -af asetrate=48000,aresample=48000,atempo=(250/215) out86_2.aac

/*
TASK:
взять звук с длинной Х и сделать длинну У



страя
----- = tempo
новая

old = tempo * new

tempo = old/new

new = old/tempo


100/86
*/

// 25/x = 1.16279
// 25 = 1.16279x
// x= 25/1.16279

// x = 25 / (wantFPS / currentFPS)

// 250/215

// 1.162790697674419
