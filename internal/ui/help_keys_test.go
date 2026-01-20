package ui

import (
	"testing"
)

func TestFormatKeyList(t *testing.T) {
	tests := []struct {
		name string
		keys []string
		want string
	}{
		{
			name: "single key",
			keys: []string{"j"},
			want: "j",
		},
		{
			name: "two keys",
			keys: []string{"j", "k"},
			want: "j/k",
		},
		{
			name: "three keys",
			keys: []string{"a", "b", "c"},
			want: "a/b/c",
		},
		{
			name: "duplicate keys",
			keys: []string{"j", "j", "k"},
			want: "j/k",
		},
		{
			name: "with esc",
			keys: []string{"h", "esc"},
			want: "h/ESC",
		},
		{
			name: "with enter",
			keys: []string{"l", "enter"},
			want: "l/Enter",
		},
		{
			name: "with space",
			keys: []string{" "},
			want: "SPACE",
		},
		{
			name: "empty keys filtered",
			keys: []string{"a", "", "b"},
			want: "a/b",
		},
		{
			name: "all empty",
			keys: []string{"", ""},
			want: "",
		},
		{
			name: "no keys",
			keys: []string{},
			want: "",
		},
		{
			name: "mixed special keys",
			keys: []string{"h", "esc", "enter", " "},
			want: "h/ESC/Enter/SPACE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatKeyList(tt.keys...)
			if got != tt.want {
				t.Errorf("formatKeyList(%v) = %q, want %q", tt.keys, got, tt.want)
			}
		})
	}
}

func TestFormatDoubleKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{
			name: "single char",
			key:  "g",
			want: "gg",
		},
		{
			name: "uppercase single char",
			key:  "G",
			want: "GG",
		},
		{
			name: "multi char",
			key:  "ctrl+a",
			want: "ctrl+a ctrl+a",
		},
		{
			name: "esc becomes ESC ESC",
			key:  "esc",
			want: "ESC ESC",
		},
		{
			name: "enter becomes Enter Enter",
			key:  "enter",
			want: "Enter Enter",
		},
		{
			name: "space becomes SPACE SPACE",
			key:  " ",
			want: "SPACE SPACE",
		},
		{
			name: "empty string",
			key:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDoubleKey(tt.key)
			if got != tt.want {
				t.Errorf("formatDoubleKey(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

func TestFormatKeyLabel(t *testing.T) {
	tests := []struct {
		key  string
		want string
	}{
		{"esc", "ESC"},
		{"enter", "Enter"},
		{" ", "SPACE"},
		{"a", "a"},
		{"j", "j"},
		{"ctrl+c", "ctrl+c"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := formatKeyLabel(tt.key)
			if got != tt.want {
				t.Errorf("formatKeyLabel(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}
