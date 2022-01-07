package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

const (
	WELCOME_PAGE = 0
	PLUGIN_PAGE  = 1
	FLAG_PAGE    = 2
)

var (
	currentPage = WELCOME_PAGE
)

type Pages interface {
	Update(tea.Model, tea.Msg) (tea.Model, tea.Cmd, int)
	View() string
}

type HelpUI struct {

	// Welcome Page (main page to get to other pages)
	// 		Welcome Tab
	// 		Tutorial Tab (Guide to how Telegraf works, explain plugins)
	// Plugin List Page (usage, all plugins listed)
	// Flag List Page (usage, all plugins listed)
	pages []Pages
}

func NewHelpUI(version string) HelpUI {
	w := NewWelcomePage(version)
	p := NewPluginPage()
	f := NewFlagsPage()

	var pages []Pages
	pages = append(pages, &w)
	pages = append(pages, &p)
	pages = append(pages, &f)

	return HelpUI{pages: pages}
}

func (m HelpUI) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m HelpUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd, pageNumber := m.pages[currentPage].Update(m, msg)
	currentPage = pageNumber
	return model, cmd
}

func (m HelpUI) View() string {
	return m.pages[currentPage].View()
}
