package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 20

var version string

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type WelcomePage struct {
	Tabs       []string
	TabContent []list.Model

	activatedTab int
}

func NewWelcomePage(v string) WelcomePage {
	tabs := []string{
		"Welcome Page",
		"Tutorial",
	}
	version = v

	items_welcome := []list.Item{
		item{title: "Show Plugins", desc: "All the plugins supported by Telegraf"},
		item{title: "Show Flags", desc: "Flags come with Telegraf"},
	}

	items_tutorial := []list.Item{
		item{title: "How Telegraf works", desc: "..."},
		item{title: "Set up", desc: "..."},
		item{title: "Showing me around", desc: "..."},
	}

	var tabcontent []list.Model
	defaultWidth = 40
	tabcontent = append(tabcontent, list.NewModel(items_welcome, list.NewDefaultDelegate(), defaultWidth, listHeight))
	tabcontent = append(tabcontent, list.NewModel(items_tutorial, list.NewDefaultDelegate(), defaultWidth, listHeight))
	return WelcomePage{Tabs: tabs, TabContent: tabcontent}
}

func (w *WelcomePage) Update(m tea.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right":
			if w.activatedTab < len(w.Tabs)-1 {
				w.activatedTab++
			}
			return m, nil
		case "left":
			if w.activatedTab > 0 {
				w.activatedTab--
			}
			return m, nil
		}
	}
	var cmd tea.Cmd
	w.TabContent[w.activatedTab], cmd = w.TabContent[w.activatedTab].Update(msg)
	return m, cmd
}

func (w *WelcomePage) View() string {
	doc := strings.Builder{}

	// Tabs
	{
		var renderedTabs []string

		for i, t := range w.Tabs {
			if i == w.activatedTab {
				renderedTabs = append(renderedTabs, activeTab.Render(t))
			} else {
				renderedTabs = append(renderedTabs, tab.Render(t))
			}
		}

		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			renderedTabs...,
		)
		gap := tabGap.Render(strings.Repeat(" ", max(0, defaultWidth-lipgloss.Width(row)-2)))
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
		_, err := doc.WriteString(row + "\n\n")
		if err != nil {
			return err.Error()
		}
	}

	// Welcome message
	{
		if w.activatedTab == 0 { // Welcome Page tab
			s := "Welcome to Telegraf! 🥳 \n\n"
			s += fmt.Sprintf("You are on %s \n\n", version)
			_, err := doc.WriteString(s)
			if err != nil {
				return err.Error()
			}
		}
	}

	// list
	_, err := doc.WriteString(w.TabContent[w.activatedTab].View())
	if err != nil {
		return err.Error()
	}

	// style
	docStyle.Foreground(lipgloss.Color("#BF2FE5"))

	return docStyle.Render(doc.String())
}
