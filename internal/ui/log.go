package ui

import (
	"fmt"
	"strings"

	"go-on-git/internal/git"

	tea "github.com/charmbracelet/bubbletea"
)


// LogModel is the bubbletea model for the log view
type LogModel struct {
	lines           []string
	scrollOffset    int
	showHelp        bool
	showVerboseHelp bool
	err             error
	width           int
	height          int
}

// NewLogModel creates a new log model
func NewLogModel() LogModel {
	return LogModel{}
}

// NewLogModelWithSize creates a new log model with dimensions
func NewLogModelWithSize(width, height int) LogModel {
	return LogModel{
		width:  width,
		height: height,
	}
}

// NewLogModelWithOptions creates a new log model with all options
func NewLogModelWithOptions(width, height int, showVerboseHelp bool) LogModel {
	return LogModel{
		width:           width,
		height:          height,
		showVerboseHelp: showVerboseHelp,
	}
}

type logMsg struct {
	content string
}

func refreshLog() tea.Msg {
	content, err := git.GetLog(100)
	if err != nil {
		return errMsg{err}
	}
	return logMsg{content}
}

// Init initializes the model
func (m LogModel) Init() tea.Cmd {
	return refreshLog
}

// Update handles messages
func (m LogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		// Handle help mode
		if m.showHelp {
			if key == Keys.Help || key == "esc" || key == Keys.Quit {
				m.showHelp = false
			}
			return m, nil
		}

		// Account for header (2 lines) and optionally help bar (2 lines)
		reservedLines := 4
		if m.showVerboseHelp {
			reservedLines = 6
		}
		visibleLines := m.height - reservedLines
		if visibleLines < 1 {
			visibleLines = 10
		}
		maxOffset := len(m.lines) - visibleLines
		if maxOffset < 0 {
			maxOffset = 0
		}

		switch key {
		case Keys.Help:
			m.showHelp = true
			return m, nil
		case Keys.VerboseHelp:
			m.showVerboseHelp = !m.showVerboseHelp
			return m, nil
		case Keys.Down, "down":
			m.scrollOffset = min(m.scrollOffset+1, maxOffset)
			return m, nil
		case Keys.Up, "up":
			m.scrollOffset = max(m.scrollOffset-1, 0)
			return m, nil
		case Keys.Bottom:
			m.scrollOffset = maxOffset
			return m, nil
		case Keys.Top:
			m.scrollOffset = 0
			return m, nil
		case "ctrl+d":
			m.scrollOffset = min(m.scrollOffset+visibleLines/2, maxOffset)
			return m, nil
		case "ctrl+u":
			m.scrollOffset = max(m.scrollOffset-visibleLines/2, 0)
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case logMsg:
		m.lines = strings.Split(msg.content, "\n")
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil
	}

	return m, nil
}

// View renders the log view
func (m LogModel) View() string {
	if m.showHelp {
		return m.renderHelp()
	}

	var content strings.Builder

	// Header
	content.WriteString(m.renderHeader())
	content.WriteString("\n\n")

	if m.err != nil {
		content.WriteString(StyleUnstaged.Render(fmt.Sprintf("Error: %v", m.err)))
		content.WriteString("\n")
		if m.showVerboseHelp {
			content.WriteString("\n")
			content.WriteString(m.renderHelpBar())
		}
		return m.anchorBottom(content.String())
	}

	if len(m.lines) == 0 {
		content.WriteString(StyleMuted.Render("Loading..."))
		content.WriteString("\n")
		if m.showVerboseHelp {
			content.WriteString("\n")
			content.WriteString(m.renderHelpBar())
		}
		return m.anchorBottom(content.String())
	}

	// Account for header (2 lines) and optionally help bar (2 lines)
	reservedLines := 4
	if m.showVerboseHelp {
		reservedLines = 6
	}
	visibleLines := m.height - reservedLines
	if visibleLines < 1 {
		visibleLines = 20
	}

	endIdx := m.scrollOffset + visibleLines
	if endIdx > len(m.lines) {
		endIdx = len(m.lines)
	}

	for i := m.scrollOffset; i < endIdx; i++ {
		line := m.lines[i]
		if strings.HasPrefix(line, "commit ") {
			content.WriteString(StyleStaged.Render(line))
		} else if strings.HasPrefix(line, "Author:") || strings.HasPrefix(line, "Date:") {
			content.WriteString(StyleMuted.Render(line))
		} else {
			content.WriteString(line)
		}
		content.WriteString("\n")
	}

	if m.showVerboseHelp {
		content.WriteString("\n")
		content.WriteString(m.renderHelpBar())
	}

	return m.anchorBottom(content.String())
}

func (m LogModel) renderHeader() string {
	return StyleMuted.Render("> git log") + "  " + StyleMuted.Render("(esc to go back)") + "\n" + StyleMuted.Render("───────────────────────────────────────────────────────────────")
}

func (m LogModel) anchorBottom(content string) string {
	lines := strings.Count(content, "\n")
	if m.height <= lines {
		return content
	}
	padding := m.height - lines - 1
	return strings.Repeat("\n", padding) + content
}

func (m LogModel) renderHelpBar() string {
	var sb strings.Builder

	sb.WriteString(StyleMuted.Render("───────────────────────────────────────────────────────────────"))
	sb.WriteString("\n")

	items := []struct{ key, desc string }{
		{formatKeyList(Keys.Down, Keys.Up), "scroll"},
		{formatKeyList(Keys.Top, Keys.Bottom), "top/bottom"},
		{"ctrl+d/u", "page down/up"},
		{Keys.Help, "help"},
		{formatKeyList(Keys.Left, "ESC"), "back"},
	}

	for _, item := range items {
		sb.WriteString(StyleHelpKey.Render(item.key))
		sb.WriteString(" ")
		sb.WriteString(StyleHelpDesc.Render(item.desc))
		sb.WriteString("  ")
	}

	return sb.String()
}

func (m LogModel) renderHelp() string {
	var sb strings.Builder

	sb.WriteString(StyleHelpTitle.Render("Log Shortcuts"))
	sb.WriteString("\n\n")

	moveKeys := formatKeyList(Keys.Down, Keys.Up, "↓", "↑")
	topKey := formatDoubleKey(Keys.Top)
	backKeys := formatKeyList(Keys.Left, "←", "ESC")

	help := []struct {
		key  string
		desc string
	}{
		{moveKeys, "Scroll down/up"},
		{topKey, "Go to top"},
		{Keys.Bottom, "Go to bottom"},
		{"ctrl+d", "Page down"},
		{"ctrl+u", "Page up"},
		{Keys.Help, "Toggle help"},
		{backKeys, "Go back"},
	}

	for _, h := range help {
		sb.WriteString(fmt.Sprintf("  %s  %s\n",
			StyleHelpKey.Render(fmt.Sprintf("%-8s", h.key)),
			StyleHelpDesc.Render(h.desc)))
	}

	return sb.String()
}
