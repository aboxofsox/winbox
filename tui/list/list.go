package tui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	TitleStyle        = lipgloss.NewStyle().MarginLeft(2)
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("150"))
	PaginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	HelpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	QuitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type Item string

func (i Item) FilterValue() string { return "" }

type ItemDelegate struct{}

func (d ItemDelegate) Height() int                             { return 1 }
func (d ItemDelegate) Spacing() int                            { return 0 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, li list.Item) {
	i, ok := li.(Item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return SelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type Model struct {
	List      list.Model
	title     string
	submitMsg string
	Choice    string
	quitting  bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.List.SelectedItem().(Item)
			if ok {
				m.Choice = string(i)
			}
			return m, tea.Quit
		}

	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.Choice != "" {
		return QuitTextStyle.Render(m.submitMsg)
	}
	if m.quitting {
		return ""
	}

	return "\n" + m.List.View()
}

func Show(title, submitMsg string, items []string) string {
	listItems := make([]list.Item, 0)
	for _, l := range items {
		listItems = append(listItems, Item(l))
	}

	const defaultWidth = 20

	l := list.New(listItems, ItemDelegate{}, defaultWidth, listHeight)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = TitleStyle
	l.Styles.PaginationStyle = PaginationStyle
	l.Styles.HelpStyle = HelpStyle

	m := Model{List: l, submitMsg: submitMsg, title: title}

	tm, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m, _ = tm.(Model)
	if m.Choice == "" {
		os.Exit(0)
	}

	return m.Choice
}

func FindAllWSBFiles() []string {
	f, err := os.Open(".")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer f.Close()

	files, err := f.Readdir(-1)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	var wsb []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".wsb") {
			wsb = append(wsb, strings.TrimSuffix(file.Name(), ".wsb"))
		}
	}

	return wsb
}
