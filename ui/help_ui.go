package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

var (
	currentPage = 0
)

type Pages interface {
	Update(tea.Model, tea.Msg) (tea.Model, tea.Cmd, int)
	View() string
}

type HelpUI struct {

	// Welcome Screen (main page to get to other pages)
	// Tutorial Screen (Guide to how Telegraf works, like the getting started page in the docs)
	// Plugin List screen (usage, all plugins listed)
	pages []Pages
}

func NewHelpUI(version string) HelpUI {
	// [ Welcome ] [ Tutorial ]
	// Telegraf, the plugin thingie
	//
	// > Show Plugins
	// > Show Flags

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
