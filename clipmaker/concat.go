package clipmaker

import "fmt"

func ConcatClips(cm map[int]clip) (string, []string) {
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

	return "", tasks
}

/*
создать файл
закидать последовательность клипов
сшить через ффмпег
присвоить имя? как?
*/
