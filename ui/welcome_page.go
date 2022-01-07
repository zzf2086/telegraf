package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type WelcomePage struct {
	Tabs       []string
	TabContent []list.Model

	activatedTab int
	version      string

	tutorialContent    string
	selectedMenuOption *Item
}

func NewWelcomePage(v string) WelcomePage {
	tabs := []string{
		"Welcome Page",
		"Tutorial",
	}

	itemsWelcome := []list.Item{
		Item{ItemTitle: "Show Plugins", Desc: "All the plugins supported by Telegraf"},
		Item{ItemTitle: "Show Flags", Desc: "Flags come with Telegraf"},
	}

	itemsTutorial := []list.Item{}

	var tabcontent []list.Model
	welcomePageOptions := list.NewModel(itemsWelcome, list.NewDefaultDelegate(), defaultWidth, listHeight)
	tutorialPageOptions := list.NewModel(itemsTutorial, list.NewDefaultDelegate(), defaultWidth, listHeight)
	tabcontent = append(tabcontent, welcomePageOptions)
	tabcontent = append(tabcontent, tutorialPageOptions)

	// Tutorial Content
	in := `# Telegraf

## Intro

Telegraf is an agent for collecting, processing, aggregating, and writing metrics. Based on a plugin system to enable developers in the community to easily add support for additional metric collection. There are *four* distinct types of plugins:

1. **Input** Plugins collect metrics from the system, services, or 3rd party APIs
2. **Processor** Plugins transform, decorate, and/or filter metrics
3. **Aggregator** Plugins create aggregate metrics (e.g. mean, min, max, quantiles, etc.)
4. **Output** Plugins write metrics to various destinations
	`
	r, _ := glamour.NewTermRenderer(
		// detect background color and pick either the default dark or light theme
		glamour.WithAutoStyle(),
	)
	out, _ := r.Render(in)

	return WelcomePage{
		Tabs:            tabs,
		TabContent:      tabcontent,
		tutorialContent: out,
		version:         v,
	}
}

func (w *WelcomePage) Update(m tea.Model, msg tea.Msg) (tea.Model, tea.Cmd, int) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit, WELCOME_PAGE
		case "right":
			if w.activatedTab < len(w.Tabs)-1 {
				w.activatedTab++
			}
			return m, nil, WELCOME_PAGE
		case "left":
			if w.activatedTab > 0 {
				w.activatedTab--
			}
			return m, nil, WELCOME_PAGE
		case "enter":
			listItem := w.TabContent[w.activatedTab].SelectedItem()
			i, ok := listItem.(Item)
			if !ok {
				return m, nil, WELCOME_PAGE
			}
			if i.ItemTitle == "Show Plugins" {
				return m, nil, PLUGIN_PAGE
			} else if i.ItemTitle == "Show Flags" {
				return m, nil, FLAG_PAGE
			}

		}
	}
	var cmd tea.Cmd
	w.TabContent[w.activatedTab], cmd = w.TabContent[w.activatedTab].Update(msg)
	return m, cmd, WELCOME_PAGE
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

	if w.activatedTab == 0 {
		// Welcome Page tab
		s := "Welcome to Telegraf! 🥳 \n\n"
		s += fmt.Sprintf("You are on %s \n\n", w.version)
		_, err := doc.WriteString(s)
		if err != nil {
			return err.Error()
		}
		// list
		_, err = doc.WriteString(w.TabContent[w.activatedTab].View())
		if err != nil {
			return err.Error()
		}
	} else {
		// Tutorial tab
		_, err := doc.WriteString(w.tutorialContent)
		if err != nil {
			return err.Error()
		}
	}

	// Set default font color to be Influx color
	// Color code is from https://influxdata.github.io/branding/visual/colors/
	docStyle.Foreground(lipgloss.Color("#BF2FE5"))

	return docStyle.Render(doc.String())
}
