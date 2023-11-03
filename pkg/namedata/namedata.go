package namedata

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/translit"
	"github.com/Galdoba/utils"
)

const (
	RESOLUTION_MP4    = "mp4"
	RESOLUTION_M4A    = "m4a"
	RESOLUTION_AAC    = "aac"
	RESOLUTION_SRT    = "srt"
	TAG_HD            = "_HD"
	TAG_SD            = "_SD"
	TAG_4K            = "_4K"
	TAG_SUB           = "_SUB"
	TAG_AUDIORUS20    = "_AUDIORUS20"
	TAG_AUDIORUS51    = "_AUDIORUS51"
	TAG_AUDIOENG20    = "_AUDIOENG20"
	TAG_AUDIOENG51    = "_AUDIOENG51"
	CONTENT_TYPE_FILM = "FILM"
	CONTENT_TYPE_TRL  = "TRL"
	CONTENT_TYPE_SER  = "SER"
)

/*
from Name:
-basename
-extention
-tags


*/

//RetrieveAll -
func RetrieveAll(path string) (string, string, []string) {
	fileName := shortFileName(path)
	tags, ext := splitName(fileName)
	base, tags2 := nameBase(tags)
	return base, ext, tags2
}

func RetrieveDirectory(path string) string {
	sep := string(filepath.Separator)
	pathData := strings.Split(path, sep)
	return strings.Join(pathData[0:len(pathData)-1], sep) + sep

}

func RetrieveShortName(path string) string {
	return shortFileName(path)
}

// func RetrieveBase(path string) string {
// 	fileName := shortFileName(path)
// 	if strings.Contains(fileName, "__") {
// 		return strings.Split(fileName, "__")[0]
// 	}
// 	tags, _ := splitName(fileName)
// 	base, _ := nameBase(tags)
// 	return base
// }

func RetrieveExtention(path string) string {
	p := strings.Split(path, ".")
	return p[len(p)-1]
}

func RetrieveTags(path string) []string {
	fileName := shortFileName(path)
	tags, _ := splitName(fileName)
	_, tags2 := nameBase(tags)
	return tags2
}

func RetrieveDrive(path string) string {
	data := strings.Split(path, "\\")
	return data[0]
}

func shortFileName(path string) string {
	path = strings.ReplaceAll(path, `\`, `/`)
	data := strings.Split(path, "/")
	fileName := path
	if len(data) > 1 {
		fileName = data[len(data)-1]
	}
	return fileName
}

func splitName(fileName string) ([]string, string) {
	data := strings.Split(fileName, "_")
	tags := []string{}
	ext := ""
	for index, val := range data {
		if index == len(data)-1 {
			p := strings.Split(val, ".")
			ext = p[len(p)-1]
			tags = append(tags, p[0])
			continue
		}
		tags = append(tags, val)
	}
	return tags, ext
}

func nameBase(tags []string) (string, []string) {
	base := ""
	tags2 := []string{}
	for i, val := range tags {
		for _, tg := range KnownTags() {
			if tg == val {
				tags2 = append(tags2, tg)
			}
			if tg != val {
				continue
			}
		}
		if i != len(tags)-1 {
			base += val + "_"
		}
	}
	base = strings.TrimSuffix(base, "_")
	return base, tags2
}

func KnownTags() []string {
	return []string{
		"HD",
		"SD",
		"43",
		"AUDIOENG20",
		"AUDIORUS20",
		"AUDIOENG51",
		"AUDIORUS51",
		"TRL",
		"Proxy",
		"ar2",
		"ar6",
		"ar2e2",
		"ar2e6",
		"ar6e2",
		"ar6e6",
		"rus51",
		"rus20",
		"eng51",
		"eng20",
	}
}

func TrimLoudnormPrefix(name string) (string, error) {
	newName := ""
	if err := validateOldname(name); err != nil {
		return newName, err
	}
	base, vid, aud, ebur := ungroupeName(name)
	if aud == "" {
		return newName, fmt.Errorf("audio tag can't be detected '%v'", name)
	}
	if ebur == "" {
		return name, nil
	}
	//fmt.Println("UNGROUPE:", base, vid, aud, ebur)
	switch {
	default:
		newName = ""
		return newName, fmt.Errorf("New name undecided for '%v'", name)

	case (vid == "hd" || vid == "4k") && strings.Contains(aud, "51") && strings.Contains(ebur, "-stereo"):
		vid = "sd"
		aud = strings.TrimSuffix(aud, "51") + "20"
		newName = base + "__" + vid + "_" + aud + ".ac3"
	case (vid == "hd" || vid == "4k") && strings.Contains(aud, "51") && !strings.Contains(ebur, "-stereo"):
		newName = base + "__" + vid + "_" + aud + ".ac3"
	case vid == "sd" && strings.Contains(aud, "51"):
		aud = strings.TrimSuffix(aud, "51") + "20"
		newName = base + "__" + vid + "_" + aud + ".ac3"
	case strings.Contains(aud, "20"):
		newName = base + "__" + vid + "_" + aud + ".ac3"
	}
	return newName, nil
}

func validateOldname(name string) error {
	if strings.TrimSuffix(name, ".ac3") == name {
		return fmt.Errorf("invalid name [%v] - is not ac3 file", name)
	}
	data := strings.Split(name, "__")
	if len(data) != 2 {
		return fmt.Errorf("invalid name [%v] - does not contain '__'", name)
	}
	if len(strings.Split(data[1], "_")) != 2 {
		return fmt.Errorf("invalid name [%v] - can not define audio and/or video tags", name)
	}
	/////HD20

	return nil
}

func ungroupeName(name string) (base, video, audio, ebur string) {
	data := strings.Split(name, "__")
	base = data[0]
	data2 := strings.Split(data[1], "_")
	video = data2[0]
	if strings.Contains(data2[1], "-ebur128-stereo.ac3") {
		audio = strings.TrimSuffix(data2[1], "-ebur128-stereo.ac3")
		ebur = "-ebur128-stereo"
		return
	}
	if strings.Contains(data2[1], "-ebur128.ac3") {
		audio = strings.TrimSuffix(data2[1], "-ebur128.ac3")
		ebur = "-ebur128"
		return
	}
	if strings.Contains(data2[1], "51.ac3") || strings.Contains(data2[1], "20.ac3") {
		audio = data2[1]
	}
	return
}

func renamerMapLocation() string {
	cu, _ := user.Current()
	fmt.Println("Call:", "c:\\Users\\"+cu.Name+"\\config\\ffstuff\\renamerMap.txt")
	return "c:\\Users\\" + cu.Name + "\\config\\ffstuff\\renamerMap.txt"
}

func AddToRenamerMap(oldName, newName string) error {
	f, err := os.OpenFile(renamerMapLocation(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.WriteString(oldName + " ==> " + newName); err != nil {
		return err
	}
	return nil
}

func RenamerMap() (map[string]string, error) {
	rnMap := make(map[string]string)
	lines := utils.LinesFromTXT(renamerMapLocation())
	for _, ln := range lines {
		data := strings.Split(ln, " ==> ")
		if len(data) != 2 {
			return rnMap, fmt.Errorf("Invalid entry '%v'", ln)
		}
		fmt.Println("FOUND:", rnMap[data[0]], "|", data[1], "||", data[0])
		rnMap[data[0]] = data[1]
	}
	return rnMap, nil
}

type textMask struct {
	//original     string
	matchPattern string
	typePattern  string
}

func (m *textMask) MatchPattern() string {
	return m.matchPattern
}

func (m *textMask) TypePattern() string {
	return m.typePattern
}

func (m *textMask) BasePatern() string {
	basePatern := ""
	base := true
	for i, l := range strings.Split(m.matchPattern, "") {
		if base {
			if l == "*" {
				base = false
			}
			switch base {
			case true:
				basePatern += l
			case false:
				typePat := strings.Split(m.typePattern, "")
				switch typePat[i] {
				case "d":
					basePatern += "+"
				case "_":
					basePatern += "_"
				default:
					return basePatern
				}
			}

		}
	}
	return basePatern
}

/*получая список имен пытаемся вывести маску из списка
Легенда выхода:
L - Letter
D - Digit
_ - space or special symbol
+ - Any symbol expected
- - Any symbol might or might not be
*/
//
func SearchMask(names []string) (textMask, error) {
	tm := textMask{}
	longest := 0
	for _, name := range names {
		l := len(strings.Split(name, ""))
		if longest <= l {
			longest = l
		}
	}
	for n := 0; n <= longest; n++ {
		symb := ""
		typeMap := make(map[string]int)
		sMap := make(map[string]int)
		for _, name := range names {
			if n > len(strings.Split(name, ""))-1 {
				symb = " "
				sMap[symb]++
				typeMap[" "]++
				continue
			} else {
				sl := strings.Split(name, "")
				symb = sl[n]
				sMap[symb]++
				symbType := glyphType(symb)
				typeMap[symbType]++

			}
		}
		switch {
		default:
			tm.matchPattern += "*"
		case len(sMap) == 1:
			tm.matchPattern += symb
		}
		switch {
		default:
			if n < longest {
				tm.typePattern += "*"
			}
		case summOfMap(typeMap) == typeMap[`d`]:
			tm.typePattern += "d"
		case summOfMap(typeMap) == typeMap[`w`]:
			tm.typePattern += "w"
		case summOfMap(typeMap) == typeMap[`_`]:
			tm.typePattern += "_"
		}
	}
	return tm, nil
}

func glyphType(s string) string {
	switch s {
	default:
		return "*"
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		return `d`
	case "A", "a", "B", "b", "C", "c", "D", "d", "E", "e", "F", "f", "G", "g", "H", "h", "I", "i", "J", "j", "K", "k", "L", "l", "M", "m", "N", "n", "O", "o", "P", "p", "Q", "q", "R", "r", "S", "s", "T", "t", "U", "u", "V", "v", "W", "w", "X", "x", "Y", "y", "Z", "z", "А", "а", "Б", "б", "В", "в", "Г", "г", "Д", "д", "Е", "е", "Ё", "ё", "Ж", "ж", "З", "з", "И", "и", "Й", "й", "К", "к", "Л", "л", "М", "м", "Н", "н", "О", "о", "П", "п", "Р", "р", "С", "с", "Т", "т", "У", "у", "Ф", "ф", "Х", "х", "Ц", "ц", "Ч", "ч", "Ш", "ш", "Щ", "щ", "Ъ", "ъ", "Ы", "ы", "Ь", "ь", "Э", "э", "Ю", "ю", "Я", "я":
		return `w`
	case `!`, `@`, `#`, `$`, `%`, `^`, `&`, `*`, `(`, `)`, `_`, `+`, ` `, `"`, `№`, `;`, `:`, `?`, `-`, `=`, `/`, `\`, `|`, `,`, `.`, ``:
		return `_`
	}
}

func summOfMap(m map[string]int) int {
	s := 0
	for _, v := range m {
		s += v
	}
	return s
}

type NameForm struct {
	name       string
	season     int
	episode    int
	seTag      string
	year       string
	tag        string
	resolution string
}

func ParseName(file string) *NameForm {
	n := NameForm{}
	n.resolution = detectResolution(file)
	n.tag = detectTag(file)
	n.season, n.episode, n.seTag = detectSeasonEpisode(file)
	nm := strings.TrimSuffix(file, "."+n.resolution)
	nm = strings.TrimSuffix(nm, n.tag)
	nms := strings.Split(nm, n.seTag)
	nm = nms[0]
	n.name = strings.ToLower(nm)
	return &n
}

func detectResolution(file string) string {
	file = reverse(file)
	data := strings.Split(file, ".")
	//validResolutions := []string{RESOLUTION_MP4, RESOLUTION_M4A, RESOLUTION_AAC}
	for _, resol := range validResolutions() {
		if reverse(resol) == data[0] {
			return resol
		}
	}
	return ""
}

func detectTag(file string) string {
	switch {
	case strings.Contains(file, TAG_HD):
		return TAG_HD
	case strings.Contains(file, TAG_SD):
		return TAG_SD
	case strings.Contains(file, TAG_4K):
		return TAG_4K
	case strings.Contains(file, TAG_SUB):
		return TAG_SUB
	case strings.Contains(file, TAG_AUDIORUS20):
		return TAG_AUDIORUS20
	case strings.Contains(file, TAG_AUDIORUS51):
		return TAG_AUDIORUS51
	case strings.Contains(file, TAG_AUDIOENG20):
		return TAG_AUDIOENG20
	case strings.Contains(file, TAG_AUDIOENG51):
		return TAG_AUDIOENG51
	}
	return "_TAG_ERROR"
}

func validYears() []string {
	tags := []string{}
	for i := 1900; i < 2025; i++ {
		yearTag := fmt.Sprintf("_%v__", i)
		yearTag = reverse(yearTag)
		tags = append(tags, yearTag)
	}
	return tags
}

func detectSeasonEpisode(file string) (int, int, string) {
	for s := 0; s < 30; s++ {
		for e := 0; e < 100; e++ {
			sT, eT := strconv.Itoa(s), strconv.Itoa(e)
			if s < 10 {
				sT = "0" + sT
			}
			if e < 10 {
				eT = "0" + eT
			}
			seTag := "_s" + sT + "e" + eT
			if strings.Contains(file, seTag) {
				return s, e, seTag
			}
		}
	}
	return -1, -1, ""
}

func validResolutions() []string {
	return []string{
		RESOLUTION_MP4,
		RESOLUTION_M4A,
		RESOLUTION_AAC,
		RESOLUTION_SRT,
	}
}

func reverse(str string) string {
	r := []rune(str)
	var res []rune
	for i := len(r) - 1; i >= 0; i-- {
		res = append(res, r[i])
	}
	return string(res)
}

func (nf *NameForm) ReconstructName() (string, error) {
	rnMap, err := RenamerMap()
	if err != nil {
		return "", err
	}
	newName := nf.name
	fmt.Println("oldName:", nf.name)
	if rnMap[nf.name] != "" {
		newName = rnMap[nf.name]
	}
	if nf.season != -1 {
		newName += "_s" + intToIndex(nf.season, 2)
	}
	if nf.episode != -1 {
		newName += "_" + intToIndex(nf.episode, 2)
	}
	newName += "_0000_"
	switch nf.tag {
	default:
		return "", fmt.Errorf("unknown tag '%v'", nf.tag)
	case TAG_HD:
		newName += "_hd" + "." + nf.resolution
	case TAG_SD:
		newName += "_sd" + "." + nf.resolution
	case TAG_4K:
		newName += "_4k" + "." + nf.resolution
	case TAG_SUB:
		newName += "_hd" + "." + nf.resolution
	case TAG_AUDIORUS20:
		newName += "_hd_rus20" + "." + nf.resolution
	case TAG_AUDIORUS51:
		newName += "_hd_rus51" + "." + nf.resolution
	case TAG_AUDIOENG20:
		newName += "_hd_eng20" + "." + nf.resolution
	case TAG_AUDIOENG51:
		newName += "_hd_eng51" + "." + nf.resolution
	}
	return newName, nil
}

func intToIndex(i, f int) string {
	index := strconv.Itoa(i)
	for len(index) < f {
		index = "0" + index
	}
	return index

}

func TransliterateForEdit(name string) string {
	clName := strings.Split(name, " (")[0]
	trName := translit.Transliterate(clName)
	nTag := ""
	sTag := ""
	for s := 0; s < 100; s++ {
		val := "_" + zeroIf(s) + fmt.Sprintf("%v", s) + "_sezon"
		if strings.Contains(trName, val) {
			nTag = strings.Split(trName, val)[0]
			sTag = fmt.Sprintf("_s%v%v", zeroIf(s), s)
			break
		}
	}
	eTag := ""
	for e := 0; e < 100; e++ {
		val := "_" + zeroIf(e) + fmt.Sprintf("%v", e) + "_seriya"
		if strings.Contains(trName, val) {
			eTag = fmt.Sprintf("e%v%v", zeroIf(e), e)
			break
		}
	}
	sTag = sTag + eTag
	res := nTag + sTag
	if res == "" {
		res = trName
	}
	return strings.Title(res)
}

func zeroIf(i int) string {
	if i < 10 {
		return "0"
	}
	return ""
}

func ValidateName(name string) string {
	letters := strings.Split(name, "")
	newName := ""
	for _, l := range letters {
		switch l {
		case "_", "(", ")", "-":
			l = " "
		}
		newName += l
	}
	return strings.Join(strings.Fields(newName), "_")
}

const (
	EDIT = "EDIT"
)

func SortFileNames(list []string, sortType string) ([]string, error) {
	sorted := []string{}
	switch sortType {
	default:
		return list, fmt.Errorf("Неизвестный тип сортировки: %v", sortType)
	case EDIT:
		sorted = sortAsEdit(list)
	}
	return sorted, nil
}

func sortAsEdit(list []string) []string {
	sList := []string{}
	sort.Strings(list)
	//sList = list
	weightMap := make(map[int]string)
	maxW := 0
	minW := 2000000
	for i, fl := range list {
		fl = strings.ToLower(fl)
		weight := i * -1
		if strings.Contains(fl, strings.ToLower(".srt")) {
			weight += 1000000
		}
		if strings.Contains(fl, strings.ToLower("_proxy")) {
			weight += 100000
		}
		if strings.Contains(fl, strings.ToLower("_audio")) {
			weight += 10000
		}
		if strings.Contains(fl, strings.ToLower("_nocens")) {
			weight += 10000
		}
		if strings.Contains(fl, strings.ToLower("_sd")) {
			weight += 1000
		}
		if strings.Contains(fl, strings.ToLower("_hd")) {
			weight += 1000
		}
		if strings.Contains(fl, strings.ToLower("_4k")) {
			weight += 100
		}
		if weight > maxW {
			maxW = weight
		}
		if weight < minW {
			minW = weight
		}
		weightMap[weight] = list[i]
	}
	for i := maxW; i >= minW; i-- {
		if fl, ok := weightMap[i]; ok {
			sList = append(sList, fl)
		}
	}
	if len(sList) != len(list) {
		panic("sortAsEdit() failed: len(sList) != len(list)")
	}
	return sList
}

func IsEditName(name string) bool {
	name = shortFileName(name)
	if len(strings.Split(name, "__")) == 2 {
		return true
	}
	return false
}

var EditTags = []string{
	"SD",
	"HD",
	"4K",
	"NOCENS",
}

type EditNameForm struct {
	source     string
	dir        string
	short      string
	base       string
	season     string
	episode    string
	prt        string
	extention  string
	tags       []string
	readyToUse bool
	editName   string
	film       bool
	trl        bool
	serial     bool
}

func (enf *EditNameForm) Source() string {
	return enf.source
}

func (enf *EditNameForm) Base() string {
	return enf.base
}

func (enf *EditNameForm) ShortName() string {
	return enf.short
}

func serialData(name string) (string, string, string) {
	re := regexp.MustCompile(`_(s[0-9]{1,}_[0-9]{1,})`)
	tag := re.FindString(name)
	data := strings.Split(tag, "_")
	if tag != "" {
		re2 := regexp.MustCompile(`(PRT[0-9]{1,})`)
		data = append(data, re2.FindString(name))
	}
	for len(data) < 4 {
		data = append(data, "")
	}
	return data[1], data[2], data[3]
}

func (ed *EditNameForm) PRT() string {
	return ed.prt
}

func RetrieveAudioTag(name string) string {
	re := regexp.MustCompile(`AUDIO(.*[0-9]{2})`)
	return re.FindString(name)
}

func EditForm(path string) *EditNameForm {
	ef := EditNameForm{}
	ef.source = path
	ef.short = shortFileName(path)
	ef.dir = RetrieveDirectory(path)
	ef.season, ef.episode, ef.prt = serialData(path)
	ef.extention = RetrieveExtention(path)
	ef.base = RetrieveBase(ef.short)
	ef.tags = tags(ef)
	editMarkers := strings.Split(ef.short, "--")
	switch len(editMarkers) {
	case 3:
		ef.editName = editMarkers[0]
		switch editMarkers[1] {
		case CONTENT_TYPE_TRL:
			ef.trl = true
		case CONTENT_TYPE_FILM:
			ef.film = true
		case CONTENT_TYPE_SER:
			ef.serial = true
		}
	}
	return &ef
}

func (ef *EditNameForm) Words() []string {
	wds := []string{}
	wds = append(wds, strings.Split(ef.base, "_")...)
	if ef.serial {
		wds = append(wds, ef.season, ef.episode)
	}
	return wds
}

func (ef *EditNameForm) EditName() string {
	return ef.editName
}

func (ef *EditNameForm) Season() string {
	return ef.season
}
func (ef *EditNameForm) Episode() string {
	return ef.episode
}

func (ef *EditNameForm) IsFilm() bool {
	return ef.film
}
func (ef *EditNameForm) IsTrl() bool {
	return ef.trl
}
func (ef *EditNameForm) IsSerial() bool {
	return ef.serial
}

func (ef *EditNameForm) ContentType() string {
	if ef.trl {
		return CONTENT_TYPE_TRL
	}
	if ef.serial {
		return CONTENT_TYPE_SER
	}
	if ef.film {
		return CONTENT_TYPE_FILM
	}
	return ""
}

func (ef *EditNameForm) AddPrefix(prefix string) error {
	oldName := ef.dir + ef.short
	ef.short = prefix + ef.short
	newName := ef.dir + ef.short
	err := os.Rename(oldName, newName)
	if err != nil {
		return fmt.Errorf("can not rename \n%v\n%v\n  reason: %v", oldName, newName, err.Error())
	}
	ef.source = newName
	return nil
}

func (ef *EditNameForm) HasTags(tags ...string) bool {
	met := 0
	for _, check := range tags {
		for _, got := range ef.tags {

			if check == got {
				met++
			}
		}
	}
	if met == len(ef.tags)-1 && met == len(tags) {
		return true
	}
	return false
}

func (ef *EditNameForm) HasExtention(ext string) bool {
	if ef.extention == ext {
		return true
	}
	return false
}

func RetrieveBase(shortName string) string {
	words := strings.Split(shortName, "_")
	base := ""
	for _, w := range words {
		for _, pref := range []string{"SD", "HD", "4K", "SD.", "HD.", "4K.", "AUDIO", "TRL", "PROXY", "proxy"} {
			if strings.HasPrefix(w, pref) {
				base = strings.TrimSuffix(base, "_")
				return base
			}
		}
		base += w + "_"
	}
	base = strings.TrimSuffix(base, "_")
	if base == shortName {
		baseWE := strings.Split(base, ".")
		base = strings.Join(baseWE[:len(baseWE)-1], ".")
	}
	return base
}

func tags(ef EditNameForm) []string {
	short := ef.short
	short = strings.TrimSuffix(short, "."+ef.extention)
	short = strings.TrimPrefix(short, ef.base)
	tags := strings.Split(short, "_")
	for i, tag := range tags {
		if strings.HasPrefix(tag, "AUDIO") {
			tags[i] = "AUDIO"
		}
	}
	return tags
}

type EditFormList struct {
	list []EditNameForm
}

// func NewEditFormList(list []string) *EditFormList {
// 	efl := EditFormList{}

// }

func have6mono(names []string) bool {
	ending := make(map[string]int)
	for _, name := range names {
		name = strings.ToLower(name)
		dots := strings.Split(name, ".")
		ext := dots[len(dots)-1]
		ending[ext]++
	}

	for _, qty := range ending {
		if qty != 6 {
			continue
		}
		for _, n := range names {

			for _, suf := range []string{
				"L.wav", "01.wav", "L.aif", "01.aif",
				"R.wav", "02.wav", "R.aif", "02.aif",
				"C.wav", "03.wav", "C.aif", "03.aif",
				"LFE.wav", "Lfe.wav", "04.wav", "LFE.aif", "Lfe.aif", "04.aif",
				"LS.wav", "Ls.wav", "05.wav", "LS.aif", "Ls.aif", "05.aif",
				"RS.wav", "Rs.wav", "06.wav", "RS.aif", "Rs.aif", "06.aif",
			} {

				if strings.HasSuffix(n, suf) {
					qty--
				}

			}

		}
		if qty == 0 {
			return true
		}
	}
	return false
}

func NormalizeSoundNames(names []string) ([]string, error) {
	if !have6mono(names) {
		return names, fmt.Errorf("6 mono sounds not detected")
	}
	newnames := []string{}
	norm := make(map[string]string)
	norm["01.wav"] = "L.wav"
	norm["01.aif"] = "L.aif"
	norm["02.wav"] = "R.wav"
	norm["02.aif"] = "R.aif"
	norm["03.wav"] = "C.wav"
	norm["03.aif"] = "C.aif"
	norm["LFE.wav"] = "Lfe.wav"
	norm["04.wav"] = "Lfe.wav"
	norm["LFE.aif"] = "Lfe.aif"
	norm["04.aif"] = "Lfe.aif"
	norm["LS.wav"] = "Ls.wav"
	norm["05.wav"] = "Ls.wav"
	norm["LS.aif"] = "Ls.aif"
	norm["05.aif"] = "Ls.aif"
	norm["RS.wav"] = "Rs.wav"
	norm["06.wav"] = "Rs.wav"
	norm["RS.aif"] = "Rs.aif"
	norm["06.aif"] = "Rs.aif"
	for _, nam := range names {
		for k, v := range norm {
			if !strings.HasSuffix(nam, k) {
				continue
			}
			newName := strings.TrimSuffix(nam, k) + v

			err := os.Rename(nam, newName)
			if err != nil {
				return names, fmt.Errorf("os.Rename(): %v", err.Error())
			}
			newnames = append(newnames, newName)
		}
	}
	return newnames, nil
}

func Words(name string) []string {
	wds := []string{}
	parts := strings.Split(name, " (")
	switch len(parts) > 1 {
	case true:
		if strings.HasSuffix(parts[1], ")") {
			name = parts[0]
		}
	case false:
	}
	for _, delim := range []string{" ", ".", ",", ":", "-", "–", "?", "!", "	"} {
		name = strings.ReplaceAll(name, delim, "_")
	}
	wds = append(wds, strings.Split(name, "_")...)
	return wds
}

func numbers() []string {
	out := []string{}
	for i := 0; i < 1000; i++ {
		n := strconv.Itoa(i)
		if i < 10 {
			n = "0" + n
		}
		out = append(out, n)
	}
	return out
}

func NameBaseOnly(tableName string) []string {
	//fmt.Println(tableName)
	out := []string{}
	for _, tag := range []string{" SD", " 4K", "_3D"} {
		tableName = strings.ReplaceAll(tableName, tag, "")
	}
	if strings.Contains(tableName, "(") {
		prebracket := strings.Split(tableName, "(")
		postBracket := strings.Join(prebracket[1:], "(")
		tailBracket := strings.Split(postBracket[len(postBracket)-1:], ")")
		tableName = prebracket[0] + tailBracket[len(tailBracket)-1]
	}
	tableNameLow := strings.ToLower(tableName)
	trNameLow := strings.ToLower(translit.Transliterate(tableNameLow))
	data := strings.Split(trNameLow, "_sezon_")
	base := ""
	switch len(data) {
	case 1:
		base := data[0]
		out = append(out, base)
		//return out
	default:
		for _, n := range numbers() {
			if strings.HasSuffix(data[0], "_"+n) {
				base = strings.TrimSuffix(data[0], "_"+n) + "_s" + n
				break
			}
		}
		for _, serTag := range []string{"_seriya", "_serii", "_seriy"} {
			if !strings.Contains(data[1], serTag) {
				continue
			}
			switch serTag {
			case "_seriya":
				data[1] = strings.TrimSuffix(data[1], serTag)
				out = append(out, base+"_"+data[1])

			default:

				data[1] = strings.TrimSuffix(data[1], serTag)
				episodes := strings.Split(data[1], "_")
				//fmt.Println(episodes)
				max := -1
				min := 1
				for _, ep := range episodes {
					eNum, _ := strconv.Atoi(ep)
					if eNum == 0 {
						continue
					}
					if max < 0 {
						max = eNum
						continue
					}
					min = max
					max = eNum
				}
				for e := min; e <= max; e++ {
					out = append(out, base+"_"+num2Str(e))
				}
			}
		}
	}
	for i, o := range out {
		l := strings.Split(o, "")
		l[0] = strings.ToUpper(l[0])
		o = strings.Join(l, "")
		out[i] = o
	}
	return out
	// if strings.Contains(tableNameLow, " сезон") {
	// 	for _, n := range numbers() {
	// 		if strings.Contains(tableNameLow, n+" сезон") {
	// 			tableName = strings.ReplaceAll(tableNameLow, n+" сезон", "")
	// 			seasonData = n
	// 			break
	// 		}
	// 	}
	// }
	// if strings.Contains(tableName, " сери") {
	// 	for _, n := range numbers() {
	// 		if strings.Contains(tableName, n+" сезон") {
	// 			tableName = strings.ReplaceAll(tableName, n+" сезон", "")
	// 			seasonData = n
	// 			break
	// 		}
	// 	}
	// }

}

func num2Str(n int) string {
	s := fmt.Sprintf("%v", n)
	if n < 0 {
		return s
	}
	if n < 10 {
		s = "0" + s
	}
	return s
}
