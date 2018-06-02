package watch

import (
	"bufio"
	"io"

	"github.com/jlucasnsilva/atog/atog"
)

// Show ...
func Show(r io.Reader) {
	scanner := bufio.NewScanner(r)

	atog.RunSimpleLoopApp(func() string {
		if scanner.Scan() {
			return atog.Highlight(scanner.Text())
		}
		return ""
	})
}
