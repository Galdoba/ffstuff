package sorting

import (
	"fmt"
	"testing"
)

type testData struct {
	startPos    int
	expectedPos int
	startSel    bool
	expectedSel bool
}

func input() [][]testData {
	test := [][]testData{}
	test = append(test, []testData{
		{0, 1, false, true},
		{1, 2, true, true},
		{2, 0, true, false},
		{3, 3, false, false},
		{4, 4, false, false},
		{5, 5, false, false},
	})

	test = append(test, []testData{
		{0, 1, false, true},
		{1, 2, true, true},
		{2, 3, true, true},
		{3, 0, true, false},
		{4, 5, false, true},
		{5, 4, true, false},
	})

	test = append(test, []testData{
		{0, 0, true, true},
		{1, 1, true, true},
		{2, 2, true, true},
		{3, 3, true, true},
		{4, 5, false, true},
		{5, 4, true, false},
	})

	test = append(test, []testData{
		{0, 0, true, true},
		{1, 1, true, true},
		{2, 3, false, true},
		{3, 4, true, true},
		{4, 5, true, true},
		{5, 2, true, false},
	})

	test = append(test, []testData{
		{0, 1, false, true},
		{1, 0, true, false},
		{2, 3, false, true},
		{3, 2, true, false},
		{4, 5, false, true},
		{5, 4, true, false},
	})

	test = append(test, []testData{
		{0, 0, false, false},
		{1, 2, false, true},
		{2, 3, true, true},
		{3, 1, true, false},
		{4, 5, false, true},
		{5, 4, true, false},
	})

	return test
}

func TestMoveUp(t *testing.T) {

	for i, test := range input() {
		ilStartList := []index{}
		ilExpectedList := []index{}
		fmt.Println("Test", i+1)
		for _, td := range test {
			ilStartList = append(ilStartList, index{td.startPos, td.startSel})
			ilExpectedList = append(ilExpectedList, index{td.expectedPos, td.expectedSel})
		}
		il := IndexList{}
		il.index = ilStartList
		il.MoveUp()
		result := il.index
		for l, td := range test {
			verdict := "PASS"
			if result[l].possition != td.expectedPos || result[l].selected != td.expectedSel {
				verdict = "FAIL"
				t.Errorf("test %v: line %v: FAIL\nhave position %v (expect %v) - selection %v (expect %v)", i, l, result[l].possition, ilExpectedList[l].possition, result[l].selected, ilExpectedList[l].selected)
			}
			fmt.Println(ilExpectedList[l], result[l], verdict, td)
		}

	}

}
