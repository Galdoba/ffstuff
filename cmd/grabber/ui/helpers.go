package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nsf/termbox-go"
)

func commandSequanceToTBKey(s string) (termbox.Key, error) {
	tk := termbox.KeyF1
	err := fmt.Errorf("no err")
	s = strings.ToUpper(s)
	switch s {

	default:
		err = fmt.Errorf("rune not found: ''", s)
	case "F1":
		tk = termbox.KeyF1
	case "F2":
		tk = termbox.KeyF2
	case "F3":
		tk = termbox.KeyF3
	case "F4":
		tk = termbox.KeyF4
	case "F5":
		tk = termbox.KeyF5
	case "F6":
		tk = termbox.KeyF6
	case "F7":
		tk = termbox.KeyF7
	case "F8":
		tk = termbox.KeyF8
	case "F9":
		tk = termbox.KeyF9
	case "F10":
		tk = termbox.KeyF10
	case "F11":
		tk = termbox.KeyF11
	case "F12":
		tk = termbox.KeyF12
	case "INSERT":
		tk = termbox.KeyInsert
	case "DELETE":
		tk = termbox.KeyDelete
	case "HOME":
		tk = termbox.KeyHome
	case "END":
		tk = termbox.KeyEnd
	case "PGUP":
		tk = termbox.KeyPgup
	case "PGDN":
		tk = termbox.KeyPgdn
	case "UP":
		tk = termbox.KeyArrowUp
	case "DOWN":
		tk = termbox.KeyArrowDown
	case "LEFT":
		tk = termbox.KeyArrowLeft
	case "RIGHT":
		tk = termbox.KeyArrowRight
	case "LMB":
		tk = termbox.MouseLeft
	case "MMB":
		tk = termbox.MouseMiddle
	case "RMB":
		tk = termbox.MouseRight
	case "MR":
		tk = termbox.MouseRelease
	case "MWUP":
		tk = termbox.MouseWheelUp
	case "MWDOWN":
		tk = termbox.MouseWheelDown
	case "SPACE":
		tk = termbox.KeySpace
	case "ENTER":
		tk = termbox.KeyEnter
	case "CTRL+~":
		tk = termbox.KeyCtrlTilde
	case "CTRL+2":
		tk = termbox.KeyCtrl2
	case "CTRL+SPACE":
		tk = termbox.KeyCtrlSpace
	case "CTRL+A":
		tk = termbox.KeyCtrlA
	case "CTRL+B":
		tk = termbox.KeyCtrlB
	case "CTRL+C":
		tk = termbox.KeyCtrlC
	case "CTRL+D":
		tk = termbox.KeyCtrlD
	case "CTRL+E":
		tk = termbox.KeyCtrlE
	case "CTRL+F":
		tk = termbox.KeyCtrlF
	case "CTRL+G":
		tk = termbox.KeyCtrlG
	case "BACKSPACE":
		tk = termbox.KeyBackspace
	case "CTRL+H":
		tk = termbox.KeyCtrlH
	case "TAB":
		tk = termbox.KeyTab
	case "CTRL+I":
		tk = termbox.KeyCtrlI
	case "CTRL+J":
		tk = termbox.KeyCtrlJ
	case "CTRL+K":
		tk = termbox.KeyCtrlK
	case "CTRL+L":
		tk = termbox.KeyCtrlL
	case "CTRL+M":
		tk = termbox.KeyCtrlM
	case "CTRL+N":
		tk = termbox.KeyCtrlN
	case "CTRL+O":
		tk = termbox.KeyCtrlO
	case "CTRL+P":
		tk = termbox.KeyCtrlP
	case "CTRL+Q":
		tk = termbox.KeyCtrlQ
	case "CTRL+R":
		tk = termbox.KeyCtrlR
	case "CTRL+S":
		tk = termbox.KeyCtrlS
	case "CTRL+T":
		tk = termbox.KeyCtrlT
	case "CTRL+U":
		tk = termbox.KeyCtrlU
	case "CTRL+V":
		tk = termbox.KeyCtrlV
	case "CTRL+W":
		tk = termbox.KeyCtrlW
	case "CTRL+X":
		tk = termbox.KeyCtrlX
	case "CTRL+Y":
		tk = termbox.KeyCtrlY
	case "CTRL+Z":
		tk = termbox.KeyCtrlZ
	case "ESC":
		tk = termbox.KeyEsc
		// case "":
		// tk = termbox.KeyCtrlLsqBracket
	case "CTRL+3":
		tk = termbox.KeyCtrl3
	case "CTRL+4":
		tk = termbox.KeyCtrl4
	case "CTRL+5":
		tk = termbox.KeyCtrl5
	case "CTRL+6":
		tk = termbox.KeyCtrl6
	case "CTRL+7":
		tk = termbox.KeyCtrl7
	case "CTRL+8":
		tk = termbox.KeyCtrl8
	}
	if err.Error() != "no err" {
		return tk, err
	}
	return tk, nil
}

func map_evCh(r rune) string {
	switch r {
	case 9:
		return "\t"
	case 32:
		return " "
	case 44, 60, 1073, 1041:
		return "<"
	case 45, 95:
		return "-"
	case 46, 62, 1102, 1070:
		return ">"
	case 47, 63:
		return "?"
	case 48, 41:
		return "0"
	case 49, 33:
		return "1"
	case 50, 64, 34:
		return "2"
	case 51, 35, 8470:
		return "3"
	case 52, 36:
		return "4"
	case 53, 37:
		return "5"
	case 54, 94:
		return "6"
	case 55, 38:
		return "7"
	case 56, 42:
		return "8"
	case 57, 40:
		return "9"
	case 59, 58, 1078, 1046:
		return ":"
	case 61, 43:
		return "="
	case 92:
		return "\\"
	case 124:
		return "|"
	case 91, 123, 1093, 1061:
		return "{"
	case 93, 125, 1098, 1066:
		return "}"
	case 96, 126, 1105, 1025:
		return "~"
	case 97, 65, 1092, 1060:
		return "A"
	case 98, 66, 1080, 1048:
		return "B"
	case 99, 67, 1089, 1057:
		return "C"
	case 100, 68, 1074, 1042:
		return "D"
	case 101, 69, 1091, 1059:
		return "E"
	case 102, 70, 1072, 1040:
		return "F"
	case 103, 71, 1087, 1055:
		return "G"
	case 104, 72, 1088, 1056:
		return "H"
	case 105, 73, 1096, 1064:
		return "I"
	case 106, 74, 1086, 1054:
		return "J"
	case 107, 75, 1083, 1051:
		return "K"
	case 108, 76, 1076, 1044:
		return "L"
	case 109, 77, 1100, 1068:
		return "M"
	case 110, 78, 1090, 1058:
		return "N"
	case 111, 79, 1097, 1065:
		return "O"
	case 112, 80, 1079, 1047:
		return "P"
	case 113, 81, 1081, 1049:
		return "Q"
	case 114, 82, 1082, 1050:
		return "R"
	case 115, 83, 1099, 1067:
		return "S"
	case 116, 84, 1077, 1045:
		return "T"
	case 117, 85, 1075, 1043:
		return "U"
	case 118, 86, 1084, 1052:
		return "V"
	case 119, 87, 1094, 1062:
		return "W"
	case 120, 88, 1095, 1063:
		return "X"
	case 121, 89, 1085, 1053:
		return "Y"
	case 122, 90, 1103, 1071:
		return "Z"

	default:
		panic(strconv.QuoteRuneToGraphic(r) + " " + strconv.QuoteRune(rune(r)) + " " + fmt.Sprintf("%v", r))
	}
}

func runeToKey(ev termbox.Key) string {
	tbR := ev
	//panic(ev)
	switch tbR {
	default:
		return ""
	case termbox.KeyF1:
		return "F1"
	case termbox.KeyF2:
		return "F2"
	case termbox.KeyF3:
		return "F3"
	case termbox.KeyF4:
		return "F4"
	case termbox.KeyF5:
		return "F5"
	case termbox.KeyF6:
		return "F6"
	case termbox.KeyF7:
		return "F7"
	case termbox.KeyF8:
		return "F8"
	case termbox.KeyF9:
		return "F9"
	case termbox.KeyF10:
		return "F10"
	case termbox.KeyF11:
		return "F11"
	case termbox.KeyF12:
		return "F12"
	case termbox.KeyInsert:
		return "Insert"
	case termbox.KeyDelete:
		return "Delete"
	case termbox.KeyHome:
		return "Home"
	case termbox.KeyEnd:
		return "End"
	case termbox.KeyPgup:
		return "PgUp"
	case termbox.KeyPgdn:
		return "PgDn"
	case termbox.KeyArrowUp:
		return "Up"
	case termbox.KeyArrowDown:
		return "Down"
	case termbox.KeyArrowLeft:
		return "Left"
	case termbox.KeyArrowRight:
		return "Right"
	case termbox.MouseLeft:
		return "LMB (unsupported)"
	case termbox.MouseMiddle:
		return "MMB (unsupported)"
	case termbox.MouseRight:
		return "RMB (unsupported)"
	case termbox.MouseRelease:
		return "MR (unsupported)"
	case termbox.MouseWheelUp:
		return "MWU (unsupported)"
	case termbox.MouseWheelDown:
		return "MWD (unsupported)"
	case termbox.KeySpace:
		return "Space"
	case termbox.KeyEnter:
		return "Enter"
	case termbox.KeyCtrlA:
		return "Ctrl+A"
	case termbox.KeyCtrlB:
		return "Ctrl+B"
	case termbox.KeyCtrlC:
		return "Ctrl+C"
	case termbox.KeyCtrlD:
		return "Ctrl+D"
	case termbox.KeyCtrlE:
		return "Ctrl+E"
	case termbox.KeyCtrlF:
		return "Ctrl+F"
	case termbox.KeyCtrlG:
		return "Ctrl+G"
	case termbox.KeyBackspace:
		return "Backspace"
	case termbox.KeyTab:
		return "Tab"
	case termbox.KeyCtrlSpace:
		return "Ctrl+Space"
	case termbox.KeyCtrlJ:
		return "Ctrl+J"
	case termbox.KeyCtrlK:
		return "Ctrl+K"
	case termbox.KeyCtrlL:
		return "Ctrl+L"
	case termbox.KeyCtrlN:
		return "Ctrl+N"
	case termbox.KeyCtrlO:
		return "Ctrl+O"
	case termbox.KeyCtrlP:
		return "Ctrl+P"
	case termbox.KeyCtrlQ:
		return "Ctrl+Q"
	case termbox.KeyCtrlR:
		return "Ctrl+R"
	case termbox.KeyCtrlS:
		return "Ctrl+S"
	case termbox.KeyCtrlT:
		return "Ctrl+T"
	case termbox.KeyCtrlU:
		return "Ctrl+U"
	case termbox.KeyCtrlV:
		return "Ctrl+V"
	case termbox.KeyCtrlW:
		return "Ctrl+W"
	case termbox.KeyCtrlX:
		return "Ctrl+X"
	case termbox.KeyCtrlY:
		return "Ctrl+Y"
	case termbox.KeyCtrlZ:
		return "Ctrl+Z"
	case termbox.KeyEsc:
		return "Esc"
	case termbox.KeyCtrlBackslash:
		return "Backslash"
	case termbox.KeyCtrlRsqBracket:
		return "Ctrl+]" //уточнить
	case termbox.KeyCtrl6:
		return "Ctrl+6"
	case termbox.KeyCtrl7:
		return "Ctrl+7"
	case termbox.KeyCtrl8:
		return "Ctrl+8"
	}
}
