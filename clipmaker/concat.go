package clipmaker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Galdoba/ffstuff/fldr"
)

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
		if len(concatIndexes[i]) == 0 || len(concatIndexes[i]) == 1 {
			continue
		}
		concatVideo(cm, concatIndexes[i])
	}

	return "", tasks
}

func concatVideo(cm map[int]clip, partIndexes []int) {
	fmt.Println("Start concat")
	file, err := os.OpenFile(fldr.MuxPath()+"temp.bat", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(2)
	}
	outName := "Out.mp4"
	file.WriteString("set root=" + fldr.MuxPath() + "\n")
	file.WriteString("pushd %root%\n")

	for p, i := range partIndexes {
		sName := strings.TrimPrefix(shortName(cm[i].targetFileName), fldr.MuxPath()) + ".mp4"
		fmt.Println("Clip:", cm[i].index, cm[i].targetFileName)
		nextStr := "echo file " + sName + " " + passToFile(p) + " " + fldr.MuxPath() + "newlist.txt\n"
		if _, err := file.WriteString(nextStr); err != nil {
			panic(err)
		}
	}
	nextStr := "ffmpeg -safe 0 -f concat -i " + "newList.txt" + " -c copy " + outName + "\n"
	if _, err := file.WriteString(nextStr); err != nil {
		panic(err)
	}
	bat := file.Name()
	file.Close()
	runBatchFile(bat)
	fmt.Println("end concat")
}

func newClipList() *os.File {
	fmt.Println("Create File: ", fldr.MuxPath()+"Clip_Parts.txt")
	file, err := os.OpenFile(fldr.MuxPath()+"Clip_Parts.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		panic(err)
	}
	return file
}

func passToFile(i int) string {
	str := ">"
	switch i {
	case 0:
		return str
	default:
		str += ">"
	}
	return str
}

func runBatchFile(path string) error {
	cmd := exec.Command(`cmd.exe`, `/C`, `Start `+path)
	return cmd.Run()
}

/*
создать файл
закидать последовательность клипов
сшить через ффмпег
присвоить имя? как?
*/
