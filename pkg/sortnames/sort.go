package sortnames

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
