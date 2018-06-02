package tabbed

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

// Show ...
func Show(r io.Reader) {
	scanner := bufio.NewScanner(r)

	runApplication(func() string {
		if scanner.Scan() {
			return strings.TrimFunc(scanner.Text(), unicode.IsSpace)
		}
		return ""
	})
}
