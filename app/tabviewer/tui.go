package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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
