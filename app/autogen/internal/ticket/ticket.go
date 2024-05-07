package ticket

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Galdoba/ffstuff/app/autogen/internal/tabledata"
	"github.com/Galdoba/ffstuff/pkg/translit"
)

const (
	UNDEFINED = iota
	Check_Status_NOT_EXPECTED
	Check_Status_PASS
	Check_Status_FAIL
	Check_Status_WAIT
	Check_Status_ERR
	FILM
	TRL
	SER
)

func categoryPrefix(i int) string {
	switch i {
	case FILM:
		return "FILM"
	case TRL:
		return "TRL"
	case SER:
		return "SER"
	default:
		return "???"
	}
}

type Ticket struct {
	Name           string    `json:"Name"`                      //uname
	Category       int       `json:"Category"`                  //FILM-TRL-SER
	Contragent     string    `json:"Contragent"`                // заказчик (Amedia или иное)
	SourceFiles    []string  `json:"Source Files"`              // файлы относящиеся к тикету
	Season         int       `json:"Season Number,omitempty"`   //Номер сезона если сериал
	Episodes       []int     `json:"Episode Numbers,omitempty"` //Номер эпизода если сериал
	Register       string    `json:"Register,omitempty"`        //Регистр серии (например 14b) если сериал
	BasicCheck     int       `json:"Basic Scan Status"`         //статус проверки
	InterlaceCheck int       `json:"Interlace Scan Status"`     //статус проверки
	LoudnessCheck  int       `json:"Loudness Scan Status"`      //статус проверки
	ReadWriteCheck int       `json:"ReadWrite Scan Status"`     //статус проверки
	TimeOpen       time.Time `json:"Opened"`                    //время открытия тикета
	TimeClose      time.Time `json:"Closed"`                    //время закрытия тикета
	TableData      []string  `json:"Table Data"`                //данные из таблицы
}

/*
нужно для создания тикета:
название из таблицы
категория (фильм/сериал/трейлер)
контрагент

тикета дает данные:
basename - транслитное имя основы
destination - папку куда складывать результат
archive - папку куда складывать исходники
*/

func FilmFromEntry(entry tabledata.TableEntry) (*Ticket, error) {
	tkt := Ticket{}
	if entry.TableAgent == "" {
		return nil, fmt.Errorf("ticket must have contragent")
	}
	tableName := translit.TransliterateLower(entry.TableTaskName)
	//fmt.Println("try", tableName)
	if strings.Contains(tableName, translit.TransliterateLower("сезон")) && strings.Contains(tableName, translit.TransliterateLower("серия")) {
		return nil, fmt.Errorf("ticket have serial numbers")
	}
	//fmt.Printf("%v: TASKStatus: '%v'\n", tableName, entry.TableTaskStatus)
	if entry.TableEditStatus == "v" || entry.TableEditStatus == "м" {
		return nil, fmt.Errorf("ticket have DOWNLOADING status")
	}
	if entry.TableEditStatus == "g" || entry.TableEditStatus == "п" {
		return nil, fmt.Errorf("ticket have DONE status")
	}
	if entry.TableEditStatus == "" {
		return nil, fmt.Errorf("ticket have NO status")
	}
	if entry.TablePath != "" {
		return nil, fmt.Errorf("ticket have PATH")
	}

	tkt.Contragent = entry.TableAgent
	tkt.Category = FILM
	tkt.Name = translit.TransliterateLower(entry.TableTaskName)
	tkt.TableData = []string{
		entry.TableTaskName,
		entry.Desination,
	}
	tkt.TimeOpen = time.Now()

	return &tkt, nil
}

func SerialFromEntry(entry tabledata.TableEntry) (*Ticket, error) {
	tkt := Ticket{}
	if entry.TableAgent == "" {
		return nil, fmt.Errorf("ticket must have contragent")
	}
	tableName := translit.TransliterateLower(entry.TableTaskName)
	//fmt.Println("try", tableName)
	if !strings.Contains(tableName, translit.TransliterateLower("сезон")) {
		return nil, fmt.Errorf("ticket have no serial numbers")
	}

	if entry.TableEditStatus == "v" || entry.TableEditStatus == "м" {
		return nil, fmt.Errorf("ticket have DOWNLOADING status")
	}
	if entry.TableEditStatus == "g" || entry.TableEditStatus == "п" {
		return nil, fmt.Errorf("ticket have DONE status")
	}
	if entry.TableEditStatus == "" {
		return nil, fmt.Errorf("ticket have NO status")
	}
	if entry.TablePath != "" {
		return nil, fmt.Errorf("ticket have PATH")
	}

	if strings.Contains(tableName, translit.TransliterateLower("серия")) || strings.Contains(tableName, translit.TransliterateLower("серий")) {
		//fmt.Println("defines series HERE")
	}

	tkt.Contragent = entry.TableAgent
	tkt.Category = SER
	tkt.Name = translit.TransliterateLower(entry.TableTaskName)
	tkt.TableData = []string{
		entry.TableTaskName,
		entry.Desination,
	}
	tkt.TimeOpen = time.Now()

	return &tkt, nil
}

func episodes(str string) []int {
	//`serial_pro_kartoshku_02_sezon_04_serii`
	episodes := []int{}
	if strings.Contains(str, "serii_") {
		words := strings.Split(str, "_")
		val := -1
		for _, w := range words {
			v, err := strconv.Atoi(w)
			if err == nil {
				val = v
			}
		}
		for i := 1; i <= val; i++ {
			episodes = append(episodes, i)
		}
		return episodes
	}
	if strings.Contains(str, "seriya_") {
		words := strings.Split(str, "_")
		val := -1
		for _, w := range words {
			v, err := strconv.Atoi(w)
			if err == nil {
				val = v
			}
		}
		episodes = append(episodes, val)
		return episodes
	}
	return episodes

}
