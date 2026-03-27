package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
)

type state int

const (
	stateLang state = iota
	stateSince
	stateUntil
)

type Model struct {
	step    state
	langIdx int
	since   textinput.Model
	until   textinput.Model
	err     string
	done    bool
	config  domain.Config
}

func New() Model {
	since := textinput.New()
	since.Prompt = "since: "
	since.Placeholder = "YYYY-MM-DD"
	since.SetValue(time.Now().AddDate(0, 0, -7).Format("2006-01-02"))
	since.Focus()
	since.CharLimit = 10
	since.Width = 16

	until := textinput.New()
	until.Prompt = "until: "
	until.Placeholder = "YYYY-MM-DD"
	until.SetValue(time.Now().Format("2006-01-02"))
	until.CharLimit = 10
	until.Width = 16

	return Model{
		step:  stateLang,
		since: since,
		until: until,
		config: domain.Config{
			Language: "ko",
		},
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}

		switch m.step {
		case stateLang:
			return m.updateLanguage(msg)
		case stateSince, stateUntil:
			return m.updateDate(msg)
		}
	}

	var cmd tea.Cmd
	if m.step == stateSince || m.step == stateUntil {
		m.until, cmd = m.until.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	if m.done {
		return fmt.Sprintf("language: %s\nsince: %s\nuntil: %s\n",
			m.config.Language,
			m.config.Since.Format("2006-01-02"),
			m.config.Until.Format("2006-01-02"),
		)
	}

	var b strings.Builder
	b.WriteString("notebook battery analyzer\n\n")

	if m.step == stateLang {
		b.WriteString("choose language\n")
		opts := []string{"ko", "en"}
		for i, opt := range opts {
			prefix := "  "
			if i == m.langIdx {
				prefix = "> "
			}
			b.WriteString(prefix + opt + "\n")
		}
		b.WriteString("\nenter to continue")
		return b.String()
	}

	if m.step == stateSince || m.step == stateUntil {
		b.WriteString(m.since.View())
		b.WriteString("\n")
		b.WriteString(m.until.View())
		b.WriteString("\nenter to run, tab to switch")
		if m.err != "" {
			b.WriteString("\n")
			b.WriteString(m.err)
		}
	}

	return b.String()
}

func (m Model) updateLanguage(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.langIdx > 0 {
			m.langIdx--
		}
	case "down", "j":
		if m.langIdx < 1 {
			m.langIdx++
		}
	case "enter":
		if m.langIdx == 0 {
			m.config.Language = "ko"
		} else {
			m.config.Language = "en"
		}
		m.step = stateSince
		m.since.Focus()
		m.until.Blur()
		return m, nil
	}
	return m, nil
}

func (m Model) updateDate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		if m.step == stateSince {
			m.step = stateUntil
			m.since.Blur()
			m.until.Focus()
		} else {
			m.step = stateSince
			m.until.Blur()
			m.since.Focus()
		}
		return m, nil
	case "enter":
		since, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(m.since.Value()), time.Local)
		if err != nil {
			m.err = "invalid since date"
			return m, nil
		}
		until, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(m.until.Value()), time.Local)
		if err != nil {
			m.err = "invalid until date"
			return m, nil
		}
		m.config.Since = since
		m.config.Until = until.Add(24*time.Hour - time.Second)
		m.done = true
		m.step = stateUntil
		return m, tea.Quit
	}

	var cmd tea.Cmd
	if m.step == stateSince {
		m.since, cmd = m.since.Update(msg)
		return m, cmd
	}
	m.until, cmd = m.until.Update(msg)
	return m, cmd
}

func Run() (domain.Config, error) {
	prog := tea.NewProgram(New(), tea.WithAltScreen())
	finalModel, err := prog.Run()
	if err != nil {
		return domain.Config{}, err
	}
	m, ok := finalModel.(Model)
	if !ok {
		return domain.Config{}, fmt.Errorf("unexpected model type")
	}
	return m.config, nil
}
