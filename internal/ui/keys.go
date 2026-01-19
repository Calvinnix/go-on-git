package ui

import "strings"

// Keymap holds all configurable key bindings
type Keymap struct {
	// Navigation
	Up       string
	Down     string
	Left     string
	Right    string
	Top      string
	Bottom   string
	Select   string
	Back     string
	Quit     string

	// Actions
	Stage      string
	StageAll   string
	Unstage    string
	UnstageAll string
	Discard    string
	Commit     string
	CommitEdit string
	Push       string
	Stash      string
	StashAll   string

	// Views
	FileDiff string
	AllDiffs string
	Branches string
	Stashes  string
	Log      string

	// Modes
	Visual      string
	Help        string
	VerboseHelp string
	NewBranch   string
	Delete      string
}

// DefaultKeymap returns the default key bindings
func DefaultKeymap() *Keymap {
	return &Keymap{
		// Navigation
		Up:       "k",
		Down:     "j",
		Left:     "h",
		Right:    "l",
		Top:      "g",
		Bottom:   "G",
		Select:   "h",
		Back:     "h",
		Quit:     "q",

		// Actions
		Stage:      "a",
		StageAll:   "A",
		Unstage:    "u",
		UnstageAll: "U",
		Discard:    "d",
		Commit:     "c",
		CommitEdit: "C",
		Push:       "p",
		Stash:      "s",
		StashAll:   "S",

		// Views
		FileDiff: "l",
		AllDiffs: "i",
		Branches: "b",
		Stashes:  "e",
		Log:      "o",

		// Modes
		Visual:      "v",
		Help:        "?",
		VerboseHelp: "/",
		NewBranch:   "n",
		Delete:      "d",
	}
}

// ParseKeymapArg parses a keymap override argument in the format "action=key"
// Returns the action name and key, or empty strings if invalid
func ParseKeymapArg(arg string) (action, key string, valid bool) {
	parts := strings.SplitN(arg, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// ApplyOverride applies a single keymap override
func (k *Keymap) ApplyOverride(action, key string) bool {
	switch action {
	case "up":
		k.Up = key
	case "down":
		k.Down = key
	case "left":
		k.Left = key
	case "right":
		k.Right = key
	case "top":
		k.Top = key
	case "bottom":
		k.Bottom = key
	case "select":
		k.Select = key
	case "back":
		k.Back = key
	case "quit":
		k.Quit = key
	case "stage":
		k.Stage = key
	case "stage-all":
		k.StageAll = key
	case "unstage":
		k.Unstage = key
	case "unstage-all":
		k.UnstageAll = key
	case "discard":
		k.Discard = key
	case "commit":
		k.Commit = key
	case "commit-edit":
		k.CommitEdit = key
	case "push":
		k.Push = key
	case "stash":
		k.Stash = key
	case "stash-all":
		k.StashAll = key
	case "file-diff":
		k.FileDiff = key
	case "all-diffs":
		k.AllDiffs = key
	case "branches":
		k.Branches = key
	case "stashes":
		k.Stashes = key
	case "log":
		k.Log = key
	case "visual":
		k.Visual = key
	case "help":
		k.Help = key
	case "verbose-help":
		k.VerboseHelp = key
	case "new-branch":
		k.NewBranch = key
	case "delete":
		k.Delete = key
	default:
		return false
	}
	return true
}

// ListActions returns all available action names for help text
func ListKeymapActions() []string {
	return []string{
		"up", "down", "left", "right", "top", "bottom", "select", "back", "quit",
		"stage", "stage-all", "unstage", "unstage-all", "discard",
		"commit", "commit-edit", "push", "stash", "stash-all",
		"file-diff", "all-diffs", "branches", "stashes", "log",
		"visual", "help", "verbose-help", "new-branch", "delete",
	}
}

// Global keymap instance
var Keys = DefaultKeymap()
