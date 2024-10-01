package sourcesort

import (
	"fmt"
	"math"
	"os"
	"slices"

	"github.com/Galdoba/ffstuff/app/grabber/internal/origin"
)

type PriorityProvider interface {
	NamesPriority() map[string]int
	DirectoryPriority() map[string]int
}

func SortByPriority(keepGroups bool, feed ...origin.Origin) []origin.Origin {
	sortedByScore := make(map[int][]origin.Origin)
	scoreMax := math.MaxInt * -1
	scoreMin := math.MaxInt
	///get score
	for _, unit := range feed {
		sortedByScore[unit.Score()] = append(sortedByScore[unit.Score()], unit)
		scoreMin, scoreMax = updateMinMax(scoreMin, scoreMax, unit.Score())
	}
	sorted := []origin.Origin{}
	//sort
	for i := scoreMax; i >= scoreMin; i-- {
		if list, ok := sortedByScore[i]; ok {
			for _, src := range list {
				sorted = append(sorted, src)
			}
		}
	}
	if keepGroups {
		sorted = groupSorted(sorted)
	}
	return sorted
}

func groupSorted(sorted []origin.Origin) []origin.Origin {
	groups := []string{}
	for _, src := range sorted {
		groups = appendUnique(groups, src.Group())
	}
	groupSorted := []origin.Origin{}
	for _, grp := range groups {
		for _, src := range sorted {
			if src.Group() == grp {
				groupSorted = append(groupSorted, src)
			}
		}
	}
	return groupSorted
}

func updateMinMax(min, max, new int) (int, int) {
	if new > max {
		max = new
	}
	if new < min {
		min = new
	}
	return min, max
}

func appendUnique(sl []string, elem string) []string {
	for _, s := range sl {
		if s == elem {
			return sl
		}
	}
	return append(sl, elem)
}

func SortBySize(keepGroups bool, feed ...origin.Origin) []origin.Origin {
	allSizes := []int64{}
	sourceSizeMap := make(map[origin.Origin]int64)
	for _, unit := range feed {
		f, err := os.Stat(unit.Path())
		if err != nil {
			fmt.Println("stats:", err.Error())
			continue
		}
		size := f.Size()
		allSizes = append(allSizes, size)
		sourceSizeMap[unit] = size
	}
	slices.Sort(allSizes)
	allSizes = removeDuplicates(allSizes)
	sorted := []origin.Origin{}
	for _, sizeKey := range allSizes {
		for source, size := range sourceSizeMap {
			if size != sizeKey {
				continue
			}
			sorted = append(sorted, source)
		}
	}
	if keepGroups {
		sorted = groupSorted(sorted)
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
