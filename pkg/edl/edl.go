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
	track             []clip //последовательность найденых клипов
	inputFilePaths    []string
	inputFileCheckMap map[string]bool
}

type clip struct {
	////ВАЖНО mix определяет с какого типа склейки клип начинается
	//технически это означает что nextclip должно стать lastclip
	nextClip     *clip  //адресс следующего клипа
	mix          string //тип склейки - если не "C" то дальше идем в nextClip
	mixx         mixType
	sourcefile   string //имя файла из которого берем данные
	fileTime     timeSegment
	sequanceTime timeSegment
	effects      []string
	fcm          string //fcm внутри clip
}

type State struct {
	currentClip  *clip
	waitFileName bool
	waitMix      bool
}

type mixType struct {
	mixCode string
	mixLen  float64
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
	//str += fmt.Sprintf("fcm = %v\n", ed.fcm)
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
	state := State{}
	for scanner.Scan() {
		// parseLine(state, &eData, scanner.Text()) (state, err)
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
		fmt.Println(state)
		switch {
		default:
			return &eData, fmt.Errorf("unknown err = %v | %v", line, parseError)
		case parseError != nil:
			return &eData, parseError
		case index == "*":
			fromFile, toFile := parseComment(line)
			if fromFile != "" && state.waitFileName {
				if !state.waitFileName {
					parseError = fmt.Errorf("unexpected 'fromFile'")
				}
				//заполняем в state поля fromFile
				state.waitFileName = false

			}

			if toFile != "" { //TODO: улучшить запись
				if state.currentClip == nil || state.currentClip.mix == "" {
					parseError = fmt.Errorf("unexpected 'toFile'")
				}
				//колдуем с mix
			}

			//eData, parseError = parseComment(eData, line)
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
		case state.currentClip == nil && effect == "C":
			state.currentClip = &clip{} //выкинуть создание объекта за пределы цикла
			switch reel {
			default:
				parseError = fmt.Errorf("Unknown err")
			case "BL":
				fmt.Printf("сегмент пустоты: %q\n", line)
				fmt.Printf("clip is BL\n")
			case "AX":
				fmt.Printf("Parse main Data:  %v\n", line)
				state.currentClip.fileTime, parseError = parseFileTime(fileIN, fileOUT)
				state.currentClip.sequanceTime, parseError = parseSequenceTime(sequenceIN, sequenceOUT)
				fmt.Println(state.currentClip.fileTime, state.currentClip.sequanceTime)
				//			state.currentclip.lenght = state.currentclip.fileOUT - state.currentclip.fileIN
				eData.track = append(eData.track, *state.currentClip)
				switch trackType {
				default:
					return nil, fmt.Errorf("clip is unknown type = %v", line)
				case "V":
					fmt.Printf("clip is video\n")
				case "A", "A2", "A3", "A4":
					fmt.Printf("clip is audio\n")
				}
				state.waitFileName = true
				state.waitMix = true

			}

		case state.currentClip != nil && state.waitMix:
			switch reel {
			default:
				parseError = fmt.Errorf("Unknown err2")
			case "AX":
				//state.currentClip  TODO: добавляем данные со следующей строки в currentClip
				state.waitMix = false
			}

		}
		if state.currentClip != nil && !state.waitFileName && !state.waitMix {
			//аппендим clip в edlData
			state = State{}
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

func parseComment(line string) (fromFile, toFile string) {
	//func parseComment( line string) (fromFile, toFile string) {
	switch {
	default:
	case strings.HasPrefix(line, "* FROM CLIP NAME: "):
		//заполняем Source file name для клипа
		fromFile = strings.TrimPrefix(line, "* FROM CLIP NAME: ")
		//fmt.Printf("Source file name: %q\n", source)
		//eData.inputFileCheckMap[source] = true
	case strings.HasPrefix(line, "* TO CLIP NAME: "):
		//заполняем Dest file name для клипа
		toFile = strings.TrimPrefix(line, "* TO CLIP NAME: ")
		//fmt.Printf("Dest file name: %q\n", source)
		//eData.inputFileCheckMap[source] = true

	}
	return
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

func isStandardStatement(line string) bool {
	if len(strings.Fields(line)) < 7 {
		return false
	}
	return true
}

type StandardStatement struct {
	index        string
	reel         string
	channels     string
	editType     string
	editDuration string
	fileIN       types.Timecode
	fileOUT      types.Timecode
	sequanceIN   types.Timecode
	sequanceOUT  types.Timecode
}

func parseFields(line string) (*StandardStatement, error) {
	ss := StandardStatement{}
	err := errors.New("Initial")
	err = nil
	fields := strings.Fields(line)
	switch len(fields) {

	case 9:
		for i, _ := range fields {
			if err != nil {
				return nil, err
			}
			switch i {
			case 0:
				ss.index = fields[i]
			case 1:
				ss.reel = fields[i]
			case 2:
				ss.channels = fields[i]
			case 3:
				ss.editType = fields[i]
			case 4:
				ss.editDuration = fields[i]
			case 5:
				ss.fileIN, err = types.ParseTimecode(fields[i])
			case 6:
				ss.fileOUT, err = types.ParseTimecode(fields[i])
			case 7:
				ss.sequanceIN, err = types.ParseTimecode(fields[i])
			case 8:
				ss.sequanceOUT, err = types.ParseTimecode(fields[i])
			}
		}
	case 8:
		for i, _ := range fields {
			if err != nil {
				return nil, err
			}
			switch i {
			case 0:
				ss.index = fields[i]
			case 1:
				ss.reel = fields[i]
			case 2:
				ss.channels = fields[i]
			case 3:
				ss.editType = fields[i]
			case 4:
				ss.fileIN, err = types.ParseTimecode(fields[i])
			case 5:
				ss.fileOUT, err = types.ParseTimecode(fields[i])
			case 6:
				ss.sequanceIN, err = types.ParseTimecode(fields[i])
			case 7:
				ss.sequanceOUT, err = types.ParseTimecode(fields[i])
			}
		}

	}
	return &ss, nil
}
