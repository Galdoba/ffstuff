package sourcesort

import (
	"fmt"
	"math"
	"os"
	"slices"
	"strings"
)

func SortByPriority(priorityMap map[string]int, feed ...string) []string {
	sortedByScore := make(map[int][]string)
	scoreMax := math.MaxInt * -1
	scoreMin := math.MaxInt
	for _, unit := range feed {
		upUnit := strings.ToUpper(unit)
		scoreTotal := 0
		for key, score := range priorityMap {
			upKey := strings.ToUpper(key)
			if strings.Contains(upUnit, upKey) {
				scoreTotal += score
			}
		}
		sortedByScore[scoreTotal] = append(sortedByScore[scoreTotal], unit)
		if scoreTotal > scoreMax {
			scoreMax = scoreTotal
		}
		if scoreTotal < scoreMin {
			scoreMin = scoreTotal
		}
	}
	sorted := []string{}
	for i := scoreMax; i >= scoreMin; i-- {
		if i < 0 {
			break
		}
		if list, ok := sortedByScore[i]; ok {
			for _, path := range list {
				sorted = append(sorted, path)
			}
		}
	}
	return sorted
}

func SortBySize(feed ...string) []string {
	allSizes := []int64{}
	sourceSizeMap := make(map[string]int64)
	for _, path := range feed {
		f, err := os.Stat(path)
		if err != nil {
			fmt.Println("bad path", path)
			continue
		}
		size := f.Size()
		allSizes = append(allSizes, size)
		sourceSizeMap[f.Name()] = size
	}
	slices.Sort(allSizes)
	allSizes = removeDuplicates(allSizes)
	sorted := []string{}
	for _, sizeKey := range allSizes {
		for source, size := range sourceSizeMap {
			if size == sizeKey {
				sorted = append(sorted, source)
				break
			}
		}
	}

	return sorted
}

func removeDuplicates(sl []int64) []int64 {
	newSl := []int64{}
	for _, s1 := range sl {
		if int64ContainedInSl(newSl, s1) {
			continue
		}
		newSl = append(newSl, s1)
	}
	return newSl
}

func int64ContainedInSl(sl []int64, el int64) bool {
	for _, n := range sl {
		if n == el {
			return true
		}
	}
	return false
}
