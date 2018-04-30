package ui

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

func format(currentText string, newText []byte) string {
	lines := strings.Split(string(newText), "\n")
	return fmt.Sprintf("%v\n%v", currentText, formatText(lines))
}

func formatText(lines []string) string {
	return joinLines(processText(lines))
}

func joinLines(lines []string) string {
	ls := make([]string, len(lines))

	for i, line := range lines {
		ls[i] = fmt.Sprintf("%v", line)
	}

	return strings.Join(ls, "\n\n")
}

func processText(lines []string) []string {
	s := make([]string, len(lines))

	for i, line := range lines {
		s[i] = processLine(line)
	}

	return s
}

func processLine(line string) string {
	regex := regexp.MustCompile(`:\s+`)
	parts := regex.Split(line, -1)
	whole := strings.Join(parts, ":\n\t")
	return highlight(whole)
}

func insertLineBreaks(line string) string {
	regex := regexp.MustCompile(`:\s+`)
	parts := regex.Split(line, -1)
	return strings.Join(parts, ":\n\t")
}

func highlight(s string) string {
	runes := []rune(s)
	runesCount := len(runes)
	out := bytes.Buffer{}
	opposite := map[rune]rune{
		'{': '}',
		'[': ']',
		'(': ')',
	}

	for i := 0; i < runesCount; i++ {
		c := runes[i]

		switch {
		case strings.ContainsRune("{[(", c):
			out.WriteString("[#c83737]")
			out.WriteRune(c)
			out.WriteString("[white:#162d50]")

			i++
			balance := 1
			o := opposite[c]
			for ; i < runesCount && balance > 0; i++ {
				x := runes[i]
				if x == o {
					balance--
					if balance > 0 {
						out.WriteRune(x)
					} else {
						out.WriteString("[#c83737:black]")
						out.WriteRune(x)
						out.WriteString("[white]")
					}
				} else {
					out.WriteRune(x)
					if x == c {
						balance++
					}
				}
			}
		case strings.ContainsRune("\"'`", c):
			out.WriteString("[#217844]")
			out.WriteRune(c)

			i++
			for ; i < runesCount; i++ {
				x := runes[i]
				out.WriteRune(x)

				if c == x {
					out.WriteString("[white]")
					break
				}
			}
		case strings.ContainsRune("/\\:-=.", c):
			out.WriteString("[#c83737]")
			out.WriteRune(c)
			out.WriteString("[white]")
		default:
			out.WriteRune(c)
		}
	}

	out.WriteString("[white:black]")
	r := regexp.MustCompile(`(line\s\d+)|(L\d+)|(col\s\d+)|(C\d+)|(column\s\d+)`)
	return r.ReplaceAllStringFunc(out.String(), func(x string) string {
		return fmt.Sprintf("[#c83737]%v[white]", x)
	})
}
