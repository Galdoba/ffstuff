package clipmaker

import "fmt"

func ConcatClips(cm map[int]clip) []string {
	tasks := []string{}
	lastClip := clip{}
	workClip := clip{}
	concatIndexes := make(map[int][]int)
	stIndex := 0
	for i := 0; i < 1000; i++ {
		if _, ok := cm[i]; !ok {
			continue
		}
		lastClip = workClip
		workClip = cm[i]

		if lastClip.seqPosEndTimeCode != workClip.seqPosStartTimeCode {
			stIndex = i
			concatIndexes[stIndex] = []int{i}
			continue
		}
		concatIndexes[stIndex] = append(concatIndexes[stIndex], i)
	}
	for i := 0; i < 1000; i++ {
		if len(concatIndexes[i]) != 0 {
			fmt.Println(i, concatIndexes[i])
		}
	}

	return tasks
}

/*
Клипы: 1,4,6,11,12,14,18
Ожидаем на выходе:
1
4,6
11
12,14,18

*/
