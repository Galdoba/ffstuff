package stamp

import (
	"strconv"
	"strings"
	"time"
)

//Seconds - returns string with timestamp in format: [HH:MM:SS]
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

//Date - return current date as a string in format: [YYYY-MM-DD]
func Date() string {
	y, m, d := time.Now().Date()
	yy := strconv.Itoa(y)
	mm := strconv.Itoa(int(m))
	dd := strconv.Itoa(d)
	if int(m) < 10 {
		mm = "0" + mm
	}
	if d < 10 {
		dd = "0" + dd
	}
	return yy + "-" + mm + "-" + dd
}

func IsDate(str string) bool {
	if len(str) != 10 {
		return false
	}
	parts := strings.Split(str, "-")
	for i, part := range parts {
		switch i {
		case 0:
			if len(part) != 4 {
				return false
			}
			yy, err := strconv.Atoi(part)
			if err != nil {
				return false
			}
			if yy < 1000 {
				return false
			}
		}
	}
	return true
}
