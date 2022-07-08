package dmcombine

import (
	"fmt"
	"strings"
)

func autoMap(input []string) map[string]string {
	tags := []string{".L.", ".R.", ".C.", ".LFE.", ".Ls.", ".Rs."}
	tagMap := make(map[string]string)
	for _, path := range input {
		for i, tag := range tags {
			if !strings.Contains(path, tag) {
				continue
			}

			tagMap[path] = fmt.Sprintf("%v:a:0#", i)
		}
	}
}
