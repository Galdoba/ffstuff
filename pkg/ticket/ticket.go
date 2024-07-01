package ticket

import (
	"fmt"
	"time"
)

const (
	DEFAULT           = iota
	TYPE_TRL          = "--TRL--"
	TYPE_FILM         = "--FILM--"
	TYPE_SER          = "--SER--"
	SOURCE_DATA       = "*source_profile"
	PROCESS_DATA      = "*Process data"
	PROCESS_REQUEST   = "*Process request"
	ARCHIVE_DEST      = "path_archive"
	PROCESS_DEST      = "path_process"
	SEASON_NUM        = "season "
	EPISODE_NUM       = "episode"
	EPISODES_EXPECTED = "expected_episodes"
	SEASON_EXPECTED   = "season"
)

type Ticket struct {
	Name        string                       `json:"Table Name"`
	TicketType  string                       `json:"Type"`
	BaseName    string                       `json:"Output Base,omitempty"`
	Agent       string                       `json:"Agent,omitempty"`
	Date        string                       `json:"Deliver Date,omitempty"`
	StartTime   time.Time                    `json:"Started at"`
	CloseTime   string                       `json:"Closed at,omitempty"`
	Info_Tags   map[string]map[string]string `json:"Info,omitempty"`                 //карта [имя_файла]набор_тэгов
	Processable bool                         `json:"Ready for processing,omitempty"` //карта [имя_файла]набор_тэгов
}

func New(name, tType string) *Ticket {
	tk := Ticket{}
	tk.Name = name
	tk.TicketType = tType
	tk.StartTime = time.Now()
	tk.Info_Tags = make(map[string]map[string]string)
	tk.Info_Tags[PROCESS_DATA] = make(map[string]string)
	return &tk
}

func (tk *Ticket) AddSource(source, mfdata string) error {
	if _, ok := tk.Info_Tags[source]; ok {
		return fmt.Errorf("source '%v' exists", source)
	}
	tk.Info_Tags[source] = make(map[string]string)
	tk.Info_Tags[source][SOURCE_DATA] = mfdata
	return nil
}

//AddTag - добавляет тэг с информацией
//менять и удалять тэг нельзя
func (tk *Ticket) AddTag(source, key, val string) error {
	if _, ok := tk.Info_Tags[source]; !ok {
		return fmt.Errorf("source '%v' was not added", source)
	}
	sourcemap := tk.Info_Tags[source]

	if _, ok := sourcemap[key]; ok {
		return fmt.Errorf("tag '%v' for source '%v' exists", key, source)
	}
	tk.Info_Tags[source][key] = val
	return nil
}

//AddRequest - задел на будущее
//менять и удалять нельзя
func (tk *Ticket) AddRequest(key, val string) error {
	if _, ok := tk.Info_Tags[PROCESS_REQUEST]; !ok {
		tk.Info_Tags[PROCESS_REQUEST] = make(map[string]string)
	}
	requestmap := tk.Info_Tags[PROCESS_REQUEST]

	if _, ok := requestmap[key]; ok {
		return fmt.Errorf("request '%v' exists", key)
	}
	tk.Info_Tags[PROCESS_REQUEST][key] = val
	return nil
}

func (tk *Ticket) ValidateWith(validator Validator) error {
	tk.Processable = false
	if _, ok := tk.Info_Tags[PROCESS_DATA]; !ok {
		return fmt.Errorf("process data absent")
	}
	procdata := tk.Info_Tags[PROCESS_DATA]
	if _, ok := procdata[PROCESS_DEST]; !ok {
		return fmt.Errorf("process destination absent")
	}
	if _, ok := procdata[ARCHIVE_DEST]; !ok {
		return fmt.Errorf("archive destination absent")
	}
	switch tk.TicketType {
	default:
		return fmt.Errorf("ticket type invalid")
	case TYPE_SER:
		if _, ok := procdata[SEASON_NUM]; !ok {
			return fmt.Errorf("season number absent")
		}
		if _, ok := procdata[EPISODE_NUM]; !ok {
			return fmt.Errorf("episode number absent")
		}
		if !validator.ValidateAgent(procdata[EPISODES_EXPECTED]) {
			return fmt.Errorf("invalid expected episodes: [%v]", procdata[EPISODES_EXPECTED])
		}
	case TYPE_FILM:
		if _, ok := procdata[SEASON_NUM]; ok {
			return fmt.Errorf("season number not expected")
		}
		if _, ok := procdata[EPISODE_NUM]; !ok {
			return fmt.Errorf("episode number not expected")
		}
	}
	if !validator.ValidateAgent(tk.Agent) {
		return fmt.Errorf("invalid agent: %v", tk.Agent)
	}
	if !validator.ValidateProcessPath(procdata[PROCESS_DEST]) {
		return fmt.Errorf("invalid processing destination: %v", procdata[PROCESS_DEST])
	}
	if !validator.ValidateArchivePath(procdata[ARCHIVE_DEST]) {
		return fmt.Errorf("invalid archive destination: %v", procdata[ARCHIVE_DEST])
	}
	tk.Processable = true
	return nil
}

//Validator - ожидаем реализацию трех валидаторов ФИЛЬМ, СЕРИАЛ, ТРЕЙЛЕР
type Validator interface {
	ValidateAgent(string) bool
	ValidateProcessPath(string) bool
	ValidateArchivePath(string) bool
	ValidateEpisodes(string) bool
}

type pseudoValidator struct{}

func (pv *pseudoValidator) ValidateAgent(s string) bool {
	return true
}
func (pv *pseudoValidator) ValidateProcessPath(s string) bool {
	return true
}
func (pv *pseudoValidator) ValidateArchivePath(s string) bool {
	return true
}
func (pv *pseudoValidator) ValidateEpisodes(s string) bool {
	return true
}

// type requester struct {
// 	conditions     map[string]string
// 	key            string
// 	val            string
// 	validationFunc func(map[string]string) bool
// }

// func (rq *requester) Condition() map[string]string {
// 	return rq.conditions
// }

// func (rq *requester) Validation() func(map[string]string) bool {
// 	return rq.validationFunc
// }

// func (rq *requester) Request(map[string]string) (string, string) {
// 	if rq.validationFunc(rq.Condition()) {
// 		return rq.key, rq.val
// 	}
// 	return "", ""
// }

// type Requester interface {
// 	Condition() map[string]string
// 	ValidationFunc() func(map[string]string) bool
// 	Request() (string, string)
// }

// func (tk *Ticket) RequestFrom(r Requester) (string, string) {
// 	for k, sourcemap := range tk.Info_Tags {
// 		switch k {
// 		case PROCESS_DATA, PROCESS_REQUEST:
// 			continue
// 		default:
// 			condLen := len(r.Condition())
// 			condMet := 0
// 			conditions := r.Condition()
// 			for haveKey, haveVal := range sourcemap {
// 				for condKey, condVal := range conditions {
// 					if condKey == haveKey && condVal != haveVal {
// 						condMet++
// 					}
// 				}
// 			}
// 		}
// 	}
// 	if r.ValidationFunc()(r.Condition()) {
// 		return r.Request()
// 	}
// 	return "", ""
// }

// func testSizeValidation(actial map[string]string) bool {
// 	key := "size"
// 	if _, ok := actial[key]; !ok {
// 		return false
// 	}
// 	requesters := make(map[string]string)
// 	expectedW := 1920
// 	expectedH := 1080
// 	switch requesters[key] {
// 	case "HD":
// 		expectedW = 1920
// 		expectedH = 1080
// 	}

// 	whStrSl := strings.Split(actial[key], "x")
// 	w, _ := strconv.Atoi(whStrSl[0])
// 	h, _ := strconv.Atoi(whStrSl[1])
// 	if w > expectedW || (w == expectedW && h < expectedH) {
// 		return true
// 	}
// 	return false
// }
