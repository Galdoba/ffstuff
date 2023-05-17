package sortnames

import (
	"fmt"
	"strings"
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

func SeacrhFileNameBases(list []string) []string {
	words := searchWords(list)
	startWords := []string{}
	for _, w := range words {
		fmt.Println(w)
		wLow := strings.ToLower(w)
		wTitle := strings.Title(wLow)
		if wTitle == w {
			startWords = append(startWords, w)
		}
	}
	fmt.Println(startWords)
	otherWords := otherWords(words, startWords)
	fmt.Println(otherWords)
	baseMap := make(map[string]int)
	bases := []string{}
	//words = append(words, "")
	for _, sw := range startWords {
		for _, second := range words {
			bases = append(bases, strings.TrimSuffix(sw+"_"+second, "_"))
		}

	}
	bases = OmitDuplicates(bases)

	for _, name := range list {
		for _, base := range bases {
			if strings.HasPrefix(name, base) {
				baseMap[base]++
			}
		}
	}

	fmt.Println(baseMap)
	return []string{}
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
