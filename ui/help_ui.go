package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Pages interface {
	Update(tea.Model, tea.Msg) (tea.Model, tea.Cmd)
	View() string
}

type HelpUI struct {

	// Welcome Screen (main page to get to other pages)
	// Tutorial Screen (Guide to how Telegraf works, like the getting started page in the docs)
	// Plugin List screen (usage, all plugins listed)
	pages []Pages

	currentPage int
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

	return HelpUI{pages: pages, currentPage: 2}
}

func (m HelpUI) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m HelpUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.pages[m.currentPage].Update(m, msg)
}

func (m HelpUI) View() string {
	return m.pages[m.currentPage].View()
}
