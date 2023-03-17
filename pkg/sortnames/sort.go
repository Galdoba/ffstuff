package sortnames

func NoDuplicates(sl []string) []string {
	newSl := []string{}
	for _, val := range sl {
		if inSlice(newSl, val) {
			continue
		}
		newSl = append(newSl, val)
	}
	return newSl
}

func inSlice(sl []string, elem string) bool {
	for _, val := range sl {
		if elem == val {
			return true
		}
	}
	return false
}
