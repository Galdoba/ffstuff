package ticket

import (
	"fmt"
	"regexp"
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

func CategoryPrefix(i int) string {
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
	Name             string    `json:"Name"`                         //uname
	Category         int       `json:"Category"`                     //FILM-TRL-SER
	Contragent       string    `json:"Contragent"`                   // заказчик (Amedia или иное)
	SourceFiles      []string  `json:"Source Files"`                 // файлы относящиеся к тикету
	SourcePrefix     string    `json:"Source File Prefix,omitempty"` // префикс рабочих файлов
	Season           int       `json:"Season Number,omitempty"`      //Номер сезона если сериал
	Episodes         []int     `json:"Episode Numbers,omitempty"`    //Номер эпизода если сериал
	EpisodeDefined   bool      `json:"Episode Defined,omitempty"`    //Номер эпизода если сериал
	Register         string    `json:"Register,omitempty"`           //Регистр серии (например 14b) если сериал
	SourceNamesCheck int       `json:"Source Names Status"`          //статус проверки
	BasicCheck       int       `json:"Basic Scan Status"`            //статус проверки
	InterlaceCheck   int       `json:"Interlace Scan Status"`        //статус проверки
	LoudnessCheck    int       `json:"Loudness Scan Status"`         //статус проверки
	ReadWriteCheck   int       `json:"ReadWrite Scan Status"`        //статус проверки
	TimeOpen         time.Time `json:"Opened"`                       //время открытия тикета
	TimeClose        time.Time `json:"Closed"`                       //время закрытия тикета
	IsClosed         bool      `json:"IsClosed"`                     //статус закрытия тикета
	TableData        []string  `json:"Table Data"`                   //данные из таблицы
	BaseWords        []string  `json:"Base Words"`                   //слова в имени файлов для авто определения
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
	// if entry.TablePath != "" {
	// 	return nil, fmt.Errorf("ticket have PATH")
	// }
	base, season, episodes := guessSerialName(entry.TableTaskName)
	//fmt.Println("guessed", base, season, episodes)
	tkt.Season = season
	tkt.Episodes = episodes
	tkt.BaseWords = base
	// if strings.Contains(tableName, translit.TransliterateLower("серия")) || strings.Contains(tableName, translit.TransliterateLower("серий")) {
	// 	//fmt.Println("defines series HERE")
	// }

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

// func episodes(str string) []int {
// 	//`serial_pro_kartoshku_02_sezon_04_serii`
// 	episodes := []int{}
// 	if strings.Contains(str, "serii_") {
// 		words := strings.Split(str, "_")
// 		val := -1
// 		for _, w := range words {
// 			v, err := strconv.Atoi(w)
// 			if err == nil {
// 				val = v
// 			}
// 		}
// 		for i := 1; i <= val; i++ {
// 			episodes = append(episodes, i)
// 		}
// 		return episodes
// 	}
// 	if strings.Contains(str, "seriya_") {
// 		words := strings.Split(str, "_")
// 		val := -1
// 		for _, w := range words {
// 			v, err := strconv.Atoi(w)
// 			if err == nil {
// 				val = v
// 			}
// 		}
// 		episodes = append(episodes, val)
// 		return episodes
// 	}
// 	return episodes

// }

func guessSerialName(name string) ([]string, int, []int) {
	base := []string{}
	seasonNum := -1
	episodes := []int{}
	if haltSerialGuess(name) {
		return base, seasonNum, episodes
	}
	seasonNum = grabSeason(name)
	switch seasonNum {
	case -1:
		//base = "can not guess"
	case 0:
	default:
		episodes = grabEpisodes(name)
		data := strings.Split(name, numToString(seasonNum)+" сезон")
		if len(data) < 2 {
			//	base = "ERROR"
			return base, seasonNum, episodes
		}
		baseWord := ensureTitle(translit.TransliterateLower(data[0]))
		baseWord = strings.TrimSuffix(baseWord, "_")
		base = strings.Split(baseWord, "_")
	}

	return base, seasonNum, episodes
}

func ensureTitle(s string) string {
	letters := strings.Split(s, "")
	if len(letters) > 0 {
		letters[0] = strings.ToUpper(letters[0])
	}
	return strings.Join(letters, "")
}

func haltSerialGuess(n string) bool {
	if strings.Contains(n, "(Трейлер)") {
		return true
	}
	if strings.Contains(n, "(Замена трейлера)") {
		return true
	}
	return false
}

func grabSeason(input string) int {

	r := regexp.MustCompile(` [0-9]?[0-9]?[0-9]? сезон`)
	smatch := r.FindString(input)
	if smatch == "" {
		return -1
	}
	data := strings.Split(smatch, " ")
	num := 0

	for _, d := range data {
		n, err := strconv.Atoi(d)
		if err != nil {
			continue
		}
		num = n

	}

	return num
}

func grabEpisodes(input string) []int {
	episodes := []int{}
	smatch := ""
	one := true
	for _, word := range []string{"серия", "серии", "серий"} {
		expression := `[0-9]?[0-9]?[0-9]? `
		r := regexp.MustCompile(expression + fmt.Sprintf("%v", word))
		smatch = r.FindString(input)
		if smatch == "" {
			one = false
			continue
		}
		break
	}
	if smatch == "" {
		for i := 1; i <= 300; i++ {
			episodes = append(episodes, i)
		}
		return episodes
	}
	data := strings.Split(smatch, " ")
	num := 0
	for _, d := range data {
		n, err := strconv.Atoi(d)
		if err != nil {
			continue
		}
		num = n

	}
	if !one {
		for i := 1; i < num; i++ {
			episodes = append(episodes, i)
		}
	}
	episodes = append(episodes, num)
	return episodes
}

func numToString(i int) string {
	s := fmt.Sprintf("%v", i)
	for len(s) < 2 {
		s = "0" + s
	}
	return s
}

func (tkt *Ticket) SplitByEpisodes() ([]*Ticket, error) {
	if tkt.Category != SER {
		return nil, fmt.Errorf("ticket is not a serial")
	}
	if len(tkt.Episodes) < 1 {
		return nil, fmt.Errorf("ticket episodes unawailable")
	}
	if len(tkt.BaseWords) < 1 {
		return nil, fmt.Errorf("ticket basewords unawailable")
	}
	splittedTickets := []*Ticket{}
	for _, e := range tkt.Episodes {
		newTicket := Ticket{}
		base := strings.Join(tkt.BaseWords, "_")
		sertag := "s" + numToString(tkt.Season) + "e" + numToString(e)
		newTicket.Name = ensureTitle(base) + "_" + sertag
		newTicket.Season = tkt.Season
		newTicket.Category = tkt.Category
		newTicket.Contragent = tkt.Contragent
		newTicket.BaseWords = tkt.BaseWords
		newTicket.Episodes = append([]int{}, e)
		newTicket.EpisodeDefined = true
		newTicket.SourcePrefix = newTicket.Name
		newTicket.SourceNamesCheck = Check_Status_WAIT
		newTicket.TimeOpen = time.Now()
		splittedTickets = append(splittedTickets, &newTicket)
	}
	return splittedTickets, nil
}

func (t *Ticket) EpisodeTag() string {
	if t.Season < 1 {
		return ""
	}
	if len(t.Episodes) != 1 {
		return "s" + numToString(t.Season)
	}
	return "s" + numToString(t.Season) + "e" + numToString(t.Episodes[0])
}
