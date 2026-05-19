package wizard

import (
	"strings"
	"testing"
)

func TestStylerColorDisabled(t *testing.T) {
	s := NewStyler(false)

	tests := []struct {
		name string
		fn   func(string) string
		want string
	}{
		{"Bold", s.Bold, "x"},
		{"Dim", s.Dim, "x"},
		{"Green", s.Green, "x"},
		{"BoldCyan", s.BoldCyan, "x"},
		{"BoldGreen", s.BoldGreen, "x"},
		{"BoldWhite", s.BoldWhite, "x"},
	}
	for _, tt := range tests {
		got := tt.fn("x")
		if got != tt.want {
			t.Errorf("%s(\"x\") with color off = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestStylerColorEnabled(t *testing.T) {
	s := NewStyler(true)

	tests := []struct {
		name string
		fn   func(string) string
		want string
	}{
		{"Bold", s.Bold, "\033[1mx\033[0m"},
		{"Dim", s.Dim, "\033[2mx\033[0m"},
		{"Green", s.Green, "\033[32mx\033[0m"},
		{"BoldCyan", s.BoldCyan, "\033[1m\033[36mx\033[0m"},
		{"BoldGreen", s.BoldGreen, "\033[1m\033[32mx\033[0m"},
		{"BoldWhite", s.BoldWhite, "\033[1m\033[37mx\033[0m"},
	}
	for _, tt := range tests {
		got := tt.fn("x")
		if got != tt.want {
			t.Errorf("%s(\"x\") with color on = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestPadBold(t *testing.T) {
	// Color disabled: should just be padded text
	s := NewStyler(false)
	got := s.PadBold("Hi:", 10)
	if len(got) != 10 {
		t.Errorf("PadBold color off: len = %d, want 10; got %q", len(got), got)
	}
	if !strings.HasPrefix(got, "Hi:") {
		t.Errorf("PadBold color off should start with 'Hi:', got %q", got)
	}

	// Color enabled: contains ANSI codes but wraps padded text
	sc := NewStyler(true)
	gotColor := sc.PadBold("Hi:", 10)
	if !strings.Contains(gotColor, "\033[1m") {
		t.Error("PadBold with color should contain bold ANSI code")
	}
	if !strings.Contains(gotColor, "\033[0m") {
		t.Error("PadBold with color should contain reset ANSI code")
	}
	if !strings.Contains(gotColor, "Hi:") {
		t.Error("PadBold should contain the label text")
	}
}

func TestBanner(t *testing.T) {
	s := NewStyler(false)
	banner := s.Banner()

	if !strings.Contains(banner, "appteam") {
		t.Error("Banner should contain 'appteam'")
	}
	if !strings.Contains(banner, "Claude Code Agent Team Generator") {
		t.Error("Banner should contain 'Claude Code Agent Team Generator'")
	}
	if !strings.Contains(banner, "┌") {
		t.Error("Banner should contain top-left box character")
	}
	if !strings.Contains(banner, "└") {
		t.Error("Banner should contain bottom-left box character")
	}
}

func TestStepHeader(t *testing.T) {
	s := NewStyler(false)
	header := s.StepHeader(2, 7, "Git Repository")

	if !strings.Contains(header, "Step 2 of 7") {
		t.Errorf("StepHeader should contain 'Step 2 of 7', got %q", header)
	}
	if !strings.Contains(header, "Git Repository") {
		t.Errorf("StepHeader should contain 'Git Repository', got %q", header)
	}
}

func TestDivider(t *testing.T) {
	s := NewStyler(false)
	divider := s.Divider()

	if !strings.Contains(divider, "─") {
		t.Error("Divider should contain '─' characters")
	}
	if len(divider) < 10 {
		t.Error("Divider should be a substantial horizontal rule")
	}
}
