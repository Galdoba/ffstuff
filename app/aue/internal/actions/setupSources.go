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
	"github.com/Galdoba/ffstuff/app/aue/logger"
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
	renamingOptions []nameTranslationOption
	sources         []*source.SourceFile
}

type nameTranslationOption struct {
	rusName           string
	engName           string
	expectedDirPrefix string
	renameTarget      string
}

/*

2024/08/30 20:16:17.428 [WARN ] >> project //192.168.31.4/buffer/IN/@AMEDIA_IN/Industry_303_PRT240830124300: no sources created
entering dormant mode:
leave dormant mode
2024/08/30 20:21:20.255 [INFO ] >> start project:%!(EXTRA string=//192.168.31.4/buffer/IN/@AMEDIA_IN/Industry_303_PRT240830124300)
SOURCE ADDED:
  path    : //192.168.31.4/buffer/IN/Industry_s03e03_PRT240830124300_SER_05052_18.RUS.srt
  name    : Industry_s03e03_PRT240830124300_SER_05052_18.RUS.srt
  purpose : Input_Subs
  profile : &{   0xc0013285a0 [0xc001368380] [file ==> 1 subtitle stream detected] map[] 0001-1 0;0;0;1;w1  []}

SOURCE ADDED:
  path    : //192.168.31.4/buffer/IN/Industry_s03e03_PRT240830124300_SER_05052_18.mp4
  name    : Industry_s03e03_PRT240830124300_SER_05052_18.mp4
  purpose : Input_Media
  profile : &{   0xc0014261b0 [0xc00005c380 0xc00005c700 0xc00005ca80 0xc00005ce00] [] map[0:a:0:#5.1#48#341 0:a:1:#stereo#48#130 0:v:0:#HD#25#[SAR=1:1_DAR=16
:9]#15983#ns] 1210-0 1v0{#HD#25#[SAR=1:1_DAR=16:9]#15983#ns};2a0{#5.1#48#341}a1{#stereo#48#130};1;0;w0 6 2 []}

LOG WARN: profile: can't consume file 'metadata.json': can't read: [tedcaptions @ 000001e49549fe00] Syntax error near offset 14.
//192.168.31.4/buffer/IN/@AMEDIA_IN/Industry_303_PRT240830124300/metadata.json: Invalid data found when processing input

2024/08/30 20:21:21.702 [INFO ] >> job creation complete: //192.168.31.4/buffer/IN/@AMEDIA_IN/Industry_303_PRT240830124300
creating targets for V1A2S1
LOG ERROR: parseComplex Failed
LOG ERROR: parseSimple Failed
LOG ERROR: parseComplex Failed
LOG ERROR: parseSimple Failed
LOG ERROR: job decide type source/target connection failed: sealing failed: conflicting data 'BASE': [Industry_s03e03_PRT240830124300_SER_05052_18] != [Indust
ry_s03e03_PRT240830124300_SER_05052_18RUS]
entering dormant mode:
leave dormant mode
2024/08/30 20:26:24.554 [INFO ] >> start project:%!(EXTRA string=//192.168.31.4/buffer/IN/@AMEDIA_IN/Industry_303_PRT240830124300)
LOG WARN: profile: can't consume file 'metadata.json': can't read: [tedcaptions @ 00000127aa51fe00] Syntax error near offset 14.
//192.168.31.4/buffer/IN/@AMEDIA_IN/Industry_303_PRT240830124300/metadata.json: Invalid data found when processing input

2024/08/30 20:26:25.516 [WARN ] >> project //192.168.31.4/buffer/IN/@AMEDIA_IN/Industry_303_PRT240830124300: no sources created
entering dormant mode:
leave dormant mode
2024/08/30 20:31:28.295 [INFO ] >> start project:%!(EXTRA string=//192.168.31.4/buffer/IN/@AMEDIA_IN/Industry_303_PRT240830124300)
LOG WARN: profile: can't consume file 'metadata.json': can't read: [tedcaptions @ 0000020874d8fe00] Syntax error near offset 14.
//192.168.31.4/buffer/IN/@AMEDIA_IN/Industry_303_PRT240830124300/metadata.json: Invalid data found when processing input

2024/08/30 20:31:29.328 [WARN ] >> project //192.168.31.4/buffer/IN/@AMEDIA_IN/Industry_303_PRT240830124300: no sources created
entering dormant mode:
^Cke up in 154 seconds
*/

func SetupSources(sourceDir, targetDir, translationFile string) ([]*source.SourceFile, error) {
	sc := sourceCollector{}
	sc.sourceDir = sourceDir
	sc.targetDir = targetDir
	sc.renamingMap = make(map[string]string)
	for _, err := range []error{
		sc.fillTranslationMap(translationFile),
		sc.assertProjectDirectory(),
		sc.projectPrefix(),
		sc.collectFiles(),
		sc.executeRenaming(),
	} {
		if err != nil {
			err = fmt.Errorf("source setup failed: %v", err)
			logger.Error(err)
			return nil, fmt.Errorf("source setup failed: %v", err)
		}
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
			logger.Info("skip %v", f.Name())
			continue
		}
		file := path(sc.sourceDir, f.Name())

		expectedSourcePath := path(sc.targetDir, sc.projectSourceFilepath(f.Name()))

		prf := ump.NewProfile()
		if err := prf.ConsumeFile(file); err != nil {
			logger.Warn("profile: can't consume file '%v': %v", f.Name(), err)
			//panic("can't consume")
			//fmt.Println("LOG WARN:", fmt.Errorf("profile: can't consume file '%v': %v", f.Name(), err))
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
		logger.Info("source added: %v", newSource.Name())
		sc.renamingMap[file] = expectedSourcePath
	}
	return nil
}

func newSourceWithProfile(expectedpath, purpose string, profile *ump.MediaProfile) (*source.SourceFile, error) {
	src := source.New(expectedpath, purpose)
	src.FillProfile(profile)
	return src, nil
}

func (sc *sourceCollector) projectPrefix() error {
	base := filepath.Base(sc.sourceDir)
	base = inputBaseCleaned(base)
	sc.base = base
	sc.prt = prtStr(sc.sourceDir)
	if sc.prt == "" {
		logger.Warn("no PRT detected")
	}
	sc.seNum = seNumConverted(seNum(sc.sourceDir))
	if sc.seNum == "" {
		logger.Warn("no serial season/episode number detected")
	}
	sc.exprctedPrefix = sc.base + sc.seNum + sc.prt
	return nil
}

func (sc *sourceCollector) executeRenaming() error {
	for source, destination := range sc.renamingMap {
		if err := os.Rename(source, destination); err != nil {
			return fmt.Errorf("renaming failed: %v", source)
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

func (sc *sourceCollector) projectSourceFilepath(name string) string {
	for _, renameOpt := range sc.renamingOptions {
		if sc.base == renameOpt.expectedDirPrefix {
			return fmt.Sprintf("%v_%v", renameOpt.renameTarget+sc.seNum+sc.prt, name)
		}
	}
	return fmt.Sprintf("%v_%v", sc.base+sc.seNum+sc.prt, name)
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

func (sc *sourceCollector) fillTranslationMap(translationFilePath string) error {
	renamingOptions, err := collectTranslationVariants(translationFilePath)
	if err != nil {
		fmt.Println("LOG ERROR:", err.Error())
	}
	sc.renamingOptions = renamingOptions
	return nil
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
			if name != nameGuessed {
				translations = append(translations, nameTranslationOption{
					rusName:           translationRus,
					engName:           name,
					expectedDirPrefix: nameGuessed,
					renameTarget:      translit(translationRus),
				})

			}
			print = false
		}

	}
	return translations, nil
}

func between(head, tail, text string) string {
	out := strings.TrimPrefix(text, head)
	out = strings.TrimSuffix(out, tail)
	return out
}

func guessFolder(name string) string {
	for _, glyph := range []string{":", "-", "_", ",", ".", "&", ";", "#", "'", `"`, "?", "/", `\`, "|", "(", ")", "’", "*", "!", "[", "]", "{", "}"} {
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
