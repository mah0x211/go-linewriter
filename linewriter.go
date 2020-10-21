package linewriter

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

type LineWriter struct {
	mu        *sync.RWMutex
	multiline bool
	buf       *bytes.Buffer
	out       io.Writer
}

func New(multiline bool, w io.Writer) *LineWriter {
	return &LineWriter{
		mu:        &sync.RWMutex{},
		multiline: multiline,
		buf:       bytes.NewBuffer(nil),
		out:       w,
	}
}

func (w *LineWriter) Len() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.buf.Len()
}

func (w *LineWriter) FlushAll() (int, error) {
	n := w.Len()
	if n == 0 {
		return 0, nil
	}

	b := w.buf.Bytes()
	n, err := w.out.Write(b)
	if n > 0 {
		m := len(b)
		if n > m {
			n = m
		}
		w.buf.Next(n)
		return n, err
	}

	return 0, err
}

func (w *LineWriter) flushLines() (int, error) {
	var tail int
	var err error
	b := w.buf.Bytes()
	n := bytes.IndexByte(b, '\n')
	for n != -1 {
		n++
		c := b[tail : tail+n]
		nw, werr := w.out.Write(c)
		if nw > 0 {
			m := len(c)
			if nw > m {
				nw = m
			}
			tail += nw
		}

		if werr != nil {
			err = werr
			break
		}
		n = bytes.IndexByte(b[tail:], '\n')
	}

	// remove the output bytes
	if tail > 0 {
		w.buf.Next(tail)
	}

	return tail, err
}

func (w *LineWriter) flushMultiline() (int, error) {
	var tail int
	b := w.buf.Bytes()
	n := bytes.IndexByte(b, '\n')
	for n != -1 {
		tail += n + 1
		n = bytes.IndexByte(b[tail:], '\n')
	}

	if tail == 0 {
		return 0, nil
	}

	c := b[:tail]
	nw, err := w.out.Write(c)
	if nw > 0 {
		m := len(c)
		if nw > m {
			nw = m
		}
		// remove the output bytes
		w.buf.Next(nw)
		return nw, err
	}

	return 0, err
}

func (w *LineWriter) Flush() (int, error) {
	if w.multiline {
		return w.flushMultiline()
	}
	return w.flushLines()
}

func (w *LineWriter) Write(b []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	n, _ := w.buf.Write(b)
	_, err := w.Flush()
	return n, err
}

func (w *LineWriter) WriteString(s string) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	n, _ := w.buf.WriteString(s)
	_, err := w.Flush()
	return n, err
}

func (w *LineWriter) Printf(s string, a ...interface{}) (int, error) {
	return w.WriteString(fmt.Sprintf(s, a...))
}
