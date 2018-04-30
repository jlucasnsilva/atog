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
	"io"
	"os"
	"sync"

	"github.com/rjeczalik/notify"
)

var (
	// ChannelSize - filesystem events come through a buffered channel
	// this sets the size of the buffer of such channels
	ChannelSize int64 = 20
	// MaxBufferSize - at most, the last BufferSize bytes are read from a
	// file whenever it needs to be read.
	MaxBufferSize int64 = 20480
)

type (
	stalker struct {
		filename string
		filepos  int64
		notifier chan notify.EventInfo
	}

	// Value ...
	Value struct {
		Source string
		Buffer []byte
		Error  error
	}
)

// Watch creates a new Stalker to watch over a file.
func Watch(filenames []string) <-chan Value {
	out := make(chan Value)
	wg := &sync.WaitGroup{}

	wg.Add(len(filenames))

	for _, fn := range filenames {
		s := &stalker{
			filename: fn,
			notifier: make(chan notify.EventInfo, ChannelSize),
		}

		go start(s, out, wg)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func start(s *stalker, writeTo chan<- Value, wg *sync.WaitGroup) {
	if err := notify.Watch(s.filename, s.notifier, notify.Write); err != nil {
		writeTo <- Value{Source: s.filename, Error: err}
		return
	}

	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		writeTo <- Value{Source: s.filename, Error: err}
		return
	}

	pos, err := initialPos(s.filename)
	buffer, count, err := readAt(s.filename, pos)

	buffer = dropFirstLine(buffer)
	s.filepos = pos + count

	writeTo <- Value{
		Source: s.filename,
		Buffer: buffer,
		Error:  err,
	}

	for range s.notifier {
		buffer, count, err := readAt(s.filename, s.filepos)
		s.filepos += count

		writeTo <- Value{
			Source: s.filename,
			Buffer: buffer,
			Error:  err,
		}
	}

	wg.Done()
}

func initialPos(filename string) (int64, error) {
	stat, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}

	if fileSize := stat.Size(); fileSize > MaxBufferSize {
		return fileSize - MaxBufferSize, nil
	}

	return 0, nil
}

func readAt(filename string, at int64) ([]byte, int64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

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

// func stop(s *stalker) {
// 	if s != nil && s.notifier != nil {
// 		notify.Stop(s.notifier)
// 	}
// }
