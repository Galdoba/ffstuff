package sorting

import "fmt"

type index struct {
	//source    string
	possition int
	selected  bool
}

type IndexList struct {
	index []index
}

func Check(il1, il2 IndexList) (bool, error) {
	if len(il1.index) != len(il2.index) {
		return false, fmt.Errorf("indexes do not match")
	}
	for i := range il1.index {
		if il1.index[i].possition != il2.index[i].possition {
			return false, nil
		}
	}
	return true, nil
}

func Import(sel []bool) *IndexList {
	il := IndexList{}
	for i, s := range sel {
		il.index = append(il.index, index{i, s})
	}
	return &il
}

func (il *IndexList) Add(pos int, sel bool) {
	il.index = append(il.index, index{pos, sel})
}

func (il *IndexList) Export() ([]int, []bool) {
	pos := []int{}
	sel := []bool{}
	for _, in := range il.index {
		pos = append(pos, in.possition)
		sel = append(sel, in.selected)
	}
	return pos, sel
}

func (il *IndexList) MoveTop() {
	sel := []int{}
	unsel := []int{}
	for i, ind := range il.index {
		switch ind.selected {
		case true:
			sel = append(sel, i)
		case false:
			unsel = append(unsel, i)
		}
	}
	newIL := []index{}
	for _, s := range sel {
		newIL = append(newIL, il.index[s])
	}
	for _, uns := range unsel {
		newIL = append(newIL, il.index[uns])
	}
	il.index = newIL
}

func (il *IndexList) MoveBottom() {
	sel := []int{}
	unsel := []int{}
	for i, ind := range il.index {
		switch ind.selected {
		case false:
			sel = append(sel, i)
		case true:
			unsel = append(unsel, i)
		}
	}
	newIL := []index{}
	for _, s := range sel {
		newIL = append(newIL, il.index[s])
	}
	for _, uns := range unsel {
		newIL = append(newIL, il.index[uns])
	}
	il.index = newIL
}

func (il *IndexList) MoveUp() {
	keep := true
	for i, ind := range il.index {
		if ind.selected && keep {
			continue
		}
		if keep && !ind.selected {
			keep = false
		}
		if i == 0 {
			continue
		}
		if ind.selected {
			il.index[i-1], il.index[i] = il.index[i], il.index[i-1]
		}
	}
}

func (il *IndexList) reverse() {

	for i, j := 0, len(il.index)-1; i < j; i, j = i+1, j-1 {
		il.index[i], il.index[j] = il.index[j], il.index[i]
	}

}

func (il *IndexList) MoveDown() {
	il.reverse()
	il.MoveUp()
	il.reverse()
}
