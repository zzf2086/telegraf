package ui

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	flagStyle         = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#22ADF6"))
	selected          = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#13002D")).Background(lipgloss.Color("#9394FF"))
	descStyle         = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#F6F6F8"))
	selectedFlagStyle = lipgloss.NewStyle().Background(lipgloss.Color("#9394FF"))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("#D6F622"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)

	detailStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Foreground(lipgloss.Color("#FAFAFA")).
			Margin(1, 3, 0, 0).
			Padding(1, 2).
			Height(19).
			Width(40)
)

type itemDelegate struct{}

func (d itemDelegate) Height() int  { return 2 }
func (d itemDelegate) Spacing() int { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	listItem := m.SelectedItem()
	i, ok := listItem.(Item)
	if !ok {
		return nil
	}
	m.Title = fmt.Sprintf("Usage: telegraf %s", i.ItemTitle)
	return nil
}
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	// fn := itemStyle.Render
	if index == m.Index() {
		// str := lipgloss.JoinVertical(lipgloss.Left, fmt.Sprintf("• %s", i.ItemTitle), i.Desc)
		str := lipgloss.JoinHorizontal(lipgloss.Top, itemStyle.Render(">"), selected.Render(i.ItemTitle))
		str = lipgloss.JoinVertical(lipgloss.Left, str, descStyle.Render(i.Desc))
		fmt.Fprintf(w, selectedFlagStyle.Render(str))
	} else {
		str := lipgloss.JoinHorizontal(lipgloss.Top, itemStyle.Render("•"), flagStyle.Render(i.ItemTitle))
		str = lipgloss.JoinVertical(lipgloss.Left, str, descStyle.Render(i.Desc))
		fmt.Fprintf(w, str)
	}
}

type FlagsPage struct {
	Tabs       []string
	TabContent []list.Model

	activatedTab int
	selectedFlag *Item
}

func NewFlagsPage() FlagsPage {
	tabs := []string{
		"Flags",
		// "Commands",
	}

	tabContent := []list.Model{
		FlagsContent(),
	}

	return FlagsPage{Tabs: tabs, TabContent: tabContent}
}

func (f *FlagsPage) Update(m tea.Model, msg tea.Msg) (tea.Model, tea.Cmd, int) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit, 2
		case "right":
			if f.activatedTab < len(f.Tabs)-1 {
				f.activatedTab++
			}
			return m, nil, 2
		case "left":
			if f.activatedTab > 0 {
				f.activatedTab--
			}
			return m, nil, 2
		case "enter":
			listItem := f.TabContent[f.activatedTab].SelectedItem()
			i, ok := listItem.(Item)
			if !ok {
				return m, nil, 2
			}
			f.selectedFlag = &i
		case "backspace":
			if f.selectedFlag != nil {
				f.selectedFlag = nil
				return m, nil, 2
			}
			return m, nil, 0
		}
	}

	var cmd tea.Cmd
	f.TabContent[f.activatedTab], cmd = f.TabContent[f.activatedTab].Update(msg)
	return m, cmd, 2
}

func (f *FlagsPage) View() string {
	doc := strings.Builder{}

	// Tabs

	if f.selectedFlag != nil {
		doc.WriteString("<- backspace")
		details := lipgloss.JoinVertical(lipgloss.Top, f.selectedFlag.ItemTitle, f.selectedFlag.Desc)
		doc.WriteString(detailStyle.Render(details))
	} else {
		{
			var renderedTabs []string

			for i, t := range f.Tabs {
				if i == f.activatedTab {
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
		_, err := doc.WriteString(f.TabContent[f.activatedTab].View())
		if err != nil {
			return err.Error()
		}
	}

	return docStyle.Render(doc.String())
}

func FlagsContent() list.Model {
	var flagInfo []flag.Flag
	flag.VisitAll(func(f *flag.Flag) {
		flagInfo = append(flagInfo, *f)
	})

	var styledFlags []list.Item
	for _, f := range flagInfo {
		styledFlags = append(styledFlags, Item{ItemTitle: "--" + f.Name, Desc: f.Usage})
	}

	l := list.NewModel(styledFlags, itemDelegate{}, 50, 14)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return l
}
