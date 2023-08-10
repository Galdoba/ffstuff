package fdf

import (
	"fmt"
	"os"
	"time"

	"github.com/Galdoba/ffstuff/pkg/mdm/inputinfo"
)

/*
small: 20-80
medium: 81-130
big: 131+


fd.IsAccesible() bool


*/

type filedata struct {
	path         string
	scandataPath string
	FileInfo     os.FileInfo
	IsDir        bool
	comment      string
}

func newFD(path string) *filedata {
	fd := filedata{}
	//tcur := time.Now()
	fd.comment = "NO DATA"
	FileInfo, _ := os.Stat(path)
	fd.IsDir = FileInfo.IsDir()
	dur := time.Since(FileInfo.ModTime())
	fmt.Println(dur.String())
	pi, _ := inputinfo.ParseFile(path)
	fmt.Println(FMP(pi))
	return &fd
	//.Format("2006-01-02 15:04:05")
}

//FMP - (FileMediaProfile) return short entry on stream data in ehex glyphs:
//Format: 0123-5
//0 - quantity of video streams
//1 - quantity of audio streams
//2 - quantity of data streams
//3 - quantity of subtitle streams
//4 - quantity of warnings
func FMP(pi *inputinfo.ParseInfo) string {
	str := fmt.Sprintf("%v%v%v%v-%v", ehex(len(pi.Video)), ehex(len(pi.Audio)), ehex(len(pi.Data)), ehex(len(pi.Subtitles)), ehex(len(pi.Warnings())))
	return str
}

func ehex(i int) string {
	switch i {
	default:
		return "-"
	case 0:
		return "0"
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "4"
	case 5:
		return "5"
	case 6:
		return "6"
	case 7:
		return "7"
	case 8:
		return "8"
	case 9:
		return "9"
	case 10:
		return "A"
	case 11:
		return "B"
	case 12:
		return "C"
	case 13:
		return "D"
	case 14:
		return "E"
	case 15:
		return "F"
	case 16:
		return "G"
	case 17:
		return "H"
	case 18:
		return "J"
	case 19:
		return "K"
	case 20:
		return "L"
	case 21:
		return "M"
	case 22:
		return "N"
	case 23:
		return "P"
	case 24:
		return "Q"
	case 25:
		return "R"
	case 26:
		return "S"
	case 27:
		return "T"
	case 28:
		return "U"
	case 29:
		return "V"
	case 30:
		return "W"
	case 31:
		return "X"
	case 32:
		return "Y"
	case 33:
		return "Z"
	}
}
