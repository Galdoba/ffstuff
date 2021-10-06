package scanner

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/Galdoba/ffstuff/pkg/namedata"
	"github.com/Galdoba/utils"
)

//ScanReady - walks through all paths under the root and returns slice of all ready files
func Scan(root string, querry string) ([]string, error) {
	var resultList []string
	//открываем корень и собираем статистику
	rootStat, errSt := os.Stat(root)
	if errSt != nil {
		return resultList, errSt
	}
	if !rootStat.IsDir() {
		return resultList, errors.New("root is not a directory")
	}
	srcInfo, errS := os.Open(root)
	if errS != nil {
		return resultList, errors.New("scan source error: " + errS.Error()) //
	}
	defer srcInfo.Close()

	//Читаем и получаем список всего находящегося в корне
	found, errR := srcInfo.Readdir(0)
	if errR != nil {
		return resultList, errors.New("scan read error: " + errR.Error())
	}
	for _, val := range found {
		if strings.Contains(val.Name(), querry) {
			//fmt.Println(root + val.Name())
			resultList = append(resultList, root+val.Name())
		}
		if val.IsDir() {
			subResults, errSub := Scan(root+val.Name()+"\\", querry)
			if errSub != nil {
				return resultList, errSub
			}
			resultList = append(resultList, subResults...)
		}
	}
	return resultList, nil
}

func ListContent(root string) ([]string, error) {
	var resultList []string
	//открываем корень и собираем статистику
	rootStat, errSt := os.Stat(root)
	if errSt != nil {
		return resultList, errors.New("specified root was not found")
	}
	if !rootStat.IsDir() {
		return resultList, errors.New("root is not a directory")
	}
	srcInfo, errS := os.Open(root)
	if errS != nil {
		return resultList, errors.New("scan source error: " + errS.Error()) //
	}
	defer srcInfo.Close()

	//Читаем и получаем список всего находящегося в корне
	found, errR := srcInfo.Readdir(0)
	if errR != nil {
		return resultList, errors.New("scan read error: " + errR.Error())
	}
	for _, val := range found {
		resultList = append(resultList, val.Name())
	}
	return resultList, nil
}

//evaluate - возвращает ошибку с путем до найденного ready файла
func evaluate(path string, f os.FileInfo, err error) error {
	if f.IsDir() {
		fmt.Println("Scanning: ", path)
		return nil
	}
	if filepath.Ext(path) != "ready" {
		return nil
	}
	return errors.New(path)
}

func ListReady(readyfiles []string) []string {
	resSl := []string{}
	for i := range readyfiles {
		//fmt.Println(i, results[i])
		name := namedata.RetrieveShortName(readyfiles[i])
		name = strings.TrimSuffix(name, ".ready")
		dir := namedata.RetrieveDirectory(readyfiles[i])
		sList, err2 := ListContent(dir)
		if err2 != nil {
			fmt.Println(err2.Error())
		}
		for f := range sList {
			if strings.Contains(sList[f], name) {
				resSl = append(resSl, dir+sList[f])
			}
		}
	}
	sorted := []string{}
	sorted = AppendIfContainsStr(sorted, resSl, ".ready")
	for _, val := range resSl {
		if strings.Contains(val, ".ready") {
			sorted = append(sorted, val)
		}
	}
	for _, val := range resSl {
		if strings.Contains(val, "_Proxy_") {
			sorted = append(sorted, val)
		}
	}
	for _, val := range resSl {
		if strings.Contains(val, ".m4a") {
			sorted = append(sorted, val)
		}
	}
	for _, val := range resSl {
		sorted = utils.AppendUniqueStr(sorted, val)
	}
	return resSl
}

func downloadingMarker(path string) (string, error) {
	trimmed := strings.TrimSuffix(path, ".ready")
	if trimmed == path {
		return path, errors.New("cannot trim marker")
	}
	cu, _ := user.Current()
	username := cu.Name

	return trimmed + ".downloadingBY" + username, nil
}

func AppendIfContainsStr(targetSl []string, selectionSl []string, marker string) []string {
	for _, val := range selectionSl {
		if strings.Contains(val, marker) {
			targetSl = append(targetSl, val)
		}
	}
	return targetSl
}

func SortPriority(rf []string) []string {
	sorted := []string{}

	keys := []string{}
	currentTime := time.Now()
	for i := -30; i < 8; i++ {
		date := utils.DateStampFrom(currentTime.AddDate(0, 0, i))
		date = strings.TrimPrefix(date, "20")
		date = strings.Replace(date, "-", "_", -1)
		keys = append(keys, date)
	}
	keys = append(keys, "_amedia", "_wb", "_disney")
	sorted = sortByKeys(rf, keys...)
	return sorted
}

func sortByKeys(sl []string, keys ...string) []string {
	newSl := []string{}
	for _, key := range keys {
		for _, val := range sl {
			if strings.Contains(val, key) {
				newSl = append(newSl, val)
			}
		}
	}
	for _, val := range sl {
		if !utils.ListContains(newSl, val) {
			newSl = append(newSl, val)
		}
	}
	resSl := []string{}
	for _, val := range newSl {
		resSl = utils.AppendUniqueStr(resSl, val)
	}
	return resSl
}
