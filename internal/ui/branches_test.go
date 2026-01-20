package ui

import (
	"fmt"
	"strings"
	"testing"

	"go-on-git/internal/git"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewBranchesModel(t *testing.T) {
	m := NewBranchesModel()

	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0", m.cursor)
	}
	if m.showHelp {
		t.Error("showHelp should be false initially")
	}
	if m.inputMode {
		t.Error("inputMode should be false initially")
	}
	if m.deleteConfirmMode {
		t.Error("deleteConfirmMode should be false initially")
	}
	if m.forceDeleteMode {
		t.Error("forceDeleteMode should be false initially")
	}
}

func TestBranchesModelInit(t *testing.T) {
	m := NewBranchesModel()
	cmd := m.Init()

	if cmd == nil {
		t.Error("Init() should return a command")
	}
}

func TestBranchesModelNavigation(t *testing.T) {
	m := NewBranchesModel()
	m.branches = []git.Branch{
		{Name: "main", IsCurrent: true},
		{Name: "feature-1"},
		{Name: "feature-2"},
	}

	// Test move down
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = newModel.(BranchesModel)
	if m.cursor != 1 {
		t.Errorf("after 'j', cursor = %d, want 1", m.cursor)
	}

	// Test move up
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = newModel.(BranchesModel)
	if m.cursor != 0 {
		t.Errorf("after 'k', cursor = %d, want 0", m.cursor)
	}

	// Test jump to bottom
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
	m = newModel.(BranchesModel)
	if m.cursor != 2 {
		t.Errorf("after 'G', cursor = %d, want 2", m.cursor)
	}

	// Test double g to top
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	m = newModel.(BranchesModel)
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	m = newModel.(BranchesModel)
	if m.cursor != 0 {
		t.Errorf("after 'gg', cursor = %d, want 0", m.cursor)
	}
}

func TestBranchesModelNavigationBounds(t *testing.T) {
	m := NewBranchesModel()
	m.branches = []git.Branch{
		{Name: "main"},
		{Name: "feature"},
	}

	// Can't go above 0
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = newModel.(BranchesModel)
	if m.cursor != 0 {
		t.Errorf("cursor should stay at 0, got %d", m.cursor)
	}

	// Move to bottom
	m.cursor = 1

	// Can't go past end
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = newModel.(BranchesModel)
	if m.cursor != 1 {
		t.Errorf("cursor should stay at 1, got %d", m.cursor)
	}
}

func TestBranchesModelHelpToggle(t *testing.T) {
	m := NewBranchesModel()

	// Toggle help
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m = newModel.(BranchesModel)
	if !m.showHelp {
		t.Error("showHelp should be true after '?'")
	}

	// Close help with '?'
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m = newModel.(BranchesModel)
	if m.showHelp {
		t.Error("showHelp should be false after pressing '?' again")
	}
}

func TestBranchesModelInputMode(t *testing.T) {
	m := NewBranchesModel()

	// Press n to enter input mode
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	m = newModel.(BranchesModel)
	if !m.inputMode {
		t.Error("should be in input mode after 'n'")
	}

	// Press esc to exit
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = newModel.(BranchesModel)
	if m.inputMode {
		t.Error("should exit input mode after esc")
	}
}

func TestBranchesModelInputModeEnter(t *testing.T) {
	m := NewBranchesModel()
	m.inputMode = true
	m.branchInput.SetValue("new-branch")

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newModel.(BranchesModel)

	if m.inputMode {
		t.Error("should exit input mode after enter")
	}
	if cmd == nil {
		t.Error("should return a command to create branch")
	}
}

func TestBranchesModelInputModeEnterEmpty(t *testing.T) {
	m := NewBranchesModel()
	m.inputMode = true
	m.branchInput.SetValue("")

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newModel.(BranchesModel)

	if m.inputMode {
		t.Error("should exit input mode after enter")
	}
	if cmd != nil {
		t.Error("should not return a command for empty branch name")
	}
}

func TestBranchesModelDeleteConfirmMode(t *testing.T) {
	m := NewBranchesModel()
	m.branches = []git.Branch{
		{Name: "main", IsCurrent: true},
		{Name: "feature"},
	}
	m.cursor = 1 // Select non-current branch

	// Press d to enter delete confirm mode
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = newModel.(BranchesModel)
	if !m.deleteConfirmMode {
		t.Error("should be in delete confirm mode after 'd' on non-current branch")
	}

	// Press esc to cancel
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = newModel.(BranchesModel)
	if m.deleteConfirmMode {
		t.Error("should exit delete confirm mode after esc")
	}
}

func TestBranchesModelDeleteConfirmModeOnCurrentBranch(t *testing.T) {
	m := NewBranchesModel()
	m.branches = []git.Branch{
		{Name: "main", IsCurrent: true},
	}
	m.cursor = 0

	// Press d on current branch - should NOT enter delete mode
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = newModel.(BranchesModel)
	if m.deleteConfirmMode {
		t.Error("should NOT enter delete confirm mode on current branch")
	}
}

func TestBranchesModelDeleteConfirmWithTyping(t *testing.T) {
	m := NewBranchesModel()
	m.branches = []git.Branch{
		{Name: "main", IsCurrent: true},
		{Name: "feature"},
	}
	m.cursor = 1
	m.deleteConfirmMode = true
	m.deleteInput.SetValue("feature")

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newModel.(BranchesModel)

	if m.deleteConfirmMode {
		t.Error("should exit delete confirm mode after correct name")
	}
	if cmd == nil {
		t.Error("should return a command to delete branch")
	}
}

func TestBranchesModelDeleteConfirmWrongName(t *testing.T) {
	m := NewBranchesModel()
	m.branches = []git.Branch{
		{Name: "main", IsCurrent: true},
		{Name: "feature"},
	}
	m.cursor = 1
	m.deleteConfirmMode = true
	m.deleteInput.SetValue("wrong-name")

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newModel.(BranchesModel)

	if m.deleteConfirmMode {
		t.Error("should exit delete confirm mode")
	}
	if cmd != nil {
		t.Error("should NOT return a command for wrong name")
	}
}

func TestBranchesModelForceDeleteMode(t *testing.T) {
	m := NewBranchesModel()
	m.forceDeleteMode = true
	m.pendingDeleteBranch = "feature"

	// Press y to confirm
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	m = newModel.(BranchesModel)

	if m.forceDeleteMode {
		t.Error("should exit force delete mode after 'y'")
	}
	if cmd == nil {
		t.Error("should return a command to force delete")
	}
}

func TestBranchesModelForceDeleteModeCancel(t *testing.T) {
	m := NewBranchesModel()
	m.forceDeleteMode = true
	m.pendingDeleteBranch = "feature"

	// Press n to cancel
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	m = newModel.(BranchesModel)

	if m.forceDeleteMode {
		t.Error("should exit force delete mode after 'n'")
	}
	if m.pendingDeleteBranch != "" {
		t.Error("pendingDeleteBranch should be cleared")
	}
}

func TestBranchesModelCheckout(t *testing.T) {
	m := NewBranchesModel()
	m.branches = []git.Branch{
		{Name: "main", IsCurrent: true},
		{Name: "feature"},
	}
	m.cursor = 1 // Select non-current branch

	// Press right/l/enter to checkout
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m = newModel.(BranchesModel)

	if cmd == nil {
		t.Error("should return a command to checkout branch")
	}
}

func TestBranchesModelCheckoutCurrentBranch(t *testing.T) {
	m := NewBranchesModel()
	m.branches = []git.Branch{
		{Name: "main", IsCurrent: true},
	}
	m.cursor = 0

	// Try to checkout current branch - should not do anything
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})

	if cmd != nil {
		t.Error("should NOT return a command to checkout current branch")
	}
}

func TestBranchesModelWindowResize(t *testing.T) {
	m := NewBranchesModel()

	newModel, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	m = newModel.(BranchesModel)

	if m.width != 100 {
		t.Errorf("width = %d, want 100", m.width)
	}
	if m.height != 50 {
		t.Errorf("height = %d, want 50", m.height)
	}
}

func TestBranchesModelBranchesMsg(t *testing.T) {
	m := NewBranchesModel()

	branches := []git.Branch{
		{Name: "main", IsCurrent: true},
		{Name: "feature-1"},
		{Name: "feature-2"},
	}

	newModel, _ := m.Update(branchesMsg{branches: branches})
	m = newModel.(BranchesModel)

	if len(m.branches) != 3 {
		t.Errorf("len(branches) = %d, want 3", len(m.branches))
	}
	// Should position cursor on current branch
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0 (current branch)", m.cursor)
	}
}

func TestBranchesModelBranchesMsgCursorOnCurrent(t *testing.T) {
	m := NewBranchesModel()

	branches := []git.Branch{
		{Name: "feature-1"},
		{Name: "main", IsCurrent: true},
		{Name: "feature-2"},
	}

	newModel, _ := m.Update(branchesMsg{branches: branches})
	m = newModel.(BranchesModel)

	// Should position cursor on current branch (index 1)
	if m.cursor != 1 {
		t.Errorf("cursor = %d, want 1 (current branch)", m.cursor)
	}
}

func TestBranchesModelBranchDeleteFailedMsg(t *testing.T) {
	m := NewBranchesModel()

	msg := branchDeleteFailedMsg{
		branchName: "feature",
		err:        fmt.Errorf("not fully merged"),
	}

	newModel, _ := m.Update(msg)
	m = newModel.(BranchesModel)

	if !m.forceDeleteMode {
		t.Error("should enter force delete mode")
	}
	if m.pendingDeleteBranch != "feature" {
		t.Errorf("pendingDeleteBranch = %q, want 'feature'", m.pendingDeleteBranch)
	}
	if m.err == nil {
		t.Error("err should be set")
	}
}

func TestBranchesModelErrMsg(t *testing.T) {
	m := NewBranchesModel()

	newModel, _ := m.Update(errMsg{err: fmt.Errorf("test error")})
	m = newModel.(BranchesModel)

	if m.err == nil {
		t.Error("err should be set")
	}
}

func TestBranchesModelView(t *testing.T) {
	m := NewBranchesModel()
	m.branches = []git.Branch{
		{Name: "main", IsCurrent: true, LastCommit: "Initial commit"},
		{Name: "feature", LastCommit: "Add feature"},
	}

	view := m.View()

	if !strings.Contains(view, "git branch") {
		t.Error("view should contain 'git branch' header")
	}
	if !strings.Contains(view, "main") {
		t.Error("view should contain 'main' branch")
	}
	if !strings.Contains(view, "feature") {
		t.Error("view should contain 'feature' branch")
	}
}

func TestBranchesModelViewEmpty(t *testing.T) {
	m := NewBranchesModel()
	m.branches = nil

	view := m.View()

	if !strings.Contains(view, "No branches") {
		t.Error("view should show 'No branches found'")
	}
}

func TestBranchesModelViewWithError(t *testing.T) {
	m := NewBranchesModel()
	m.err = fmt.Errorf("test error")
	m.branches = []git.Branch{{Name: "main"}}

	view := m.View()

	if !strings.Contains(view, "Error:") {
		t.Error("view should show error")
	}
}

func TestBranchesModelViewHelp(t *testing.T) {
	m := NewBranchesModel()
	m.showHelp = true

	view := m.View()

	if !strings.Contains(view, "Branches Shortcuts") {
		t.Error("help view should contain 'Branches Shortcuts'")
	}
}

func TestBranchesModelViewInputMode(t *testing.T) {
	m := NewBranchesModel()
	m.branches = []git.Branch{{Name: "main"}}
	m.inputMode = true
	m.branchInput.Focus()

	view := m.View()

	if !strings.Contains(view, "New branch name") {
		t.Error("view should show input prompt")
	}
}

func TestBranchesModelViewDeleteConfirmMode(t *testing.T) {
	m := NewBranchesModel()
	m.branches = []git.Branch{
		{Name: "main", IsCurrent: true},
		{Name: "feature"},
	}
	m.cursor = 1
	m.deleteConfirmMode = true
	m.deleteInput.Focus()

	view := m.View()

	if !strings.Contains(view, "Type") && !strings.Contains(view, "delete") {
		t.Error("view should show delete confirmation prompt")
	}
}

func TestBranchesModelViewForceDeleteMode(t *testing.T) {
	m := NewBranchesModel()
	m.branches = []git.Branch{{Name: "main"}}
	m.forceDeleteMode = true
	m.pendingDeleteBranch = "feature"
	m.err = fmt.Errorf("not fully merged")

	view := m.View()

	if !strings.Contains(view, "Force delete") {
		t.Error("view should show force delete prompt")
	}
}

func TestBranchesModelViewWithTracking(t *testing.T) {
	m := NewBranchesModel()
	m.branches = []git.Branch{
		{Name: "main", IsCurrent: true, Upstream: "origin/main", Ahead: 2, Behind: 1},
	}

	view := m.View()

	if !strings.Contains(view, "+2") && !strings.Contains(view, "-1") {
		t.Error("view should show ahead/behind counts")
	}
}

func TestBranchesModelViewTruncatesLongCommit(t *testing.T) {
	m := NewBranchesModel()
	longMessage := strings.Repeat("a", 100)
	m.branches = []git.Branch{
		{Name: "main", IsCurrent: true, LastCommit: longMessage},
	}

	view := m.View()

	// Should contain truncated message with ellipsis
	if !strings.Contains(view, "...") {
		t.Error("view should truncate long commit messages")
	}
}

func TestBranchesModelArrowKeys(t *testing.T) {
	m := NewBranchesModel()
	m.branches = []git.Branch{
		{Name: "main"},
		{Name: "feature"},
	}

	// Test down arrow
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = newModel.(BranchesModel)
	if m.cursor != 1 {
		t.Errorf("after down arrow, cursor = %d, want 1", m.cursor)
	}

	// Test up arrow
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = newModel.(BranchesModel)
	if m.cursor != 0 {
		t.Errorf("after up arrow, cursor = %d, want 0", m.cursor)
	}
}

func TestBranchesModelEnterKey(t *testing.T) {
	m := NewBranchesModel()
	m.branches = []git.Branch{
		{Name: "main", IsCurrent: true},
		{Name: "feature"},
	}
	m.cursor = 1

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Error("enter should checkout branch")
	}
}

func TestBranchesModelHelpModeBlocksNavigation(t *testing.T) {
	m := NewBranchesModel()
	m.branches = []git.Branch{
		{Name: "main"},
		{Name: "feature"},
	}
	m.showHelp = true

	// Navigation should be blocked in help mode
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = newModel.(BranchesModel)
	if m.cursor != 0 {
		t.Error("navigation should be blocked in help mode")
	}
}

func TestBranchesModelInputModeTyping(t *testing.T) {
	m := NewBranchesModel()
	m.inputMode = true
	m.branchInput.Focus()

	// Type a character
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = newModel.(BranchesModel)

	// The input should have been updated
	if !m.inputMode {
		t.Error("should still be in input mode")
	}
}
