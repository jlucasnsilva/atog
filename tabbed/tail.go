package tabbed

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	fileTagRE = regexp.MustCompile(`^==> .* <==$`)
)

func isFileTag(line string) bool {
	return fileTagRE.MatchString(strings.TrimFunc(line, unicode.IsSpace))
}

func getFilename(line string) string {
	s := strings.TrimPrefix(line, "==> ")
	return strings.TrimSuffix(s, " <==\n")
}
