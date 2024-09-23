package ui

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/Galdoba/ffstuff/cmd/grabber/download"
// 	"github.com/Galdoba/ffstuff/pkg/namedata"
// 	"github.com/Galdoba/utils"
// 	"github.com/nsf/termbox-go"
// )

// func statusData(ap *allProc) (int, int, int, int) {
// 	wi, wa, dn, tt := 0, 0, 0, 0
// 	for _, stream := range ap.stream {
// 		tt++
// 		if stream.warning != "" {
// 			wa++
// 			continue
// 		}
// 		if stream.handler == nil {
// 			wi++
// 			continue
// 		}
// 		if stream.handler.Status() == download.STATUS_COMPLETED {
// 			dn++
// 			continue
// 		}
// 		wi++
// 	}
// 	return wi, wa, dn, tt
// }

// func (ib *InfoBox) Update(ap *allProc) error {
// 	newData := []string{}

// 	switch ib.inputMode {
// 	default:
// 		return fmt.Errorf("unknown input mode: %v", ib.inputMode)
// 	case input_mode_NORMAL:
// 		for _, pr := range ap.stream {
// 			if pr.warning == "duplicate" {
// 				newData = append(newData, pr.ErrString())
// 				continue
// 			}
// 			if pr.handler == nil {
// 				newData = append(newData, pr.QueueString())
// 				continue
// 			}
// 			switch pr.handler.Status() {
// 			default:
// 				panic(pr.handler.Status()) //не должно срабатывать
// 			case download.STATUS_COMPLETED:
// 				newData = append(newData, pr.CompleteString())
// 				// if ib.cursor <= i {
// 				// 	ib.cursor++
// 				// }
// 			case download.STATUS_ERR:
// 				newData = append(newData, pr.ErrString())
// 			case download.STATUS_TRANSFERING, download.STATUS_PAUSED, download.STATUS_NIL:
// 				newData = append(newData, pr.QueueString())
// 			}

// 		}

// 	case input_mode_WAIT_CONFIRM:
// 		newData = append(newData, "Press Enter to confirm or Esc to deny")
// 		for _, pr := range ap.stream {
// 			newData = append(newData, pr.QueueString())
// 		}
// 		//panic("not expecting confirm mode")
// 		//pos, sel := ap.streamDataProposition.Export()
// 		// for i, _ := range sel {
// 		// 	pr := ap.stream[pos[i]].String()
// 		// 	newData = append(newData, pr)
// 		// }
// 		// case input_mode_CONFIRM_RECEIVED, input_mode_DENIAL_RECEIVED:
// 		// 	switchToNORMALMode(ap, ib)
// 	}
// 	ib.data = newData
// 	return nil
// }

// func (ib *InfoBox) Draw(ap *allProc) {
// 	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
// 	w, h := termbox.Size()
// 	fg := termbox.ColorWhite
// 	bg := termbox.ColorBlack

// 	ib.drawLen = h - len(helpBlock) - 6
// 	//tkr := tickerImage(ib.ticker / 5)
// 	prog := ""
// 	dest := ""
// 	if ap.activeStream != nil {
// 		szStr := ""
// 		if ap.activeStream.handler != nil {

// 			szStr = ap.activeStream.ProgressData()
// 		}
// 		prog = ap.activeStream.Progress() + szStr
// 		dest = namedata.RetrieveShortName(ap.activeStream.source) + "  ----->  " + ap.activeStream.dest
// 	}
// 	//tbprint(0, 0, fg, bg, "Last Key Pressed:"+ib.lastKeysPressed+"__: "+fmt.Sprintf("%v", len(ap.indexBuf.Set))+" ccl:"+fmt.Sprintf("%v", ib.cursor))
// 	wi, wa, dn, tt := statusData(ap)
// 	statusData := fmt.Sprintf("List Status: [Wait]=%v [Warn]=%v [Done]=%v [Total]=%v", wi, wa, dn, tt)
// 	tbprint(0, 0, fg, bg, "Active Transfert: "+prog)
// 	tbprint(0, 1, fg, bg, dest)
// 	tbprint(0, 2, fg, bg, statusData)
// 	if ib.lowBorder == 0 && ib.highBorder == 0 {
// 		ib.lowBorder = utils.Max(0, ib.cursor-ib.drawLen)
// 		ib.highBorder = utils.Min(len(ap.stream)-1, ib.lowBorder+ib.drawLen)
// 	}

// 	switch ib.lastScroll { //0 = down/ 1 = up //ЧТО-ТО НЕ РАБОТАЕТ ОТ ЗПвт
// 	case 0:

// 		if ib.cursor > ib.highBorder {
// 			ib.lowBorder = utils.Max(0, ib.cursor-ib.drawLen)
// 			ib.highBorder = utils.Min(len(ap.stream)-1, ib.lowBorder+ib.drawLen)
// 		}
// 	case 1:

// 		for ib.cursor < ib.lowBorder {
// 			ib.lowBorder--
// 			ib.highBorder--
// 		}
// 	default:
// 		panic("Этого не должно случиться:\nib.lastScroll")
// 	}

// 	uptag := ""
// 	downtag := ""
// 	if ib.lowBorder > 0 {
// 		uptag = fmt.Sprintf("%v files up", ib.lowBorder)
// 	}
// 	uptag = "+---" + uptag
// 	if ib.highBorder < len(ap.stream)-1 {
// 		downtag = fmt.Sprintf("%v files down", len(ap.stream)-1-ib.highBorder)
// 	}
// 	mld := w - 5
// 	for len(uptag) < mld {
// 		uptag += "-"
// 	}
// 	downtag = "+---" + downtag
// 	for len(downtag) < mld {
// 		downtag += "-"
// 	}

// 	uptag += "+"
// 	downtag += "+"
// 	tbprint(0, 3, fg, bg, uptag)
// 	lastDrowed := 0
// 	for i, data := range ib.data {
// 		if i < ib.lowBorder {
// 			continue
// 		}
// 		if i > ib.highBorder {
// 			continue
// 		}
// 		fg = activeColor(data)
// 		bg = termbox.ColorBlack
// 		if i == ib.cursor {
// 			fg = termbox.ColorBlack
// 			bg = activeColor(data)
// 		}
// 		tbprint(0, i+4-ib.lowBorder, termbox.ColorWhite, termbox.ColorBlack, "|")
// 		tbprint(2, i+4-ib.lowBorder, fg, bg, strings.TrimSuffix(data, " $"))
// 		tbprint(mld, i+4-ib.lowBorder, termbox.ColorWhite, termbox.ColorBlack, "|")
// 		lastDrowed = i + 4 - ib.lowBorder
// 		fg = activeColor(data)
// 		bg = termbox.ColorBlack
// 	}
// 	tbprint(0, lastDrowed+1, termbox.ColorWhite, termbox.ColorBlack, downtag)

// 	for i, line := range helpBlock {
// 		tbprint(0, h-len(helpBlock)+i, termbox.ColorWhite, termbox.ColorBlack, line)
// 	}
// 	termbox.Flush()
// }
