package scanner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//ScanReady - walks through all paths under the root and returns slice of all ready files
func Scan(root string, querry string) ([]string, error) {
	var resultList []string
	//открываем корень и собираем статистику
	rootStat, errSt := os.Stat(root)
	if errSt != nil {
		return resultList, errors.New("Specified root was not found")
	}
	if !rootStat.IsDir() {
		return resultList, errors.New("Root is not a directory")
	}
	srcInfo, errS := os.Open(root)
	if errS != nil {
		return resultList, errors.New("Scan source error: " + errS.Error()) //
	}
	defer srcInfo.Close()

	//Читаем и получаем список всего находящегося в корне
	found, errR := srcInfo.Readdir(0)
	if errR != nil {
		return resultList, errors.New("Scan read error: " + errR.Error())
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
		return resultList, errors.New("Specified root was not found")
	}
	if !rootStat.IsDir() {
		return resultList, errors.New("Root is not a directory")
	}
	srcInfo, errS := os.Open(root)
	if errS != nil {
		return resultList, errors.New("Scan source error: " + errS.Error()) //
	}
	defer srcInfo.Close()

	//Читаем и получаем список всего находящегося в корне
	found, errR := srcInfo.Readdir(0)
	if errR != nil {
		return resultList, errors.New("Scan read error: " + errR.Error())
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
