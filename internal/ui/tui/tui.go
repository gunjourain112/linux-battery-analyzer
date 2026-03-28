package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/domain"
	"github.com/gunjourain112/notebook-battery-analyzer/internal/ui/i18n"
)

type state int

const (
	stateLang state = iota
	stateRange
	stateSince
	stateUntil
)

type Model struct {
	step     state
	langIdx  int
	rangeIdx int
	since    textinput.Model
	until    textinput.Model
	err      string
	done     bool
	config   domain.Config
}

func New() Model {
	lang := detectLanguage()
	since := textinput.New()
	since.Prompt = i18n.New(lang).Get(i18n.SincePrompt)
	since.Placeholder = "YYYY-MM-DD"
	since.SetValue(time.Now().AddDate(0, 0, -7).Format("2006-01-02"))
	since.CharLimit = 10
	since.Width = 16

	until := textinput.New()
	until.Prompt = i18n.New(lang).Get(i18n.UntilPrompt)
	until.Placeholder = "YYYY-MM-DD"
	until.SetValue(time.Now().Format("2006-01-02"))
	until.CharLimit = 10
	until.Width = 16

	return Model{
		step:     stateLang,
		langIdx:  languageIndex(lang),
		rangeIdx: 0,
		since:    since,
		until:    until,
		config: domain.Config{
			Language: lang,
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
		case stateRange:
			return m.updateRange(msg)
		case stateSince, stateUntil:
			return m.updateDate(msg)
		}
	}

	var cmd tea.Cmd
	switch m.step {
	case stateSince:
		m.since, cmd = m.since.Update(msg)
	case stateUntil:
		m.until, cmd = m.until.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	m = m.localize()
	tr := i18n.New(m.config.Language)

	if m.done {
		return fmt.Sprintf("%s: %s\n%s: %s\n%s: %s\n",
			tr.Get(i18n.LanguageLabel), m.config.Language,
			tr.Get(i18n.SinceLabel), m.config.Since.Format("2006-01-02"),
			tr.Get(i18n.UntilLabel), m.config.Until.Format("2006-01-02"),
		)
	}

	var b strings.Builder
	b.WriteString(tr.Get(i18n.AppTitle))
	b.WriteString("\n\n")

	if m.step == stateLang {
		b.WriteString(tr.Get(i18n.ChooseLanguage))
		b.WriteString("\n")
		opts := []string{tr.Get(i18n.LanguageKo), tr.Get(i18n.LanguageEn)}
		for i, opt := range opts {
			prefix := "  "
			if i == m.langIdx {
				prefix = "> "
			}
			b.WriteString(prefix + opt + "\n")
		}
		b.WriteString("\n")
		b.WriteString(tr.Get(i18n.EnterToContinue))
		b.WriteString(" | ")
		b.WriteString(tr.Get(i18n.ChooseRange))
		return b.String()
	}

	if m.step == stateRange {
		b.WriteString(tr.Get(i18n.ChooseRange))
		b.WriteString("\n")
		opts := []string{tr.Get(i18n.RangeLast7), tr.Get(i18n.RangeLast30), tr.Get(i18n.RangeCustom)}
		for i, opt := range opts {
			prefix := "  "
			if i == m.rangeIdx {
				prefix = "> "
			}
			b.WriteString(prefix + opt + "\n")
		}
		b.WriteString("\n")
		b.WriteString(tr.Get(i18n.EnterToSelect))
		if m.err != "" {
			b.WriteString("\n")
			b.WriteString(m.err)
		}
		return b.String()
	}

	if m.step == stateSince || m.step == stateUntil {
		b.WriteString(m.since.View())
		b.WriteString("\n")
		b.WriteString(m.until.View())
		b.WriteString("\n")
		b.WriteString(tr.Get(i18n.EnterToRunTabToSwitch))
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
		m = m.localize()
		m.step = stateRange
		m.err = ""
		m.since.Blur()
		m.until.Blur()
		return m, nil
	}
	return m, nil
}

func (m Model) updateRange(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.rangeIdx > 0 {
			m.rangeIdx--
		}
	case "down", "j":
		if m.rangeIdx < 2 {
			m.rangeIdx++
		}
	case "enter":
		switch m.rangeIdx {
		case 0:
			m.applyPresetRange(7)
			m.done = true
			return m, tea.Quit
		case 1:
			m.applyPresetRange(30)
			m.done = true
			return m, tea.Quit
		default:
			m.step = stateSince
			m.err = ""
			m.since.Focus()
			m.until.Blur()
			return m, nil
		}
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
			m.err = i18n.New(m.config.Language).Get(i18n.InvalidSinceDate)
			return m, nil
		}
		until, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(m.until.Value()), time.Local)
		if err != nil {
			m.err = i18n.New(m.config.Language).Get(i18n.InvalidUntilDate)
			return m, nil
		}
		m.config.Since = since
		m.config.Until = until.Add(24*time.Hour - time.Second)
		m.done = true
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

func (m Model) localize() Model {
	tr := i18n.New(m.config.Language)
	m.since.Prompt = tr.Get(i18n.SincePrompt)
	m.until.Prompt = tr.Get(i18n.UntilPrompt)
	return m
}

func (m *Model) applyPresetRange(days int) {
	now := time.Now()
	m.config.Since = now.AddDate(0, 0, -days)
	m.config.Until = now
	m.since.SetValue(m.config.Since.Format("2006-01-02"))
	m.until.SetValue(m.config.Until.Format("2006-01-02"))
}

func detectLanguage() string {
	envs := []string{
		os.Getenv("LC_ALL"),
		os.Getenv("LC_MESSAGES"),
		os.Getenv("LANG"),
	}
	for _, v := range envs {
		v = strings.ToLower(v)
		if strings.Contains(v, "ko") {
			return "ko"
		}
		if strings.Contains(v, "en") {
			return "en"
		}
	}
	return "ko"
}

func languageIndex(lang string) int {
	if strings.ToLower(strings.TrimSpace(lang)) == "en" {
		return 1
	}
	return 0
}
