package buffer

import (
	"io"
	"math"
)

// Buffer represents a buffer.
type Buffer struct {
	rrs   []readerRange
	index int64
}

// ReadSeekCloser is the interface that groups the basic Read, Seek and Close methods.
type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}

type readerRange struct {
	r    ReadSeekCloser
	min  int64
	max  int64
	diff int64
}

// NewBuffer creates a new buffer.
func NewBuffer(r ReadSeekCloser) *Buffer {
	return &Buffer{
		rrs:   []readerRange{{r: r, min: 0, max: math.MaxInt64, diff: 0}},
		index: 0,
	}
}

// Read reads bytes.
func (b *Buffer) Read(p []byte) (int, error) {
	if _, err := b.rrs[0].r.Seek(b.index, io.SeekStart); err != nil {
		return 0, err
	}
	return b.rrs[0].r.Read(p)
}

// Seek sets the offset.
func (b *Buffer) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		b.index = offset
	case io.SeekCurrent:
		b.index += offset
	case io.SeekEnd:
		if l, err := b.Len(); err != nil {
			return 0, err
		} else {
			b.index = l + offset
		}
	}
	return b.index, nil
}

// Close the buffer.
func (b *Buffer) Close() (err error) {
	for _, rr := range b.rrs {
		if e := rr.r.Close(); e != nil {
			err = e
		}
	}
	return
}

// Len returns the total size of the buffer.
func (b *Buffer) Len() (int64, error) {
	rr := b.rrs[len(b.rrs)-1]
	l, err := rr.r.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	return l - rr.diff, nil
}