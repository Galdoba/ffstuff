package actions

import (
	"fmt"
	"strings"
	"testing"
)

var translationFilePathTest = `//192.168.31.4/buffer/IN/@AMEDIA_IN/amedia_tv_series.xml`

func testInputs() []string {
	return []string{
		"Don_t_Give_Up_109_PRT240830000000",
		"Mythical_Creatures_Are_My_Dinners_113_PRT240830000123",
	}
}

func TestTranslations(t *testing.T) {

	translations, err := collectTranslationVariants(translationFilePathTest)
	if err != nil {
		t.Errorf(err.Error())
	}

	for _, inp := range testInputs() {
		inputClean := inputBaseCleaned(inp)
		fmt.Println("INPUT:", inputClean)
		for i, trans := range translations {

			if strings.HasPrefix(inputClean, trans.expectedDirPrefix) {
				fmt.Println(inputClean, i, "====>", trans.renameTarget)
				continue
			}
		}
		t.Errorf("no translation found: %v", inputClean)
	}

}
