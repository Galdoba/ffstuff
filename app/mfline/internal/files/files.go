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
	BadNameMarker = "BAD_NAME--"
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
