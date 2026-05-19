package wizard

import (
	"fmt"
	"os"
	"strings"
)

// ANSI escape codes.
const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiDim    = "\033[2m"
	ansiCyan   = "\033[36m"
	ansiGreen  = "\033[32m"
	ansiWhite = "\033[37m"
)

// Styler applies ANSI styling when color is enabled.
type Styler struct {
	Enabled bool
}

// NewStyler creates a Styler with the given color state.
func NewStyler(color bool) *Styler {
	return &Styler{Enabled: color}
}

// IsTTY returns true if stdout is a terminal (not piped).
func IsTTY() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

func (s *Styler) wrap(code, text string) string {
	if !s.Enabled {
		return text
	}
	return code + text + ansiReset
}

func (s *Styler) Bold(text string) string  { return s.wrap(ansiBold, text) }
func (s *Styler) Dim(text string) string   { return s.wrap(ansiDim, text) }
func (s *Styler) Green(text string) string { return s.wrap(ansiGreen, text) }

func (s *Styler) BoldCyan(text string) string  { return s.wrap(ansiBold+ansiCyan, text) }
func (s *Styler) BoldGreen(text string) string { return s.wrap(ansiBold+ansiGreen, text) }
func (s *Styler) BoldWhite(text string) string { return s.wrap(ansiBold+ansiWhite, text) }

// Banner returns a styled welcome banner.
func (s *Styler) Banner() string {
	const w = 50
	top := "┌" + strings.Repeat("─", w) + "┐"
	bot := "└" + strings.Repeat("─", w) + "┘"
	blank := "│" + strings.Repeat(" ", w) + "│"
	title := boxLine("appteam", w, s.BoldCyan)
	subtitle := boxLine("Claude Code Agent Team Generator", w, s.Dim)
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s", top, blank, title, subtitle, bot)
}

// boxLine centers text inside a box row of the given width, applying styleFn.
func boxLine(text string, width int, styleFn func(string) string) string {
	visLen := len(text)
	left := (width - visLen) / 2
	right := width - visLen - left
	return "│" + strings.Repeat(" ", left) + styleFn(text) + strings.Repeat(" ", right) + "│"
}

// StepHeader returns a styled step header line.
func (s *Styler) StepHeader(step, total int, title string) string {
	return fmt.Sprintf("%s %s %s %s",
		s.Dim("━━"),
		s.Bold(fmt.Sprintf("Step %d of %d", step, total)),
		s.Dim("━━"),
		s.BoldCyan(title))
}

// PadBold pads a label to the given width, then applies bold styling.
// This ensures fmt %-Ns alignment works on visible width, not ANSI byte length.
func (s *Styler) PadBold(label string, width int) string {
	padded := fmt.Sprintf("%-*s", width, label)
	return s.Bold(padded)
}

// Divider returns a dimmed horizontal rule.
func (s *Styler) Divider() string {
	return s.Dim("──────────────────────────────────────────────────")
}
