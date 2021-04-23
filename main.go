package main

func main() {
	// err := termbox.Init()
	// w, h := termbox.Size()
	// termbox.Close()
	// termbox.SetCursor(0, 0)
	// termbox.HideCursor()

	// fmt.Print(w, h, err, "\n")

}

/*




===USER INPUT REQUIRED============
Q: Question:
0	APPLY
1	[x] - Answer1
15	[ ] - Answer2
21	[x] - Answer3
29	[ ] - Answer4
125	[ ] - Answer5
[Error message]
==================================



*/

// func progressBar(now, all int64) string {

// }

/*
10
██████████ 10%

20
████████████████████ 5%
25
█████████████████████████ 4%
30
██████████████████████████████ 4%
40
[█████████████████████████████████████▌  ] 2.5%
50
██████████████████████████████████████████████████ 2%

[1234567890123456789012345678901234567890]
[ progress: 100.000% ]
[>>>>>>              ]
[++++++--------------]

if now%4 >= 2 {
				s += string(rune(9612))
				continue
*/

/*
[]

*/

//DurationStamp - return duration (float64 - seconds) as a string in format: [HH:MM:SS.ms]
// func DurationStamp(dur float64) string {
// 	if dur < 0 {
// 		return "NEGATIVE"
// 	}
// 	stamp := ""
// 	hh := int(dur) / int(3600)
// 	hLeft := int(dur) % int(3600)
// 	mm := hLeft / 60
// 	ss := hLeft % 60
// 	sLeft := dur - (float64(hh*3600) + float64(mm*60) + float64(ss))
// 	ms := int(sLeft * 1000)
// 	////////
// 	hhs := strconv.Itoa(int(hh))
// 	for len(hhs) < 2 {
// 		hhs = "0" + hhs
// 	}
// 	mms := strconv.Itoa(int(mm))
// 	for len(mms) < 2 {
// 		mms = "0" + mms
// 	}
// 	sss := strconv.Itoa(int(ss))
// 	for len(sss) < 2 {
// 		sss = "0" + sss
// 	}
// 	mss := strconv.Itoa(int(ms))
// 	for len(mss) < 3 {
// 		mss = "0" + mss
// 	}
// 	stamp = hhs + ":" + mms + ":" + sss + "." + mss
// 	return stamp

// }
