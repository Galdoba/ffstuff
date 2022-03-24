package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"os"

	"github.com/olekukonko/tablewriter"
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
	LoudnormScan  = "loudnormSCAN"
	SounddcanFAST = "sounscanFAST"
	SounddcanFULL = "sounscanFULL"
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
					ir.Print()

					if !ir.allFilesReported() {
						for _, fl := range ir.files {
							if fl.loudnormReport == false {
								source := fl.filepath
								repPath := reportPath(source, LoudnormScan)
								fmt.Println(fl.filepath)
								fmt.Println("Scanning...")
								info.MakeLoudnormReport(source, repPath)
								break
							}

						}
					} else {
						done = true
					}

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

func (ir *issuesReport) allFilesReported() bool {
	for _, fl := range ir.files {
		if !fl.loudnormReport {
			return false
		}
	}
	return true
}

func reportPath(source, reportType string) string {
	filename := strings.TrimSuffix(namedata.RetrieveShortName(source), ".m4a")
	repDirectory := namedata.RetrieveDirectory(source)
	return repDirectory + "reports\\" + filename + "." + reportType
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

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (ir *issuesReport) String() string {

	rep := fmt.Sprintf("Tracking %v file(s):%v ", len(ir.files), ir.Summary())
	rep += "--------------------------------------------------------------------------------\n"
	for _, fr := range ir.files {
		rep += fr.String()
	}
	rep += "--------------------------------------------------------------------------------"
	return rep
}

func (ir *issuesReport) Print() {
	clearScreen()
	table := tablewriter.NewWriter(os.Stdout)
	for _, fl := range ir.files {
		flData := []string{namedata.RetrieveShortName(fl.filepath)}
		switch fl.loudnormReport {
		case true:
			for _, d := range info.LoudnormData(reportPath(fl.filepath, LoudnormScan)) {
				flData = append(flData, d)
			}
		case false:
			flData = append(flData, "No Data")
			flData = append(flData, "No Data")
			flData = append(flData, "No Data")
			flData = append(flData, "No Data")

		}
		colors := assignColors(flData)
		table.Rich(flData, colors)
	}
	table.SetHeader([]string{"File", "IL", "RA", "channels", "stats"})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})
	table.SetColumnSeparator("|")
	table.Render()

}

func assignColors(dataRow []string) []tablewriter.Colors {
	colors := []tablewriter.Colors{}
	for i, dt := range dataRow {
		chosenColor := tablewriter.Colors{}
		switch {
		default:
			switch i {
			case 2:
				chosenColor = processRAdata(dt)
			case 3:
				chosenColor = processChannelsData(dt)
			case 4:
				chosenColor = processStatsData(dt)
			}
		case dt == "data corrupted":
			chosenColor = tablewriter.Colors{tablewriter.FgRedColor}
		case dt == "No Data":
			chosenColor = tablewriter.Colors{tablewriter.FgYellowColor}
		}
		colors = append(colors, chosenColor)
	}
	return colors
}

func processRAdata(data string) tablewriter.Colors {
	col := tablewriter.Colors{}
	ra, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return tablewriter.Colors{tablewriter.FgYellowColor}
	}
	if ra >= 19 {
		col = tablewriter.Colors{tablewriter.FgYellowColor}
	}
	if ra >= 20 {
		col = tablewriter.Colors{tablewriter.FgHiRedColor}
	}
	return col
}

func processChannelsData(data string) tablewriter.Colors {
	col := tablewriter.Colors{}
	channel := []int{}
	chanStr := strings.Fields(data)
	for _, v := range chanStr {
		ch, err := strconv.Atoi(v)
		if err != nil {
			return tablewriter.Colors{tablewriter.FgRedColor}
		}
		channel = append(channel, ch)
	}
	switch len(channel) {
	default:
		return tablewriter.Colors{tablewriter.FgRedColor}
	case 2:
		if difference(channel[0], channel[1]) == 1 {
			return tablewriter.Colors{tablewriter.FgYellowColor}
		}
		if difference(channel[0], channel[1]) > 1 {
			return tablewriter.Colors{tablewriter.FgHiRedColor}
		}
	case 6:
		if difference(channel[2], 0) > 0 {
			col = tablewriter.Colors{tablewriter.FgYellowColor}
		}
		if difference(channel[2], 0) > 1 {
			col = tablewriter.Colors{tablewriter.FgHiRedColor}
		}
		if difference(channel[0], channel[1]) > 1 {
			return tablewriter.Colors{tablewriter.FgHiRedColor}
		}
		if difference(channel[4], channel[5]) > 1 {
			return tablewriter.Colors{tablewriter.FgHiRedColor}
		}
	}

	return col
}

func processStatsData(data string) tablewriter.Colors {
	col := tablewriter.Colors{}
	stat := []int{}
	chanStr := strings.Split(data, "<")
	for _, v := range chanStr {
		ch, err := strconv.Atoi(v)
		if err != nil {
			return tablewriter.Colors{tablewriter.FgHiRedColor}
		}
		stat = append(stat, ch)
	}
	warn := 0
	if stat[0] >= 75 {
		warn++
	}
	if stat[0] >= 80 {
		warn++
	}
	if stat[1] <= 25 {
		warn++
	}
	if stat[1] <= 20 {
		warn++
	}
	if stat[2] <= 3 {
		warn++
	}
	if stat[2] <= 1 {
		warn++
	}
	if warn > 0 {
		col = tablewriter.Colors{tablewriter.FgYellowColor}
	}
	if warn > 2 {
		col = tablewriter.Colors{tablewriter.FgHiRedColor}
	}

	return col
}

func difference(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

/*
colorData1 := []string{"TestCOLOR1Merge", "HelloCol2 - COLOR1", "HelloCol3 - COLOR1", "HelloCol4 - COLOR1"}
for i, row := range data {
	if i == 4 {
		table.Rich(colorData1,
			[]tablewriter.Colors{
				tablewriter.Colors{},
				tablewriter.Colors{tablewriter.Normal, tablewriter.FgCyanColor},
				tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor},
				tablewriter.Colors{},
			})

	}
	table.Append(row)
}

*/

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
	filepath       string
	loudnormReport bool
	issues         map[int]string
	data           []string
}

func newFileReport(path string) *fileReport {
	fr := fileReport{}
	fr.filepath = path
	lnPath := reportPath(fr.filepath, LoudnormScan)
	if _, err := os.Stat(lnPath); err == nil {
		fr.loudnormReport = true
	}
	fr.issues = make(map[int]string)
	for _, checkType := range allChecks {
		//fr.issues[checkType] = "No data provided"
		if checkType == checkByLoudnorm {

			if _, err := os.Stat(lnPath); err != nil {
				continue
			}
			//	fr.issues[checkType] = info.LoudnormReportToString(lnPath)
			//	info.LoudnormData(lnPath)
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
