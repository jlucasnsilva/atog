package ui

import (
	"fmt"

	"github.com/jlucasnsilva/atog/stalker"

	"github.com/rivo/tview"
)

const atogTag = "# ATOG"

type controller struct {
	textView   *tview.TextView
	fileList   *tview.List
	indexes    map[string]int
	invIndexes []string
	logViews   map[string]*logView
}

func shortcut(i int) rune {
	shortcuts := []rune{
		'a', 'b', 'c', 'd', 'e',
		'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o',
		'p', 'q', 'r', 's', 't',
		'u', 'v', 'w', 'x', 'y',
		'z',
	}

	if i >= len(shortcuts) {
		return ' '
	}
	return shortcuts[i]
}

func newRouter(filenames []string, textView *tview.TextView, fileList *tview.List) *controller {
	fns := append([]string{atogTag}, filenames...)
	con := &controller{}
	con.textView = textView
	con.fileList = fileList
	con.indexes = make(map[string]int)
	con.logViews = make(map[string]*logView)
	con.invIndexes = make([]string, len(fns))

	for i, fn := range fns {
		if fn != atogTag {
			con.fileList.AddItem(fn, "", shortcut(i), nil)
		} else {
			con.fileList.AddItem(fn+" (0)", "Errors log", shortcut(i), nil)
		}
		con.logViews[fn] = newLogView()
		con.indexes[fn] = i
		con.invIndexes[i] = fn
	}

	fileList.SetChangedFunc(func(i int, title string, description string, shortcut rune) {
		fn := con.invIndexes[i]
		con.update(fn)
	}).SetCurrentItem(1)

	return con
}

func (c *controller) updateListItem(filename string) {
	i := c.indexes[filename]
	view := c.logViews[filename]
	text := fmt.Sprintf("%v (%v)", filename, view.dirt())

	if filename != atogTag {
		c.fileList.SetItemText(i, text, "")
	} else {
		c.fileList.SetItemText(i, text, "Errors log")
	}
}

func (c *controller) updateView() {
	tv := c.textView
	idx := c.fileList.GetCurrentItem()
	fn := c.invIndexes[idx]

	tv.Clear()
	c.logViews[fn].flush(tv)
}

func (c *controller) update(source string) {
	c.updateView()
	c.updateListItem(source)
}

func (c *controller) handle(v stalker.Value) {
	var source string
	var buffer []byte

	if v.Error != nil {
		source = atogTag
		buffer = []byte(fmt.Sprintf("@{{ %v }}\n\t%v\n", source, v.Error))
	} else {
		source = v.Source
		buffer = v.Buffer
	}

	view := c.logViews[source]

	view.write(buffer)
	c.update(source)
}
