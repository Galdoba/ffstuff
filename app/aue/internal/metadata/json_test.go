package metadata

import (
	"fmt"
	"testing"
)

func Test_json(t *testing.T) {

	trMap, err := TranslationsMap()
	fmt.Println(err)
	//	for k, v := range trMap {
	//		fmt.Println(k, "--", words(k), "---", v)
	//	}
	fmt.Println(len(trMap))
}
