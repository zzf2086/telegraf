package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/influxdata/telegraf/plugins/aggregators"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/outputs"
	"github.com/influxdata/telegraf/plugins/processors"
)

type Item struct {
	ItemTitle, Desc string
}

func (i Item) Title() string       { return i.ItemTitle }
func (i Item) Description() string { return i.Desc }
func (i Item) FilterValue() string { return i.ItemTitle }

type PluginPage struct {
	Tabs       []string
	TabContent []list.Model

	activatedTab int
}

func createPluginList(content []list.Item) list.Model {
	pluginList := list.NewModel(content, list.NewDefaultDelegate(), 50, listHeight)
	pluginList.SetShowStatusBar(false)
	pluginList.SetShowTitle(false)

	return pluginList
}

func NewPluginPage() PluginPage {
	tabs := []string{
		"Inputs",
		"Outputs",
		"Aggregators",
		"Processors",
	}

	var inputContent, outputContent, aggregatorContent, processorContent []list.Item

	for name, creator := range inputs.Inputs {
		inputContent = append(inputContent, Item{ItemTitle: name, Desc: creator().Description()})
	}

	for name, creator := range outputs.Outputs {
		outputContent = append(outputContent, Item{ItemTitle: name, Desc: creator().Description()})
	}

	for name, creator := range aggregators.Aggregators {
		aggregatorContent = append(aggregatorContent, Item{ItemTitle: name, Desc: creator().Description()})
	}

	for name, creator := range processors.Processors {
		processorContent = append(processorContent, Item{ItemTitle: name, Desc: creator().Description()})
	}

	var t []list.Model
	t = append(t, createPluginList(inputContent))
	t = append(t, createPluginList(outputContent))
	t = append(t, createPluginList(aggregatorContent))
	t = append(t, createPluginList(processorContent))

	return PluginPage{Tabs: tabs, TabContent: t}
}

func (p *PluginPage) Update(m tea.Model, msg tea.Msg) (tea.Model, tea.Cmd, int) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit, 1
		case "right":
			if p.activatedTab < len(p.Tabs)-1 {
				p.activatedTab++
			}
			return m, nil, 1
		case "left":
			if p.activatedTab > 0 {
				p.activatedTab--
			}
			return m, nil, 1
		case "backspace":
			if p.TabContent[p.activatedTab].FilterState() != list.Filtering {
				return m, nil, 0
			}
		}
	}

	var cmd tea.Cmd
	p.TabContent[p.activatedTab], cmd = p.TabContent[p.activatedTab].Update(msg)
	return m, cmd, 1
}

func (p *PluginPage) View() string {
	doc := strings.Builder{}

	// Tabs
	{
		var renderedTabs []string

		for i, t := range p.Tabs {
			if i == p.activatedTab {
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

	// list
	_, err := doc.WriteString(p.TabContent[p.activatedTab].View())
	if err != nil {
		return err.Error()
	}

	return docStyle.Render(doc.String())
}
