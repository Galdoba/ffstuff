package main

import (
	"fmt"
	"os/exec"
	"strings"

	"os"

	"github.com/urfave/cli"

	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/ffstuff/pkg/config"
	"github.com/Galdoba/ffstuff/pkg/glog"
	"github.com/Galdoba/ffstuff/pkg/info"
	"github.com/Galdoba/ffstuff/pkg/namedata"
	"github.com/Galdoba/ffstuff/pkg/scanner"
)

const (
	checkByLoudnorm = iota
	checkBySoundscanFAST
	checkBySoundscanFULL
)

var configMap map[string]string

var logger glog.Logger
var logLocation string
var issueFilePath string
var root string
var allChecks []int

func init() {
	conf, err := config.ReadProgramConfig("ffstuff")
	if err != nil {
		fmt.Println(err)
	}
	configMap = conf.Field
	if err != nil {
		switch err.Error() {
		case "Config file not found":
			fmt.Print("Expecting config file in:\n", conf.Path)
			os.Exit(1)
		}
	}
	root = fldr.InPath()
	issueFilePath = root + "issues.txt"
	if _, err := os.Stat(issueFilePath); err != nil {
		os.Create(issueFilePath)
		fmt.Println("File created:", issueFilePath)
	}
	allChecks = []int{checkByLoudnorm} //, checkBySoundscanFAST, checkBySoundscanFULL}

}

func main() {
	//root := configMap[constant.SearchRoot]
	//marker := configMap[constant.SearchMarker]
	// if configMap[constant.LogDirectory] == "default" {
	// 	logLocation = fldr.MuxPath() + "logfile.txt"
	// }
	//logger = glog.New(logLocation, glog.LogLevelINFO)
	logger = glog.New(glog.LogPathDEFAULT, glog.LogLevelINFO)

	app := cli.NewApp()
	app.Version = "v 0.0.1"
	app.Name = "issuetracker"
	app.Usage = "Scans audio with Loudnorm and Soundscan and creates report files for analisys"
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "vocal",
			Usage: "If flag is active soundscan will print data on terminal",
		},
	}
	app.Commands = []cli.Command{
		//////////////////////////////////////
		{
			Name:  "track",
			Usage: "Listens all files and checks silence in audio stream",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "repeat, r ",
					Usage: "If used creates loop to run issuetraker every x seconds",
					Value: "180",
				},
			},
			Action: func(c *cli.Context) error {
				done := false
				for !done {
					ir := newIssueReport()
					found, err := scanner.Scan(root, ".m4a")
					fmt.Println(err)
					for _, f := range found {
						_, err := os.Stat(f)
						if err != nil {
							fmt.Println(err.Error())
						}

						ir.files = append(ir.files, *newFileReport(f))
					}
					fmt.Println(ir.String())
					ungraceful := false
					for _, fl := range ir.files {
						if len(ir.files) == 0 {
							os.Exit(0)
						}

						if fl.issues[checkByLoudnorm] == "No data provided" {
							source := fl.filepath
							repPath := reportPathLN(source)
							fmt.Println(fl.filepath)
							fmt.Println("Scanning...")
							info.MakeLoudnormReport(source, repPath)
							ungraceful = true
							break
						} else {
							continue
						}
					}
					if ungraceful {
						continue
					}
					ir.Print()
					done = true
					// if len(validSources) > 0 {
					// 	source := validSources[0]

					// 	fmt.Println(source)
					// 	fmt.Println("Scanning...")
					// 	err := info.MakeLoudnormReport(source, reportPathLN(source))
					// 	if err != nil {
					// 		fmt.Println(err)
					// 		panic(err)
					// 	}
					// 	summary := info.LoudnormReportToString(reportPathLN(source))
					// 	for i, fl := range ir.files {
					// 		if fl.filepath == validSources[0] {
					// 			ir.files[i].issues[checkByLoudnorm] = summary
					// 			break
					// 		}
					// 	}
					// }
					// if len(validSources) == 0 {
					// 	done = true
					// }

				}
				return nil
			},
		},
	}
	args := os.Args
	if len(args) < 2 {
		args = append(args, "help") //Принудительно зовем помощь если нет других аргументов
	}
	if err := app.Run(args); err != nil {
		fmt.Println(err.Error())
	}
}

/*
issuetracker run --repeat 42
init()
1. определяем папку IN
2. Создаем issues.txt
Run()
1. Создаем список Аудио в Папке IN
2. Удоставеряемся что для каждого аудио есть Report
3. Публикуем issues.txt

track - запускает треккер
--repeat x  - повторяет процесс каждые x секунд

*/

func reportPathLN(source string) string {
	filename := strings.TrimSuffix(namedata.RetrieveShortName(source), ".m4a")
	repDirectory := namedata.RetrieveDirectory(source)
	return repDirectory + "proxy\\" + filename + "_loudnorm_report.txt"
}

func createIssueFile() {
	os.Create(issueFilePath)
	f, err := os.Open(issueFilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
}

type issuesReport struct {
	reportPath string
	files      []fileReport
}

func newIssueReport() *issuesReport {
	return &issuesReport{reportPath: issueFilePath}
}

func (ir *issuesReport) String() string {
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
	rep := fmt.Sprintf("Tracking %v file(s):\n%v", len(ir.files), ir.Summary())
	rep += "--------------------------------------------------------------------------------\n"
	for _, fr := range ir.files {
		rep += fr.String()
	}
	rep += "--------------------------------------------------------------------------------"
	return rep
}

func (ir *issuesReport) Print() {
	fmt.Println("Here be Total Summary Table...")

}

func (ir *issuesReport) Summary() string {
	noData := 0
	ok := 0
	warnings := 0
	total := 0
	for _, fr := range ir.files {
		for _, issue := range fr.issues {
			total++
			switch issue {
			case "ok":
				ok++
			case "No data provided":
				noData++
			default:
				warnings++
			}
		}
	}
	return fmt.Sprintf("Tests passed: %v/%v (%v warnings; No data on %v tests)\n", ok, total, warnings, noData)
}

type fileReport struct {
	filepath string
	issues   map[int]string
}

func newFileReport(path string) *fileReport {
	fr := fileReport{}
	fr.filepath = path
	fr.issues = make(map[int]string)
	for _, checkType := range allChecks {
		fr.issues[checkType] = "No data provided"
		if checkType == checkByLoudnorm {
			lnPath := reportPathLN(fr.filepath)
			if _, err := os.Stat(lnPath); err != nil {
				continue
			}
			fr.issues[checkType] = info.LoudnormReportToString(lnPath)
		}
	}
	return &fr
}

func (fr *fileReport) String() string {
	rep := namedata.RetrieveShortName(fr.filepath) + ":"
	rep += "   Loudnorm (Scan) : " + fr.issues[checkByLoudnorm] + "\n"
	//rep += "   Soundscan (FAST): " + fr.issues[checkBySoundscanFAST] + "\n"
	//rep += "   Soundscan (FULL): " + fr.issues[checkBySoundscanFULL] + "\n"
	return rep
}

/*
issueFile Exapmle:
START
File [file1.m4a]:
 Loudnorm warning:
 warning 1
 warning 2
 ...
 warning n
 Loudnorm report end.
 --------------------
 Soundscan report start:
 [report text line 1]
 [report text line 2]
 ...
 [report text line n]
 Soundscan report end.
 --------------------
[blank line]
File [file2.m4a]:
...
[blank line]
END
*/
