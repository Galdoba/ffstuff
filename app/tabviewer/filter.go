package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	funcMATCH    = "FUNC_MATCH"
	funcCONTAINS = "FUNC_CONTAINS"
	funcREG      = "FUNC_REG"
)

type filter struct {
	RowWiseScope bool
	FilterFunc   string
	RHD          []string
}

func (f *filter) run(cnt *content, target string) (bool, error) {
	pass := false
	cellCoords := []string{}
	r, c := toNum(target)
	if f.RowWiseScope {
		for i := 0; i < c; i++ {
			trg := coord(r, i)
			cellCoords = append(cellCoords, trg.String())
		}
	}
	if len(cellCoords) == 0 {
		trg := coord(r, c)
		cellCoords = append(cellCoords, trg.String())
	}
	for _, crd := range cellCoords {
		if cell, ok := cnt.cells[crd]; ok {
			switch f.FilterFunc {
			case funcMATCH:
				for _, check := range f.RHD {
					passLoc, err := filterFunc_Match(cell.rawText, check)
					if err != nil {
						return false, err
					}
					if passLoc {
						fmt.Println("PASS", cell.rawText, check)
						return true, nil
					}
				}
			case funcCONTAINS:
				for _, check := range f.RHD {
					passLoc, err := filterFunc_Contains(cell.rawText, check)
					if err != nil {
						return false, err
					}
					if passLoc {
						return true, nil
					}
				}
			default:
				panic(999)
			}
		}

	}
	return pass, nil
}

func toNum(coord string) (int, int) {
	dat := strings.TrimPrefix(coord, "R")
	data := strings.Split(dat, "C")
	fmt.Println(data)
	i1, _ := strconv.Atoi(data[0])
	i2, _ := strconv.Atoi(data[1])
	return i1, i2
}

/*
match
contains
reg
*/

func filterFunc_Match(lhd string, rhd string) (bool, error) {
	fmt.Println(lhd == rhd, lhd, rhd)
	return lhd == rhd, nil
}

func filterFunc_Contains(lhd string, rhd string) (bool, error) {
	return strings.Contains(lhd, rhd), nil
}

func filterFunc_Regular(lhd string, maskKey string) (bool, error) {
	mask := maskKey
	rg, err := regexp.Compile(mask)
	if err != nil {
		return true, err
	}
	found := rg.FindString(lhd)
	if lhd == "" && lhd == found {
		return true, nil
	}
	if lhd != "" && found == "" {
		return false, nil
	}
	return true, nil
}
