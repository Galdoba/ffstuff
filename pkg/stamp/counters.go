package stamp

import "strconv"

func Seconds(seconds int64) string {
	hh := seconds / 3600
	mm := (seconds - (hh * 3600)) / 60
	for mm > 60 {
		hh++
		mm -= 60
	}
	ss := seconds % 60
	hStr := strconv.Itoa(int(hh))
	mStr := strconv.Itoa(int(mm))
	sStr := strconv.Itoa(int(ss))
	if len(hStr) < 2 {
		hStr = "0" + hStr
	}
	if len(mStr) < 2 {
		mStr = "0" + mStr
	}
	if len(sStr) < 2 {
		sStr = "0" + sStr
	}
	return hStr + ":" + mStr + ":" + sStr
}
