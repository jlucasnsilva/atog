package tabbed

import (
	"bufio"
	"io"
)

// Show ...
func Show(r io.Reader) {
	scanner := bufio.NewScanner(r)

	runApplication(func() string {
		if scanner.Scan() {
			return scanner.Text()
		}
		return ""
	})
}
