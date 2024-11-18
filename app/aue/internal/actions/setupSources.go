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
	"github.com/Galdoba/ffstuff/app/aue/internal/metadata"
	"github.com/Galdoba/ffstuff/pkg/logman"
	logger "github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/Galdoba/ffstuff/pkg/ump"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type sourceCollector struct {
	sourceDir       string
	targetDir       string
	base            string
	seNum           string
	prt             string
	exprctedPrefix  string
	containedFiles  []string
	renamingMap     map[string]string
	projectName     string
	renamingOptions []nameTranslationOption
	sources         []*source.SourceFile
}

type nameTranslationOption struct {
	rusName           string
	engName           string
	expectedDirPrefix string
	renameTarget      string
}

func SetupSources(sourceDir, targetDir, translationFile string) ([]*source.SourceFile, error) {
	// serialFile := ``
	// serialInfo, err := xmlparse.ParseSeriesVideoData(translationFile)

	sc := sourceCollector{}
	sc.sourceDir = sourceDir
	sc.projectName = filepath.Base(sourceDir)
	sc.targetDir = targetDir
	sc.renamingMap = make(map[string]string)
	for _, err := range []error{
		//sc.fillTranslationMap(translationFile),
		sc.assertProjectDirectory(),
		sc.projectPrefix(),
		sc.collectFiles(),
	} {
		if err != nil {
			err = fmt.Errorf("project %v: source setup failed", sc.projectName)
			return nil, logger.Error(err)
		}
	}
	if err := sc.executeRenaming(); err != nil {
		return nil, logger.Error(err)
	}
	return sc.sources, nil
}

func (sc *sourceCollector) assertProjectDirectory() error {
	f, err := os.Stat(sc.sourceDir)
	if err != nil {
		return fmt.Errorf("stat: %v", err)
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
	for _, f := range fi {
		if f.Name() == "metadata.json" {
			logger.Debug(logger.NewMessage(fmt.Sprintf("skip %v", f.Name())))
			continue
		}
		file := path(sc.sourceDir, f.Name())
		expectedSourcePath := path(sc.targetDir, sc.exprctedPrefix+"_"+f.Name())
		prf := ump.NewProfile()
		if err := prf.ConsumeFile(file); err != nil {
			err := fmt.Errorf("failed to consume file '%v': %v", f.Name(), err)

			return logger.Error(err)
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
		logger.Info("source added: %v", f.Name())
		sc.renamingMap[file] = expectedSourcePath

	}
	if len(sc.sources) == 0 {
		return logger.Warn("project '%v': no sources detected", sc.projectName)
	}
	return nil
}

func newSourceWithProfile(expectedpath, purpose string, profile *ump.MediaProfile) (*source.SourceFile, error) {
	src := source.New(expectedpath, purpose)
	src.FillProfile(profile)
	return src, nil
}

// func (sc *sourceCollector) projectPrefix() error {
// 	base := filepath.Base(sc.sourceDir)
// 	base = inputBaseCleaned(base)
// 	sc.base = base
// 	sc.prt = prtStr(sc.sourceDir)
// 	if sc.prt == "" {
// 		logger.Warn("no PRT detected")
// 	}
// 	sc.seNum = seNumConverted(seNum(sc.sourceDir))
// 	if sc.seNum == "" {
// 		logger.Warn("no serial season/episode number detected")
// 	}
// 	sc.exprctedPrefix = sc.base + sc.seNum + sc.prt
// 	return nil
// }

func (sc *sourceCollector) projectPrefix() error {
	base := filepath.Base(sc.sourceDir)
	se := seNum(base)
	if se != "" {
		se = seNumConverted(se)
	}
	prt := prtStr(base)
	baseCleaned := inputBaseCleaned(base)
	sc.base = baseCleaned
	prefix, err := translate(baseCleaned)
	if err != nil {
		logman.Warn("translation failed: %v", baseCleaned)
	}
	sc.exprctedPrefix = prefix + se + prt

	return nil
}

func (sc *sourceCollector) executeRenaming() error {
	for source, destination := range sc.renamingMap {
		if err := os.Rename(source, destination); err != nil {
			return fmt.Errorf("renaming failed: %v", source)
		}
		//		fmt.Printf("os.Rename(%v, %v)", source, destination)
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

// func (sc *sourceCollector) projectSourceFilepath(name string) string {
// 	for _, renameOpt := range sc.renamingOptions {
// 		if sc.exprctedPrefix == renameOpt.expectedDirPrefix {
// 			fmt.Println("========", sc.exprctedPrefix)
// 			return fmt.Sprintf("%v_%v", renameOpt.renameTarget+sc.seNum+sc.prt, name)
// 		}
// 	}
// 	return fmt.Sprintf("%v_%v", sc.base+sc.seNum+sc.prt, name)
// }

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

// func baseToSource(base, translation string) string {
// 	seNum := seNum(base)
// 	seNumConv := "_"
// 	prt := prtStr(base)
// 	if seNum != "" {
// 		seNumConv = "s" + seNumConverted(seNum)
// 	}
// 	sourceBase := translation + strings.TrimSuffix(seNumConv, "_") + prt
// 	return sourceBase
// }

func seNum(str string) string {
	re := regexp.MustCompile(`(_[0-9]{3,})`)
	return re.FindString(str)
}

func prtStr(str string) string {
	re := regexp.MustCompile(`(_PRT[0-9]{12,})`)
	return re.FindString(str)
}

func seNumConverted(seNum string) string {
	if seNum == "" {
		return ""
	}
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
	lines := []string{}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
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

// func (sc *sourceCollector) fillTranslationMap(translationFilePath string) error {
// 	renamingOptions, err := collectTranslationVariants(translationFilePath)
// 	if err != nil {
// 		fmt.Println("LOG ERROR:", err.Error())
// 	}
// 	sc.renamingOptions = renamingOptions
// 	return nil
// }

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
	case "а", "б", "в", "г", "д", "е", "ё", "ж", "з", "и",
		"й", "к", "л", "м", "н", "о", "п", "р", "с", "т",
		"у", "ф", "х", "ц", "ч", "ш", "щ", "ъ", "ы", "ь",
		"э", "ю", "я":
		lMap := make(map[string]string)
		lMap = map[string]string{
			"а": "a", "б": "b", "в": "v", "г": "g", "д": "d",
			"е": "e", "ё": "e", "ж": "zh", "з": "z", "и": "i",
			"й": "y", "к": "k", "л": "l", "м": "m", "н": "n",
			"о": "o", "п": "p", "р": "r", "с": "s", "т": "t",
			"у": "u", "ф": "f", "х": "h", "ц": "c", "ч": "ch",
			"ш": "sh", "щ": "sh", "ъ": "", "ы": "y", "ь": "",
			"э": "e", "ю": "yu", "я": "ya"}
		return lMap[a]
	case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
		"k", "l", "m", "n", "o", "p", "q", "r", "s", "t",
		"u", "v", "w", "x", "y", "z",
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "0",
		`/`, `\`:
		return a
	}
}

func collectTranslationVariants(path string) ([]nameTranslationOption, error) {
	counter := 0
	lines := encodeToUTF8(path)
	print := false
	name := ""
	translationRus := ""
	depth := 0
	translations := []nameTranslationOption{}
	for _, line := range lines {
		if strings.HasPrefix(line, `<group type="`) {
			depth++
		}
		if strings.HasPrefix(line, `</group>`) {
			depth--
		}
		if strings.HasPrefix(line, `<title type="original">`) {
			if depth != 1 {
				continue
			}
			counter++
			name = between(`<title type="original">`, `</title>`, line)
		}
		if strings.HasPrefix(line, `<title type="translated">`) {
			if depth != 1 {
				continue
			}
			translationRus = between(`<title type="translated">`, `</title>`, line)
			print = true
		}
		if print {
			nameGuessed := guessFolder(name)

			//if name != nameGuessed {
			translations = append(translations, nameTranslationOption{
				rusName:           translationRus,
				engName:           name,
				expectedDirPrefix: nameGuessed,
				renameTarget:      translit(translationRus),
			})

			//}
			print = false
		}

	}
	return translations, nil
}

func translate(base string) (string, error) {
	trMap, err := metadata.TranslationsMap()
	if err != nil {
		return base, fmt.Errorf("failed to compose translation map: %v", err)
	}
	keyWords := words(strings.ToLower(inputBaseCleaned(base)))

	for key, value := range trMap {
		if equal(keyWords, words(key)) {
			return value, nil
		}
	}
	return base, fmt.Errorf("words not found: %v", keyWords)
}

func equal(sl1, sl2 []string) bool {
	if len(sl1) != len(sl2) {
		return false
	}
	for i, s := range sl1 {
		if s != sl2[i] {
			return false
		}
	}
	return true
}

func between(head, tail, text string) string {
	out := strings.TrimPrefix(text, head)
	out = strings.TrimSuffix(out, tail)
	return out
}

func guessFolder(name string) string {
	for _, glyph := range []string{":", "-", "_", ",", "&", ".", ";", "#", "'", `"`, "?", "/", `\`, "|", "(", ")", "’", "*", "!", "[", "]", "{", "}"} {
		name = strings.ReplaceAll(name, glyph, " ")
	}
	words := strings.Fields(name)
	name = strings.Join(words, "_")
	return name
}

func inputBaseCleaned(inputBase string) string {
	prt := prtStr(inputBase)
	inputBase = strings.TrimSuffix(inputBase, prt)
	seNum := seNum(inputBase)
	inputBase = strings.TrimSuffix(inputBase, seNum)
	return inputBase
}
