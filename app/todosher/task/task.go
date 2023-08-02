package task

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

const (
	KEY_receiver = "Receiver"
	KEY_sender   = "Sender"
	KEY_category = "Category"
	KEY_title    = "Title"
	KEY_descr    = "Descr"
	//KEY_created         = "Created"
	KEY_start      = "Start_after"
	KEY_deadline   = "Deadline"
	KEY_comment    = "Comment"
	KEY_importance = "Importance"
	KEY_private    = "Private"
	timeformat1    = "2006-01-02 15:04:05"
	timeformat2    = "2006-01-02 15:04"
	timeformat3    = "2006-01-02"
	timeformat4    = "15:04:05"
	timeformat5    = "15:04"

	defaultVal = iota
	timing_irrelevant
	timing_active
	timing_missed
)

type Task struct {
	//id              int64 `xml:"name,omitempty"`
	MXLName     xml.Name `xml:"task"`                     //задача
	Sender      string   `xml:"sender"`                   //тот кто поставил задачу (юзер или программа)
	Receiver    string   `xml:"receiver"`                 //тот кто должен выполнить задачу (юзер или программа)
	Category    string   `xml:"category,omitempty"`       //
	Title       string   `xml:"title"`                    //краткое описание задачи
	Descr       string   `xml:"descr,omitempty"`          //описываем что как и зачем делать для человека
	ResolveArgs []string `xml:"args,omitempty"`           //описываем что как и зачем делать для машины
	Created     string   `xml:"created"`                  //Дата создания файла
	Start_after string   `xml:"relevant after,omitempty"` //отметка времени после которого задача считается актуальной
	Deadline    string   `xml:"deadline,omitempty"`       //отметка времени после которого задача перестает быть актуальной
	Comment     string   `xml:"comment,omitempty"`        //дополнительный коментарий в любой форме
	Importance  int      `xml:"importance,omitempty"`     //степень важности
	timing      int
}

/*
Importance:
0 неопределено
1 Критично			Critical		Red
2 Очень Важно		Very Important	Magenta
3 Важно				Important		Yellow
4 Средней Важности	Average			Green
5 Рутина			Routine			White/Gray
6 Не важно			Unimportant		Cyan
7 Факультативно		Optional		Blue
*/

func NewTask(inputs ...*inputpair) (*Task, error) {
	tm := time.Now()
	ts := Task{}
	//ts.id = tm.UnixNano()
	ts.Created = timeFormatted(tm)
	errors := []error{}
	for _, input := range inputs {
		if err := ts.Change(input); err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		errText := "task initiation errors:"
		for _, err := range errors {
			errText += "\n" + err.Error()
		}
		return &ts, fmt.Errorf(errText)
	}
	ts.evaluateTiming()
	return &ts, nil
}

func FromFile(file string) (*Task, error) {
	ts := Task{}
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadFile(file): %v", err.Error())
	}
	err := xml.Unmarshal(f, &tsk2)
}

// func (ts *Task) String() string {
// 	str := fmt.Sprintf("From: %v\n", ts.Sender)
// 	str += fmt.Sprintf("To: %v\n", ts.Receiver)
// 	str += fmt.Sprintf("Task: %v\n", ts.Title)
// 	str += fmt.Sprintf("Description: %v\n", ts.Descr)
// 	if ts.Deadline == "" {
// 		str += fmt.Sprintf("Status: %v\n", ts.timing)
// 	}
// 	return str
// }

func (ts *Task) evaluateTiming() {

	tn := time.Now()
	//ts := time.Parse(timeformat1, ts.Created)

	startAfter, errS := time.Parse(timeformat1, ts.Start_after)

	deadline, errA := time.Parse(timeformat1, ts.Deadline)

	if errS == nil && !tn.After(startAfter) {
		ts.timing = timing_irrelevant
		return
	}
	if errA == nil && tn.After(deadline) {
		ts.timing = timing_missed
		return
	}
	ts.timing = timing_active

}

type inputpair struct {
	key  string
	val  string
	from string
}

func Input(key, val string) *inputpair {
	input := inputpair{}
	input.key = key
	input.val = val
	return &input
}

func (inp *inputpair) From(from string) *inputpair {
	inp.from = from
	return inp
}

func (ts *Task) Change(input *inputpair) error {
	if ts.Sender != "" && (ts.Sender != input.from) {
		return fmt.Errorf("can not change task of different Sender")
	}
	switch input.key {
	default:
		return fmt.Errorf("can not change task '%v': key '%v' unrecognized", ts.Title, input.key)
	case KEY_sender:
		if ts.Sender != "" {
			return fmt.Errorf("can not change Sender if already set")
		}
		ts.Sender = input.val
	case KEY_receiver:
		ts.Receiver = input.val
	case KEY_category:
		ts.Category = input.val
	case KEY_title:
		ts.Title = input.val
	case KEY_descr:
		ts.Descr = input.val
	case KEY_start:
		_, err := time.Parse("2006-01-02 15:04:05", input.val)
		if err != nil {
			return fmt.Errorf("can not change start: %v", err.Error())
		}
		ts.Start_after = input.val
	case KEY_deadline:
		_, err := time.Parse("2006-01-02 15:04:05", input.val)
		if err != nil {
			return fmt.Errorf("can not change Deadline: %v", err.Error())
		}
		ts.Deadline = input.val
	case KEY_comment:
		ts.Comment = input.val
	case KEY_importance:
		switch input.val {
		default:
			return fmt.Errorf("can not change Importance: '%v' value unknown", input.val)
		case "1", "2", "3", "4", "5", "6", "7":
			ts.Importance, _ = strconv.Atoi(input.val)
		}
	}
	return nil
}

func formats() []string {
	return []string{
		timeformat1,
		timeformat2,
		timeformat3,
		timeformat4,
		timeformat5,
	}
}

func timeFormatted(t time.Time) string {
	for _, f := range []string{timeformat1, timeformat2, timeformat3} {
		str := t.Format(f)
		if str != "" {
			return str
		}
	}
	return ""
}

/*
Importance:
1 Критично			Critical		Red
2 Очень Важно		Very Important	Mreceivera
3 Важно				Important		Yellow
4 Средней Важности	Average			Green
5 Рутина			Routine			White/Gray
6 Не важно			Unimportant		Cyan
7 Факультативно		Optional		Blue
*/

// type progress struct {
// 	done   int
// 	steps  int
// 	status string
// }

// func newProgress(steps int) *progress {
// 	pr := progress{}
// 	pr.steps = steps
// 	return &pr
// }

// func (pr *progress) Bar() string {

// 	//12345678901234567890
// 	//100% [+++++++++---]

// /*
// 10/10
// 9/10
// 9/9
// 100%
// */

// 	return ""
// }

///////////////////////////////////////////////////
// func (ts *Task) ID() int64 {
// 	return ts.id
// }
