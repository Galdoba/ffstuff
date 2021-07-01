package edl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/macroblock/imed/pkg/types"
)

type edlData struct {
	filepath          string
	timeline          []clip
	inputFilePaths    []string
	inputFileCheckMap map[string]bool
}

type clip struct {
	nextClip    *clip  //адресс следующего клипа
	mix         string //тип склейки - если не пуст то дальше идем в nextClip
	fileName    string //имя файла из которого берем данные
	sequanceIN  types.Timecode
	sequanceOUT types.Timecode
	fileIN      types.Timecode
	fileOUT     types.Timecode
	effects     []string
}

func New() (edlData, error) {
	ed := edlData{}
	ed.inputFileCheckMap = make(map[string]bool)
	return ed, nil
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
	scanner := bufio.NewScanner(r)
	i := 0
	for scanner.Scan() {
		i++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var index string
		var reel string
		var mediaType string
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
			mediaType = fields[2]
			effect = fields[3]
			fileIN = fields[4]
			fileOUT = fields[5]
			sequenceIN = fields[6]
			sequenceOUT = fields[7]
		}
		switch {
		default:
			return nil, fmt.Errorf("unknown err = %v", line)
		case index == "*":
			switch {
			default:
				return nil, fmt.Errorf("index field err = %v", line)
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
		case index == "TITLE:":
			//skip
		case index == "EFFECTS":
			//заполняем effect name для клипа
			fmt.Printf("Effect name: %q\n", strings.TrimPrefix(line, "EFFECTS NAME IS "))
		case index == "FCM:":
			fmt.Printf("TODO: разобраться что это %q\n", strings.TrimPrefix(line, "FCM: "))
		case isIndex(index):
			fmt.Printf("TODO: присвоить индекс '%v' к клипу\n", index)
		case reel == "BL":
			fmt.Printf("сегмент пустоты: %q", line)
			continue

		}
		////////////
		fmt.Printf("Source IN = %q | Source OUT = %q\n", fileIN, fileOUT)
		if 5 > 3 {
			fmt.Println(i, line)
			fmt.Print(index)
			fmt.Print(reel)
			fmt.Print(mediaType)
			fmt.Print(effect)
			fmt.Print(fileIN)
			fmt.Print(fileOUT)
			fmt.Print(sequenceIN)
			fmt.Print(sequenceOUT)
			fmt.Print("\n")
		}
	}

	return &eData, nil
}

func isIndex(s string) bool {
	n, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	if n < 1 || n > 999 {
		return false
	}
	return true
}

/*
index
reel
mediaType
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
