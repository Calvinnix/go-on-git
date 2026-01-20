package main

import (
	"strings"
	"testing"

	"go-on-git/internal/ui"
)

// TestKeymapArgParsing tests the --key.action=key argument parsing
func TestKeymapArgParsing(t *testing.T) {
	tests := []struct {
		name       string
		arg        string
		wantAction string
		wantKey    string
		wantValid  bool
	}{
		{
			name:       "valid override",
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
			name:       "valid special key",
			arg:        "quit=ctrl+c",
			wantAction: "quit",
			wantKey:    "ctrl+c",
			wantValid:  true,
		},
		{
			name:       "invalid no equals",
			arg:        "upw",
			wantAction: "",
			wantKey:    "",
			wantValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, key, valid := ui.ParseKeymapArg(tt.arg)
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

// TestKeymapOverrideApplication tests applying keymap overrides
func TestKeymapOverrideApplication(t *testing.T) {
	// Save and restore global Keys
	originalKeys := ui.Keys
	defer func() { ui.Keys = originalKeys }()

	// Reset to default
	ui.Keys = ui.DefaultKeymap()

	// Apply valid override
	if !ui.Keys.ApplyOverride("up", "w") {
		t.Error("ApplyOverride should return true for valid action")
	}
	if ui.Keys.Up != "w" {
		t.Errorf("Up = %q, want 'w'", ui.Keys.Up)
	}

	// Apply invalid override
	if ui.Keys.ApplyOverride("invalid-action", "x") {
		t.Error("ApplyOverride should return false for invalid action")
	}
}

// TestKeymapConflictDetection tests detecting keymap conflicts
func TestKeymapConflictDetection(t *testing.T) {
	defaults := ui.DefaultKeymap()
	current := ui.DefaultKeymap()

	// No conflicts with default
	conflicts := ui.FindKeymapOverrideConflicts(defaults, current)
	if len(conflicts) != 0 {
		t.Errorf("expected 0 conflicts with defaults, got %d", len(conflicts))
	}

	// Create a conflict
	current.ApplyOverride("up", "a") // 'a' is used by stage

	conflicts = ui.FindKeymapOverrideConflicts(defaults, current)
	if len(conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(conflicts))
	}
	if len(conflicts) > 0 && conflicts[0].Key != "a" {
		t.Errorf("conflict key = %q, want 'a'", conflicts[0].Key)
	}
}

// TestListKeymapActions tests listing all keymap actions
func TestListKeymapActions(t *testing.T) {
	actions := ui.ListKeymapActions()

	// Should contain all expected actions
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
			t.Errorf("expected action %q not in list", expected)
		}
	}
}

// TestVersion tests the version constant
func TestVersion(t *testing.T) {
	if version == "" {
		t.Error("version should not be empty")
	}
	// Version should be in semver format (e.g., "0.7.0")
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		t.Errorf("version %q is not in semver format", version)
	}
}

// TestCommandLineArgsSimulation tests simulating command line argument processing
func TestCommandLineArgsSimulation(t *testing.T) {
	// Test argument parsing patterns used in main()
	tests := []struct {
		name       string
		arg        string
		isHelp     bool
		isVersion  bool
		isHideHelp bool
		isKeymap   bool
		isUnknown  bool
	}{
		{"help short", "-h", true, false, false, false, false},
		{"help long", "--help", true, false, false, false, false},
		{"version short", "-v", false, true, false, false, false},
		{"version long", "--version", false, true, false, false, false},
		{"hide help", "--hide-help", false, false, true, false, false},
		{"keymap override", "--key.up=w", false, false, false, true, false},
		{"unknown", "--unknown", false, false, false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arg := tt.arg
			isHelp := arg == "--help" || arg == "-h"
			isVersion := arg == "--version" || arg == "-v"
			isHideHelp := arg == "--hide-help"
			isKeymap := strings.HasPrefix(arg, "--key.")

			if isHelp != tt.isHelp {
				t.Errorf("isHelp = %v, want %v", isHelp, tt.isHelp)
			}
			if isVersion != tt.isVersion {
				t.Errorf("isVersion = %v, want %v", isVersion, tt.isVersion)
			}
			if isHideHelp != tt.isHideHelp {
				t.Errorf("isHideHelp = %v, want %v", isHideHelp, tt.isHideHelp)
			}
			if isKeymap != tt.isKeymap {
				t.Errorf("isKeymap = %v, want %v", isKeymap, tt.isKeymap)
			}
		})
	}
}

// TestKeyPrefixTrim tests the --key. prefix trimming
func TestKeyPrefixTrim(t *testing.T) {
	arg := "--key.up=w"
	override := strings.TrimPrefix(arg, "--key.")

	if override != "up=w" {
		t.Errorf("override = %q, want 'up=w'", override)
	}

	action, key, valid := ui.ParseKeymapArg(override)
	if !valid {
		t.Error("should be valid")
	}
	if action != "up" {
		t.Errorf("action = %q, want 'up'", action)
	}
	if key != "w" {
		t.Errorf("key = %q, want 'w'", key)
	}
}

// TestMultipleKeymapOverrides tests applying multiple keymap overrides
func TestMultipleKeymapOverrides(t *testing.T) {
	km := ui.DefaultKeymap()

	overrides := []struct {
		action string
		key    string
	}{
		{"up", "w"},
		{"down", "s"},
		{"left", "a"},
		{"right", "d"},
	}

	for _, o := range overrides {
		if !km.ApplyOverride(o.action, o.key) {
			t.Errorf("failed to apply override %s=%s", o.action, o.key)
		}
	}

	if km.Up != "w" {
		t.Errorf("Up = %q, want 'w'", km.Up)
	}
	if km.Down != "s" {
		t.Errorf("Down = %q, want 's'", km.Down)
	}
	if km.Left != "a" {
		t.Errorf("Left = %q, want 'a'", km.Left)
	}
	if km.Right != "d" {
		t.Errorf("Right = %q, want 'd'", km.Right)
	}
}

// TestKeymapConflictDetails tests the details of keymap conflicts
func TestKeymapConflictDetails(t *testing.T) {
	defaults := ui.DefaultKeymap()
	current := ui.DefaultKeymap()

	// Create multiple actions mapped to same key
	current.ApplyOverride("up", "x")
	current.ApplyOverride("down", "x")

	conflicts := ui.FindKeymapOverrideConflicts(defaults, current)

	if len(conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(conflicts))
	}

	conflict := conflicts[0]
	if conflict.Key != "x" {
		t.Errorf("conflict key = %q, want 'x'", conflict.Key)
	}

	// Should have both up and down in actions
	hasUp := false
	hasDown := false
	for _, action := range conflict.Actions {
		if action == "up" {
			hasUp = true
		}
		if action == "down" {
			hasDown = true
		}
	}

	if !hasUp {
		t.Error("expected 'up' in conflict actions")
	}
	if !hasDown {
		t.Error("expected 'down' in conflict actions")
	}
}

// TestAllKeymapActionsApplicable tests that all listed actions can be applied
func TestAllKeymapActionsApplicable(t *testing.T) {
	actions := ui.ListKeymapActions()

	for _, action := range actions {
		km := ui.DefaultKeymap()
		if !km.ApplyOverride(action, "TEST") {
			t.Errorf("action %q could not be applied", action)
		}
	}
}

// TestDefaultKeymapValues tests all default keymap values
func TestDefaultKeymapValues(t *testing.T) {
	km := ui.DefaultKeymap()

	// Test all default values match expected
	expectations := map[string]string{
		"Up":          "k",
		"Down":        "j",
		"Left":        "h",
		"Right":       "l",
		"Top":         "g",
		"Bottom":      "G",
		"Select":      "h",
		"Back":        "h",
		"Quit":        "q",
		"Stage":       "a",
		"StageAll":    "A",
		"Unstage":     "u",
		"UnstageAll":  "U",
		"Discard":     "d",
		"Commit":      "c",
		"CommitEdit":  "C",
		"Push":        "p",
		"Stash":       "s",
		"StashAll":    "S",
		"FileDiff":    "l",
		"AllDiffs":    "i",
		"Branches":    "b",
		"Stashes":     "e",
		"Log":         "o",
		"Visual":      "v",
		"Help":        "?",
		"VerboseHelp": "/",
		"NewBranch":   "n",
		"Delete":      "d",
	}

	values := map[string]string{
		"Up":          km.Up,
		"Down":        km.Down,
		"Left":        km.Left,
		"Right":       km.Right,
		"Top":         km.Top,
		"Bottom":      km.Bottom,
		"Select":      km.Select,
		"Back":        km.Back,
		"Quit":        km.Quit,
		"Stage":       km.Stage,
		"StageAll":    km.StageAll,
		"Unstage":     km.Unstage,
		"UnstageAll":  km.UnstageAll,
		"Discard":     km.Discard,
		"Commit":      km.Commit,
		"CommitEdit":  km.CommitEdit,
		"Push":        km.Push,
		"Stash":       km.Stash,
		"StashAll":    km.StashAll,
		"FileDiff":    km.FileDiff,
		"AllDiffs":    km.AllDiffs,
		"Branches":    km.Branches,
		"Stashes":     km.Stashes,
		"Log":         km.Log,
		"Visual":      km.Visual,
		"Help":        km.Help,
		"VerboseHelp": km.VerboseHelp,
		"NewBranch":   km.NewBranch,
		"Delete":      km.Delete,
	}

	for field, expected := range expectations {
		if values[field] != expected {
			t.Errorf("%s = %q, want %q", field, values[field], expected)
		}
	}
}
