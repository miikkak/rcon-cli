package rcon

import (
	"runtime"
	"testing"
)

func TestColorize(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("colorize is a no-op on windows")
	}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "plain text",
			input: "hello",
			want:  "hello",
		},
		{
			name:  "single legacy code",
			input: "§ahello",
			want:  "\x1b[92mhello" + reset,
		},
		{
			name:  "newline reset",
			input: "§cline1\nline2",
			want:  "\x1b[91mline1\n" + reset + "line2" + reset,
		},
		{
			name:  "truecolor sequence",
			input: "§x§1§2§3§4§5§6hello",
			want:  "\x1b[38;2;18;52;86mhello" + reset,
		},
		{
			name:  "truecolor uppercase hex",
			input: "§x§A§B§C§D§E§Fhello",
			want:  "\x1b[38;2;171;205;239mhello" + reset,
		},
		{
			name:  "mixed legacy and truecolor",
			input: "§cred §x§0§0§F§F§0§0cyan-ish§rreset",
			want:  "\x1b[91mred \x1b[38;2;0;255;0mcyan-ish" + reset + "reset" + reset,
		},
		{
			name:  "multiple truecolor runs",
			input: "§x§1§1§1§1§1§1a§x§2§2§2§2§2§2b",
			want:  "\x1b[38;2;17;17;17ma\x1b[38;2;34;34;34mb" + reset,
		},
		{
			name:  "uppercase legacy code",
			input: "§Chello",
			want:  "\x1b[91mhello" + reset,
		},
		{
			name:  "truncated truecolor falls back to legacy processing",
			input: "§x§1§2§3",
			want:  "§x\x1b[34m\x1b[32m\x1b[36m" + reset,
		},
		{
			name:  "section-x immediately followed by non-section char",
			input: "§xhello",
			want:  "§xhello",
		},
		{
			name:  "section-x at end of string",
			input: "foo§x",
			want:  "foo§x",
		},
		{
			name:  "unrecognized code passes through literally",
			input: "§zfoo",
			want:  "§zfoo",
		},
		{
			name:  "raw escape byte is stripped",
			input: "foo\x1b]0;pwned\x07bar",
			want:  "foo]0;pwnedbar",
		},
		{
			name:  "newline and tab preserved",
			input: "foo\n\tbar",
			want:  "foo\n\tbar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := colorize(tt.input)
			if got != tt.want {
				t.Errorf("colorize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestColorizeAllLegacyCodes(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("colorize is a no-op on windows")
	}

	for code, ansi := range colors {
		t.Run(string(code), func(t *testing.T) {
			input := "§" + string(code) + "text"
			want := ansi + "text" + reset
			got := colorize(input)
			if got != want {
				t.Errorf("colorize(%q) = %q, want %q", input, got, want)
			}
		})
	}
}
