package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bad-noodles/kv-store/pkg/client"
	typesystem "github.com/bad-noodles/kv-store/pkg/type_system"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type command struct {
	query    string
	response string
}

type model struct {
	client   *client.Client
	input    textinput.Model
	viewport viewport.Model
	commands []command
}

func initialModel() model {
	input := textinput.New()
	input.Prompt = ">> "
	input.Focus()
	input.ShowSuggestions = true

	input.SetSuggestions([]string{"get ", "GET ", "set ", "SET "})

	vp := viewport.New(30, 5)

	return model{
		client:   client.NewClient(),
		input:    input,
		viewport: vp,
	}
}

type connected bool

func connect(client *client.Client) tea.Cmd {
	return func() tea.Msg {
		err := client.Connect("localhost:1337")
		if err != nil {
			log.Fatal(err)
		}

		return connected(true)
	}
}

func execute(client *client.Client, input string) tea.Cmd {
	return func() tea.Msg {
		err := client.Execute(input)
		if err != nil {
			return command{
				query:    input,
				response: typesystem.NewStatus(false, err.Error()).Pretty(),
			}
		}

		response, err := client.Read()
		if err != nil {
			return command{
				query:    input,
				response: typesystem.NewStatus(false, err.Error()).Pretty(),
			}
		}

		return command{
			query:    input,
			response: response.Pretty(),
		}
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("KV"), textinput.Blink, connect(m.client))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			val := m.input.Value()
			m.input.Reset()
			return m, execute(m.client, val)
		}
	case tea.WindowSizeMsg:
		m.viewport.Height = msg.Height - lipgloss.Height(m.input.View())
		m.viewport.Width = msg.Width
		m.input.Width = msg.Width
	case command:
		m.commands = append(m.commands, msg)

		var b strings.Builder
		for _, cmd := range m.commands {
			b.WriteString(">> ")
			b.WriteString(cmd.query)
			b.WriteRune('\n')
			b.WriteString(cmd.response)
			b.WriteRune('\n')
		}

		content := strings.Trim(b.String(), "\n")
		textPadding := m.viewport.Height - lipgloss.Height(content)

		if textPadding > 0 {
			content = fmt.Sprintf("%s%s", strings.Repeat("\n", textPadding), content)
		}

		m.viewport.SetContent(content)
		m.viewport.GotoBottom()
		return m, nil
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.viewport.View(),
		m.input.View(),
	)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
