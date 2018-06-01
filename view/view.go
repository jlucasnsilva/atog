package view

import (
	"io"
	"io/ioutil"

	"github.com/jlucasnsilva/atog/atog"
)

// Show ...
func Show(r io.Reader) {
	atog.RunSimpleApp(func() string {
		bytes, _ := ioutil.ReadAll(r)
		return string(bytes)
	})
}
