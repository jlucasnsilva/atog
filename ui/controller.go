/*
 * Copyright (c) 2018, João Lucas Nunes e Silva
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *     * Neither the name of the <organization> nor the
 *       names of its contributors may be used to endorse or promote products
 *       derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL JOÃO LUCAS NUNES E SILVA BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

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
