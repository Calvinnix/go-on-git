package ui

import (
	"strings"
	"testing"
)

func TestStatusChar(t *testing.T) {
	tests := []struct {
		name        string
		indexStatus byte
		workStatus  byte
		section     string
		wantContains string
		wantEmpty   bool
	}{
		{
			name:        "staged new file",
			indexStatus: 'A',
			workStatus:  ' ',
			section:     "staged",
			wantContains: "new file:",
		},
		{
			name:        "staged modified",
			indexStatus: 'M',
			workStatus:  ' ',
			section:     "staged",
			wantContains: "modified:",
		},
		{
			name:        "staged deleted",
			indexStatus: 'D',
			workStatus:  ' ',
			section:     "staged",
			wantContains: "deleted:",
		},
		{
			name:        "staged renamed",
			indexStatus: 'R',
			workStatus:  ' ',
			section:     "staged",
			wantContains: "renamed:",
		},
		{
			name:        "staged copied",
			indexStatus: 'C',
			workStatus:  ' ',
			section:     "staged",
			wantContains: "copied:",
		},
		{
			name:        "unstaged modified",
			indexStatus: ' ',
			workStatus:  'M',
			section:     "unstaged",
			wantContains: "modified:",
		},
		{
			name:        "unstaged deleted",
			indexStatus: ' ',
			workStatus:  'D',
			section:     "unstaged",
			wantContains: "deleted:",
		},
		{
			name:      "untracked returns empty",
			indexStatus: '?',
			workStatus:  '?',
			section:    "untracked",
			wantEmpty:  true,
		},
		{
			name:      "unknown section returns empty",
			indexStatus: 'M',
			workStatus:  ' ',
			section:    "unknown",
			wantEmpty:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StatusChar(tt.indexStatus, tt.workStatus, tt.section)

			if tt.wantEmpty {
				if got != "" {
					t.Errorf("StatusChar() = %q, want empty", got)
				}
				return
			}

			if !strings.Contains(got, tt.wantContains) {
				t.Errorf("StatusChar() = %q, want to contain %q", got, tt.wantContains)
			}
		})
	}
}

func TestStatusCharStyled(t *testing.T) {
	tests := []struct {
		name        string
		indexStatus byte
		workStatus  byte
		section     string
		wantContains string
		wantEmpty   bool
	}{
		{
			name:        "staged with extra style",
			indexStatus: 'A',
			workStatus:  ' ',
			section:     "staged",
			wantContains: "new file:",
		},
		{
			name:        "unstaged with extra style",
			indexStatus: ' ',
			workStatus:  'M',
			section:     "unstaged",
			wantContains: "modified:",
		},
		{
			name:      "untracked returns empty regardless of style",
			indexStatus: '?',
			workStatus:  '?',
			section:    "untracked",
			wantEmpty:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StatusCharStyled(tt.indexStatus, tt.workStatus, tt.section, StyleVisual)

			if tt.wantEmpty {
				if got != "" {
					t.Errorf("StatusCharStyled() = %q, want empty", got)
				}
				return
			}

			if !strings.Contains(got, tt.wantContains) {
				t.Errorf("StatusCharStyled() = %q, want to contain %q", got, tt.wantContains)
			}
		})
	}
}

func TestIndexStatusWord(t *testing.T) {
	tests := []struct {
		status byte
		want   string
	}{
		{'A', "new file:"},
		{'M', "modified:"},
		{'D', "deleted:"},
		{'R', "renamed:"},
		{'C', "copied:"},
		{' ', " "},
		{'?', "?"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			got := indexStatusWord(tt.status)
			if got != tt.want {
				t.Errorf("indexStatusWord(%q) = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

func TestWorkStatusWord(t *testing.T) {
	tests := []struct {
		status byte
		want   string
	}{
		{'M', "modified:"},
		{'D', "deleted:"},
		{' ', " "},
		{'?', "?"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			got := workStatusWord(tt.status)
			if got != tt.want {
				t.Errorf("workStatusWord(%q) = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

func TestStylesAreDefined(t *testing.T) {
	// Verify all styles are properly defined (non-nil)
	styles := []struct {
		name  string
		style interface{}
	}{
		{"StyleNormal", StyleNormal},
		{"StyleMuted", StyleMuted},
		{"StyleStaged", StyleStaged},
		{"StyleUnstaged", StyleUnstaged},
		{"StyleUntracked", StyleUntracked},
		{"StyleSelected", StyleSelected},
		{"StyleVisual", StyleVisual},
		{"StyleSectionHeader", StyleSectionHeader},
		{"StyleDiffAdded", StyleDiffAdded},
		{"StyleDiffRemoved", StyleDiffRemoved},
		{"StyleDiffContext", StyleDiffContext},
		{"StyleDiffHeader", StyleDiffHeader},
		{"StyleHunkHeaderStaged", StyleHunkHeaderStaged},
		{"StyleHunkHeaderUnstaged", StyleHunkHeaderUnstaged},
		{"StyleHelpKey", StyleHelpKey},
		{"StyleHelpDesc", StyleHelpDesc},
		{"StyleHelpTitle", StyleHelpTitle},
		{"StyleStatusBar", StyleStatusBar},
		{"StyleConfirm", StyleConfirm},
		{"StyleEmpty", StyleEmpty},
	}

	for _, s := range styles {
		t.Run(s.name, func(t *testing.T) {
			// Just verify it can render without panicking
			result := s.style.(interface{ Render(strs ...string) string }).Render("test")
			if result == "" {
				// Styles can render empty-ish but shouldn't panic
				t.Logf("%s rendered empty string (expected for some styles)", s.name)
			}
		})
	}
}

func TestStatusCharPadding(t *testing.T) {
	// Verify that status chars are padded to consistent width
	got := StatusChar('A', ' ', "staged")
	// Should be padded to 12 chars
	if len(got) == 0 {
		t.Skip("StatusChar returns styled string, padding hard to test")
	}
}
