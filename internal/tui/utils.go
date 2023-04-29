package tui

import (
	"bytes"
	"fmt"

	"github.com/alecthomas/chroma/quick"
)

const nbFormat = 2
const (
	YAML format = iota
	JSON
)

func Highlight(content, extension, syntaxTheme string) (string, error) {
	buf := new(bytes.Buffer)
	if err := quick.Highlight(buf, content, extension, "terminal256", syntaxTheme); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return buf.String(), nil
}
