package edl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"github.com/macroblock/imed/pkg/types"
)

type edlData struct {
	edlSource         string //источник самого edl (не обязательно файл)
	title             string //title внутри edl
	fcm               string //fcm внутри edl
	track             []clip //последовательность найденых клипов
	inputFilePaths    []string
	inputFileCheckMap map[string]bool
}

type clip struct {
	////ВАЖНО mix определяет с какого типа склейки клип начинается
	//технически это означает что nextclip должно стать lastclip
	nextClip     *clip  //адресс следующего клипа
	mix          string //тип склейки - если не "C" то дальше идем в nextClip
	sourcefile   string //имя файла из которого берем данные
	fileTime     timeSegment
	sequanceTime timeSegment
	effects      []string
}

type timeSegment struct {
	in     types.Timecode
	out    types.Timecode
	lenght types.Timecode
}

func New() (edlData, error) {
	ed := edlData{}
	ed.inputFileCheckMap = make(map[string]bool)
	return ed, nil
}

func (ed edlData) String() string {
	str := ""
	str += fmt.Sprintf("edlSource = %v\n", ed.edlSource)
	str += fmt.Sprintf("title = %v\n", ed.title)
	str += fmt.Sprintf("fcm = %v\n", ed.fcm)
	for i, val := range ed.track {
		str += fmt.Sprintf("track %v = %v\n", i, val)
	}
	str += fmt.Sprintf("inputFilePaths = %v\n", ed.inputFilePaths)
	str += fmt.Sprintf("inputFileCheckMap:\n")
	for k, v := range ed.inputFileCheckMap {
		str += fmt.Sprintf("%v  %v\n", k, v)
	}
	return str
}

func (ts timeSegment) String() string {
	return fmt.Sprintf("[ IN %v | OUT %v | LEN %v ]\n", ts.in.HHMMSSMs(), ts.out.HHMMSSMs(), ts.lenght.HHMMSSMs())
}

// edl.Parse("file.edl") (edlData, error)

func ParseFile(path string) (*edlData, error) {
	fmt.Println("Start Parse File")

	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer f.Close()

	return Parse(f)
}

func Parse(r io.Reader) (*edlData, error) {
	fmt.Println("Start Parse Reader")
	eData := edlData{}
	eData.inputFileCheckMap = make(map[string]bool)
	//eData, parseError = parseLine()
	scanner := bufio.NewScanner(r)
	parseError := errors.New("Initial")
	parseError = nil
	i := 0
	for scanner.Scan() {
		i++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var index string
		var reel string
		var trackType string
		var effect string
		var fileIN string
		var fileOUT string
		var sequenceIN string
		var sequenceOUT string
		//рабочая зона
		fields := strings.Fields(line)
		index = fields[0]
		if len(fields) > 7 {
			reel = fields[1]
			trackType = fields[2]
			effect = fields[3]
			fileIN = fields[4]
			fileOUT = fields[5]
			sequenceIN = fields[6]
			sequenceOUT = fields[7]
		}

		switch {
		default:
			return &eData, fmt.Errorf("unknown err = %v", line)
		case parseError != nil:
			return &eData, parseError
		case index == "*":
			eData, parseError = parseComment(eData, line)
		case index == "TITLE:":
			eData, parseError = parseTitle(eData, line)
		case index == "FCM:":
			eData, parseError = parseFCM(eData, line)
		case index == "EFFECTS":
			//заполняем effect name для клипа
			fmt.Printf("TODO:   EFFECTS not implemented\n")
			fmt.Printf("Effect name: %q\n", strings.TrimPrefix(line, "EFFECTS NAME IS "))
		case reel == "BL":
			fmt.Printf("сегмент пустоты: %q\n", line)
			fmt.Printf("clip is BL\n")
		case reel == "AX":
			fmt.Printf("Parse main Data:  %v\n", line)
			newclip := clip{} //выкинуть создание объекта за пределы цикла
			newclip.fileTime, parseError = parseFileTime(fileIN, fileOUT)
			newclip.sequanceTime, parseError = parseSequenceTime(sequenceIN, sequenceOUT)
			fmt.Println(newclip.fileTime, newclip.sequanceTime)
			//			newclip.lenght = newclip.fileOUT - newclip.fileIN
			eData.track = append(eData.track, newclip)
			switch trackType {
			default:
				return nil, fmt.Errorf("clip is unknown type = %v", line)
			case "V":
				fmt.Printf("clip is video\n")
			case "A", "A2", "A3", "A4":
				fmt.Printf("clip is audio\n")
			}

		}

		////////////
		//fmt.Printf("Source IN = %q | Source OUT = %q\n", fileIN, fileOUT)

		//fmt.Printf("Sequance IN = %q | Sequance OUT = %q\n", sequenceIN, sequenceOUT)
		if 5 < 3 {
			fmt.Println(i, line)
			fmt.Print(index)
			fmt.Print(reel)
			fmt.Print(trackType)
			fmt.Print(effect)
			fmt.Print(fileIN)
			fmt.Print(fileOUT)
			fmt.Print(sequenceIN)
			fmt.Print(sequenceOUT)
			fmt.Print("\n")
		}
	}
	fmt.Println("End Parse Reader")
	return &eData, nil
}

func parseLine(eData edlData, line string) (edlData, error) {
	fmt.Println("parseLine(eData edlData, line string) (edlData, error) - not implemented")
	return eData, nil
}

func parseTitle(eData edlData, line string) (edlData, error) {
	title := strings.TrimPrefix(line, "TITLE: ")
	if title == line {
		return eData, fmt.Errorf("title cannot be parsed %v", line)
	}
	eData.title = title
	return eData, nil
}

func parseFCM(eData edlData, line string) (edlData, error) {
	fmt.Printf("TODO: разобраться что это %q\n", strings.TrimPrefix(line, "FCM: "))
	fcm := strings.TrimPrefix(line, "FCM: ")
	if fcm == line {
		return eData, fmt.Errorf("fcm cannot be parsed %v", fcm)
	}
	return eData, nil
}

func parseComment(eData edlData, line string) (edlData, error) {
	switch {
	default:
		return eData, fmt.Errorf("index field err = %v", line)
	case strings.HasPrefix(line, "* FROM CLIP NAME: "):
		//заполняем Source file name для клипа
		source := strings.TrimPrefix(line, "* FROM CLIP NAME: ")
		fmt.Printf("Source file name: %q\n", source)
		eData.inputFileCheckMap[source] = true
	case strings.HasPrefix(line, "* TO CLIP NAME: "):
		//заполняем Dest file name для клипа
		source := strings.TrimPrefix(line, "* TO CLIP NAME: ")
		fmt.Printf("Dest file name: %q\n", source)
		eData.inputFileCheckMap[source] = true
	}
	return eData, nil
}

func parseFileTime(fileIN, fileOUT string) (timeSegment, error) {
	fileTime := timeSegment{}
	in, errIN := types.ParseTimecode(fileIN)
	if errIN != nil {
		return fileTime, fmt.Errorf("can't parse fileIN: %v", errIN.Error())
	}
	out, errOUT := types.ParseTimecode(fileOUT)
	if errOUT != nil {
		return fileTime, fmt.Errorf("can't parse fileOUT: %v", errIN.Error())
	}
	fileTime.in = in
	fileTime.out = out
	fileTime.lenght = fileTime.out - fileTime.in
	return fileTime, nil
}

func roundFloatTo(fl float64, digit float64) float64 {
	return math.Round(fl/digit) * digit
}

//TODO: посмотреть как можно адекватно избавиться от дублирующей функции
//сейчас она нужна только, чтобы отличать ошибки fileTime от sequenceTime
func parseSequenceTime(sequanceIN, sequanceOUT string) (timeSegment, error) {
	sequanceTime := timeSegment{}
	in, errIN := types.ParseTimecode(sequanceIN)
	if errIN != nil {
		return sequanceTime, fmt.Errorf("can't parse sequanceIN: %v", errIN.Error())
	}
	out, errOUT := types.ParseTimecode(sequanceOUT)
	if errOUT != nil {
		return sequanceTime, fmt.Errorf("can't parse sequanceOUT: %v", errIN.Error())
	}
	sequanceTime.in = in
	sequanceTime.out = out
	sequanceTime.lenght = sequanceTime.out - sequanceTime.in
	return sequanceTime, nil
}

func parseClip(eData edlData, line string) (edlData, error) {

	return eData, nil
}

func isFolower(cl clip) bool {
	return false
}

/*
index
reel
trackType
effect
fileIN
fileOUT
sequenceIN
sequenceOUT


*/

// func readFile(fPath string) error {
// 	fmt.Println("readFileWithReadString")
// 	file, err := os.Open(fPath)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	// Start reading from the file with a reader.
// 	reader := bufio.NewReader(file)
// 	var line string
// 	for {
// 		line, err = reader.ReadString('\n')
// 		if err != nil && err != io.EOF {
// 			break
// 		}

// 		// Process the line here.
// 		fmt.Printf(" > Read %d characters\n", len(line))

// 		if err != nil {
// 			break
// 		}
// 	}
// 	if err != io.EOF {
// 		fmt.Printf(" > Failed with error: %v\n", err)
// 		return err
// 	}
// 	return nil
// }
