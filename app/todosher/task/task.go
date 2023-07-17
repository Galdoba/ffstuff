package task

import (
	"fmt"
	"strconv"
	"time"
)

const (
	KEY_agent    = "agent"
	KEY_director = "director"
	KEY_category = "category"
	KEY_title    = "title"
	KEY_descr    = "descr"
	//KEY_created         = "created"
	KEY_start      = "start_after"
	KEY_deadline   = "deadline"
	KEY_comment    = "comment"
	KEY_importance = "importance"
	KEY_private    = "private"
)

type tsk struct {
	//id              int64
	director    string
	agent       string
	category    string
	title       string
	descr       string
	created     time.Time
	start_after time.Time
	deadline    time.Time
	comment     string
	importance  int
	private     bool
}

func NewTask(inputs ...*inputpair) (*tsk, error) {
	tm := time.Now()
	ts := tsk{}
	//ts.id = tm.UnixNano()
	ts.created = tm
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
	return &ts, nil
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

func (ts *tsk) Change(input *inputpair) error {
	if ts.private && ts.director != "" && (ts.director != input.from) {
		return fmt.Errorf("can not change private task of different agent")
	}
	switch input.key {
	default:
		return fmt.Errorf("can not change task '%v': key '%v' unrecognized", ts.title, input.key)
	case KEY_director:
		if ts.director != "" {
			return fmt.Errorf("can not change director if already set")
		}
		ts.director = input.val
	case KEY_agent:
		ts.agent = input.val
	case KEY_category:
		ts.category = input.val
	case KEY_title:
		ts.title = input.val
	case KEY_descr:
		ts.descr = input.val
	case KEY_start:
		tm, err := time.Parse("2006-01-02 15:04:05", input.val)
		if err != nil {
			return fmt.Errorf("can not change start: %v", err.Error())
		}
		ts.start_after = tm
	case KEY_deadline:
		tm, err := time.Parse("2006-01-02 15:04:05", input.val)
		if err != nil {
			return fmt.Errorf("can not change deadline: %v", err.Error())
		}
		ts.deadline = tm
	case KEY_comment:
		ts.comment = input.val
	case KEY_importance:
		switch input.val {
		default:
			return fmt.Errorf("can not change importance: '%v' value unknown", input.val)
		case "1", "2", "3", "4", "5", "6", "7":
			ts.importance, _ = strconv.Atoi(input.val)
		}
	case KEY_private:
		v, err := strconv.ParseBool(input.val)
		if err != nil {
			return fmt.Errorf("can not change privacy: %v", err.Error())
		}
		ts.private = v
	}
	return nil
}

/*
IMPORTANCE:
1 Критично			Critical		Red
2 Очень Важно		Very Important	Magenta
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
// func (ts *tsk) ID() int64 {
// 	return ts.id
// }
func (ts *tsk) Director() string {
	return ts.director
}
func (ts *tsk) Agent() string {
	return ts.agent
}
func (ts *tsk) Category() string {
	return ts.category
}
func (ts *tsk) Title() string {
	return ts.title
}
func (ts *tsk) Descr() string {
	return ts.descr
}
func (ts *tsk) Created() time.Time {
	return ts.created
}
func (ts *tsk) StartTime() time.Time {
	return ts.start_after
}
func (ts *tsk) Deadline() time.Time {
	return ts.deadline
}
func (ts *tsk) Comment() string {
	return ts.comment
}
func (ts *tsk) Importance() int {
	return ts.importance
}
func (ts *tsk) IsPrivate() bool {
	return ts.private
}

// func (ts *tsk) SetFrom(val string) {
// 	ts.from = val
// }
func (ts *tsk) SetAgent(val string) {
	ts.agent = val
}
func (ts *tsk) SetCategory(val string) {
	ts.category = val
}
func (ts *tsk) SetTitle(val string) {
	ts.title = val
}
func (ts *tsk) SetDescr(val string) {
	ts.descr = val
}
func (ts *tsk) SetCreated(val time.Time) {
	ts.created = val
}
func (ts *tsk) Setstart_after(val time.Time) {
	ts.start_after = val
}
func (ts *tsk) SetDeadline(val time.Time) {
	ts.deadline = val
}
func (ts *tsk) SetComment(val string) {
	ts.comment = val
}
func (ts *tsk) SetImportance(val int) {
	ts.importance = val
}
func (ts *tsk) SetPrivate(val bool) {
	ts.private = val
}
