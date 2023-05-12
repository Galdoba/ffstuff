package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Galdoba/ffstuff/cmd/grabber/download"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const (
	commandPAUSE                                 = "pause"
	commandCONTINUE                              = "continue"
	commandNONE                                  = "NONE"
	ACTION_MOVE_CURSOR_UP                        = "MOVE_CURSOR_UP"
	ACTION_MOVE_CURSOR_DOWN                      = "MOVE_CURSOR_DOWN"
	ACTION_MOVE_CURSOR_DOWN_AND_TOGGLE_SELECTION = "MOVE_CURSOR_DOWN_AND_TOGGLE_SELECTION"
	ACTION_TOGGLE_SELECTION_STATE                = "TOGGLE_SELECTION_STATE"
	ACTION_SELECT_ALL_WITH_SAME_EXTENTION        = "SELECT_ALL_WITH_SAME_EXTENTION"
	ACTION_DROP_SELECTIONS                       = "DROP_SELECTIONS"
	ACTION_MOVE_SELECTED_TOP                     = "MOVE_SELECTED_TOP"
	ACTION_MOVE_SELECTED_BOTTOM                  = "MOVE_SELECTED_BOTTOM"
	ACTION_MOVE_SELECTED_UP                      = "MOVE_SELECTED_UP"
	ACTION_MOVE_SELECTED_DOWN                    = "MOVE_SELECTED_DOWN"
	DECIDION_CONFIRM                             = "DECIDION_CONFIRM"
	DECIDION_DENY                                = "DECIDION_DENY"
	DELETE_SELECTED                              = "DELETE_SELECTED"
	DOWNLOAD_PAUSE                               = "DOWNLOAD_PAUSE"
	UNDO_MOVEMENT                                = "UNDO_MOVEMENT"
	input_mode_NORMAL                            = 1000
	input_mode_WAIT_CONFIRM                      = 1001
	input_mode_CONFIRM_RECEIVED                  = 1002
	input_mode_DENIAL_RECEIVED                   = 1003
)

func actionMap() map[string]func(*allProc, *InfoBox) error {
	mp := make(map[string]func(*allProc, *InfoBox) error)

	mp[ACTION_MOVE_CURSOR_UP] = Action_MoveCursorUP
	mp[ACTION_MOVE_CURSOR_DOWN] = Action_MoveCursorDOWN
	mp[ACTION_MOVE_CURSOR_DOWN_AND_TOGGLE_SELECTION] = Action_MoveCursorDOWNandSELECT
	mp[ACTION_TOGGLE_SELECTION_STATE] = Action_ToggleSelection
	mp[ACTION_SELECT_ALL_WITH_SAME_EXTENTION] = Action_SelectAllWithSameExtention
	mp[ACTION_DROP_SELECTIONS] = Action_DropSelection
	mp[ACTION_MOVE_SELECTED_TOP] = Action_MoveSelectedTop
	mp[ACTION_MOVE_SELECTED_BOTTOM] = Action_MoveSelectedBottom
	mp[ACTION_MOVE_SELECTED_UP] = Action_MoveSelectedUp
	mp[ACTION_MOVE_SELECTED_DOWN] = Action_MoveSelectedDown
	mp[DECIDION_CONFIRM] = DesidionConfirm
	mp[DECIDION_DENY] = DesidionDeny
	mp[DOWNLOAD_PAUSE] = Action_TogglePause
	mp[UNDO_MOVEMENT] = Action_UndoMovement
	mp[DELETE_SELECTED] = Action_DeleteSelected

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
	stream            []*stream
	globalStop        bool
	activeHandlerChan chan download.Response
	activeHandler     download.Handler
	activeStream      *stream
	//streamDataBak     sorting.IndexList
	indexBuf *IndexBuffer
	warnings []warning
	//endEvent bool
	//cursorSelection int

}

type stream struct {
	source string
	//temp             string
	dest             string
	baseName         string
	progress         int64
	expected         int64
	handler          download.Handler
	lastResponse     string
	lastResponseTime time.Time
	lastCommand      string
	isSelected       bool
	size             int64
	warning          string
}

func NewStream(source, dest, baseName string) *stream {
	st := stream{source, dest, baseName, 0, 0, nil, "NONE", time.Now(), "NONE", false, 0, ""}
	return &st
}

func (st *stream) Start() {
	st.handler = download.StartNew(st.source, st.dest+"\\"+st.baseName+".gdf")
}

func (st *stream) QueueString() string {
	str := "["
	switch st.isSelected {
	case true:
		str += "x]"
	case false:
		str += " ]"
	}
	progr := " "
	if st.handler != nil {
		progr = st.Progress()
	}
	str += " " + st.source + "|" + st.lastCommand + "|" + progr + "|"
	if st.handler != nil {
		str += fmt.Sprintf("handler:%v", st.handler)
	} else {
		str += fmt.Sprintf("no handler")
	}
	return str
}

func (st *stream) CompleteString() string {
	str := "["
	switch st.isSelected {
	case true:
		str += "x]"
	case false:
		str += " ]"
	}
	progr := " "
	str += " " + st.source + "|" + st.lastCommand + "|" + progr + "|" + "done"
	return str
}

func (st *stream) ErrString() string {
	str := "["
	switch st.isSelected {
	case true:
		str += "x]"
	case false:
		str += " ]"
	}
	progr := " "
	str += " " + st.source + "|" + st.lastCommand + "|" + progr + "|" + st.warning
	return str
}

func (st *stream) Progress() string {

	switch st.handler.Status() {
	default:
		proc := int64(0)
		size := st.handler.FileSize()
		if size != 0 {
			proc = st.handler.Progress() / (st.handler.FileSize() / 100)
		}
		return fmt.Sprintf(" %v", proc) + "% "
	case download.STATUS_COMPLETED:
		return "completed"
	case download.STATUS_TERMINATED:
		return "terminated"
	case download.STATUS_ERR:
		return "error"
	}
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
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	fg := termbox.ColorWhite
	bg := termbox.ColorBlack
	tkr := tickerImage(ib.ticker / 5)
	tbprint(0, 0, fg, bg, "Last Key Pressed:"+ib.lastKeysPressed+"__: "+fmt.Sprintf("%v", len(ap.indexBuf.Set)))
	tbprint(0, 1, fg, bg, "Tiker:"+tkr)
	tbprint(0, 2, fg, bg, "Grabber Dowloading: "+fmt.Sprintf("%v", ap.activeHandler))
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
	newData := []string{}
	switch ib.inputMode {
	default:
		return fmt.Errorf("unknown input mode: %v", ib.inputMode)
	case input_mode_NORMAL:
		// completeList := []string{}
		// queueList := []string{}
		// errList := []string{}
		for _, pr := range ap.stream {
			if pr.warning == "duplicate" {
				newData = append(newData, pr.ErrString())
				continue
			}
			if pr.handler == nil {
				newData = append(newData, pr.QueueString())
				continue
			}
			switch pr.handler.Status() {
			default:
				panic(pr.handler.Status()) //не должно срабатывать
			case download.STATUS_COMPLETED:
				newData = append(newData, pr.CompleteString())
				// if ib.cursor <= i {
				// 	ib.cursor++
				// }
			case download.STATUS_ERR:
				newData = append(newData, pr.ErrString())
			case download.STATUS_TRANSFERING, download.STATUS_PAUSED, download.STATUS_NIL:
				newData = append(newData, pr.QueueString())
			}

		}
		//newData = append(newData, completeList...)
		//newData = append(newData, queueList...)
		//newData = append(newData, errList...)

	case input_mode_WAIT_CONFIRM:
		newData = append(newData, "Press Enter to confirm or Esc to deny")
		for _, pr := range ap.stream {
			newData = append(newData, pr.QueueString())
		}
		//panic("not expecting confirm mode")
		//pos, sel := ap.streamDataProposition.Export()
		// for i, _ := range sel {
		// 	pr := ap.stream[pos[i]].String()
		// 	newData = append(newData, pr)
		// }
		// case input_mode_CONFIRM_RECEIVED, input_mode_DENIAL_RECEIVED:
		// 	switchToNORMALMode(ap, ib)
	}
	ib.data = newData
	return nil
}

func NewAllProcesses(dest string, paths ...string) *allProc {
	ap := allProc{}

	for i, path := range paths {
		file := filepath.Base(path)
		ap.stream = append(ap.stream, NewStream(path, dest, file))
		ap.stream[i].warning = "?"
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
	function       func(*allProc, *InfoBox) error
	validInputMode int
}

func (act *action) setValidInputMode() {
	act.validInputMode = input_mode_NORMAL
	if strings.Contains(act.eventName, "DECIDION") {
		act.validInputMode = input_mode_WAIT_CONFIRM
	}
	// if strings.Contains(act.eventName, "DELETE") {
	// 	act.validInputMode = input_mode_CONFIRM_RECEIVED
	// }
}

func setupAction(key string, configMap map[string]string, function func(*allProc, *InfoBox) error) ([]action, error) {
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
	включить/выключить ограниченную скорость закачки
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

	ap := NewAllProcesses(configMap["dest"], paths...)
	//activeStream := (*stream)(nil)

	ib := &InfoBox{}
	ib.data = []string{}
	actionPool := ActionPool{}
	//competedFilesHash := make(map[string]bool)
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
	ap.indexBuf = CreateIndexBuffer()
	ap.SaveState()
	// started := false
loop:
	for {
		ap.initialCheck()
		if len(ap.stream) == 0 && len(ap.warnings) == 0 {
			break
		}
		// if !started {
		// 	err := ap.stream[0].start()
		// 	if err != nil {
		// 		panic("start dowload stream: " + err.Error())
		// 	}
		// 	ap.stream[0].lastCommand = commandCONTINUE
		// 	ap.activeHandler = ap.stream[0].handler
		// 	started = true
		// }

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
								switch action.eventName {
								case ACTION_MOVE_SELECTED_UP, ACTION_MOVE_SELECTED_TOP, ACTION_MOVE_SELECTED_DOWN, ACTION_MOVE_SELECTED_BOTTOM:
									ap.SaveState()
								}
								break
							} else {
								ib.lastKeysPressed += " skip action" + " " + action.eventName
							}
						}

					}
				case 'q', 'й':
					break loop
				default:
					key := map_evCh(ev.Ch) + fmt.Sprintf("_%v", ib.inputMode)
					actions := actionPool.byKBKey[key]
					for _, action := range actions {
						if action.validInputMode != ib.inputMode {
							continue
						}

						err := action.function(ap, ib)
						if err != nil {
							panic(err.Error())
						}
						switch action.eventName {
						case ACTION_MOVE_SELECTED_UP, ACTION_MOVE_SELECTED_TOP, ACTION_MOVE_SELECTED_DOWN, ACTION_MOVE_SELECTED_BOTTOM:
							ap.SaveState()
						}
						ib.lastKeysPressed = key + " do action  " + action.eventName

					}
					//ib.lastKeysPressed = key //fmt.Sprintf("%v", string(ev.Ch))
				}
				ib.ticker = 1
			}
		case <-draw_tick.C:
			ib.ticker++
			if len(ap.stream) == 0 && len(ap.warnings) == 0 {
				break loop
			}
			if len(ap.stream) > 0 && ap.stream[0].handler != nil && ap.stream[0].lastCommand == commandCONTINUE {

			}
			ap.confirmStreams()
			//		case ev := <-handlerEvents:
			// case ev, ok := <-ap.activeHandler.Listen():

			// 	if ok {
			// 		ap.stream[0].lastResponse = ev.String()
			// 		if ev.String() == "completed" {

			// 			ib.ticker = 0

			// 			//err := ap.CloseStream()
			// 			// if err := ap.CloseStream(); err != nil {
			// 			// 	panic("CLOSE STREAM: " + err.Error())
			// 			// }

			// 			// err := Action_StartNext(ap, ib)
			// 			// if err != nil {
			// 			// 	return err
			// 			// }
			// 			//Action_Continue(ap, ib)
			// 			//handlerEvents = make(chan download.Response)
			// 			//ap.activeHandlerChan = make(chan download.Response)
			// 			//ap.activeHandler.Kill()
			// 		}
			// 	} else {
			// 		if ap.activeHandler == nil {
			// 			fmt.Print("chennel closed\r")
			// 			ap.activeHandlerChan = nil
			// 		} else {
			// 			fmt.Print(ev.String() + "  channel closed (not realy) " + fmt.Sprintf("%v", ap.activeHandler) + "\r")
			// 			for i, stream := range ap.stream {
			// 				if stream.handler != nil {
			// 					continue
			// 				}
			// 				ap.stream[i].start()
			// 				ap.activeHandler = ap.stream[i].handler
			// 			}

			// 			time.Sleep(time.Second)
			// 		}

			// 	}

		}
		if ap.activeStream == nil {
			err := ap.ActivateStream()
			if err == ErrAllCompleted {
				return nil
			}
		}
		if ap.activeStream != nil && ap.activeStream.handler.Status() == download.STATUS_COMPLETED {
			ap.activeStream = nil
		}

		ib.Update(ap)
		ib.Draw(ap)

	}

	fmt.Println("END")

	return nil
}

var ErrAllCompleted = fmt.Errorf("All streams complete")

func (ap *allProc) ActivateStream() error {
	if ap.globalStop {
		return nil
	}
	haveNIL := 0
	statusMap := make(map[download.Status]int)
	for i := range ap.stream {
		if ap.stream[i].handler == nil {
			haveNIL++
			continue
		}
		statusMap[ap.stream[i].handler.Status()]++
	}
	if statusMap[download.STATUS_TRANSFERING] > 0 {
		return nil
	}
	if len(ap.stream) == statusMap[download.STATUS_COMPLETED] {
		return ErrAllCompleted
	}

	for i := range ap.stream {
		if ap.stream[i].warning != "" {
			continue
		}
		if ap.stream[i].handler != nil && ap.stream[i].handler.Status() == download.STATUS_PAUSED {
			ap.stream[i].handler.Continue()
			ap.activeStream = ap.stream[i]
			return nil
		}
		if ap.stream[i].handler == nil {
			ap.stream[i].Start()
			ap.activeStream = ap.stream[i]
			return nil
		}
	}
	// switch {
	// case statusMap[download.STATUS_TRANSFERING] > 0:
	// 	return
	// case statusMap[download.STATUS_PAUSED] > 0:
	// 	for i := range ap.stream {
	// 		if ap.stream[i].handler == nil {
	// 			ap.activeStream = ap.stream[i]
	// 			ap.activeStream.Start()
	// 			return
	// 		}
	// 		if ap.stream[i].handler.Status() == download.STATUS_PAUSED {
	// 			ap.activeStream = ap.stream[i]
	// 			ap.activeStream.handler.Continue()
	// 			return
	// 		}
	// 	}
	// default:
	// 	for i := range ap.stream {
	// 		if ap.stream[i].handler == nil {
	// 			ap.activeStream = ap.stream[i]
	// 			ap.activeStream.Start()
	// 			return
	// 		}
	// 	}
	// }

	// for i := range ap.stream {
	// 	if ap.stream[i].handler == nil || ap.stream[i].handler.Status() == download.STATUS_PAUSED {
	// 		ap.activeStream = ap.stream[i]
	// 		ap.activeStream.Start()
	// 		return
	// 	}
	// }
	ap.activeStream = nil
	//panic("Do not Know what to do")
	return nil
}

////////////////////

func (ap *allProc) ExportSelected() []bool {
	sel := []bool{}
	for _, stream := range ap.stream {
		sel = append(sel, stream.isSelected)
	}
	return sel
}

func renameFile(stream *stream) error {
	//panic(stream.temp + stream.baseName + "===>" + stream.dest + stream.baseName)
	return os.Rename(stream.dest+stream.baseName+".gdf", stream.dest+stream.baseName)
}

func (ap *allProc) CloseStream() error {
	if len(ap.stream) < 1 {
		return fmt.Errorf(" CloseStream(): no streams to close")
	}
	stream := ap.stream[0]
	ap.addWarning(newWarning(stream.baseName, stream.dest, "transfert not confirmed"))
	//time.Sleep(time.Millisecond * 500)

	// if _, err := os.Stat(stream.dest + stream.baseName); os.IsNotExist(err) {
	// 	go renameFile(stream)
	// }
	// time.Sleep(time.Millisecond * 200)
	//ap.warnings = append(ap.warnings, stream.baseName+"|"+stream.temp+"|"+stream.dest)
	/*
		The process cannot access the file because it is being used by an  tobot_s01_12_2010__hd_rus20.m4a: rename d:\IN\IN_2022-05-11\proxy\temp\tobot_s01_12_2010__hd_rus20.m4a d:\IN\IN_2022-05-11\proxy\tobot_s01_12_2010__hd_rus20.m4a: The system cannot find the file specified.
	*/
	//ap.indexBuf.Remove(stream.source)
	ap.activeHandlerChan = nil
	// if len(ap.stream) > 0 {
	// 	ap.stream = ap.stream[1:]
	// }
	return nil
}

func (ap *allProc) StreamString() string {
	s := ""
	for _, str := range ap.stream {
		s += str.lastCommand + "=" + str.lastResponse + "=" + str.source + "\n"
	}
	return s
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

type IndexData struct {
	//InitialPos int
	SavedPos int
	Selected bool
}

type IndexState struct {
	state map[string]IndexData
}

type IndexBuffer struct {
	SavedStates int
	Set         []IndexState
}

func (ap *allProc) IndexStateCurrent() IndexState {
	is := IndexState{}
	is.state = map[string]IndexData{}
	for i, stream := range ap.stream {
		is.state[stream.source] = IndexData{i, stream.isSelected}
	}
	return is
}

func CreateIndexBuffer() *IndexBuffer {
	indBuf := IndexBuffer{}
	return &indBuf
}

func (indBuf *IndexBuffer) LastState() IndexState {
	indexLen := len(indBuf.Set)
	if indexLen == 0 {
		return IndexState{}
	}
	return indBuf.Set[len(indBuf.Set)-1]
}

func (ap *allProc) Selected() []bool {
	sel := []bool{}
	lastState := ap.indexBuf.LastState()
	for _, stream := range ap.stream {
		sel = append(sel, lastState.state[stream.source].Selected)
	}
	return sel
}

func statesEqual(index1, index2 IndexState) bool {
	l1 := len(index1.state)
	l2 := len(index2.state)
	if l1 == 0 {
		return false
	}
	if l2 == 0 {
		return false
	}

	for k := range index1.state {
		//if index1.state[k].SavedPos != index2.state[k].SavedPos || index1.state[k].Selected != index2.state[k].Selected {
		if index1.state[k].SavedPos != index2.state[k].SavedPos || index1.state[k].Selected != index2.state[k].Selected {
			//			panic(fmt.Sprintf("%v: %v = %v   %v = %v", k, index1.state[k].SavedPos, index2.state[k].SavedPos, index1.state[k].Selected, index2.state[k].Selected))
			return false
		}
	}
	return true
}

func (is *IndexState) String() string {
	str := ""
	for k, v := range is.state {
		str += k + fmt.Sprintf("%v\n", v)
	}
	return str
}

func (ap *allProc) SaveState() {

	current := ap.IndexStateCurrent()
	last := ap.indexBuf.LastState()
	if !statesEqual(current, last) {
		//panic(current.String() + "===" + current.String())
		ap.indexBuf.Set = append(ap.indexBuf.Set, ap.IndexStateCurrent())
	}

}

func (indBuf *IndexBuffer) DeleteLast() {
	if len(indBuf.Set) < 2 {
		return
	}
	indBuf.Set = indBuf.Set[:len(indBuf.Set)-1]
}

func (indBuf *IndexBuffer) Remove(source string) {
	for i, _ := range indBuf.Set {
		positionWas := indBuf.Set[i].state[source].SavedPos
		delete(indBuf.Set[i].state, source)
		for k, v := range indBuf.Set[i].state {
			if v.SavedPos > positionWas {
				v.SavedPos--
			}
			indBuf.Set[i].state[k] = v
		}
	}
}

func (ap *allProc) arrangeStreamsBy(index IndexState) {
	newOrder := ap.stream
	for i, stream := range ap.stream {
		recall := index.state[stream.source]
		if recall.SavedPos >= len(newOrder) {
			panic(fmt.Sprintf("DEBUG must not happen: saved=%v; len=%v", recall.SavedPos, len(newOrder)))
		}
		newOrder[i].isSelected = recall.Selected
		newOrder[i], newOrder[recall.SavedPos] = newOrder[recall.SavedPos], newOrder[i]
	}
	ap.stream = newOrder
}

const (
	stat_haveDuplicate = 10
	stat_notConfirmed  = 11
)

type warning struct {
	base string
	temp string
	dest string
	//status int
	text string
}

func newWarning(base, dest, msg string) warning {
	wrn := warning{}
	wrn.base = base
	wrn.dest = dest
	wrn.text = msg
	return wrn
}

func (ap *allProc) addWarning(wrn warning) {
	for i, haveW := range ap.warnings {
		if haveW.base == wrn.base {
			ap.warnings[i] = wrn
			return
		}
	}
	ap.warnings = append(ap.warnings, wrn)
}

func (ap *allProc) removeWarning(wrnBase string) {
	newW := []warning{}
	for _, w := range ap.warnings {
		if wrnBase == w.base {
			continue
		}
		newW = append(newW, w)
	}
	ap.warnings = newW
}

func (ap *allProc) confirmStreams() {
	// for _, wrn := range ap.warnings {
	// 	if err := renameFileName(wrn.temp+wrn.base, wrn.dest+wrn.base); err != nil {
	// 		if strings.Contains(err.Error(), "The system cannot find the file specified") {
	// 			ap.addWarning(newWarning(wrn.base, wrn.dest, "The system cannot find the file specified"))
	// 		}
	// 		if strings.Contains(err.Error(), "The process cannot access the file") {
	// 			ap.addWarning(newWarning(wrn.base, wrn.dest, "The process cannot access the file"))
	// 		}
	// 		ap.addWarning(newWarning(wrn.base, wrn.dest, err.Error()))
	// 	} else {
	// 		ap.removeWarning(wrn.base)
	// 	}
	// }
	for _, stream := range ap.stream {
		if stream.warning == "done" {
			continue
		}
		if stream.handler == nil {
			exist, _ := exists(stream.dest + stream.baseName)
			if exist {
				stream.warning = "duplicate"
			} else {
				stream.warning = ""
			}

			//if exist {}
		}
		if stream.handler != nil && stream.handler.Status() == download.STATUS_COMPLETED {
			renameFileName(stream.dest+stream.baseName+".gdf", stream.dest+stream.baseName)
			stream.warning = "done"
		}
	}
}

func (ap *allProc) initialCheck() {
	for _, stream := range ap.stream {
		exist, err := exists(stream.dest + stream.baseName)
		if err != nil {
			ap.addWarning(newWarning(stream.baseName, stream.dest, err.Error()))
			continue
		}
		if exist {
			ap.addWarning(newWarning(stream.baseName, stream.dest, "duplicate found"))
		}
	}
}

func renameFileName(file1, file2 string) error {
	exist, err := exists(file2)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("duplicate found") //./IN/@SCRIPTS
	}
	return os.Rename(file1, file2)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
