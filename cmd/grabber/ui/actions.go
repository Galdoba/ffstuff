package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Galdoba/ffstuff/cmd/grabber/sorting"
	"github.com/Galdoba/ffstuff/pkg/namedata"
)

func Action_ToggleSelection(ap *allProc, ib *InfoBox) error {
	if ap.stream[ib.cursor].isSelected {
		ap.stream[ib.cursor].isSelected = false
	} else {
		ap.stream[ib.cursor].isSelected = true
	}
	return nil
}

func Action_MoveCursorUP(ap *allProc, ib *InfoBox) error {
	ib.cursor--
	if ib.cursor < 0 {
		ib.cursor = 0
	}
	// //проскакиваем вверх все готовые
	// for ib.cursor > 0 && ap.stream[ib.cursor].handler != nil && ap.stream[ib.cursor].handler.Status() == download.STATUS_COMPLETED {
	// 	ib.cursor--
	// }
	// //проскакиваем вниз все готовые
	// for ib.cursor < len(ap.stream) && ap.stream[ib.cursor].handler != nil && ap.stream[ib.cursor].handler.Status() == download.STATUS_COMPLETED {
	// 	ib.cursor++
	// }
	// for ib.cursor < len(ap.stream) && ap.stream[ib.cursor].handler != nil && ap.stream[ib.cursor].handler.Status() == download.STATUS_COMPLETED {
	// 	ib.cursor++
	// }
	// for ib.cursor >= len(ap.stream) {
	// 	ib.cursor--
	// }
	return nil
}

func Action_MoveCursorDOWN(ap *allProc, ib *InfoBox) error {
	ib.cursor++
	if ib.cursor >= len(ap.stream) {
		ib.cursor = len(ap.stream) - 1
	}
	// //проскакиваем вниз все готовые
	// for ib.cursor < len(ap.stream) && ap.stream[ib.cursor].handler != nil && ap.stream[ib.cursor].handler.Status() == download.STATUS_COMPLETED {
	// 	ib.cursor++
	// }
	// //проскакиваем вверх все готовые
	// for ib.cursor > 0 && ib.cursor < len(ap.stream) && ap.stream[ib.cursor].handler != nil && ap.stream[ib.cursor].handler.Status() == download.STATUS_COMPLETED {
	// 	ib.cursor--
	// }

	return nil
}

func Action_MoveCursorDOWNandSELECT(ap *allProc, ib *InfoBox) error {
	if ap.stream[ib.cursor].isSelected {
		ap.stream[ib.cursor].isSelected = false
	} else {
		ap.stream[ib.cursor].isSelected = true
	}
	ib.cursor++
	for ib.cursor >= len(ap.stream) {
		ib.cursor--
	}
	return nil
}

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

func Action_MoveSelectedTop(ap *allProc, ib *InfoBox) error {
	cursorPos := saveCursor(ap, ib)
	Action_Pause(ap, ib)
	if nothingSelected(ap.ExportSelected()) {
		ap.stream[ib.cursor].isSelected = true
	}
	il := sorting.Import(ap.ExportSelected())
	il.MoveTop()
	ap.loadStreamOrderFromIndexList(*il)
	Action_Continue(ap, ib)
	Action_DropSelection(ap, ib)
	ib.cursor = restoreCursor(ap.stream, cursorPos)
	return nil
}

func (ap *allProc) loadStreamOrderFromIndexList(il sorting.IndexList) {
	newIndex, newSelected := il.Export()
	newStreamOrder := []*stream{}
	for i, newInd := range newIndex {
		newStreamOrder = append(newStreamOrder, ap.stream[newInd])
		newStreamOrder[i].isSelected = newSelected[i]
	}
	ap.stream = newStreamOrder
}

func copyAP(ap *allProc) *allProc {
	newAP := allProc{
		stream:     ap.stream,
		globalStop: true,
	}
	return &newAP
}

func Action_MoveSelectedUp(ap *allProc, ib *InfoBox) error {
	cursorPos := saveCursor(ap, ib)
	Action_Pause(ap, ib)
	drop := false
	if nothingSelected(ap.ExportSelected()) {
		ap.stream[ib.cursor].isSelected = true
		drop = true
	}
	il := sorting.Import(ap.ExportSelected())
	il.MoveUp()
	ap.loadStreamOrderFromIndexList(*il)
	if drop {
		Action_DropSelection(ap, ib)
	}
	Action_Continue(ap, ib)
	ib.cursor = restoreCursor(ap.stream, cursorPos)
	return nil
}

func Action_MoveSelectedDown(ap *allProc, ib *InfoBox) error {
	cursorPos := saveCursor(ap, ib)
	Action_Pause(ap, ib)
	drop := false
	if nothingSelected(ap.ExportSelected()) {
		ap.stream[ib.cursor].isSelected = true
		drop = true
	}
	il := sorting.Import(ap.ExportSelected())
	il.MoveDown()
	ap.loadStreamOrderFromIndexList(*il)
	if drop {
		Action_DropSelection(ap, ib)
	}
	Action_Continue(ap, ib)
	ib.cursor = restoreCursor(ap.stream, cursorPos)
	return nil
}

func Action_MoveSelectedBottom(ap *allProc, ib *InfoBox) error {
	cursorPos := saveCursor(ap, ib)
	Action_Pause(ap, ib)
	if nothingSelected(ap.ExportSelected()) {
		ap.stream[ib.cursor].isSelected = true
	}
	il := sorting.Import(ap.ExportSelected())
	il.MoveBottom()
	ap.loadStreamOrderFromIndexList(*il)
	Action_Continue(ap, ib)
	Action_DropSelection(ap, ib)
	ib.cursor = restoreCursor(ap.stream, cursorPos)
	return nil
}

func Action_DropSelection(ap *allProc, ib *InfoBox) error {
	for i := range ap.stream {
		ap.stream[i].isSelected = false
	}
	return nil
}

func Action_TogglePause(ap *allProc, ib *InfoBox) error {
	if ap.globalStop {
		Action_Continue(ap, ib)
	} else {
		Action_Pause(ap, ib)
	}
	return nil
}

func Action_Pause(ap *allProc, ib *InfoBox) error {
	ap.globalStop = true
	for _, stream := range ap.stream {
		if stream.handler != nil {
			stream.handler.Pause()
		}
	}
	//ap.activeStream = nil

	// if ap.activeStream != nil {
	// 	ap.activeStream.handler.Pause()
	// }
	// if ap.stream[0].handler != nil {

	// 	ap.stream[0].handler.Pause()
	// 	panic(fmt.Sprintf("%v", ap.stream[0].handler))
	// 	ap.stream[0].lastCommand = commandPAUSE
	// 	time.Sleep(time.Millisecond * 200)
	// } else {
	// 	panic("have not")
	// }
	return nil
}

func Action_Continue(ap *allProc, ib *InfoBox) error {
	ap.globalStop = false
	// for _, stream := range ap.stream {
	// 	if stream.handler != nil && stream.handler.Status() == download.STATUS_PAUSED {
	// 		stream.handler.Continue()
	// 		return nil
	// 	}
	// }
	ap.activeStream = nil
	// if ap.activeStream != nil {
	// 	ap.activeStream.handler.Continue()
	// }

	// if ap.stream[0].handler != nil {
	// 	ap.activeHandler = ap.stream[0].handler //.Listen()
	// 	ap.activeHandler.Continue()
	// 	ap.stream[0].lastCommand = commandCONTINUE
	// }
	return nil
}

// func Action_StartNext(ap *allProc, ib *InfoBox) error {
// 	if ap.stream == nil {
// 		return fmt.Errorf(" Action_StartNext(): no streams to start")
// 	}
// 	if len(ap.stream) == 0 {
// 		return fmt.Errorf("no streams")
// 	}

// 	for i, stream := range ap.stream {
// 		if stream.lastResponse == "completed" {
// 			continue
// 		}
// 		if ap.stream[i].handler == nil {

// 			ap.stream[i].start()

// 		}

// 		//ap.activeHandlerChan = ap.stream[i].handler.Listen()
// 		// switch ap.stream[i].lastCommand {
// 		// case commandPAUSE, commandNONE:
// 		// 	Action_Continue(ap, ib)
// 		// }
// 		return nil
// 	}

// 	return nil //Action_MoveCursorUP(ap, ib)
// }

func switchToWaitConfirmMode(ap *allProc, ib *InfoBox) {
	Action_Pause(ap, ib)
	ib.inputMode = input_mode_WAIT_CONFIRM
}

func switchToNORMALMode(ap *allProc, ib *InfoBox) {
	ib.inputMode = input_mode_NORMAL
	if len(ap.stream) == 0 {
		return
	}
	if ap.stream[0].lastCommand == commandPAUSE {
		Action_Continue(ap, ib)
	}
	if ap.stream[0].lastCommand == commandNONE {
		Action_Continue(ap, ib)
	}
}

func DesidionConfirm(ap *allProc, ib *InfoBox) error {
	switch ib.inputMode {
	case input_mode_WAIT_CONFIRM:
		ib.inputMode = input_mode_CONFIRM_RECEIVED
	}
	return nil
}

func DesidionDeny(ap *allProc, ib *InfoBox) error {
	switch ib.inputMode {
	case input_mode_WAIT_CONFIRM:
		ib.inputMode = input_mode_DENIAL_RECEIVED
	}
	return nil
}

func Action_UndoMovement(ap *allProc, ib *InfoBox) error {
	Action_Pause(ap, ib)
	curs := saveCursor(ap, ib)
	ap.indexBuf.DeleteLast()
	lastState := ap.indexBuf.LastState()
	ap.arrangeStreamsBy(lastState)
	ib.cursor = restoreCursor(ap.stream, curs)
	Action_Continue(ap, ib)
	return nil
}

func Action_SelectAllWithSameExtention(ap *allProc, ib *InfoBox) error {
	curs := saveCursor(ap, ib)
	ext := namedata.RetrieveExtention(curs)
	for i, stream := range ap.stream {
		if strings.HasSuffix(stream.source, ext) {
			ap.stream[i].isSelected = true
		}
	}
	return nil
}

func Action_DeleteSelected(ap *allProc, ib *InfoBox) error {
	Action_Pause(ap, ib)
	toDelete := []int{}
	for i, stream := range ap.stream {
		if stream.isSelected {
			toDelete = append(toDelete, i)
		}
	}
	positionsProjected := []int{}
	for j, i := range toDelete {
		positionsProjected = append(positionsProjected, i-j)
	}
	switchToWaitConfirmMode(ap, ib)
	go func() {
		err := ap.deleteStreams(positionsProjected, ib)
		if err != nil {
			panic(err.Error())
		}
		if err == nil {
			switchToNORMALMode(ap, ib)
		}
	}()
	return nil
}

func (ap *allProc) deleteStreams(streams []int, ib *InfoBox) error {
	concluded := false
	for !concluded {
		switch ib.inputMode {
		case input_mode_CONFIRM_RECEIVED:
			concluded = true
		case input_mode_DENIAL_RECEIVED:
			return nil
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
	time.Sleep(time.Millisecond * 500)
	for _, i := range streams {
		if err := ap.DeleteStream(i); err != nil {
			panic("DELETE ERR (я хз откуда она может взяться)" + err.Error())
			return err
		}
	}
	time.Sleep(time.Millisecond * 200)
	return nil
}

func (ap *allProc) DeleteStream(i int) error {
	if len(ap.stream) < i {
		return fmt.Errorf(" DeleteStream(): can not delete stream: i > len(ap.stream)")
	}
	stream := ap.stream[i]
	ap.stream = append(ap.stream[:i], ap.stream[i+1:]...)
	ap.indexBuf.Remove(stream.source)
	ap.activeHandlerChan = nil
	return nil
}
