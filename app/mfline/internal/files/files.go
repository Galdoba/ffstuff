package files

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ListDir(dir string) []string {
	fls, _ := os.ReadDir(dir)
	names := []string{}
	for _, fl := range fls {
		if fl.IsDir() {
			continue
		}
		names = append(names, dir+fl.Name())
	}
	return names
}

const (
	BadNameMarker         = "BAD_NAME--"
	RW_Marker             = "SCANNING_RW--"
	Interlace_Marker      = "SCANNING_INTERLACE--"
	SCANS_COMPLETE_MARKER = "READY--"
)

func BadName(name string) bool {
	re := regexp.MustCompile(`(\ |А|а|Б|б|В|в|Г|г|Д|д|Е|е|Ё|ё|Ж|ж|З|з|И|и|Й|й|К|к|Л|л|М|м|Н|н|О|о|П|п|Р|р|С|с|Т|т|У|у|Ф|ф|Х|х|Ц|ц|Ч|ч|Ш|ш|Щ|щ|Ъ|ъ|Ы|ы|Ь|ь|Э|э|Ю|ю|Я|я)`)
	find := re.FindString(name)
	return find != ""
}

func MarkAsBad(name string) error {
	if strings.Contains(name, BadNameMarker) {
		return nil
	}
	dir := filepath.Dir(name)
	base := filepath.Base(name)
	return os.Rename(name, dir+string(filepath.Separator)+BadNameMarker+base)
}

func MarkScan(name, marker string) (string, error) {
	if strings.Contains(name, marker) {
		return name, nil
	}
	dir := filepath.Dir(name)
	base := filepath.Base(name)
	return dir + string(filepath.Separator) + marker + base, os.Rename(name, dir+string(filepath.Separator)+marker+base)
}

func ClearMarkers(name string) error {
	dir := filepath.Dir(name)
	base := filepath.Base(name)
	markers := strings.Split(base, "--")
	base = markers[len(markers)-1]
	return os.Rename(name, dir+string(filepath.Separator)+base)
}

func hasMarker(name string) string {
	base := filepath.Base(name)
	for _, marker := range []string{
		RW_Marker,
	} {
		if strings.Contains(base, marker) {
			return marker
		}
	}
	return ""
}
