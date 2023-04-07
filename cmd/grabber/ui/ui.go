package ui

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Galdoba/ffstuff/cmd/grabber/download"
	"github.com/Galdoba/ffstuff/cmd/grabber/sorting"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const (
	commandPAUSE                  = "pause"
	commandCONTINUE               = "continue"
	commandNONE                   = "NONE"
	ACTION_MOVE_CURSOR_UP         = "MOVE_CURSOR_UP"
	ACTION_MOVE_CURSOR_DOWN       = "MOVE_CURSOR_DOWN"
	ACTION_TOGGLE_SELECTION_STATE = "TOGGLE_SELECTION_STATE"
	ACTION_DROP_SELECTIONS        = "DROP_SELECTIONS"
	ACTION_MOVE_SELECTED_TOP      = "MOVE_SELECTED_TOP"
	ACTION_MOVE_SELECTED_BOTTOM   = "MOVE_SELECTED_BOTTOM"
	ACTION_MOVE_SELECTED_UP       = "MOVE_SELECTED_UP"
	ACTION_MOVE_SELECTED_DOWN     = "MOVE_SELECTED_DOWN"
	DECIDION_CONFIRM              = "DECIDION_CONFIRM"
	DECIDION_DENY                 = "DECIDION_DENY"
	DOWNLOAD_PAUSE                = "DOWNLOAD_PAUSE"
	input_mode_NORMAL             = 1000
	input_mode_WAIT_CONFIRM       = 1001
)

func actionMap() map[string]func(*allProc, *InfoBox) {
	mp := make(map[string]func(*allProc, *InfoBox))

	mp[ACTION_MOVE_CURSOR_UP] = Action_MoveCursorUP
	mp[ACTION_MOVE_CURSOR_DOWN] = Action_MoveCursorDOWN
	mp[ACTION_TOGGLE_SELECTION_STATE] = Action_ToggleSelection
	mp[ACTION_DROP_SELECTIONS] = Action_DropSelection
	mp[ACTION_MOVE_SELECTED_TOP] = Action_MoveSelectedTop
	mp[DECIDION_CONFIRM] = DesidionConfirm
	mp[DECIDION_DENY] = DesidionDeny
	mp[DOWNLOAD_PAUSE] = Action_TogglePause
	//mp[ACTION_MOVE_SELECTED_BOTTOM] = Action_MoveSelectedBottom
	//mp[ACTION_MOVE_SELECTED_UP] = Action_MoveSelectedUp
	//mp[ACTION_MOVE_SELECTED_DOWN] = Action_MoveSelectedDown
	return mp
}

// func main() {
// 	inputinfo.CleanScanData()
// }

type allProc struct {
	//procs  []*proc
	stream                []*stream
	globalStop            bool
	activeHandlerChan     chan download.Response
	streamDataBak         sorting.IndexList
	streamDataProposition sorting.IndexList
	//endEvent bool
	//cursorSelection int
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

// func (ap *allProc) update() {
// 	if len(ap.stream) == 0 {
// 		return
// 	}
// }

func (ap *allProc) reverseStreamOrder() {
	for i, j := 0, len(ap.stream)-1; i < j; i, j = i+1, j-1 {
		ap.stream[i], ap.stream[j] = ap.stream[j], ap.stream[i]
	}
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
	data            []string
	cursor          int
	ticker          int
	lastKeysPressed string
	inputMode       int
}

func tickerImage(i int) string {
	s := ""
	for len(s) < i && len(s) < 25 {
		s += "="
	}
	return s
}

func (ib *InfoBox) Draw(ap *allProc) {
	//ib.Update(ap)
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	fg := termbox.ColorWhite
	bg := termbox.ColorBlack
	tkr := tickerImage(ib.ticker / 5)
	tbprint(0, 0, fg, bg, "Last Key Pressed:"+ib.lastKeysPressed)
	tbprint(0, 1, fg, bg, "Tiker:"+tkr)
	tbprint(0, 2, fg, bg, "Grabber Dowloading:")
	for i, data := range ib.data {

		if i == ib.cursor {
			fg = termbox.ColorBlack
			bg = termbox.ColorWhite
		}
		tbprint(2, i+3, fg, bg, data)
		fg = termbox.ColorWhite
		bg = termbox.ColorBlack
	}
	termbox.Flush()
}

// This function is often useful:
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {

	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func (ib *InfoBox) Update(ap *allProc) error {
	match, err := sorting.Check(ap.streamDataProposition, ap.streamDataBak)
	if err != nil {
		return fmt.Errorf("func (ib *InfoBox) Update(ap *allProc): %v", err.Error()) //DEBUG
	}
	if !match && ib.inputMode == input_mode_NORMAL {
		switchToWaitConfirmMode(ap, ib)
	}
	newData := []string{}
	switch ib.inputMode {
	default:
		return fmt.Errorf("unknown input mode: %v", ib.inputMode)
	case input_mode_NORMAL:
		for _, pr := range ap.stream {
			newData = append(newData, pr.String())
		}
	case input_mode_WAIT_CONFIRM:
		pos, sel := ap.streamDataProposition.Export()
		for i, _ := range sel {
			pr := ap.stream[pos[i]].String()
			newData = append(newData, pr)
		}
	}
	ib.data = newData
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
	eventName      string
	key            string
	termbxKey      []termbox.Key
	termbxRunes    []rune
	function       func(*allProc, *InfoBox)
	validInputMode int
}

func (act *action) setValidInputMode() {
	act.validInputMode = input_mode_NORMAL
	if strings.Contains(act.eventName, "DECIDION") {
		act.validInputMode = input_mode_WAIT_CONFIRM
	}
}

func setupAction(key string, configMap map[string]string, function func(*allProc, *InfoBox)) ([]action, error) {
	indexList := []string{}
	for k := range configMap {
		if strings.Contains(k, key) {
			indexList = append(indexList, strings.TrimPrefix(k, key+"_"))
		}
	}
	if len(indexList) == 0 {
		return nil, fmt.Errorf("key '%v' was not found in configMap \n'%v'", key, configMap)
	}
	actions := []action{}
	for _, index := range indexList {
		keyIndexed := key + "_" + index
		act := action{eventName: key, key: keyIndexed, function: function}

		tk, err := commandSequanceToTBKey(configMap[keyIndexed])
		if err != nil {
			runeArray := []rune(configMap[keyIndexed])
			act.termbxRunes = runeArray
		} else {
			act.termbxKey = append(act.termbxKey, tk)
		}
		act.setValidInputMode()
		actions = append(actions, act)
	}

	return actions, nil
}

type ActionPool struct {
	acmap       map[string][]action
	byTermbxKey map[termbox.Key][]action
	byKBKey     map[string][]action
}

type Action interface {
	commence(string) error
}

func (actpl *ActionPool) fillCommandActionMap(configMap map[string]string) error {
	actpl.acmap = make(map[string][]action)
	actpl.byTermbxKey = make(map[termbox.Key][]action)
	actpl.byKBKey = make(map[string][]action)

	for actName, actFunc := range actionMap() {
		action, err := setupAction(actName, configMap, actFunc)
		if err != nil {
			return fmt.Errorf("can not setup '%v' action: %v", actName, err.Error())
		}

		actpl.acmap[actName] = action
	}

	/*TODO: Прописать действия
	переместить выделенное по списку вверх на 1 позицию
	переместить выделенное по списку вниз на 1 позицию
	переместить выделенное по списку вверх до предела . . . ok
	переместить выделенное по списку вниз до предела
	удаление из очереди
	активная пауза (для всех процессов)
	включить/выключить ограниченную скорость закачки
	сброс всех выделений . . .  . .  . .  . .  .  . .  .  . .  ok
	инсерт (переключение выделение со сдвигом курсора вниз на 1 позицию)
	*/

	for k, actions := range actpl.acmap {
		for i, a := range actions {
			if len(a.termbxKey) != 0 {
				for _, tk := range a.termbxKey {
					actpl.byTermbxKey[tk] = append(actpl.byTermbxKey[tk], actpl.acmap[k][i])
				}
				continue
			}
			if len(a.termbxRunes) != 0 {
				for _, r := range a.termbxRunes {
					//проверяем совпадает ли руна(-ы) полученные
					//от эвента с картой кнопок в map_evCh(rune)
					if map_evCh(r) != configMap[a.key] {
						panic("mismatch: " + map_evCh(r) + " and " + configMap[a.key])
					}
				}
				keyWithIM := configMap[a.key] + fmt.Sprintf("_%v", a.validInputMode)
				actpl.byKBKey[keyWithIM] = append(actpl.byKBKey[keyWithIM], actpl.acmap[k][i])

			}
		}
	}
	//panic("DONE")
	return nil
}

func (actpl *ActionPool) AddAction(key string, act action) {
	counter := 0
	added := false
	for !added {
		indexedKey := act.key + fmt.Sprintf("_%v", counter)
		if _, ok := actpl.acmap[indexedKey]; ok {
			counter++
			continue
		} else {

		}
	}
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
	switchToNORMALMode(ap, ib)
	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	draw_tick := time.NewTicker(200 * time.Millisecond)
	handlerEvents := make(chan download.Response)
	//Tick := 0

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
			//Action_StartNext(ap, ib)
		}

		select {
		case ev := <-event_queue:

			switch ev.Type {

			case termbox.EventKey:
				switch ev.Ch {
				case 0:
					ib.lastKeysPressed = runeToKey(ev.Key) + fmt.Sprintf("_%v", ib.inputMode)
					if actions, ok := actionPool.byTermbxKey[ev.Key]; ok {
						for _, action := range actions {
							if action.validInputMode == ib.inputMode {
								ib.lastKeysPressed += " do action" + " " + action.eventName
								action.function(ap, ib)
								break
							} else {
								ib.lastKeysPressed += " skip action" + " " + action.eventName
							}
						}

					}
					// if ap.stream[0].lastCommand == commandPAUSE && ap.stream[0].handler != nil {
					// 	ap.stream[0].handler.Continue()
					// }
				case 'q', 'й':
					break loop
				default:
					key := map_evCh(ev.Ch) + fmt.Sprintf("_%v", ib.inputMode)
					actions := actionPool.byKBKey[key]
					for _, action := range actions {
						if action.validInputMode != ib.inputMode {
							continue
						}
						//if action, ok := actionPool.byKBKey[key]; ok {
						action.function(ap, ib)
						ib.lastKeysPressed += "  " + action.eventName
						//}
					}
					ib.lastKeysPressed = key //fmt.Sprintf("%v", string(ev.Ch))
					//ib.lastKeysPressed = string(ev.Ch)
				}
				ib.ticker = 1
			}
		case <-draw_tick.C:
			ib.ticker++
			if ap.stream[0].handler != nil && ap.stream[0].lastCommand == commandCONTINUE {

				//	ap.stream[0].handler.Continue()
				handlerEvents = ap.stream[0].handler.Listen()
			}
		case ev := <-handlerEvents:
			ap.stream[0].lastResponse = ev.String()
			if ev.String() == "completed" {
				ib.ticker = 0
				Action_StartNext(ap, ib)
				handlerEvents = make(chan download.Response)
			}

		}
		//ap.update()
		ib.Update(ap)
		ib.Draw(ap)

	}
	fmt.Println("END")
	return nil
}

////////////////////

func Action_ToggleSelection(ap *allProc, ib *InfoBox) {
	ap.stream[ib.cursor].isSelected = !ap.stream[ib.cursor].isSelected
}

func Action_MoveCursorUP(ap *allProc, ib *InfoBox) {
	ib.cursor--
	if ib.cursor < 0 {
		ib.cursor = 0
	}
}

func Action_MoveCursorDOWN(ap *allProc, ib *InfoBox) {
	ib.cursor++
	for ib.cursor >= len(ap.stream) {
		ib.cursor--
	}
}

// func Action_MoveSelectionTO(ap *allProc, ib *InfoBox, pos int) {
// 	ib.cursor = pos
// 	if ib.cursor < 0 {
// 		ib.cursor = 0
// 	}
// 	for ib.cursor >= len(ap.stream) {
// 		ib.cursor--
// 	}
// }

func saveCursor(ap *allProc, ib *InfoBox) string {
	return ap.stream[ib.cursor].source
}

func restoreCursor(ap []*stream, src string) int {
	for i, stream := range ap {
		if src != stream.source {
			continue
		}
		return i
	}
	return 0
}

func nothingSelected(bArray []bool) bool {
	for _, b := range bArray {
		if b {
			return false
		}
	}
	return true
}

func Action_MoveSelectedTop(ap *allProc, ib *InfoBox) {
	if nothingSelected(ap.ExportSelected()) {
		return
	}
	Action_Pause(ap, ib)

	selected := ap.ExportSelected()
	il := sorting.Import(selected)
	ap.streamDataBak = *il
	il.MoveTop()
	ap.streamDataProposition = *il
	//switchToWaitConfirmMode(ap, ib)

	// newIndex, newSelected := il.Export()
	// newStreamOrder := []*stream{}
	// for i, newInd := range newIndex {
	// 	newStreamOrder = append(newStreamOrder, ap.stream[newInd])
	// 	newStreamOrder[i].isSelected = newSelected[i]
	// }
	//cursorName := saveCursor(ap, ib)
	//termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)

	// pseudoAP := copyAP(ap)
	// pseudoAP.stream = newStreamOrder
	// ib.Update(pseudoAP)
	// ib.Draw(pseudoAP)
	// time.Sleep(time.Second * 5)
	// pseudoAP = &allProc{}
	//ap.stream = newStreamOrder
	//ib.cursor = restoreCursor(ap.stream, cursorName)
}

func copyAP(ap *allProc) *allProc {
	newAP := allProc{
		stream:     ap.stream,
		globalStop: true,
	}
	return &newAP
}

func Action_MoveSelectedUp(ap *allProc, ib *InfoBox) {
	globalStop(ap)
	cursorName := saveCursor(ap, ib)
	sel := []int{}
	unsel := []int{}
	for i, str := range ap.stream {
		switch str.isSelected {
		case true:
			sel = append(sel, i)
		case false:
			unsel = append(unsel, i)
		}
	}
	newStreamOrder := []*stream{}
	for _, s := range sel {
		newStreamOrder = append(newStreamOrder, ap.stream[s])
	}
	for _, uns := range unsel {
		newStreamOrder = append(newStreamOrder, ap.stream[uns])
	}
	ap.stream = newStreamOrder
	ib.cursor = restoreCursor(ap.stream, cursorName)
}

func globalStop(ap *allProc) {
	for i := range ap.stream {
		if ap.stream[i].handler == nil {
			continue
		}
		ap.stream[i].handler.Pause()
		ap.stream[i].lastCommand = commandPAUSE
	}
	time.Sleep(time.Millisecond * 200)
}

func Action_MoveSelectedBottom(ap *allProc, ib *InfoBox) {
	switch ib.inputMode {
	case input_mode_NORMAL:
		cursorName := saveCursor(ap, ib)
		ap.reverseStreamOrder()
		Action_MoveSelectedTop(ap, ib)
		ap.reverseStreamOrder()
		ib.cursor = restoreCursor(ap.stream, cursorName)
	}

}

func Action_DropSelection(ap *allProc, ib *InfoBox) {
	for i := range ap.stream {
		ap.stream[i].isSelected = false
	}
}

func Action_TogglePause(ap *allProc, ib *InfoBox) {

	if ap.stream[0].lastCommand == commandPAUSE {
		Action_Continue(ap, ib)
		return
	}
	if ap.stream[0].lastCommand == commandNONE {
		Action_Continue(ap, ib)
		return
	}
	Action_Pause(ap, ib)

}

func Action_Pause(ap *allProc, ib *InfoBox) {
	if ap.stream[0].handler != nil {
		ap.stream[0].handler.Pause()
		ap.stream[0].lastCommand = commandPAUSE
		time.Sleep(time.Millisecond * 200)
	}
}

func Action_Continue(ap *allProc, ib *InfoBox) {
	if ap.stream[0].handler != nil {
		ap.stream[0].handler.Continue()
		ap.stream[0].lastCommand = commandCONTINUE
	}
}

func Action_StartNext(ap *allProc, ib *InfoBox) {
	if ib.inputMode != input_mode_NORMAL {
		return
	}
	if len(ap.stream) > 1 {
		ap.stream = ap.stream[1:]
		Action_MoveCursorUP(ap, ib)
		if ap.stream[0].lastCommand == commandPAUSE {
			Action_Continue(ap, ib)
		}
	} else {
		ap.stream = nil
		return
	}
}

func switchToWaitConfirmMode(ap *allProc, ib *InfoBox) {
	Action_Pause(ap, ib)
	ib.inputMode = input_mode_WAIT_CONFIRM
}

func switchToNORMALMode(ap *allProc, ib *InfoBox) {
	ib.inputMode = input_mode_NORMAL
	Action_DropSelection(ap, ib)
	if len(ap.stream) == 0 {
		return
	}
	if ap.stream[0].lastCommand == commandPAUSE {
		Action_Continue(ap, ib)
	}
}

func DesidionConfirm(ap *allProc, ib *InfoBox) {
	switch ib.inputMode {
	case input_mode_WAIT_CONFIRM:
		newIndex, newSelected := ap.streamDataProposition.Export()
		newStreamOrder := []*stream{}
		for i, newInd := range newIndex {
			newStreamOrder = append(newStreamOrder, ap.stream[newInd])
			newStreamOrder[i].isSelected = newSelected[i]
		}
		ap.stream = newStreamOrder
		ap.streamDataBak = sorting.IndexList{}
		ap.streamDataProposition = sorting.IndexList{}
		switchToNORMALMode(ap, ib)
	}
}

func DesidionDeny(ap *allProc, ib *InfoBox) {
	switch ib.inputMode {
	case input_mode_WAIT_CONFIRM:
		newIndex, newSelected := ap.streamDataBak.Export()
		newStreamOrder := []*stream{}
		for i, newInd := range newIndex {
			newStreamOrder = append(newStreamOrder, ap.stream[newInd])
			newStreamOrder[i].isSelected = newSelected[i]
		}
		ap.stream = newStreamOrder
		ap.streamDataBak = sorting.IndexList{}
		ap.streamDataProposition = sorting.IndexList{}
		switchToNORMALMode(ap, ib)
	}
}

func (ap *allProc) ExportSelected() []bool {
	sel := []bool{}
	for _, stream := range ap.stream {
		sel = append(sel, stream.isSelected)
	}
	return sel
}

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
