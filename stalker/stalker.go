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

package stalker

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/rjeczalik/notify"
)

var (
	// ChannelSize - filesystem events come through a buffered channel
	// this sets the size of the buffer of such channels
	ChannelSize int64 = 20
	// MaxBufferSize - at most, the last BufferSize bytes are read from a
	// file whenever it needs to be read.
	MaxBufferSize int64 = 20480

	waitEvent = make(chan string)
)

// Params ...
type Params struct {
	Target io.Writer
	Err    io.Writer
	Empty  bool
}

// Watch ...
func Watch(filename string, params Params) {
	go func(fname string, args Params) {
		file, err := os.Open(fname)
		if err != nil {
			args.Err.Write(makeError(fname, err))
		}
		defer file.Close()

		notifier := make(chan notify.EventInfo, ChannelSize)
		if err := notify.Watch(fname, notifier, notify.Write); err != nil {
			args.Err.Write(makeError(fname, err))
		}
		defer notify.Stop(notifier)

		run(file, notifier, &args)
	}(filename, params)
}

// WaitEvent ...
func WaitEvent() string {
	return <-waitEvent
}

func run(file *os.File, notifier chan notify.EventInfo, args *Params) {
	fp, err := initialPos(file, args.Empty)
	if err != nil {
		args.Err.Write(makeError(file.Name(), err))
	}

	buffer, count, err := readAt(file, fp)
	fp += count
	if err == nil {
		buffer = dropFirstLine(buffer)
		args.Target.Write(buffer)
	} else {
		args.Err.Write(makeError(file.Name(), err))
	}

	for range notifier {
		buffer, count, err := readAt(file, fp)
		fp += count

		if err == nil {
			args.Target.Write(buffer)
		} else {
			args.Err.Write(makeError(file.Name(), err))
		}

		waitEvent <- file.Name()
	}
}

func initialPos(file *os.File, empty bool) (int64, error) {
	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}

	if empty {
		return stat.Size(), nil
	}

	if fileSize := stat.Size(); fileSize > MaxBufferSize {
		return fileSize - MaxBufferSize, nil
	}

	return 0, nil
}

func readAt(file *os.File, at int64) ([]byte, int64, error) {
	buffer := make([]byte, MaxBufferSize+1)
	count, err := file.ReadAt(buffer, at)

	if err != nil && err != io.EOF {
		return nil, int64(count), err
	}

	return buffer[:count], int64(count), nil
}

func dropFirstLine(bs []byte) []byte {
	idx := bytes.Index(bs, []byte("\n"))
	if idx >= 0 {
		return bs[idx+1:]
	}
	return bs
}

func makeError(filename string, err error) []byte {
	return []byte(fmt.Sprintf("{{ %v }}\n%v", filename, err.Error()))
}
