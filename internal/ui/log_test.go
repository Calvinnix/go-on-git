package ui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewLogModel(t *testing.T) {
	m := NewLogModel()

	if m.scrollOffset != 0 {
		t.Errorf("scrollOffset = %d, want 0", m.scrollOffset)
	}
	if m.width != 0 {
		t.Errorf("width = %d, want 0", m.width)
	}
	if m.height != 0 {
		t.Errorf("height = %d, want 0", m.height)
	}
	if len(m.lines) != 0 {
		t.Errorf("lines should be empty, got %d", len(m.lines))
	}
}

func TestNewLogModelWithSize(t *testing.T) {
	m := NewLogModelWithSize(100, 50)

	if m.width != 100 {
		t.Errorf("width = %d, want 100", m.width)
	}
	if m.height != 50 {
		t.Errorf("height = %d, want 50", m.height)
	}
}

func TestLogModelInit(t *testing.T) {
	m := NewLogModel()
	cmd := m.Init()

	if cmd == nil {
		t.Error("Init() should return a command")
	}
}

func TestLogModelNavigation(t *testing.T) {
	m := NewLogModelWithSize(100, 10)
	m.lines = make([]string, 100)

	// Test scroll down
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = newModel.(LogModel)
	if m.scrollOffset != 1 {
		t.Errorf("after 'j', scrollOffset = %d, want 1", m.scrollOffset)
	}

	// Test scroll up
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = newModel.(LogModel)
	if m.scrollOffset != 0 {
		t.Errorf("after 'k', scrollOffset = %d, want 0", m.scrollOffset)
	}

	// Test can't scroll above 0
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = newModel.(LogModel)
	if m.scrollOffset != 0 {
		t.Errorf("scrollOffset should stay at 0, got %d", m.scrollOffset)
	}

	// Test jump to bottom
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
	m = newModel.(LogModel)
	if m.scrollOffset == 0 {
		t.Error("scrollOffset should be > 0 after 'G'")
	}

	// Test jump to top
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	m = newModel.(LogModel)
	if m.scrollOffset != 0 {
		t.Errorf("after 'g', scrollOffset = %d, want 0", m.scrollOffset)
	}
}

func TestLogModelPageNavigation(t *testing.T) {
	m := NewLogModelWithSize(100, 20)
	m.lines = make([]string, 100)

	// Test ctrl+d (half page down)
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}, Alt: false})
	// Need to use the actual key type for ctrl+d
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
	m = newModel.(LogModel)
	if m.scrollOffset == 0 {
		t.Error("scrollOffset should be > 0 after ctrl+d")
	}

	// Save current offset
	currentOffset := m.scrollOffset

	// Test ctrl+u (half page up)
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyCtrlU})
	m = newModel.(LogModel)
	if m.scrollOffset >= currentOffset {
		t.Error("scrollOffset should decrease after ctrl+u")
	}
}

func TestLogModelArrowKeys(t *testing.T) {
	m := NewLogModelWithSize(100, 10)
	m.lines = make([]string, 50)

	// Test down arrow
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = newModel.(LogModel)
	if m.scrollOffset != 1 {
		t.Errorf("after down arrow, scrollOffset = %d, want 1", m.scrollOffset)
	}

	// Test up arrow
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = newModel.(LogModel)
	if m.scrollOffset != 0 {
		t.Errorf("after up arrow, scrollOffset = %d, want 0", m.scrollOffset)
	}
}

func TestLogModelWindowResize(t *testing.T) {
	m := NewLogModel()

	newModel, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 60})
	m = newModel.(LogModel)

	if m.width != 120 {
		t.Errorf("width = %d, want 120", m.width)
	}
	if m.height != 60 {
		t.Errorf("height = %d, want 60", m.height)
	}
}

func TestLogModelLogMsg(t *testing.T) {
	m := NewLogModel()

	content := "commit abc123\nAuthor: Test\nDate: Today\n\n    Message"

	newModel, _ := m.Update(logMsg{content: content})
	m = newModel.(LogModel)

	if len(m.lines) != 5 {
		t.Errorf("len(lines) = %d, want 5", len(m.lines))
	}
}

func TestLogModelErrMsg(t *testing.T) {
	m := NewLogModel()

	newModel, _ := m.Update(errMsg{err: fmt.Errorf("test error")})
	m = newModel.(LogModel)

	if m.err == nil {
		t.Error("err should be set")
	}
}

func TestLogModelView(t *testing.T) {
	m := NewLogModelWithSize(100, 20)
	m.lines = []string{
		"commit abc123",
		"Author: Test User <test@example.com>",
		"Date:   Mon Jan 1 00:00:00 2024",
		"",
		"    Initial commit",
	}

	view := m.View()

	if !strings.Contains(view, "commit") {
		t.Error("view should contain 'commit'")
	}
	if !strings.Contains(view, "Author:") {
		t.Error("view should contain 'Author:'")
	}
}

func TestLogModelViewLoading(t *testing.T) {
	m := NewLogModelWithSize(100, 20)
	m.lines = nil

	view := m.View()

	if !strings.Contains(view, "Loading") {
		t.Error("view should show 'Loading' when lines is empty")
	}
}

func TestLogModelViewWithError(t *testing.T) {
	m := NewLogModelWithSize(100, 20)
	m.err = fmt.Errorf("test error")

	view := m.View()

	if !strings.Contains(view, "Error:") {
		t.Error("view should show error")
	}
}

func TestLogModelViewStyling(t *testing.T) {
	m := NewLogModelWithSize(100, 30)
	m.lines = []string{
		"commit abc123",
		"Author: Test User",
		"Date: Today",
		"Regular line",
	}

	view := m.View()

	// Just verify it doesn't panic and renders something
	if len(view) == 0 {
		t.Error("view should not be empty")
	}
}

func TestLogModelScrollBounds(t *testing.T) {
	m := NewLogModelWithSize(100, 10)
	m.lines = make([]string, 5) // Fewer lines than visible

	// Scroll down should not go negative
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
	m = newModel.(LogModel)

	// maxOffset should be 0 when lines fit in view
	if m.scrollOffset < 0 {
		t.Errorf("scrollOffset should not be negative, got %d", m.scrollOffset)
	}
}

func TestLogModelAnchorBottom(t *testing.T) {
	m := NewLogModelWithSize(100, 20)

	content := "line1\nline2\n"
	anchored := m.anchorBottom(content)

	// Should add padding
	if !strings.HasPrefix(anchored, "\n") {
		t.Error("anchored content should have leading newlines")
	}
}

func TestLogModelSmallHeight(t *testing.T) {
	m := NewLogModelWithSize(100, 0)
	m.lines = make([]string, 50)

	// Should use default visible lines
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = newModel.(LogModel)

	// Should not panic
	view := m.View()
	if len(view) == 0 {
		t.Error("view should render even with small height")
	}
}

func TestLogModelVisibleLinesCalculation(t *testing.T) {
	// When height is very small, visibleLines should have a minimum
	m := NewLogModelWithSize(100, 2)
	m.lines = make([]string, 100)

	// visibleLines = height - 2, minimum 1 for display
	// This is handled in the Update function
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
	m = newModel.(LogModel)

	// Should not panic and scrollOffset should be reasonable
	if m.scrollOffset < 0 {
		t.Error("scrollOffset should not be negative")
	}
}

func TestLogModelViewWithManyLines(t *testing.T) {
	m := NewLogModelWithSize(100, 20)
	lines := make([]string, 100)
	for i := 0; i < 100; i++ {
		lines[i] = fmt.Sprintf("Line %d", i)
	}
	m.lines = lines
	m.scrollOffset = 50

	view := m.View()

	// Should show lines around offset 50
	if !strings.Contains(view, "Line 50") && !strings.Contains(view, "Line 51") {
		t.Error("view should show lines near scroll offset")
	}
}
