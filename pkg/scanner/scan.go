package scanner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/namedata"
	"github.com/Galdoba/utils"
)

//ScanReady - walks through all paths under the root and returns slice of all ready files
func Scan(root string, querry string) ([]string, error) {
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
