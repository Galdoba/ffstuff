package actions

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Galdoba/ffstuff/app/aue/internal/define"
	source "github.com/Galdoba/ffstuff/app/aue/internal/files/sourcefile"
	"github.com/Galdoba/ffstuff/pkg/ump"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type sourceCollector struct {
	sourceDir      string
	targetDir      string
	translationMap map[string]string
	renamingMap    map[string]string
	sources        []*source.SourceFile
}

func SetupSources(sourceDir, targetDir, translationFile string) ([]*source.SourceFile, error) {
	sc := sourceCollector{}
	sc.sourceDir = sourceDir
	sc.targetDir = targetDir
	sc.renamingMap = make(map[string]string)
	sc.translationMap = make(map[string]string)
	for _, err := range []error{
		sc.fillTranslationMap(translationFile),
		sc.assertProjectDirectory(),
		sc.collectFiles(),
		sc.executeRenaming(),
	} {
		if err != nil {
			return nil, fmt.Errorf("source setup failed: %v", err)
		}
	}

	return sc.sources, nil
}

func (sc *sourceCollector) assertProjectDirectory() error {
	f, err := os.Stat(sc.sourceDir)
	if err != nil {
		return fmt.Errorf("stat: %v")
	}
	if !f.IsDir() {
		return fmt.Errorf("'%v' is not a directory", sc.sourceDir)
	}
	return nil
}

func (sc *sourceCollector) collectFiles() error {
	fi, err := os.ReadDir(sc.sourceDir)
	if err != nil {
		return fmt.Errorf("can't read parent dir: %v", err)
	}

	base := filepath.Base(sc.sourceDir)
	translations := sugestTranslation(base, sc.translationMap)
	longer := ""
	for _, translation := range translations {
		if len(translation.translation) > len(longer) {
			longer = translation.translation
		}
	}

	base = baseToSource(base, longer)
	for _, f := range fi {
		file := path(sc.sourceDir, f.Name())
		expectedSourcePath := path(sc.targetDir, sourceNameProjected(base, f.Name()))
		prf := ump.NewProfile()
		if err := prf.ConsumeFile(file); err != nil {
			fmt.Println("LOG:", fmt.Errorf("profile: can't consume file '%v': %v", f.Name(), err))
			continue
		}
		strComp := streamComposition(prf)
		purpose := define.PURPOSE_Input_Media
		switch strComp {
		case 0:
			fmt.Println("LOG: skip", f.Name())
			continue
		case 1:
			purpose = define.PURPOSE_Input_Subs
		default:
			purpose = define.PURPOSE_Input_Media
		}

		newSource, err := newSourceWithProfile(expectedSourcePath, purpose, prf)
		if err != nil {
			return fmt.Errorf("collectFiles: %v", err)
		}
		sc.sources = append(sc.sources, newSource)
		fmt.Println("SOURCE ADDED:")
		fmt.Println(newSource.Details())

		sc.renamingMap[file] = expectedSourcePath
	}
	return nil
}

func newSourceWithProfile(expectedpath, purpose string, profile *ump.MediaProfile) (*source.SourceFile, error) {
	src := source.New(expectedpath, purpose)
	src.FillProfile(profile)
	return src, nil
}

func (sc *sourceCollector) executeRenaming() error {
	for source, destination := range sc.renamingMap {
		//fmt.Printf("rename %v to %v\n", source, destination)
		if err := os.Rename(source, destination); err != nil {
			return fmt.Errorf("renaming failed: %v")
		}
	}
	return nil
}

func path(dir, file string) string {
	dir = strings.TrimSuffix(dir, `\`)
	dir = strings.TrimSuffix(dir, "/")
	return fmt.Sprintf("%v%v%v", dir, "/", file)
}

func streamComposition(prf *ump.MediaProfile) int {
	cmp := 0
	for _, stream := range prf.Streams {
		switch stream.Codec_type {
		case define.STREAM_VIDEO:
			cmp += 100
		case define.STREAM_AUDIO:
			cmp += 10
		case define.STREAM_SUBTITLE:
			cmp += 1
		}
	}
	return cmp
}

func sourceNameProjected(base, name string) string {
	return fmt.Sprintf("%v_%v", base, name)
}

type translations struct {
	origin      string
	translation string
}

func sliceMatch(sl1, sl2 []string) bool {
	if len(sl1) != len(sl2) {
		return false
	}
	for i := range sl1 {
		if sl1[i] != sl2[i] {
			return false
		}
	}
	return true
}

func sugestTranslation(base string, translationMap map[string]string) []translations {
	basewords := strings.Split(base, "_")
	phrase := []string{}
	transl := []translations{}
	for _, word := range basewords {
		phrase = append(phrase, word)
		origin := strings.Join(phrase, "_")
		for k, v := range translationMap {
			if sliceMatch(words(k), words(origin)) {
				transl = append(transl, translations{k, v})
			}
		}
	}
	return transl
}

func baseToSource(base, translation string) string {
	seNum := seNum(base)
	seNumConv := "_"
	prt := prtStr(base)
	if seNum != "" {
		seNumConv = "s" + seNumConverted(seNum)
	}

	//base = strings.ReplaceAll(base, seNum, seNumConverted(seNum))
	sourceBase := translation + strings.TrimSuffix(seNumConv, "_") + prt
	return sourceBase
}

func seNum(str string) string {
	re := regexp.MustCompile(`(_[0-9]{3,})`)
	return re.FindString(str)
}

func prtStr(str string) string {
	re := regexp.MustCompile(`(_PRT[0-9]{12,})`)
	return re.FindString(str)
}

func seNumConverted(seNum string) string {
	seNum = strings.TrimPrefix(seNum, "_")
	val, err := strconv.Atoi(seNum)
	if err != nil {
		return ""
	}
	s := numToStr(val / 100)
	e := numToStr(val % 100)
	return fmt.Sprintf("_s%ve%v", s, e)
}

func numToStr(n int) string {
	s := fmt.Sprintf("%v", n)
	for len(s) < 2 {
		s = "0" + s
	}
	return s
}

var SearchOrigin = 0
var SearchTranslation = 1

func encodeToUTF8(filename string) []string {
	// Read UTF-8 from a GBK encoded file.
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	r := transform.NewReader(f, charmap.Windows1251.NewDecoder())

	// Read converted UTF-8 from `r` as needed.
	// As an example we'll read line-by-line showing what was read:
	lines := []string{}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		//fmt.Println(sc.Text())
		lines = append(lines, sc.Text())
	}
	if err = sc.Err(); err != nil {
		log.Fatal(err)
	}

	if err = f.Close(); err != nil {
		log.Fatal(err)
	}
	return lines
}

func (sc *sourceCollector) fillTranslationMap(translationFilePath string) error {
	// exampleReadGBK(translationFilePath)

	// panic(0)
	lines := encodeToUTF8(translationFilePath)
	//lines := strings.Split(fileStr, "\n")
	scanStatus := SearchOrigin
	title := ""
	translated := ""
	for _, line := range lines {
		switch scanStatus {
		case SearchOrigin:
			if strings.Contains(line, "Series ") {
				continue
			}
			if strings.Contains(line, `<title type="original">`) {
				title = getTitle(line)
				// title = strings.ReplaceAll(title, " ", "_")
				// title = strings.ReplaceAll(title, "'", "_")
				// title = strings.ReplaceAll(title, ",", "")
				// title = strings.ReplaceAll(title, ":", "")
				scanStatus = SearchTranslation
			}
		case SearchTranslation:
			if strings.Contains(line, `<title type="translated">`) && !strings.Contains(line, "Серия ") && !strings.Contains(line, "Эпизод ") {
				translated = getTranslation(line)
				sc.translationMap[title] = translated
				scanStatus = SearchOrigin

				title = ""
				translated = ""

				//fmt.Printf("detected: %v ==> %v\n", title, translated)
				// if strings.Contains(title, "Legend_") {
				// 	panic(2)
				// }
			}
		}

	}

	return nil
}

func getTitle(line string) string {
	//<title type="original">The Slide</title>
	title := strings.TrimPrefix(line, `<title type="original">`)
	title = strings.TrimSuffix(title, `</title>`)

	return strings.Join(words(title), "_")
}

func getTranslation(line string) string {
	//<title type="translated">Горка</title>
	translation := strings.TrimPrefix(line, `<title type="translated">`)
	translation = strings.TrimSuffix(translation, `</title>`)
	translation = translit(translation)
	return translation
}

func translit(origin string) string {
	letters := strings.Split(origin, "")
	changed := ""
	result := ""
	for _, l := range letters {
		changed += change(l)
	}
	words := strings.Split(changed, "_")
	for _, w := range words {
		if w == "" {
			continue
		}
		result += w + "_"
	}
	result = strings.TrimSuffix(result, "_")
	out := ""
	for i, l := range strings.Split(result, "") {
		switch i {
		case 0:
			out += strings.ToUpper(l)
		default:
			out += l
		}
	}
	return out
}

func words(text string) []string {
	lowText := ""
	for _, letter := range strings.Split(text, "") {
		lowText += change(letter)
	}
	lowText = strings.ReplaceAll(lowText, "_", " ")
	return strings.Fields(lowText)
}

func change(a string) string {
	a = strings.ToLower(a)
	switch a {
	default:
		return "_"
	case "а", "б", "в", "г", "д", "е", "ё", "ж", "з", "и", "й", "к", "л", "м", "н", "о", "п", "р", "с", "т", "у", "ф", "х", "ц", "ч", "ш", "щ", "ъ", "ы", "ь", "э", "ю", "я":
		lMap := make(map[string]string)
		lMap = map[string]string{
			"а": "a",
			"б": "b",
			"в": "v",
			"г": "g",
			"д": "d",
			"е": "e",
			"ё": "e",
			"ж": "zh",
			"з": "z",
			"и": "i",
			"й": "y",
			"к": "k",
			"л": "l",
			"м": "m",
			"н": "n",
			"о": "o",
			"п": "p",
			"р": "r",
			"с": "s",
			"т": "t",
			"у": "u",
			"ф": "f",
			"х": "h",
			"ц": "c",
			"ч": "ch",
			"ш": "sh",
			"щ": "sh",
			"ъ": "",
			"ы": "y",
			"ь": "",
			"э": "e",
			"ю": "yu",
			"я": "ya"}
		return lMap[a]
	case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", `/`, `\`:
		return a
	}
}
