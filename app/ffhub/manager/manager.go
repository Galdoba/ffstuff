package main

import (
	"fmt"

	"github.com/Galdoba/devtools/keyval"
	"github.com/Galdoba/ffstuff/pkg/namedata"
)

func main() {
	//проходим по задачам
	//	исходя из текущего статуса предлогаем команду

	kv, err := keyval.Load("fftasks_status")
	if err != nil {
		fmt.Println(err.Error())
	}

	for i, vals := range kv.Keys() {
		v, e := kv.GetSingle(vals)
		fmt.Println(i, vals, v, e)
		switch v {
		default:
			fmt.Println("Skip")
		case "0":
			fmt.Println("Нужен файл проекта")
			edName := namedata.EditForm(vals)
			fmt.Println(edName)
			if !keyval.KVlistPresent(vals) {
				_, err := keyval.NewKVlist("fftasks_inputfiles")
				if err != nil {
					fmt.Println(err.Error())
				}
			}
			kv, _ := keyval.Load("fftasks_inputfiles")
			kv.Set(edName.EditName(), vals)

		}

	}
}
