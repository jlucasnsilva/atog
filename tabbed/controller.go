package tabbed

import (
	"bytes"
	"fmt"

	"github.com/jlucasnsilva/atog/atog"
	"github.com/rivo/tview"
)

type (
	logBuffer struct {
		buff  *bytes.Buffer
		dirty bool
	}

	controller struct {
		currentFile string
		fileList    *tview.List
		textView    *tview.TextView
		indexOf     map[string]int
		invIndexOf  []string
		buffers     map[string]*logBuffer
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

func newController(fileList *tview.List, textView *tview.TextView) *controller {
	con := &controller{}

	con.fileList = fileList
	con.textView = textView
	con.indexOf = make(map[string]int)
	con.buffers = make(map[string]*logBuffer)
	con.invIndexOf = make([]string, 0, 10)

	fileList.SetChangedFunc(func(i int, title string, description string, shortcut rune) {
		filename := con.invIndexOf[i]
		b := con.buffers[filename]
		fmt.Fprintln(con.textView, b.buff.String())
		con.fileList.SetItemText(i, filename, "")
	}).SetCurrentItem(1)

	return con
}

func (c *controller) switchTo(filename string) {
	c.currentFile = filename

	if _, ok := c.indexOf[filename]; !ok {
		i := len(c.invIndexOf)
		c.buffers[filename] = &logBuffer{buff: &bytes.Buffer{}}
		c.indexOf[filename] = i
		c.invIndexOf = append(c.invIndexOf, filename)
		c.fileList.AddItem(filename, "", shortcut(i), nil)
	}
}

func (c *controller) write(s string) {
	filename := c.currentFile
	b := c.buffers[filename]
	i := c.indexOf[filename]
	b.buff.WriteString(atog.Highlight(s))
	c.fileList.SetItemText(i, fmt.Sprintf("[*] %v", filename), "")
}

func (b *logBuffer) Write(p []byte) (n int, err error) {
	pc := len(p)

	if pc < 1 {
		return 0, nil
	}

	b.dirty = true
	return b.buff.Write(p)
}

func (b *logBuffer) flush(tv *tview.TextView) {
	fmt.Fprint(tv, b.buff.String())
	tv.ScrollToEnd()
	b.dirty = false
}
