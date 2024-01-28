package add

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMap struct {
	Confirm key.Binding
	Quit    key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Confirm, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Confirm, k.Quit}, // second column
	}
}

var keys = keyMap{
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "create task"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc/ctrl+c", "quit"),
	),
}

type (
	errMsg error
)

type Model struct {
	keys      keyMap
	Value     string
	textInput textinput.Model
	err       error
	help      help.Model
}

func AddInputModel() Model {
	ti := textinput.New()
	ti.Placeholder = "My new todo"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return Model{
		keys:      keys,
		Value:     "",
		textInput: ti,
		err:       nil,
		help:      help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Confirm):
			m.Value = m.textInput.Value()
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	helpView := m.help.View(m.keys)
	return fmt.Sprintf(
		"Input here your new task:\n\n%s\n\n%s",
		m.textInput.View(),
		helpView,
	) + "\n"
}
