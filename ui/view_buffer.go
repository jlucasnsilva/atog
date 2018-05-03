package ui

import "strings"

type viewBuffer struct {
	maxSize uint
	buff    []string
}

func newViewBuffer(size uint) *viewBuffer {
	return &viewBuffer{
		maxSize: size,
		buff:    make([]string, 0, size),
	}
}

func (vb *viewBuffer) Add(s string) {
	l := uint(len(vb.buff))
	vb.buff = append(vb.buff, s)

	if l >= vb.maxSize {
		vb.buff = vb.buff[1:]
	}
}

func (vb *viewBuffer) String(sep string) string {
	return strings.Join(vb.buff, sep)
}
