package tabledata

// import (
// 	"fmt"
// 	"os"

// 	"github.com/Galdoba/devtools/cli/command"
// 	"github.com/Galdoba/devtools/csvp"
// 	"github.com/Galdoba/ffstuff/app/autogen/config"
// )

// const (
// 	request = `-s --use-ascii --proxy http://proxy.local:3128 https://docs.google.com/spreadsheets/d/1Waa58usrgEal2Da6tyayaowiWujpm0rzd06P5ASYlsg/gviz/tq?tqx=out:csv -k --output`
// 	//file    = `c:\Users\pemaltynov\.ffstuff\data\taskSpreadsheet.csv`
// )

// //Update -
// func Update(cfg config.Config) error {

// 	if !cfg.TableDataAutoUpdate() {
// 		return nil
// 	}
// 	file := cfg.TableDataFile()

// 	_, _, err := command.Execute("curl "+request+" "+file+".tmp", command.Set(command.TERMINAL_ON))
// 	if err != nil {
// 		return err
// 	}
// 	newPath := file + ".tmp"
// 	oldPath := file
// 	if err := os.Rename(newPath, oldPath); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func FullTableData(path string) ([]TableEntry, error) {
// 	bt, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, fmt.Errorf("can't read file: %v")
// 	}

// 	tablDat, err := csvp.FromString(string(bt))
// 	if err != nil {
// 		return nil, fmt.Errorf("can't create container: %v")
// 	}
// 	dest := ""
// 	AllEntries := []TableEntry{}
// 	for i, e := range tablDat.Entries() {
// 		flds := e.Fields()
// 		if i != 0 && flds[14] == "" {
// 			dest = flds[8]
// 		}

// 		entry := TableEntry{}
// 		entry.TableComment = flds[0]
// 		entry.TablePath = flds[1]
// 		entry.TableReadyTrailer = flds[2]
// 		entry.TableTrailerStatus = flds[3]
// 		entry.TableTrailerMaker = flds[4]
// 		entry.TablePosterStatus = flds[5]
// 		entry.TablePosterMaker = flds[6]
// 		entry.TableLineData0 = flds[7]
// 		entry.TableTaskName = flds[8]
// 		entry.TableEditStatus = flds[9]
// 		entry.TableTaskStatus = flds[10]
// 		entry.TableTaskOutputStatus = flds[11]
// 		entry.TableTroubleStatus = flds[12]
// 		entry.TableAgent = flds[13]
// 		entry.TablePublicationDate = flds[14]
// 		entry.Desination = dest
// 		AllEntries = append(AllEntries, entry)
// 	}
// 	return AllEntries, nil
// }

// type TableEntry struct {
// 	TableComment          string
// 	TablePath             string
// 	TableReadyTrailer     string
// 	TableTrailerStatus    string
// 	TableTrailerMaker     string
// 	TablePosterStatus     string
// 	TablePosterMaker      string
// 	TableLineData0        string
// 	TableTaskName         string
// 	TableEditStatus       string
// 	TableTaskStatus       string
// 	TableTaskOutputStatus string
// 	TableTroubleStatus    string
// 	TableAgent            string
// 	TablePublicationDate  string
// 	Desination            string
// }
