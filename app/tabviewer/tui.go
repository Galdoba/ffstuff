package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gookit/color"
)

func (td tableData) Init() tea.Cmd {
	return nil
}

func FormatData(td tableData) []string {
	viewLen := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	out := ""
	outText := []string{}
	for row, line := range td.data {
		if td.hiddenRows[row] {
			continue
		}
		for col, text := range line {
			switch col {
			default:
				panic("не знаю что делать с колонкой " + fmt.Sprintf("%v", col))
			case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14:
				out += text + " "
				trueLen := len(strings.Split(text, ""))
				if viewLen[col] < trueLen {
					viewLen[col] = trueLen
				}
			}
		}
		out += "\n"
		outText = append(outText, out)
	}
	return outText
}

func (td tableData) View() string {

	return strings.Join(FormatData(td), "\n")
}

func (td tableData) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q":
			return td, tea.Quit

		}

	}
	return td, nil
}

// func FormatLine(line []string, width int) string {
// 	textByLetter := [][]string{}
// 	for _, cell := range line {
// 		textByLetter = append(textByLetter, strings.Split(cell, ""))
// 	}
// 	//coment 0
// 	//path 1
// 	//Trailer 4
// 	//Poster 6
// 	//Name 8
// 	//Agent 13
// 	//Date 14
// 	out := "|"
// 	for i, cellByLetter := range textByLetter {
// 		datatype := ""
// 		switch i {
// 		case 4:
// 			datatype = "t"
// 		case 6:
// 			datatype = "p"
// 		case 0, 1, 8, 13, 14:
// 			datatype = "f"
// 		default:
// 			continue

// 		}
// 		switch len(cellByLetter) {

// 		default:
// 			text := strings.Join(cellByLetter, "")
// 			col := colorCode(textByLetter, datatype)
// 			text = colorText(text, col)
// 			out += text + "|"
// 		case 0:
// 			out += "???|"
// 		}
// 	}
// 	out += fmt.Sprintf("%v", width)
// 	//fmt.Println(out)
// 	return out
// }

func columnSizes(data [][]string) []int {
	columnLen := []int{}
	for i, line := range data {
		if i < 2 {
			continue
		}
		for _, cell := range line {
			l := strings.Split(cell, "")
			columnLen = append(columnLen, len(l))
		}
		break
	}
	for i, line := range data {
		if i < 2 {
			continue
		}
		for j, cell := range line {
			text := strings.Split(cell, "")
			if len(text) > columnLen[j] {
				columnLen[j] = len(text)
			}
		}
	}
	return columnLen
}

func FormatLineSize(line []string, widths []int) []string {
	out := []string{}

	for r, word := range line {
		word = widen(word, widths[r])
		out = append(out, word)
	}

	return out
}

// const (
// 	colorWhite = iota
// 	colorGreen
// 	colorYellow
// 	colorRed
// 	colorCyan
// 	colorBlue
// )

// func colorText(text string, colorCode int) string {
// 	switch colorCode {
// 	default:
// 		return text
// 	case colorGreen:
// 		return color.GreenString(text)
// 	case colorYellow:
// 		return color.YellowString(text)
// 	case colorRed:
// 		return color.RedString(text)
// 	case colorCyan:
// 		return color.CyanString(text)
// 	}
// }

// func colorCode(textByLetter [][]string, datatype string) int {
// 	switch datatype {
// 	default:
// 		return colorWhite
// 	case "t":
// 		col := colorWhite
// 		if strings.Join(textByLetter[4], "") != "" {
// 			col = colorRed
// 		}
// 		codeByCell := colorDataToCode(strings.Join(textByLetter[3], ""))
// 		switch codeByCell {
// 		case colorWhite:
// 			return col
// 		default:
// 			return codeByCell
// 		}
// 	case "p":
// 		return colorDataToCode(strings.Join(textByLetter[5], ""))
// 	case "f":
// 		return colorDataToCode(strings.Join(textByLetter[10], ""))
// 	}
// }

// func colorDataToCode(colorData string) int {
// 	switch colorData {
// 	default:
// 		return colorWhite
// 	case "R", "r", "К", "к":
// 		return colorRed
// 	case "Y", "y", "Н", "н":
// 		return colorYellow
// 	case "G", "g", "П", "п":
// 		return colorGreen
// 	case "B", "b", "И", "и":
// 		return colorCyan

// 	}
// }

type outPreset struct {
	Name string
}

type columnData struct {
	Key                string                       //A
	Comment            string                       //What is it do?
	Hidden             bool                         //true
	AllowBrake         bool                         //false
	ColorRule          func(string) *color.Style256 //DefaultRule(coord.String()) *color.Style256
	MaxWidth           int                          //999
	OutSourceScript    string                       //path/to/script
	OutSourceArguments []string                     //script arguments
}
