package namedata

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/Galdoba/utils"
	"github.com/macroblock/imed/pkg/translit"
)

const (
	RESOLUTION_MP4 = "mp4"
	RESOLUTION_M4A = "m4a"
	RESOLUTION_AAC = "aac"
	RESOLUTION_SRT = "srt"
	TAG_HD         = "_HD"
	TAG_SD         = "_SD"
	TAG_4K         = "_4K"
	TAG_SUB        = "_SUB"
	TAG_AUDIORUS20 = "_AUDIORUS20"
	TAG_AUDIORUS51 = "_AUDIORUS51"
	TAG_AUDIOENG20 = "_AUDIOENG20"
	TAG_AUDIOENG51 = "_AUDIOENG51"
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
	pathData := strings.Split(path, "\\")
	return strings.Join(pathData[0:len(pathData)-1], "\\") + "\\"
}

func RetrieveShortName(path string) string {
	return shortFileName(path)
}

func RetrieveBase(path string) string {
	fileName := shortFileName(path)
	if strings.Contains(fileName, "__") {
		return strings.Split(fileName, "__")[0]
	}
	tags, _ := splitName(fileName)
	base, _ := nameBase(tags)
	return base
}

func RetrieveExtention(path string) string {
	fileName := shortFileName(path)
	_, ext := splitName(fileName)
	return ext
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
	data := strings.Split(path, "\\")
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
	trName, _ := translit.Do(clName)
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
			eTag = fmt.Sprintf("_%v%v", zeroIf(e), e)
			break
		}
	}
	res := nTag + sTag + eTag
	return strings.Title(res)
}

func zeroIf(i int) string {
	if i < 10 {
		return "0"
	}
	return ""
}
