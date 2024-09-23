package ui

// import (
// 	"fmt"
// 	"strings"
// 	"time"

// 	"github.com/Galdoba/ffstuff/cmd/grabber/sorting"
// 	"github.com/Galdoba/ffstuff/pkg/namedata"
// 	"github.com/Galdoba/ffstuff/pkg/sortnames"
// 	"github.com/Galdoba/utils"
// 	"golang.design/x/clipboard"
// )

// func Action_ToggleSelection(ap *allProc, ib *InfoBox) error {
// 	if ap.stream[ib.cursor].isSelected {
// 		ap.stream[ib.cursor].isSelected = false
// 	} else {
// 		ap.stream[ib.cursor].isSelected = true
// 	}
// 	return nil
// }

// func Action_MoveCursorUP(ap *allProc, ib *InfoBox) error {
// 	ib.cursor--
// 	if ib.cursor < 0 {
// 		ib.cursor = 0
// 	}
// 	ib.lastScroll = 1

// 	return nil
// }

// func Action_MoveCursorTOP(ap *allProc, ib *InfoBox) error {
// 	ib.cursor = 0
// 	ib.lastScroll = 1
// 	ib.lowBorder = 0
// 	ib.highBorder = 0
// 	return nil
// }

// func Action_MoveCursorPU(ap *allProc, ib *InfoBox) error {
// 	switch ib.cursor {
// 	default:
// 		ib.cursor = ib.lowBorder
// 		return nil
// 	case ib.lowBorder:
// 		ib.lowBorder = utils.Max(0, ib.lowBorder-ib.drawLen)
// 		ib.highBorder = ib.lowBorder + ib.drawLen
// 		ib.cursor = ib.lowBorder

// 	}
// 	ib.lastScroll = 1 //0 = down/ 1 = up
// 	//ib.lowBorder = 0
// 	//ib.highBorder = 0
// 	return nil
// }

// func Action_MoveCursorDOWN(ap *allProc, ib *InfoBox) error {
// 	ib.cursor++
// 	if ib.cursor >= len(ap.stream) {
// 		ib.cursor = len(ap.stream) - 1
// 	}
// 	ib.lastScroll = 0
// 	return nil
// }

// func Action_MoveCursorPD(ap *allProc, ib *InfoBox) error {
// 	switch ib.cursor {
// 	default:
// 		ib.cursor = ib.highBorder
// 	case ib.highBorder:
// 		ib.cursor += ib.drawLen
// 	}

// 	if ib.cursor >= len(ap.stream) {
// 		ib.cursor = len(ap.stream) - 1
// 	}
// 	ib.lastScroll = 0
// 	ib.lowBorder = 0
// 	ib.highBorder = 0 //0 = down/ 1 = up
// 	return nil
// }

// func Action_MoveCursorBOTTOM(ap *allProc, ib *InfoBox) error {
// 	ib.cursor = len(ap.stream) - 1
// 	ib.lastScroll = 0
// 	ib.lowBorder = 0
// 	ib.highBorder = 0
// 	return nil
// }

// func Action_MoveCursorDOWNandSELECT(ap *allProc, ib *InfoBox) error {
// 	ap.stream[ib.cursor].isSelected = !ap.stream[ib.cursor].isSelected
// 	return Action_MoveCursorDOWN(ap, ib)
// }

// func saveCursor(ap *allProc, ib *InfoBox) string {
// 	return ap.stream[ib.cursor].source
// }

// func restoreCursor(ap []*stream, src string) int {
// 	for i, stream := range ap {
// 		if src != stream.source {
// 			continue
// 		}
// 		return i
// 	}
// 	return 0
// }

// func nothingSelected(bArray []bool) bool {
// 	for _, b := range bArray {
// 		if b {
// 			return false
// 		}
// 	}
// 	return true
// }

// func Action_MoveSelectedTop(ap *allProc, ib *InfoBox) error {
// 	cursorPos := saveCursor(ap, ib)
// 	Action_Pause(ap, ib)
// 	if nothingSelected(ap.ExportSelected()) {
// 		ap.stream[ib.cursor].isSelected = true
// 	}
// 	il := sorting.Import(ap.ExportSelected())
// 	il.MoveTop()
// 	ap.loadStreamOrderFromIndexList(*il)
// 	Action_Continue(ap, ib)
// 	Action_DropSelection(ap, ib)
// 	ib.cursor = restoreCursor(ap.stream, cursorPos)
// 	return nil
// }

// func (ap *allProc) loadStreamOrderFromIndexList(il sorting.IndexList) {
// 	newIndex, newSelected := il.Export()
// 	newStreamOrder := []*stream{}
// 	for i, newInd := range newIndex {
// 		newStreamOrder = append(newStreamOrder, ap.stream[newInd])
// 		newStreamOrder[i].isSelected = newSelected[i]
// 	}
// 	ap.stream = newStreamOrder
// }

// func copyAP(ap *allProc) *allProc {
// 	newAP := allProc{
// 		stream:     ap.stream,
// 		globalStop: true,
// 	}
// 	return &newAP
// }

// func Action_MoveSelectedUp(ap *allProc, ib *InfoBox) error {
// 	cursorPos := saveCursor(ap, ib)
// 	Action_Pause(ap, ib)
// 	drop := false
// 	if nothingSelected(ap.ExportSelected()) {
// 		ap.stream[ib.cursor].isSelected = true
// 		drop = true
// 	}
// 	il := sorting.Import(ap.ExportSelected())
// 	il.MoveUp()
// 	ap.loadStreamOrderFromIndexList(*il)
// 	if drop {
// 		Action_DropSelection(ap, ib)
// 	}
// 	Action_Continue(ap, ib)
// 	ib.cursor = restoreCursor(ap.stream, cursorPos)
// 	return nil
// }

// func Action_MoveSelectedDown(ap *allProc, ib *InfoBox) error {
// 	cursorPos := saveCursor(ap, ib)
// 	Action_Pause(ap, ib)
// 	drop := false
// 	if nothingSelected(ap.ExportSelected()) {
// 		ap.stream[ib.cursor].isSelected = true
// 		drop = true
// 	}
// 	il := sorting.Import(ap.ExportSelected())
// 	il.MoveDown()
// 	ap.loadStreamOrderFromIndexList(*il)
// 	if drop {
// 		Action_DropSelection(ap, ib)
// 	}
// 	Action_Continue(ap, ib)
// 	ib.cursor = restoreCursor(ap.stream, cursorPos)
// 	return nil
// }

// func Action_MoveSelectedBottom(ap *allProc, ib *InfoBox) error {
// 	cursorPos := saveCursor(ap, ib)
// 	Action_Pause(ap, ib)
// 	if nothingSelected(ap.ExportSelected()) {
// 		ap.stream[ib.cursor].isSelected = true
// 	}
// 	il := sorting.Import(ap.ExportSelected())
// 	il.MoveBottom()
// 	ap.loadStreamOrderFromIndexList(*il)
// 	Action_Continue(ap, ib)
// 	Action_DropSelection(ap, ib)
// 	ib.cursor = restoreCursor(ap.stream, cursorPos)
// 	return nil
// }

// func Action_DropSelection(ap *allProc, ib *InfoBox) error {
// 	for i := range ap.stream {
// 		ap.stream[i].isSelected = false
// 	}
// 	return nil
// }

// func Action_TogglePause(ap *allProc, ib *InfoBox) error {
// 	if ap.globalStop {
// 		Action_Continue(ap, ib)
// 	} else {
// 		Action_Pause(ap, ib)
// 	}
// 	return nil
// }

// func Action_Pause(ap *allProc, ib *InfoBox) error {
// 	ap.globalStop = true
// 	for _, stream := range ap.stream {
// 		if stream.handler != nil {
// 			stream.handler.Pause()
// 		}
// 	}
// 	return nil
// }

// func Action_Continue(ap *allProc, ib *InfoBox) error {
// 	ap.globalStop = false
// 	ap.activeStream = nil
// 	return nil
// }

// func switchToWaitConfirmMode(ap *allProc, ib *InfoBox) {
// 	Action_Pause(ap, ib)
// 	ib.inputMode = input_mode_WAIT_CONFIRM
// }

// func switchToNORMALMode(ap *allProc, ib *InfoBox) {
// 	ib.inputMode = input_mode_NORMAL
// 	if len(ap.stream) == 0 {
// 		return
// 	}
// 	if ap.stream[0].lastCommand == commandPAUSE {
// 		Action_Continue(ap, ib)
// 	}
// 	if ap.stream[0].lastCommand == commandNONE {
// 		Action_Continue(ap, ib)
// 	}
// }

// func DesidionConfirm(ap *allProc, ib *InfoBox) error {
// 	switch ib.inputMode {
// 	case input_mode_WAIT_CONFIRM:
// 		ib.inputMode = input_mode_CONFIRM_RECEIVED
// 	}
// 	return nil
// }

// func DesidionDeny(ap *allProc, ib *InfoBox) error {
// 	switch ib.inputMode {
// 	case input_mode_WAIT_CONFIRM:
// 		ib.inputMode = input_mode_DENIAL_RECEIVED
// 	}
// 	return nil
// }

// func Action_UndoMovement(ap *allProc, ib *InfoBox) error {
// 	Action_Pause(ap, ib)
// 	curs := saveCursor(ap, ib)
// 	ap.indexBuf.DeleteLast()
// 	lastState := ap.indexBuf.LastState()
// 	ap.arrangeStreamsBy(lastState)
// 	ib.cursor = restoreCursor(ap.stream, curs)
// 	Action_Continue(ap, ib)
// 	return nil
// }

// func Action_SelectAllWithSameExtention(ap *allProc, ib *InfoBox) error {
// 	curs := saveCursor(ap, ib)
// 	ext := namedata.RetrieveExtention(curs)
// 	for i, stream := range ap.stream {
// 		if strings.HasSuffix(stream.source, ext) {
// 			ap.stream[i].isSelected = true
// 		}
// 	}
// 	return nil
// }

// func Action_DeleteSelected(ap *allProc, ib *InfoBox) error {
// 	Action_Pause(ap, ib)
// 	toDelete := []int{}
// 	for i, stream := range ap.stream {
// 		if stream.isSelected {
// 			toDelete = append(toDelete, i)
// 		}
// 	}
// 	positionsProjected := []int{}
// 	for j, i := range toDelete {
// 		positionsProjected = append(positionsProjected, i-j)
// 	}
// 	switchToWaitConfirmMode(ap, ib)
// 	go func() {
// 		err := ap.deleteStreams(positionsProjected, ib)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		if err == nil {
// 			switchToNORMALMode(ap, ib)
// 		}
// 	}()
// 	return nil
// }

// func (ap *allProc) deleteStreams(streams []int, ib *InfoBox) error {
// 	concluded := false
// 	for !concluded {
// 		switch ib.inputMode {
// 		case input_mode_CONFIRM_RECEIVED:
// 			concluded = true
// 		case input_mode_DENIAL_RECEIVED:
// 			return nil
// 		default:
// 			time.Sleep(time.Millisecond * 50)
// 		}
// 	}
// 	time.Sleep(time.Millisecond * 500)
// 	for _, i := range streams {
// 		if err := ap.DeleteStream(i); err != nil {
// 			panic("DELETE ERR (я хз откуда она может взяться)" + err.Error())
// 			return err
// 		}
// 	}
// 	time.Sleep(time.Millisecond * 200)
// 	return nil
// }

// func (ap *allProc) DeleteStream(i int) error {
// 	if len(ap.stream) < i {
// 		return fmt.Errorf(" DeleteStream(): can not delete stream: i > len(ap.stream)")
// 	}
// 	stream := ap.stream[i]
// 	ap.stream = append(ap.stream[:i], ap.stream[i+1:]...)
// 	ap.indexBuf.Remove(stream.source)
// 	ap.activeHandlerChan = nil
// 	return nil
// }

// func Action_AddNewProcess(ap *allProc, ib *InfoBox) error {
// 	cb := string(clipboard.Read(clipboard.FmtText))

// 	paths := strings.Split(cb, "\r\n")
// 	if strings.Join(paths, "") == cb {
// 		paths = strings.Split(cb, "\n")
// 	}
// 	list := []string{}
// 	for _, path := range paths {
// 		exists, err := exists(path)
// 		if err != nil {
// 			return err
// 		}

// 		if exists {
// 			list = append(list, path)
// 		}
// 	}
// 	list = sortnames.GrabberOrder(list)
// 	ap.NewProcesses(dest_gl, list...)

// 	ib.lowBorder = 0
// 	ib.highBorder = 0
// 	return nil
// }

// func Action_QUIT_PROGRAM(ap *allProc, ib *InfoBox) error {
// 	return fmt.Errorf("Quit action called by user")
// }

// /*
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_03_PRT230421102814_AUDIOENG20_proxy.ac3
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_03_PRT230421102814_AUDIORUS51.m4a
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_03_PRT230421102814_AUDIORUS51_proxy.ac3
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_03_PRT230421102814_HD.mp4
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_03_PRT230421102814_HD_proxy.mp4
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_04_PRT230421104337.srt
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_04_PRT230421104337_AUDIOENG20.m4a
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_04_PRT230421104337_AUDIOENG20_proxy.ac3
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_04_PRT230421104337_AUDIORUS51.m4a
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_04_PRT230421104337_AUDIORUS51_proxy.ac3
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_04_PRT230421104337_HD.mp4

// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_03_PRT230421102814_AUDIOENG20.m4a
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_03_PRT230421102814_AUDIOENG20_proxy.ac3
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_03_PRT230421102814_AUDIORUS51.m4a
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_03_PRT230421102814_AUDIORUS51_proxy.ac3
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_03_PRT230421102814_HD.mp4
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_07_PRT230421103310_AUDIOENG20.m4a
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_07_PRT230421103310_AUDIOENG20_proxy.ac3
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_07_PRT230421103310_AUDIORUS51.m4a
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_07_PRT230421103310_AUDIORUS51_proxy.ac3
// \\nas\ROOT\EDIT\_amedia\Barri_s01\Barri_s01_07_PRT230421103310_HD.mp4
// +---5 file down----------------------------------------+
// [F1=HELP] [F2=STATS] [F10=QUIT]===================================================
// */
