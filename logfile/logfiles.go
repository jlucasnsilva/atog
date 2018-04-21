package logfiles

import (
	"os"
	"strings"
)

// LastLines returns the nOfLines last lines of the given file.
func LastLines(file *os.File, nOfLines int) ([]string, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	maxBufferSize := int64(20480) // 20 kbytes
	fileSize := stat.Size()
	bufferSize := int64(fileSize)

	if fileSize > maxBufferSize {
		bufferSize = maxBufferSize
	}

	buffer := make([]byte, bufferSize)
	_, err = file.ReadAt(buffer, fileSize-bufferSize)
	if err != nil {
		return nil, err
	}

	return splitBufferIntoLines(buffer, nOfLines), nil
}

func splitBufferIntoLines(b []byte, nOfLines int) []string {
	i := len(b) - 1
	newline := byte('\n')
	count := nOfLines

	if b[i] == newline {
		// if the last bytes is a newline character
		// it need to count nOfLines + 1 lines.
		count++
	}

	for ; i >= 0 && count > 0; i-- {
		if b[i] == newline {
			count--
		}
	}

	return strings.Split(string(b[i:]), "\n")
}
