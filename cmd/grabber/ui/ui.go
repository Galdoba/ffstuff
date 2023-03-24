package ui

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Galdoba/ffstuff/cmd/grabber/download"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

// func main() {
// 	inputinfo.CleanScanData()
// }

type allProc struct {
	//procs  []*proc
	stream     []*stream
	globalStop bool
	//endEvent bool
	//cursorSelection int
}

const (
	s1 = `d:\IN\IN_2022-05-11\tobot_s01_02_2010__hd_rus20.m4a`
	d1 = `d:\IN\IN_2022-05-11\COPY_tobot_s01_02_2010__hd.mp4`
	s2 = `d:\IN\IN_2022-05-11\tobot_s01_04_2010__hd_rus20.m4a`
	d2 = `d:\IN\IN_2022-05-11\COPY_tobot_s01_04_2010__hd.mp4`
	s3 = `d:\IN\IN_2022-05-11\tobot_s01_03_2010__hd_rus20.m4a`
	d3 = `d:\IN\IN_2022-05-11\COPY_tobot_s01_03_2010__hd.mp4`
)

func testFiles() []string {
	return []string{
		`d:\IN\IN_2022-05-11\tobot_s01_02_2010__hd.mp4`,
		`d:\IN\IN_2022-05-11\tobot_s01_03_2010__hd.mp4`,
		`d:\IN\IN_2022-05-11\tobot_s01_04_2010__hd.mp4`,
		`d:\IN\IN_2022-05-11\tobot_s01_06_2010__hd.mp4`,
		`d:\IN\IN_2022-05-11\tobot_s01_07_2010__hd.mp4`,
		`d:\IN\IN_2022-05-11\tobot_s01_12_2010__hd.mp4`,
		`d:\IN\IN_2022-05-11\tobot_s01_16_2010__hd.mp4`,
		`d:\IN\IN_2022-05-11\tobot_s01_22_2010__hd.mp4`,
	}
}

type stream struct {
	source       string
	dest         string
	progress     int64
	expected     int64
	handler      download.Handler
	lastResponse string
	lastCommand  string
	isSelected   bool
}

func (ap *allProc) newStream(source, dest string) {
	ap.stream = append(ap.stream, &stream{source, dest, 0, 0, nil, "NONE", "NONE", false})
}

func (st *stream) start() error {
	h, e := download.StartNew(st.source, st.dest)
	if e != nil {
		return e
	}
	st.handler = h
	return nil
}

func (st *stream) String() string {
	str := "["
	switch st.isSelected {
	case true:
		str += "x]"
	case false:
		str += " ]"
	}
	str += " " + st.source + "|" + st.lastCommand + "|" + st.lastResponse
	return str
}

// func (ap *allProc) String() string {
// 	str := ""
// 	for _ := range ap.stream {
// 		str += fmt.Sprintf("[%v] %v\n", selected, s.lastResponse)
// 	}
// 	return str
// }

func newAllProcs() *allProc {
	ap := allProc{}
	for _, tf := range testFiles() {
		ap.newStream(tf, tf+"COPY.m4a")
	}
	return &ap
}

func (ap *allProc) update() {
	if len(ap.stream) == 0 {
		//	ap.endEvent = true
		return
	}

	// switch ap.stream[0].done {
	// case false:
	// 	ap.procs[0].run()
	// case true:
	// 	if len(ap.procs) > 1 {
	// 		ap.procs = ap.procs[1:]
	// 		//ap.cursorSelection--
	// 		ap.sync()
	// 	} else {
	// 		ap.procs = nil
	// 		//ap.endEvent = true
	// 		return
	// 	}
	// }
}

func (ap *allProc) bumpToTop(i int) {
	//func namesort.BumpToTopIndex(slInt []int, index int) []int
	if i < 1 || i > len(ap.stream)-1 {
		return
	}
	current := i
	for current > 0 {
		ap.bumpUpByOne(current)
		current--
	}
}

func (ap *allProc) bumpUpByOne(i int) {
	//func namesort.BumpToTopIndex(slInt []int, index int) []int
	if i < 1 || i > len(ap.stream)-1 {
		return
	}
	ap.stream[i-1], ap.stream[i] = ap.stream[i], ap.stream[i-1]
}

func prc(i int) string {
	n := fmt.Sprintf("%v ", i)
	for len(n) < 4 {
		n = " " + n
	}
	return fmt.Sprintf("[%v", n) + "%]"
}

type InfoBox struct {
	data   []string
	cursor int
}

func (ib *InfoBox) Draw(ac *allProc) {
	fg := termbox.ColorWhite
	bg := termbox.ColorBlack
	tbprint(0, 1, fg, bg, "Grabber Dowloading:")
	for i, data := range ib.data {

		if i == ib.cursor {
			fg = termbox.ColorBlack
			bg = termbox.ColorWhite
		}
		tbprint(2, i+2, fg, bg, data)
		fg = termbox.ColorWhite
		bg = termbox.ColorBlack
	}
}

// This function is often useful:
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {

	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func (ib *InfoBox) Update(ac *allProc) error {
	newData := []string{}
	for _, pr := range ac.stream {
		newData = append(newData, pr.String())
	}
	ib.data = newData
	//ib.cursor = ac.cursorSelection

	return nil
}

func SetupProcesses(dest string, paths ...string) *allProc {
	ap := allProc{}

	for _, path := range paths {
		file := filepath.Base(path)
		ap.newStream(path, dest+file)
	}
	return &ap
}

func SetupInfoBox() *InfoBox {
	return &InfoBox{}
}

type action struct {
	eventName  string
	key        string
	termbxKey  termbox.Key
	termbxRune rune
	function   func(*allProc, *InfoBox)
}

func setupAction(eventName, key string, function func(*allProc, *InfoBox)) (action, error) {
	act := action{eventName: eventName, key: key, function: function}
	runeArray := []rune(key)
	if len(runeArray) != 1 || runeArray[0] == 0 {
		act.termbxKey = commandSequanceToTBKey(key)
	} else {
		act.termbxRune = runeArray[0]
	}
	// act.termbxKey = 20000
	// switch act.key {
	// case "UP":
	// 	act.termbxKey = termbox.KeyArrowUp
	// 	act.termbxRune = 0
	// case "DOWN":
	// 	act.termbxKey = termbox.KeyArrowDown
	// 	act.termbxRune = 0
	// case "SPACE":
	// 	act.termbxKey = termbox.KeySpace
	// 	act.termbxRune = 0
	// }
	if act.termbxKey == 20000 {
		return act, fmt.Errorf("inpossible termboxKey")
	}
	return act, nil
}

type ActionPool struct {
	acmap       map[string]action
	byTermbxKey map[termbox.Key]action
	byRune      map[rune]action
}

type Action interface {
	commence(string) error
}

func (actpl *ActionPool) fillCommandActionMap(configMap map[string]string) error {
	actpl.acmap = make(map[string]action)
	actpl.byTermbxKey = make(map[termbox.Key]action)
	actpl.byRune = make(map[rune]action)

	/////////////////////////////////
	moveCursorUp, err := setupAction("MOVE_CURSOR_UP", configMap["MOVE_CURSOR_UP"], Action_MoveSelectionUP)
	if err != nil {
		return fmt.Errorf("can not setup 'moveCursorUp' action: %v", err.Error())
	}
	actpl.acmap["MOVE_CURSOR_UP"] = moveCursorUp
	/////////////////////////////////
	moveCursorDown, err := setupAction("MOVE_CURSOR_DOWN", configMap["MOVE_CURSOR_DOWN"], Action_MoveSelectionDOWN)
	if err != nil {
		return fmt.Errorf("can not setup 'moveCursorDown' action: %v", err.Error())
	}
	actpl.acmap["MOVE_CURSOR_DOWN"] = moveCursorDown
	/////////////////////////////////
	toggleSelectionState, err := setupAction("TOGGLE_SELECTION_STATE", configMap["TOGGLE_SELECTION_STATE"], Action_ToggleSelection)
	if err != nil {
		return fmt.Errorf("can not setup 'toggleSelectionState' action: %v", err.Error())
	}
	actpl.acmap["TOGGLE_SELECTION_STATE"] = toggleSelectionState
	/////////////////////////////////
	/*TODO: Прописать действия
	переместить выделенное по списку вверх на 1 позицию
	переместить выделенное по списку вниз на 1 позицию
	переместить выделенное по списку вверх до предела
	переместить выделенное по списку вниз до предела
	удаление из очереди
	активная пауза (для всех процессов)
	включить/выключить ограниченную скорость закачки
	сброс всех выделений
	инсерт (переключение выделение со сдвигом курсора вниз на 1 позицию)
	*/

	for k, a := range actpl.acmap {
		if a.termbxRune == 0 {
			actpl.byTermbxKey[a.termbxKey] = actpl.acmap[k]
		}
	}

	return nil
}

func StartMainloop(configMap map[string]string, paths []string) error {
	ap := SetupProcesses(configMap["dest"], paths...)
	ib := &InfoBox{}
	ib.data = []string{}
	actionPool := ActionPool{}
	if err := actionPool.fillCommandActionMap(configMap); err != nil {
		return err
	}

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	//termbox.SetInputMode(termbox.InputEsc)
	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()
	//ib.Draw(ap)
	//termbox.Flush()
	//termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	//ap.update()
	//ib.Update(ap)
	//ib.Draw(ap)
	draw_tick := time.NewTicker(100 * time.Millisecond)
	handlerEvents := make(chan download.Response)

loop:
	for {
		if len(ap.stream) == 0 {
			break
		}
		if ap.stream[0].handler == nil {
			err := ap.stream[0].start()
			if err != nil {
				panic("start dowload stream: " + err.Error())
			}
			handlerEvents = ap.stream[0].handler.Listen()
		}

		//ib.Update(ap)
		// if len(ap.procs) == 0 {
		// 	ib.Update(ap)
		// 	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
		// 	tbprint(0, 1, 0, 0, "DONE")
		// 	termbox.Flush()
		// 	fmt.Println("All GRABBED!")
		// 	break
		// }

		select {
		case ev := <-event_queue:

			switch ev.Type {

			case termbox.EventKey:
				//panic(eventKeyToString(ev))
				switch ev.Ch {

				case 0:
					if action, ok := actionPool.byTermbxKey[ev.Key]; ok {
						action.function(ap, ib)
					}

					// switch ev.Key {
					// case termbox.KeyEsc:
					// 	//os.Exit(1)
					// 	break loop
					// case termbox.KeyArrowDown:
					// 	Action_MoveSelectionDOWN(ap, ib)
					// case termbox.KeyArrowUp:
					// 	Action_MoveSelectionUP(ap, ib)
					// case termbox.KeyEnter:
					// 	//ap.bumpToTop(ib.cursor)
					// 	//Action_BumpStreamToTOP(ap, ib)
					// case termbox.KeySpace:
					// 	Action_ToggleSelection(ap, ib)
					// }
				case 'p':
					//Action_TogglePause(ap, ib)
				case 'q':
					break loop
				}

			}
		case <-draw_tick.C:
		case ev := <-handlerEvents:
			ap.stream[0].lastResponse = ev.String()

			if ev.String() == "completed" {
				Action_StartNext(ap, ib)
				handlerEvents = make(chan download.Response)
			}

		}
		ap.update()
		ib.Update(ap)
		termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
		ib.Draw(ap)
		termbox.Flush()
		//time.Sleep(time.Second * 2)
	}
	fmt.Println("END")
	return nil
}

// func (ap *allProc) Options() []string {
// 	opt := []string{}
// 	for i := range ap.procs {
// 		opt = append(opt, ap.procs[i].String())
// 	}
// 	return opt
// }

var cursor int

// func qqq() {
// 	event_queue := make(chan termbox.Event)
// 	go func() {
// 		for {
// 			event_queue <- termbox.PollEvent()
// 		}
// 	}()
// 	posx, posy := -1, 0
// 	color := termbox.ColorDefault
// 	color_change_tick := time.NewTicker(1 * time.Second)
// 	draw_tick := time.NewTicker(30 * time.Millisecond)

// loop:
// 	for {
// 		select {
// 		case ev := <-event_queue:
// 			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
// 				break loop
// 			}
// 		case <-color_change_tick.C:
// 			color++
// 			if color >= 8 {
// 				color = 0
// 			}
// 		case <-draw_tick.C:
// 			w, h := termbox.Size()
// 			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
// 			posx++
// 			if posx >= w {
// 				posx = 0
// 				posy++
// 				if posy >= h {
// 					posy = 0
// 				}
// 			}
// 			termbox.SetCell(posx, posy, '_', termbox.ColorDefault, color)
// 			termbox.Flush()
// 		}
// 	}
// }

/*












 */
////////////////////

func Action_ToggleSelection(ap *allProc, ib *InfoBox) {
	ap.stream[ib.cursor].isSelected = !ap.stream[ib.cursor].isSelected
}

func Action_MoveSelectionUP(ap *allProc, ib *InfoBox) {
	ib.cursor--
	if ib.cursor < 0 {
		ib.cursor = 0
	}
}

func Action_MoveSelectionDOWN(ap *allProc, ib *InfoBox) {
	ib.cursor++
	for ib.cursor >= len(ap.stream) {
		ib.cursor--
	}
}

func Action_MoveSelectionTO(ap *allProc, ib *InfoBox, pos int) {
	ib.cursor = pos
	if ib.cursor < 0 {
		ib.cursor = 0
	}
	for ib.cursor >= len(ap.stream) {
		ib.cursor--
	}
}

func Action_BumpStreamToTOP(ap *allProc, ib *InfoBox) {
	Action_Pause(ap, ib)
	ap.bumpToTop(ib.cursor)
	Action_MoveSelectionTO(ap, ib, 0)
}

func Action_BumpStreamByOne(ap *allProc, ib *InfoBox) {
	ap.bumpUpByOne(ib.cursor)
	Action_MoveSelectionUP(ap, ib)
}

func Action_TogglePause(ap *allProc, ib *InfoBox) {

	if ap.stream[0].lastCommand == "pause" {
		Action_Continue(ap, ib)
	} else {
		Action_Pause(ap, ib)
	}
}

func Action_Pause(ap *allProc, ib *InfoBox) {
	ap.stream[0].handler.Pause()
	ap.stream[0].lastCommand = "pause"
}

func Action_Continue(ap *allProc, ib *InfoBox) {
	ap.stream[0].handler.Continue()
	ap.stream[0].lastCommand = "continue"
}

func Action_StartNext(ap *allProc, ib *InfoBox) {
	if len(ap.stream) > 1 {
		ap.stream = ap.stream[1:]
		Action_MoveSelectionUP(ap, ib)
	} else {
		ap.stream = nil
		return
	}

}

/*








 */

func eventKeyToString(ev termbox.Event) string {
	out := fmt.Sprintf("%v (%v) ", ev.Ch, ev.Key) + ": key not mapped"
	out = strconv.QuoteRuneToGraphic(ev.Ch) + " " + strconv.QuoteRune(rune(ev.Key))
	/*
	   TODO:
	   подружить экшены и сигналы от нажатия клавиш в map для синхронизации
	*/

	// switch ev.Ch {
	// case 0:
	// 	switch ev.Key {
	// 	default:
	// 		out = strconv.QuoteRune(rune(ev.Key))
	// 		//case termbox.KeyTab
	// 	}
	// 	//спец клавиши

	// // case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
	// // 	out = strconv.QuoteRune(ev.Ch)
	// // case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
	// // 	out = strconv.QuoteRune(ev.Ch)
	// // case 'ф', 'и', 'с', 'в', 'у', 'а', 'п', 'р', 'ш', 'о', 'л', 'д', 'ь', 'т', 'щ', 'з', 'й', 'к', 'ы', 'е', 'г', 'м', 'ц', 'ч', 'н', 'я':
	// // 	out = strconv.QuoteRune(ev.Ch)
	// default:
	// 	out = strconv.QuoteRuneToGraphic(ev.Ch) + " " + strconv.QuoteRune(rune(ev.Key))
	// }

	return out
}

func commandSequanceToTBKey(s string) termbox.Key {
	switch s {
	case "F1":
		return termbox.KeyF1
	case "F2":
		return termbox.KeyF2
	case "F3":
		return termbox.KeyF3
	case "F4":
		return termbox.KeyF4
	case "F5":
		return termbox.KeyF5
	case "F6":
		return termbox.KeyF6
	case "F7":
		return termbox.KeyF7
	case "F8":
		return termbox.KeyF8
	case "F9":
		return termbox.KeyF9
	case "F10":
		return termbox.KeyF10
	case "F11":
		return termbox.KeyF11
	case "F12":
		return termbox.KeyF12
	case "INS":
		return termbox.KeyInsert
	case "DEL":
		return termbox.KeyDelete
	case "HOME":
		return termbox.KeyHome
	case "END":
		return termbox.KeyEnd
	case "PGUP":
		return termbox.KeyPgup
	case "PGDN":
		return termbox.KeyPgdn
	case "UP":
		return termbox.KeyArrowUp
	case "DOWN":
		return termbox.KeyArrowDown
	case "LEFT":
		return termbox.KeyArrowLeft
	case "RIGHT":
		return termbox.KeyArrowRight
	case "LMB":
		return termbox.MouseLeft
	case "MMB":
		return termbox.MouseMiddle
	case "RMB":
		return termbox.MouseRight
	case "MR":
		return termbox.MouseRelease
	case "MWUP":
		return termbox.MouseWheelUp
	case "MWDOWN":
		return termbox.MouseWheelDown
	case "SPACE":
		return termbox.KeySpace
	}
	return 20000
}
