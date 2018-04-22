package stalker

import (
	"os"
	"strings"

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

var (
	stalkers = make(map[string]*stalker)
)

type (
	stalker struct {
		filename string
		filepos  int64
		notifier chan notify.EventInfo
		out      chan Value
	}

	// Value ...
	Value struct {
		Lines []string
		Error error
	}
)

// Watch creates a new Stalker to watch over a file.
func Watch(filename string) (<-chan Value, error) {
	s := &stalker{
		filename: filename,
		notifier: make(chan notify.EventInfo, ChannelSize),
		out:      make(chan Value),
	}

	if err := s.start(); err != nil {
		return nil, err
	}

	stalkers[filename] = s
	return s.out, nil
}

// Unwatch ...
func Unwatch(filename string) {
	s, ok := stalkers[filename]

	if ok {
		cleanUp(s)
		delete(stalkers, filename)
	}
}

func (s *stalker) start() error {
	if err := notify.Watch(s.filename, s.notifier, notify.Write); err != nil {
		return err
	}

	ip, err := initialPos(s.filename)
	if err != nil {
		cleanUp(s)
		return err
	}

	lines, _, err := readLines(s.filename, ip)
	if err != nil {
		cleanUp(s)
		return err
	}

	s.filepos = ip
	s.out <- Value{Lines: lines}

	go func() {
		for range s.notifier {
			v := Value{}
			ls, count, err := readLines(s.filename, s.filepos)
			if err != nil {
				v.Error = err
			} else {
				v.Lines = ls
				s.filepos += int64(count)
			}
			s.out <- v
		}
	}()

	return nil
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

func readLines(filename string, pos int64) ([]string, int, error) {
	b, nOfBytes, err := readAt(filename, pos)
	if err != nil {
		return nil, 0, err
	}

	return splitBufferIntoLines(b), nOfBytes, nil
}

func readAt(filename string, at int64) ([]byte, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	buffer := make([]byte, MaxBufferSize)
	count, err := file.ReadAt(buffer, at)
	if err != nil {
		return nil, 0, err
	}

	return buffer, count, nil
}

func splitBufferIntoLines(b []byte) []string {
	lines := strings.Split(string(b), "\n")
	// the first line is dropped in case we read a
	// incomplete line
	return lines[1:]
}

func cleanUp(s *stalker) {
	close(s.out)
	notify.Stop(s.notifier)
}
