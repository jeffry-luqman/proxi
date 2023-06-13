package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
)

const (
	Reset uint8 = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
)

const (
	Black uint8 = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
)

func Fmt(text any, attribute ...uint8) string {
	s := fmt.Sprintf("%v", text)
	if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return s
	}
	format := make([]string, len(attribute))
	for i, v := range attribute {
		format[i] = fmt.Sprintf("%v", v)
	}
	return "\x1b[" + strings.Join(format, ";") + "m" + s + "\x1b[0m"
}
