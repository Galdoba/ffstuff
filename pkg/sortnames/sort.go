package sortnames

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/namedata"
)

func OmitDuplicates(sl []string) []string {
	newSl := []string{}
	for _, val := range sl {
		if inSlice(newSl, val) {
			continue
		}
		newSl = append(newSl, val)
	}
	return newSl
}

func ExpelDuplicates(sl []string, words ...string) []string {
	sl2 := []string{}
	for _, sample := range sl {
		met := false
		for _, word := range words {
			if sample == word && !met {
				met = true
				break
			}
		}
		if met {
			continue
		}
		sl2 = append(sl2, sample)
	}
	return sl2
}

func Prepend(elem string, sl []string) []string {
	slRes := []string{elem}
	slRes = append(slRes, sl...)
	return slRes
}

func inSlice(sl []string, elem string) bool {
	for _, val := range sl {
		if elem == val {
			return true
		}
	}
	return false
}

type Type struct {
	in   int
	info string
}

func BumpToTopIndex(slInt []int, index int) []int {
	if index < 1 || index > len(slInt)-1 {
		return slInt
	}
	newSl := []int{}
	for i := range slInt {
		switch {
		case i == 0:
			newSl = append(newSl, slInt[index])
		case i <= index:
			newSl = append(newSl, slInt[i-1])
		case i > index:
			newSl = append(newSl, slInt[i])
		}
	}
	return newSl
}

func BumpIndexUpByOne(slInt []int, index int) []int {
	if index < 1 || index > len(slInt)-1 {
		return slInt
	}
	slInt[index-1], slInt[index] = slInt[index], slInt[index-1]
	return slInt
}

func emulatedNames() []string {
	var emulatedNames []string
	base := []string{"Balet", "Buffalo", "Charlie_shin"}
	sTag := []string{"", "_s01", "_s02", "_s03"}
	epTag := []string{"", "_01", "_02", "_03"}
	revTag := []string{"", "_R1", "_R2"}
	vidTag := []string{"_sd", "_hd", "_4k", ""}
	audTag := []string{"_AUDIORUS20", "_AUDIORUS51", "_AUDIOENG20", "_AUDIOENG51", ""}
	proxyTag := []string{"", "_proxy"}
	extn := []string{".txt", ".m4a", ".srt", ".mp4", ".ready", ".ac3"}

	//	total := len(base) * len(sTag) * len(epTag) * len(revTag) * len(vidTag) * len(audTag) * len(proxyTag) * len(extn)
	for _, ba := range base {
		for _, sT := range sTag {
			for _, ep := range epTag {
				for _, re := range revTag {
					for _, vi := range vidTag {
						for _, au := range audTag {
							for _, pt := range proxyTag {
								for _, ex := range extn {
									if au != "" && vi != "" {
										continue
									}
									emulatedNames = append(emulatedNames, ba+sT+ep+re+vi+au+pt+ex)
								}
							}
						}
					}
				}
			}
		}
	}
	return emulatedNames
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

func otherWords(words, startWords []string) []string {
	ow := []string{}
	for _, word := range words {
		isOther := true
		for _, start := range startWords {
			if word == start {
				isOther = false
				break
			}
		}
		if isOther {
			ow = append(ow, word)
		}

	}
	return ow
}

func searchWords(phrases []string) []string {
	allWords := []string{}
	for _, phrase := range phrases {
		for _, s := range listSybols() {
			phrase = strings.ReplaceAll(phrase, s, "_")
		}
		words := strings.Split(phrase, "_")
		ssn, ep, prt := serialData(phrase)
		words = ExpelDuplicates(words, ssn, ep, prt)
		allWords = append(allWords, words...)
	}
	allWords = OmitDuplicates(allWords)
	return allWords
}

func listSybols() []string {
	return []string{`!`, `@`, `#`, `$`, `%`, `^`, `&`, `*`, `(`, `)`, `_`, `+`, ` `, `"`, `№`, `;`, `:`, `?`, `-`, `=`, `/`, `\`, `|`, `,`, `.`}
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

func names() []string {
	return []string{
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_01_PRT230516135638.srt`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_01_PRT230516135638_AUDIORUS51.m4a`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_01_PRT230516135638_AUDIORUS51_proxy.ac3`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_01_PRT230516135638_HD.mp4`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_01_PRT230516135638_HD_proxy.mp4`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_01_PRT230519162218.srt`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_01_PRT230519162218_AUDIOENG20.m4a`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_01_PRT230519162218_AUDIOENG20_proxy.ac3`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_01_PRT230519162218_AUDIORUS51.m4a`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_01_PRT230519162218_AUDIORUS51_proxy.ac3`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_01_PRT230519162218_HD.mp4`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_01_PRT230519162218_HD_proxy.mp4`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_02_PRT230516135330.srt`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_02_PRT230516135330_AUDIORUS51.m4a`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_02_PRT230516135330_AUDIORUS51_proxy.ac3`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_02_PRT230516135330_HD.mp4`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_02_PRT230516135330_HD_proxy.mp4`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_04_PRT230516134753.srt`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_04_PRT230516134753_AUDIORUS51.m4a`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_04_PRT230516134753_AUDIORUS51_proxy.ac3`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_04_PRT230516134753_HD.mp4`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_04_PRT230516134753_HD_proxy.mp4`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_04_PRT230519165811.srt`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_04_PRT230519165811_AUDIOENG20.m4a`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_04_PRT230519165811_AUDIOENG20_proxy.ac3`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_04_PRT230519165811_AUDIORUS51.m4a`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_04_PRT230519165811_AUDIORUS51_proxy.ac3`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_04_PRT230519165811_HD.mp4`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_04_PRT230519165811_HD_proxy.mp4`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_07_PRT230516135929.srt`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_07_PRT230516135929_AUDIORUS51.m4a`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_07_PRT230516135929_AUDIORUS51_proxy.ac3`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_07_PRT230516135929_HD.mp4`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_07_PRT230516135929_HD_proxy.mp4`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_07_PRT230519161740.srt`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_07_PRT230519161740_AUDIOENG20.m4a`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_07_PRT230519161740_AUDIOENG20_proxy.ac3`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_07_PRT230519161740_AUDIORUS51.m4a`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_07_PRT230519161740_AUDIORUS51_proxy.ac3`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_07_PRT230519161740_HD.mp4`,
		`\\nas\ROOT\EDIT\_amedia\Ne_ostavlyay_menya_s01\Ne_ostavlyay_menya_s01_07_PRT230519161740_HD_proxy.mp4`,
	}
}

func GrabberOrder(list []string) []string {
	edNam := []*namedata.EditNameForm{}
	for _, name := range list {
		edNam = append(edNam, namedata.EditForm(name))
	}
	baseMap := make(map[string]string)
	baseList := []string{}
	for _, edit := range edNam {
		base := edit.Base()
		baseList = append(baseList, base)
		baseMap[edit.Source()] = base
	}
	sort.Strings(baseList)
	baseList = OmitDuplicates(baseList)
	sorted := [][]string{}
	for _, base := range baseList {
		sources := []string{}
		for k, v := range baseMap {
			if v == base {
				sources = append(sources, k)
			}
		}
		sorted = append(sorted, sources)
	}
	fmt.Println(sorted)
	editOrd := editOrder(sorted)
	editOrd = append(editOrd, list...)
	editOrd = OmitDuplicates(editOrd)
	return editOrd
}

func editOrder(list [][]string) []string {
	grabbed := []string{}
	for _, base := range list {
		editOrderTags := editOrderTags()
		for _, tags := range editOrderTags {
			for _, inBase := range base {
				ef := namedata.EditForm(inBase)
				switch {
				case ef.HasExtention(tags[0]):
					grabbed = append(grabbed, ef.Source())
				case ef.HasTags(tags...):
					grabbed = append(grabbed, ef.Source())

				}
			}

		}

	}
	return grabbed
}

func editOrderTags() [][]string {
	return [][]string{
		{"srt"},
		{"AUDIO", "proxy"},
		{"SD", "proxy"},
		{"HD", "proxy"},
		{"4K", "proxy"},
		{"AUDIO"},
		{"SD"},
		{"HD"},
		{"4K"},
		{""},
	}
}
