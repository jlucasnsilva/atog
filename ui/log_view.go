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
	"strings"
	"sync"

	"github.com/rivo/tview"
)

type logView struct {
	lock  *sync.Mutex
	buff  *viewBuffer
	dirtv int
}

func newLogView(size uint) *logView {
	return &logView{
		lock: &sync.Mutex{},
		buff: newViewBuffer(size),
	}
}

// Write ...
func (v *logView) Write(p []byte) (n int, err error) {
	pc := len(p)

	if pc < 1 {
		return 0, nil
	}

	v.lock.Lock()
	defer v.lock.Unlock()

	ls := strings.Split(string(p), "\n")
	lines := make([]string, 0, len(ls))
	for _, line := range ls {
		if line != "\n" {
			lines = append(lines, line)
		}
	}

	for _, line := range lines {
		v.buff.Add(highlight(line))
	}

	v.dirtv += len(lines)
	return pc, nil
}

func (v *logView) flush(tv *tview.TextView) {
	v.lock.Lock()
	defer v.lock.Unlock()

	fmt.Fprint(tv, v.buff.String("\n"))
	tv.ScrollToEnd()
	v.dirtv = 0
}

func (v *logView) dirt() int {
	v.lock.Lock()
	defer v.lock.Unlock()

	return v.dirtv
}
