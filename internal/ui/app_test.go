package ui

import (
	"strings"
	"testing"

	"go-on-git/internal/git"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewAppModel(t *testing.T) {
	m := NewAppModel()

	if m.mode != viewStatus {
		t.Errorf("mode = %v, want viewStatus", m.mode)
	}
}

func TestNewAppModelWithOptions(t *testing.T) {
	m := NewAppModelWithOptions(true)

	if !m.status.showVerboseHelp {
		t.Error("status.showVerboseHelp should be true when showHelp=true")
	}
}

func TestAppModelInit(t *testing.T) {
	m := NewAppModel()
	cmd := m.Init()

	if cmd == nil {
		t.Error("Init() should return a command")
	}
}

func TestAppModelWindowResize(t *testing.T) {
	m := NewAppModel()

	newModel, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	m = newModel.(AppModel)

	if m.width != 100 {
		t.Errorf("width = %d, want 100", m.width)
	}
	if m.height != 50 {
		t.Errorf("height = %d, want 50", m.height)
	}
	if m.status.width != 100 {
		t.Errorf("status.width = %d, want 100", m.status.width)
	}
	if m.diff.width != 100 {
		t.Errorf("diff.width = %d, want 100", m.diff.width)
	}
	if m.branches.width != 100 {
		t.Errorf("branches.width = %d, want 100", m.branches.width)
	}
	if m.stashes.width != 100 {
		t.Errorf("stashes.width = %d, want 100", m.stashes.width)
	}
	if m.log.width != 100 {
		t.Errorf("log.width = %d, want 100", m.log.width)
	}
}

func TestAppModelCtrlCQuits(t *testing.T) {
	m := NewAppModel()

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	if cmd == nil {
		t.Error("ctrl+c should return a quit command")
	}
}

func TestAppModelNavigateToBranches(t *testing.T) {
	m := NewAppModel()
	m.mode = viewStatus

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	m = newModel.(AppModel)

	if m.mode != viewBranches {
		t.Errorf("mode = %v, want viewBranches", m.mode)
	}
	if cmd == nil {
		t.Error("should return commands for entering branches view")
	}
}

func TestAppModelNavigateToStashes(t *testing.T) {
	m := NewAppModel()
	m.mode = viewStatus

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	m = newModel.(AppModel)

	if m.mode != viewStashes {
		t.Errorf("mode = %v, want viewStashes", m.mode)
	}
	if cmd == nil {
		t.Error("should return commands for entering stashes view")
	}
}

func TestAppModelNavigateToLog(t *testing.T) {
	m := NewAppModel()
	m.mode = viewStatus

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
	m = newModel.(AppModel)

	if m.mode != viewLog {
		t.Errorf("mode = %v, want viewLog", m.mode)
	}
	if cmd == nil {
		t.Error("should return commands for entering log view")
	}
}

func TestAppModelNavigateToAllDiffs(t *testing.T) {
	m := NewAppModel()
	m.mode = viewStatus

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	m = newModel.(AppModel)

	if m.mode != viewFullDiff {
		t.Errorf("mode = %v, want viewFullDiff", m.mode)
	}
	if cmd == nil {
		t.Error("should return commands for entering full diff view")
	}
}

func TestAppModelNavigateToFileDiff(t *testing.T) {
	m := NewAppModel()
	m.mode = viewStatus
	m.status.items = []StatusItem{
		{File: git.FileStatus{Path: "test.txt"}, Section: "unstaged"},
	}

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m = newModel.(AppModel)

	if m.mode != viewFileDiff {
		t.Errorf("mode = %v, want viewFileDiff", m.mode)
	}
	if cmd == nil {
		t.Error("should return commands for entering file diff view")
	}
	if len(m.currentFiles) != 1 {
		t.Errorf("currentFiles = %d, want 1", len(m.currentFiles))
	}
}

func TestAppModelNavigateToFileDiffNoFiles(t *testing.T) {
	m := NewAppModel()
	m.mode = viewStatus
	m.status.items = nil

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m = newModel.(AppModel)

	if m.mode != viewStatus {
		t.Errorf("mode should stay viewStatus when no files, got %v", m.mode)
	}
	if cmd != nil {
		t.Error("should not return a command when no files")
	}
}

func TestAppModelBackFromFileDiff(t *testing.T) {
	m := NewAppModel()
	m.mode = viewFileDiff
	m.diff.viewingHunk = false

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = newModel.(AppModel)

	if m.mode != viewStatus {
		t.Errorf("mode = %v, want viewStatus", m.mode)
	}
	if cmd == nil {
		t.Error("should return commands for exiting to status")
	}
}

func TestAppModelBackFromFullDiff(t *testing.T) {
	m := NewAppModel()
	m.mode = viewFullDiff
	m.diff.viewingHunk = false

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = newModel.(AppModel)

	if m.mode != viewStatus {
		t.Errorf("mode = %v, want viewStatus", m.mode)
	}
	if cmd == nil {
		t.Error("should return commands for exiting to status")
	}
}

func TestAppModelBackFromBranches(t *testing.T) {
	m := NewAppModel()
	m.mode = viewBranches

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = newModel.(AppModel)

	if m.mode != viewStatus {
		t.Errorf("mode = %v, want viewStatus", m.mode)
	}
	if cmd == nil {
		t.Error("should return commands for exiting to status")
	}
}

func TestAppModelQuitFromBranchesGoesBack(t *testing.T) {
	m := NewAppModel()
	m.mode = viewBranches

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m = newModel.(AppModel)

	if m.mode != viewStatus {
		t.Errorf("mode = %v, want viewStatus (q goes back from branches)", m.mode)
	}
	if cmd == nil {
		t.Error("should return commands for exiting to status")
	}
}

func TestAppModelBackFromStashes(t *testing.T) {
	m := NewAppModel()
	m.mode = viewStashes

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = newModel.(AppModel)

	if m.mode != viewStatus {
		t.Errorf("mode = %v, want viewStatus", m.mode)
	}
	if cmd == nil {
		t.Error("should return commands for exiting to status")
	}
}

func TestAppModelQuitFromStashesGoesBack(t *testing.T) {
	m := NewAppModel()
	m.mode = viewStashes

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m = newModel.(AppModel)

	if m.mode != viewStatus {
		t.Errorf("mode = %v, want viewStatus (q goes back from stashes)", m.mode)
	}
	if cmd == nil {
		t.Error("should return commands for exiting to status")
	}
}

func TestAppModelBackFromLog(t *testing.T) {
	m := NewAppModel()
	m.mode = viewLog

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = newModel.(AppModel)

	if m.mode != viewStatus {
		t.Errorf("mode = %v, want viewStatus", m.mode)
	}
	if cmd == nil {
		t.Error("should return commands for exiting to status")
	}
}

func TestAppModelQuitFromLogGoesBack(t *testing.T) {
	m := NewAppModel()
	m.mode = viewLog

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m = newModel.(AppModel)

	if m.mode != viewStatus {
		t.Errorf("mode = %v, want viewStatus (q goes back from log)", m.mode)
	}
	if cmd == nil {
		t.Error("should return commands for exiting to status")
	}
}

func TestAppModelLogKeyToggle(t *testing.T) {
	m := NewAppModel()
	m.mode = viewLog

	// Pressing 'o' again should go back
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
	m = newModel.(AppModel)

	if m.mode != viewStatus {
		t.Errorf("mode = %v, want viewStatus (o toggles log)", m.mode)
	}
	if cmd == nil {
		t.Error("should return commands for exiting to status")
	}
}

func TestAppModelNavigationBlockedInInputModes(t *testing.T) {
	m := NewAppModel()
	m.mode = viewStatus
	m.status.commitMode = true

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	m = newModel.(AppModel)

	if m.mode != viewStatus {
		t.Error("navigation should be blocked when in commit mode")
	}
}

func TestAppModelNavigationBlockedInStashMode(t *testing.T) {
	m := NewAppModel()
	m.mode = viewStatus
	m.status.stashMode = stashFiles

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	m = newModel.(AppModel)

	if m.mode != viewStatus {
		t.Error("navigation should be blocked when in stash mode")
	}
}

func TestAppModelNavigationBlockedInConfirmMode(t *testing.T) {
	m := NewAppModel()
	m.mode = viewStatus
	m.status.confirmMode = confirmDiscard

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	m = newModel.(AppModel)

	if m.mode != viewStatus {
		t.Error("navigation should be blocked when in confirm mode")
	}
}

func TestAppModelDrillDownToStashDiff(t *testing.T) {
	m := NewAppModel()
	m.mode = viewStashes
	m.stashes.stashes = []git.Stash{
		{Index: 0, Message: "test stash"},
	}

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m = newModel.(AppModel)

	if m.mode != viewStashDiff {
		t.Errorf("mode = %v, want viewStashDiff", m.mode)
	}
	if cmd == nil {
		t.Error("should return a command to fetch stash diff")
	}
}

func TestAppModelBackFromStashDiff(t *testing.T) {
	m := NewAppModel()
	m.mode = viewStashDiff
	m.stashes.diffModel.viewingHunk = false

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = newModel.(AppModel)

	if m.mode != viewStashes {
		t.Errorf("mode = %v, want viewStashes", m.mode)
	}
}

func TestAppModelBackFromStashDiffHunkDetail(t *testing.T) {
	m := NewAppModel()
	m.mode = viewStashDiff
	m.stashes.diffModel.viewingHunk = true

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = newModel.(AppModel)

	// Should exit hunk detail first, not go back to stashes
	if m.stashes.diffModel.viewingHunk {
		t.Error("should exit hunk detail first")
	}
	if m.mode != viewStashDiff {
		t.Errorf("mode = %v, should still be viewStashDiff", m.mode)
	}
}

func TestAppModelQuitFromStashDiff(t *testing.T) {
	m := NewAppModel()
	m.mode = viewStashDiff
	m.stashes.diffModel.viewingHunk = false

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m = newModel.(AppModel)

	if m.mode != viewStashes {
		t.Errorf("mode = %v, want viewStashes (q goes back from stash diff)", m.mode)
	}
}

func TestAppModelDiffViewHunkDetailBack(t *testing.T) {
	m := NewAppModel()
	m.mode = viewFileDiff
	m.diff.viewingHunk = true
	m.diff.hunks = []git.Hunk{
		{FilePath: "file1.txt"},
		{FilePath: "file2.txt"},
	}

	// With multiple hunks, back should exit hunk detail first
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = newModel.(AppModel)

	// Should delegate to diff model
	if m.mode != viewFileDiff {
		t.Errorf("mode = %v, should still be viewFileDiff", m.mode)
	}
}

func TestAppModelDiffViewSingleHunkBack(t *testing.T) {
	m := NewAppModel()
	m.mode = viewFileDiff
	m.diff.viewingHunk = true
	m.diff.hunks = []git.Hunk{
		{FilePath: "file1.txt"},
	}
	m.diff.showHelp = false
	m.diff.confirmMode = false

	// With single hunk and no overlays, back should go to status
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = newModel.(AppModel)

	if m.mode != viewStatus {
		t.Errorf("mode = %v, want viewStatus", m.mode)
	}
	if cmd == nil {
		t.Error("should return commands for exiting to status")
	}
}

func TestAppModelBranchesBackBlockedInModes(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*BranchesModel)
	}{
		{
			name: "showHelp",
			setup: func(m *BranchesModel) { m.showHelp = true },
		},
		{
			name: "deleteConfirmMode",
			setup: func(m *BranchesModel) { m.deleteConfirmMode = true },
		},
		{
			name: "inputMode",
			setup: func(m *BranchesModel) { m.inputMode = true },
		},
		{
			name: "forceDeleteMode",
			setup: func(m *BranchesModel) { m.forceDeleteMode = true },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewAppModel()
			m.mode = viewBranches
			tt.setup(&m.branches)

			newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
			m = newModel.(AppModel)

			if m.mode != viewBranches {
				t.Errorf("back should be blocked in %s mode", tt.name)
			}
		})
	}
}

func TestAppModelStashesBackBlockedInModes(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*StashesModel)
	}{
		{
			name: "showHelp",
			setup: func(m *StashesModel) { m.showHelp = true },
		},
		{
			name: "confirmMode",
			setup: func(m *StashesModel) { m.confirmMode = true },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewAppModel()
			m.mode = viewStashes
			tt.setup(&m.stashes)

			newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
			m = newModel.(AppModel)

			if m.mode != viewStashes {
				t.Errorf("back should be blocked in %s mode", tt.name)
			}
		})
	}
}

func TestAppModelView(t *testing.T) {
	m := NewAppModel()
	m.status.status = &git.StatusResult{}
	m.status.branchStatus = git.BranchStatus{Name: "main"}

	view := m.View()

	// Should render status view
	if len(view) == 0 {
		t.Error("view should not be empty")
	}
}

func TestAppModelViewBranches(t *testing.T) {
	m := NewAppModel()
	m.mode = viewBranches
	m.branches.branches = []git.Branch{{Name: "main"}}

	view := m.View()

	if !strings.Contains(view, "Branches") {
		t.Error("view should show Branches header")
	}
}

func TestAppModelViewStashes(t *testing.T) {
	m := NewAppModel()
	m.mode = viewStashes
	m.stashes.stashes = []git.Stash{{Index: 0, Message: "test"}}

	view := m.View()

	if !strings.Contains(view, "Stashes") {
		t.Error("view should show Stashes header")
	}
}

func TestAppModelViewStashDiff(t *testing.T) {
	m := NewAppModel()
	m.mode = viewStashDiff
	m.stashes.diffModel = NewStashDiffModel(100, 50)

	view := m.View()

	// Should render stash diff view
	if len(view) == 0 {
		t.Error("view should not be empty")
	}
}

func TestAppModelViewLog(t *testing.T) {
	m := NewAppModel()
	m.mode = viewLog
	m.log = NewLogModelWithSize(100, 50)
	m.log.lines = []string{"commit abc123"}

	view := m.View()

	if !strings.Contains(view, "commit") {
		t.Error("view should show log content")
	}
}

func TestAppModelViewDiff(t *testing.T) {
	m := NewAppModel()
	m.mode = viewFileDiff
	m.diff = NewDiffModel(nil)
	m.diff.diff = &git.CombinedDiffResult{}

	view := m.View()

	// Should render diff view
	if len(view) == 0 {
		t.Error("view should not be empty")
	}
}

func TestAppModelEscapeKey(t *testing.T) {
	m := NewAppModel()
	m.mode = viewBranches

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = newModel.(AppModel)

	if m.mode != viewStatus {
		t.Errorf("mode = %v, want viewStatus after esc", m.mode)
	}
	if cmd == nil {
		t.Error("should return commands for exiting to status")
	}
}

func TestAppModelArrowKeyNavigation(t *testing.T) {
	m := NewAppModel()
	m.mode = viewStatus
	m.status.items = []StatusItem{
		{File: git.FileStatus{Path: "test.txt"}, Section: "unstaged"},
	}

	// Right arrow should enter file diff
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = newModel.(AppModel)

	if m.mode != viewFileDiff {
		t.Errorf("mode = %v, want viewFileDiff after right arrow", m.mode)
	}
	if cmd == nil {
		t.Error("should return commands for entering file diff")
	}
}

func TestAppModelEnterKeyNavigation(t *testing.T) {
	m := NewAppModel()
	m.mode = viewStatus
	m.status.items = []StatusItem{
		{File: git.FileStatus{Path: "test.txt"}, Section: "unstaged"},
	}

	// Enter should enter file diff
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newModel.(AppModel)

	if m.mode != viewFileDiff {
		t.Errorf("mode = %v, want viewFileDiff after enter", m.mode)
	}
	if cmd == nil {
		t.Error("should return commands for entering file diff")
	}
}

func TestAppModelLeftKeyBack(t *testing.T) {
	m := NewAppModel()
	m.mode = viewLog

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	m = newModel.(AppModel)

	if m.mode != viewStatus {
		t.Errorf("mode = %v, want viewStatus after left arrow", m.mode)
	}
	if cmd == nil {
		t.Error("should return commands for exiting to status")
	}
}

func TestRefreshStatus(t *testing.T) {
	// This test would require mocking git operations
	// Just verify it doesn't panic when git commands fail
	// In a real git repo, it would work
	msg := refreshStatus()

	// Should return either statusMsg or errMsg
	switch msg.(type) {
	case statusMsg:
		// OK
	case errMsg:
		// Also OK - expected if not in a git repo
	default:
		t.Errorf("unexpected message type: %T", msg)
	}
}

func TestFileFilter(t *testing.T) {
	filter := FileFilter{
		Path:       "test.txt",
		ShowStaged: true,
		Untracked:  false,
	}

	if filter.Path != "test.txt" {
		t.Errorf("Path = %q, want 'test.txt'", filter.Path)
	}
	if !filter.ShowStaged {
		t.Error("ShowStaged should be true")
	}
	if filter.Untracked {
		t.Error("Untracked should be false")
	}
}
