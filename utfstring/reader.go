package utfstrings

import (
	"errors"
	"io"
	"unicode/utf8"
)

const forceBufferSize = 3

// RuneReader provides functions for buffered rune-reading.
type RuneReader struct {
	r   io.Reader
	buf []byte
	pos int
	end int
}

func (rr *RuneReader) nextRune() (rune, error) {
	if rr.size() < forceBufferSize {
		// try to ensure that the buffer at-least has 'forceBufferSize'
		// number of bytes.
		err := rr.buffer()
		if err != nil && (err != io.EOF || rr.size() == 0) {
			return -1, err
		}
	}

	r, size := utf8.DecodeRune(rr.buf[rr.pos:])
	if r == utf8.RuneError && size == 0 {
		return -1, errors.New("no data")
	} else if r == utf8.RuneError && size == 1 {
		return -1, errors.New("invalid character")
	}
	rr.pos += size

	return r, nil
}

func (rr *RuneReader) size() int {
	return rr.end - rr.pos
}

func (rr *RuneReader) buffer() error {
	if rr.end > 0 && rr.pos > rr.end {
		// extend buffer by doubling the size
		rr.buf = append(rr.buf, make([]byte, len(rr.buf))...)
	} else if rr.end > 0 && rr.pos == rr.end {
		// all data in buffer is consumed, reset buffer
		rr.end = 0
		rr.pos = 0
	}

	n, err := rr.r.Read(rr.buf[rr.end:])
	if err != nil {
		return err
	}
	rr.end += n
	return nil
}
