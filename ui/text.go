package ui

import (
	"fmt"
	"regexp"
	"strings"
)

func processText(text string) string {
	lines := strings.Split(text, "\n")
	s := make([]string, len(lines))

	for i, line := range lines {
		s[i] = processLine(line)
	}

	return strings.Join(s, "\n\n\n")
}

func processLine(line string) string {
	regex := regexp.MustCompile(`:\s+`)
	parts := regex.Split(line, -1)
	whole := strings.Join(parts, ":\n\t")
	return highlight(whole)
}

func highlight(s string) string {
	runes := []rune(s)
	out := make([]rune, 0, len(runes))

	for _, c := range runes {
		if strings.ContainsRune("{[(", c) {
			add := []rune(fmt.Sprintf("[#c83737]%v[white:#162d50]", string(c)))
			out = append(out, add...)
		} else if strings.ContainsRune("}])", c) {
			add := []rune(fmt.Sprintf("[#c83737:black]%v[white]", string(c)))
			out = append(out, add...)
		} else if strings.ContainsRune("/\\:-=", c) {
			add := []rune(fmt.Sprintf("[#c83737]%v[white]", string(c)))
			out = append(out, add...)
		} else {
			out = append(out, c)
		}
	}

	return string(out)
}
