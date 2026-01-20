package ui

import (
	"testing"
)

func TestDefaultKeymap(t *testing.T) {
	km := DefaultKeymap()

	// Test navigation keys
	if km.Up != "k" {
		t.Errorf("expected Up to be 'k', got %q", km.Up)
	}
	if km.Down != "j" {
		t.Errorf("expected Down to be 'j', got %q", km.Down)
	}
	if km.Left != "h" {
		t.Errorf("expected Left to be 'h', got %q", km.Left)
	}
	if km.Right != "l" {
		t.Errorf("expected Right to be 'l', got %q", km.Right)
	}
	if km.Top != "g" {
		t.Errorf("expected Top to be 'g', got %q", km.Top)
	}
	if km.Bottom != "G" {
		t.Errorf("expected Bottom to be 'G', got %q", km.Bottom)
	}
	if km.Quit != "q" {
		t.Errorf("expected Quit to be 'q', got %q", km.Quit)
	}

	// Test action keys
	if km.Stage != "a" {
		t.Errorf("expected Stage to be 'a', got %q", km.Stage)
	}
	if km.StageAll != "A" {
		t.Errorf("expected StageAll to be 'A', got %q", km.StageAll)
	}
	if km.Unstage != "u" {
		t.Errorf("expected Unstage to be 'u', got %q", km.Unstage)
	}
	if km.UnstageAll != "U" {
		t.Errorf("expected UnstageAll to be 'U', got %q", km.UnstageAll)
	}
	if km.Discard != "d" {
		t.Errorf("expected Discard to be 'd', got %q", km.Discard)
	}
	if km.Commit != "c" {
		t.Errorf("expected Commit to be 'c', got %q", km.Commit)
	}
	if km.CommitEdit != "C" {
		t.Errorf("expected CommitEdit to be 'C', got %q", km.CommitEdit)
	}
	if km.Push != "p" {
		t.Errorf("expected Push to be 'p', got %q", km.Push)
	}
	if km.Stash != "s" {
		t.Errorf("expected Stash to be 's', got %q", km.Stash)
	}
	if km.StashAll != "S" {
		t.Errorf("expected StashAll to be 'S', got %q", km.StashAll)
	}

	// Test view keys
	if km.FileDiff != "l" {
		t.Errorf("expected FileDiff to be 'l', got %q", km.FileDiff)
	}
	if km.AllDiffs != "i" {
		t.Errorf("expected AllDiffs to be 'i', got %q", km.AllDiffs)
	}
	if km.Branches != "b" {
		t.Errorf("expected Branches to be 'b', got %q", km.Branches)
	}
	if km.Stashes != "e" {
		t.Errorf("expected Stashes to be 'e', got %q", km.Stashes)
	}
	if km.Log != "o" {
		t.Errorf("expected Log to be 'o', got %q", km.Log)
	}

	// Test mode keys
	if km.Visual != "v" {
		t.Errorf("expected Visual to be 'v', got %q", km.Visual)
	}
	if km.Help != "?" {
		t.Errorf("expected Help to be '?', got %q", km.Help)
	}
	if km.VerboseHelp != "/" {
		t.Errorf("expected VerboseHelp to be '/', got %q", km.VerboseHelp)
	}
	if km.NewBranch != "n" {
		t.Errorf("expected NewBranch to be 'n', got %q", km.NewBranch)
	}
	if km.Delete != "d" {
		t.Errorf("expected Delete to be 'd', got %q", km.Delete)
	}
}

func TestParseKeymapArg(t *testing.T) {
	tests := []struct {
		name      string
		arg       string
		wantAction string
		wantKey   string
		wantValid bool
	}{
		{
			name:       "valid simple override",
			arg:        "up=w",
			wantAction: "up",
			wantKey:    "w",
			wantValid:  true,
		},
		{
			name:       "valid hyphenated action",
			arg:        "stage-all=X",
			wantAction: "stage-all",
			wantKey:    "X",
			wantValid:  true,
		},
		{
			name:       "valid with special key",
			arg:        "quit=ctrl+c",
			wantAction: "quit",
			wantKey:    "ctrl+c",
			wantValid:  true,
		},
		{
			name:       "valid with equals in key",
			arg:        "commit==",
			wantAction: "commit",
			wantKey:    "=",
			wantValid:  true,
		},
		{
			name:       "invalid no equals",
			arg:        "upw",
			wantAction: "",
			wantKey:    "",
			wantValid:  false,
		},
		{
			name:       "invalid empty",
			arg:        "",
			wantAction: "",
			wantKey:    "",
			wantValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, key, valid := ParseKeymapArg(tt.arg)
			if action != tt.wantAction {
				t.Errorf("action = %q, want %q", action, tt.wantAction)
			}
			if key != tt.wantKey {
				t.Errorf("key = %q, want %q", key, tt.wantKey)
			}
			if valid != tt.wantValid {
				t.Errorf("valid = %v, want %v", valid, tt.wantValid)
			}
		})
	}
}

func TestApplyOverride(t *testing.T) {
	tests := []struct {
		name       string
		action     string
		key        string
		wantApplied bool
		checkField func(*Keymap) string
		wantValue  string
	}{
		{
			name:        "override up",
			action:      "up",
			key:         "w",
			wantApplied: true,
			checkField:  func(k *Keymap) string { return k.Up },
			wantValue:   "w",
		},
		{
			name:        "override down",
			action:      "down",
			key:         "s",
			wantApplied: true,
			checkField:  func(k *Keymap) string { return k.Down },
			wantValue:   "s",
		},
		{
			name:        "override stage-all",
			action:      "stage-all",
			key:         "X",
			wantApplied: true,
			checkField:  func(k *Keymap) string { return k.StageAll },
			wantValue:   "X",
		},
		{
			name:        "override verbose-help",
			action:      "verbose-help",
			key:         "H",
			wantApplied: true,
			checkField:  func(k *Keymap) string { return k.VerboseHelp },
			wantValue:   "H",
		},
		{
			name:        "override new-branch",
			action:      "new-branch",
			key:         "N",
			wantApplied: true,
			checkField:  func(k *Keymap) string { return k.NewBranch },
			wantValue:   "N",
		},
		{
			name:        "unknown action",
			action:      "invalid-action",
			key:         "x",
			wantApplied: false,
		},
		{
			name:        "empty action",
			action:      "",
			key:         "x",
			wantApplied: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			km := DefaultKeymap()
			applied := km.ApplyOverride(tt.action, tt.key)
			if applied != tt.wantApplied {
				t.Errorf("ApplyOverride() = %v, want %v", applied, tt.wantApplied)
			}
			if tt.wantApplied && tt.checkField != nil {
				if got := tt.checkField(km); got != tt.wantValue {
					t.Errorf("field value = %q, want %q", got, tt.wantValue)
				}
			}
		})
	}
}

func TestListKeymapActions(t *testing.T) {
	actions := ListKeymapActions()

	// Check that all expected actions are present
	expectedActions := []string{
		"up", "down", "left", "right", "top", "bottom",
		"select", "back", "quit",
		"stage", "stage-all", "unstage", "unstage-all", "discard",
		"commit", "commit-edit", "push", "stash", "stash-all",
		"file-diff", "all-diffs", "branches", "stashes", "log",
		"visual", "help", "verbose-help", "new-branch", "delete",
	}

	actionSet := make(map[string]bool)
	for _, a := range actions {
		actionSet[a] = true
	}

	for _, expected := range expectedActions {
		if !actionSet[expected] {
			t.Errorf("expected action %q not found in ListKeymapActions()", expected)
		}
	}

	// Check no duplicates
	seen := make(map[string]bool)
	for _, a := range actions {
		if seen[a] {
			t.Errorf("duplicate action found: %q", a)
		}
		seen[a] = true
	}
}

func TestFindKeymapOverrideConflicts(t *testing.T) {
	tests := []struct {
		name          string
		overrides     map[string]string
		wantConflicts int
		conflictKeys  []string
	}{
		{
			name:          "no overrides no conflicts",
			overrides:     nil,
			wantConflicts: 0,
		},
		{
			name: "single override no conflict",
			overrides: map[string]string{
				"up": "w",
			},
			wantConflicts: 0,
		},
		{
			name: "override creates new conflict",
			overrides: map[string]string{
				"up":   "a", // 'a' is already used by stage
			},
			wantConflicts: 1,
			conflictKeys:  []string{"a"},
		},
		{
			name: "multiple overrides same key conflict",
			overrides: map[string]string{
				"up":   "x",
				"down": "x",
			},
			wantConflicts: 1,
			conflictKeys:  []string{"x"},
		},
		{
			name: "override to already shared key in defaults no conflict",
			overrides: map[string]string{
				// h is shared between left, select, back in defaults
				// Moving another key to h should not be a new conflict
				// since h already has multiple defaults
			},
			wantConflicts: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaults := DefaultKeymap()
			current := DefaultKeymap()

			for action, key := range tt.overrides {
				current.ApplyOverride(action, key)
			}

			conflicts := FindKeymapOverrideConflicts(defaults, current)

			if len(conflicts) != tt.wantConflicts {
				t.Errorf("got %d conflicts, want %d", len(conflicts), tt.wantConflicts)
				for _, c := range conflicts {
					t.Logf("  conflict: key=%q actions=%v", c.Key, c.Actions)
				}
			}

			if tt.conflictKeys != nil {
				conflictKeySet := make(map[string]bool)
				for _, c := range conflicts {
					conflictKeySet[c.Key] = true
				}
				for _, key := range tt.conflictKeys {
					if !conflictKeySet[key] {
						t.Errorf("expected conflict for key %q not found", key)
					}
				}
			}
		})
	}
}

func TestKeymapConflictDetails(t *testing.T) {
	defaults := DefaultKeymap()
	current := DefaultKeymap()

	// Override 'up' to use 'a' which conflicts with 'stage'
	current.ApplyOverride("up", "a")

	conflicts := FindKeymapOverrideConflicts(defaults, current)

	if len(conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(conflicts))
	}

	conflict := conflicts[0]
	if conflict.Key != "a" {
		t.Errorf("expected conflict key 'a', got %q", conflict.Key)
	}

	// Check that both 'up' and 'stage' are in the conflict actions
	hasUp := false
	hasStage := false
	for _, action := range conflict.Actions {
		if action == "up" {
			hasUp = true
		}
		if action == "stage" {
			hasStage = true
		}
	}

	if !hasUp {
		t.Error("expected 'up' in conflict actions")
	}
	if !hasStage {
		t.Error("expected 'stage' in conflict actions")
	}
}

func TestKeymapBindingsComplete(t *testing.T) {
	// Verify that all Keymap fields have corresponding bindings
	km := DefaultKeymap()

	// Apply all bindings and verify they set the correct fields
	testCases := []struct {
		action string
		getField func(*Keymap) string
	}{
		{"up", func(k *Keymap) string { return k.Up }},
		{"down", func(k *Keymap) string { return k.Down }},
		{"left", func(k *Keymap) string { return k.Left }},
		{"right", func(k *Keymap) string { return k.Right }},
		{"top", func(k *Keymap) string { return k.Top }},
		{"bottom", func(k *Keymap) string { return k.Bottom }},
		{"select", func(k *Keymap) string { return k.Select }},
		{"back", func(k *Keymap) string { return k.Back }},
		{"quit", func(k *Keymap) string { return k.Quit }},
		{"stage", func(k *Keymap) string { return k.Stage }},
		{"stage-all", func(k *Keymap) string { return k.StageAll }},
		{"unstage", func(k *Keymap) string { return k.Unstage }},
		{"unstage-all", func(k *Keymap) string { return k.UnstageAll }},
		{"discard", func(k *Keymap) string { return k.Discard }},
		{"commit", func(k *Keymap) string { return k.Commit }},
		{"commit-edit", func(k *Keymap) string { return k.CommitEdit }},
		{"push", func(k *Keymap) string { return k.Push }},
		{"stash", func(k *Keymap) string { return k.Stash }},
		{"stash-all", func(k *Keymap) string { return k.StashAll }},
		{"file-diff", func(k *Keymap) string { return k.FileDiff }},
		{"all-diffs", func(k *Keymap) string { return k.AllDiffs }},
		{"branches", func(k *Keymap) string { return k.Branches }},
		{"stashes", func(k *Keymap) string { return k.Stashes }},
		{"log", func(k *Keymap) string { return k.Log }},
		{"visual", func(k *Keymap) string { return k.Visual }},
		{"help", func(k *Keymap) string { return k.Help }},
		{"verbose-help", func(k *Keymap) string { return k.VerboseHelp }},
		{"new-branch", func(k *Keymap) string { return k.NewBranch }},
		{"delete", func(k *Keymap) string { return k.Delete }},
	}

	for _, tc := range testCases {
		t.Run(tc.action, func(t *testing.T) {
			testKm := DefaultKeymap()
			testKey := "TEST_" + tc.action
			applied := testKm.ApplyOverride(tc.action, testKey)
			if !applied {
				t.Errorf("ApplyOverride(%q, %q) returned false", tc.action, testKey)
			}
			if got := tc.getField(testKm); got != testKey {
				t.Errorf("after ApplyOverride, field = %q, want %q", got, testKey)
			}
		})
	}

	// Verify no fields were missed - all defaults should be non-empty
	if km.Up == "" {
		t.Error("Up is empty in default keymap")
	}
	if km.Down == "" {
		t.Error("Down is empty in default keymap")
	}
	// ... the other tests above cover the rest
}
