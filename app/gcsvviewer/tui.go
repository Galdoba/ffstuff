package main

import tea "github.com/charmbracelet/bubbletea"

func (td tableData) Init() tea.Cmd {
	return nil
}

func (td tableData) View() string {
	return "nil"
}

func (td tableData) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return td, nil
}
