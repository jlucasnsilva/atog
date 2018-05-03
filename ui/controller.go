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

func newController(args Params, textView *tview.TextView, fileList *tview.List) *controller {
	fns := append([]string{atogTag}, args.Filenames...)
	con := &controller{}

	con.textView = textView
	con.fileList = fileList
	con.indexes = make(map[string]int)
	con.logViews = make(map[string]*logView)
	con.invIndexes = make([]string, len(fns))

	for i, fn := range fns {
		if fn != atogTag {
			con.fileList.AddItem(fn, "\t(0)", shortcut(i), nil)
		} else {
			con.fileList.AddItem(fn, "Errors log (0)", shortcut(i), nil)
		}

		con.logViews[fn] = newLogView(args.BufferSize)
		con.indexes[fn] = i
		con.invIndexes[i] = fn
	}

	fileList.SetChangedFunc(func(i int, title string, description string, shortcut rune) {
		fn := con.invIndexes[i]
		con.update(fn)
	}).SetCurrentItem(1)

	atogv := con.logViews[atogTag]
	for _, fn := range args.Filenames {
		stalker.Watch(fn, stalker.Params{
			Empty:  args.Empty,
			Err:    atogv,
			Target: con.logViews[fn],
		})
	}

	return con
}

func (c *controller) updateListItem(filename string) {
	i := c.indexes[filename]
	view := c.logViews[filename]

	if filename != atogTag {
		c.fileList.SetItemText(i, filename, fmt.Sprintf("\t(%v)", view.dirt()))
	} else {
		c.fileList.SetItemText(i, filename, fmt.Sprintf("Errors log (%v)", view.dirt()))
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

func (c *controller) handle(app *tview.Application) {
	app.Draw()

	for {
		source := stalker.WaitEvent()
		c.update(source)
		app.Draw()
	}
}
