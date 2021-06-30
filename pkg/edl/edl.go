package edl

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
}

func New(path string) (edlData, error) {
	ed := edlData{}
	ed.inputFileCheckMap = make(map[string]bool)
	return ed, nil
}

// edl.Parse("file.edl") (edlData, error)

func ParseFile(path string) (*edlData, error) {
	fmt.Println("Start Parse")

	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer f.Close()

	return Parse(f)
}

func Parse(r io.Reader) (*edlData, error) {
	fmt.Println("Start Parse")

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

		switch {
		default:
			return nil, fmt.Errorf("unknown err = %v", line)
		case index == "*":
			switch {
			default:
				return nil, fmt.Errorf("index field err = %v", line)
			case strings.HasPrefix(line, "* FROM CLIP NAME: "):
				//заполняем Source file name для клипа
				fmt.Printf("Source file name: %q\n", strings.TrimPrefix(line, "* FROM CLIP NAME: "))
			case strings.HasPrefix(line, "* TO CLIP NAME: "):
				//заполняем Dest file name для клипа
				fmt.Printf("Dest file name: %q\n", strings.TrimPrefix(line, "* TO CLIP NAME: "))
			}
		case index == "TITLE:":
			//skip
		case index == "FCM:":
			fmt.Printf("TODO: разобраться что это %q\n", strings.TrimPrefix(line, "FCM: "))
		case reel == "BL":
			fmt.Printf("сегмент пустоты: %q", line)
		}
		////////////

	}

	return &edlData{}, nil
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
