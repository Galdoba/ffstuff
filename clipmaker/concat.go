package clipmaker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Galdoba/ffstuff/fldr"
)

//ffmpeg -i test_Input_AUDIO51_ACLIP_001.m4a -i test_Input_AUDIO51_ACLIP_002.m4a -filter_complex acrossfade=d=0.1:o=0:c1=exp:c2=exp output5.ac3
//ffmpeg -i test_Input_AUDIO51.m4a -vn -af "afade=t=in:ss=0:d=2,afade=t=out:ss=8:d=2" -acodec copy  -t 10 outFaded.m4a  - ДЕЛАЕТ ФЕЙДИН И ФЕЙДАУТ В ФАЙЛЕ

/*
Продседура сшивания звука:
0 по едл понять какие куски станут единым выводом
1 порезать куски с запасом для кросфейда
2 наложить фейды
3 сшить с оверлапом

0:
замуксить звук и видео:
если старт и конец видео и аудио клипа совпадают - они единая компазиция


*/

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
	if !strings.Contains(cm[partIndexes[0]].sourceFileName, "VCLIP") {
		return
	}
	file, err := os.OpenFile(fldr.MuxPath()+"temp.bat", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(2)
	}
	outName := ""
	file.WriteString("set root=" + fldr.MuxPath() + "\n")
	file.WriteString("pushd %root%\n")
	for p, i := range partIndexes {
		if p == 0 {
			outName = "concated_" + strings.TrimPrefix(shortName(cm[i].targetFileName), fldr.MuxPath()) + ".mp4"
		}
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
	file.WriteString("exit\n")
	file.Close()
	err2 := runBatchFile(bat)
	if err2 != nil {
		fmt.Println(err2)
	}
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
	cmd := exec.Command(`cmd.exe`, `/C`, path)
	cmd.Stdout = os.Stdout
	//output, _ := cmd.CombinedOutput()
	//sOUT := string(output)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

/*
создать файл
закидать последовательность клипов
сшить через ффмпег
присвоить имя? как?
*/
