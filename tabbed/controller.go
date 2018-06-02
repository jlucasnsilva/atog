package tabbed

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/gdamore/tcell"

	"github.com/jlucasnsilva/atog/atog"
	"github.com/rivo/tview"
)

type (
	controller struct {
		currentFile   string
		app           *tview.Application
		fileList      *tview.List
		textView      *tview.TextView
		fileByIndex   []string
		bufferByIndex []*bytes.Buffer
		buffers       map[string]*bytes.Buffer
	}
)

func shortcut(i int) rune {
	shortcuts := []rune{
		'a', 'b', 'c', 'd', 'e',
		'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o',
		'p', 'q', 'r', 's', 't',
		'u', 'v', 'w', 'x', 'y',
		'z',
	}

	if i >= len(shortcuts) || i < 0 {
		return ' '
	}
	return shortcuts[i]
}

func newController(fileList *tview.List, textView *tview.TextView, app *tview.Application) *controller {
	controller := &controller{
		buffers:       make(map[string]*bytes.Buffer),
		bufferByIndex: make([]*bytes.Buffer, 0, 10),
		fileByIndex:   make([]string, 0, 10),
		textView:      textView,
		fileList:      fileList,
		app:           app,
	}

	app.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyCtrlL {
			if i := fileList.GetCurrentItem(); i < len(controller.bufferByIndex) {
				b := controller.bufferByIndex[i]
				b.Reset()
			}
			textView.Clear()
			app.Draw()
			return nil
		}
		return e
	})

	fileList.SetSelectedBackgroundColor(tcell.ColorDarkSlateBlue)

	fileList.SetChangedFunc(func(i int, title, subtitle string, sc rune) {
		if i < len(controller.bufferByIndex) {
			textView.Clear()
			b := controller.bufferByIndex[i]
			fmt.Fprintln(textView, b.String())
			controller.currentFile = controller.fileByIndex[i]
		}
	})
	return controller
}

func (c *controller) write(filename, text string) {
	if len(text) < 1 {
		return
	}

	b := c.buffers[filename]
	htxt := atog.Highlight(strings.TrimFunc(text, unicode.IsSpace))

	fmt.Fprintln(b, htxt)
	if c.currentFile == filename {
		fmt.Fprintln(c.textView, htxt)
	}
}

func (c *controller) createComponentsIfNeeded(filename string) {
	if _, ok := c.buffers[filename]; !ok {
		i := len(c.bufferByIndex)
		buffer := &bytes.Buffer{}

		c.buffers[filename] = buffer
		c.fileByIndex = append(c.fileByIndex, filename)
		c.bufferByIndex = append(c.bufferByIndex, buffer)
		if c.currentFile == "" {
			c.currentFile = filename
		}

		c.fileList.AddItem(filename, "", shortcut(i), nil)
	}
}
