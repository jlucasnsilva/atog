package ui

import (
	"fmt"
	"sync"

	"github.com/rivo/tview"
)

type logView struct {
	lock  *sync.Mutex
	text  string
	dirtv uint
}

func newLogView() *logView {
	return &logView{lock: &sync.Mutex{}}
}

func (st *logView) write(buff []byte) {
	if len(buff) > 0 {
		st.lock.Lock()
		defer st.lock.Unlock()

		st.text = format(st.text, buff)
		st.dirtv++
	}
}

func (st *logView) flush(tv *tview.TextView) {
	st.lock.Lock()
	defer st.lock.Unlock()

	fmt.Fprint(tv, st.text)
	tv.ScrollToEnd()
	st.dirtv = 0
}

func (st *logView) dirt() uint {
	st.lock.Lock()
	defer st.lock.Unlock()

	return st.dirtv
}
