package linewriter

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testWriter struct {
	w func([]byte) (int, error)
}

func (tw testWriter) Write(b []byte) (int, error) {
	return tw.w(b)
}

func TestNew(t *testing.T) {
	b := bytes.NewBuffer(nil)

	// test that returns *LineWriter
	assert.NotNil(t, New(false, b))

	// test that returns *LineWriter
	assert.NotNil(t, New(true, b))
}

func TestLineWriter_Len(t *testing.T) {
	b := bytes.NewBuffer(nil)
	w := New(false, b)

	// test that returns 0
	assert.Equal(t, 0, w.Len())

	// test that returns 3
	w.buf.WriteString("foo")
	assert.Equal(t, 3, w.Len())
}

func TestLineWriter_FlushAll(t *testing.T) {
	b := bytes.NewBuffer(nil)
	ncall := 0
	tw := &testWriter{
		w: func(p []byte) (int, error) {
			ncall++
			return b.Write(p)
		},
	}
	w := New(false, tw)

	// test that no flush the data if there is no buffered data
	n, err := w.FlushAll()
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.Equal(t, 0, ncall)

	// test that flush the buffered data
	w.buf.WriteString("foo")
	n, err = w.FlushAll()
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, 1, ncall)
	// confirm
	assert.EqualValues(t, "foo", b.Bytes())

	// test that returns the return values of writer
	tw.w = func([]byte) (int, error) {
		return 0, fmt.Errorf("error of writer")
	}
	w.buf.WriteString("bar")
	n, err = w.FlushAll()
	assert.Equal(t, fmt.Errorf("error of writer"), err)
	assert.Equal(t, 0, n)
	assert.Equal(t, 3, w.buf.Len())

	// test that returns the return values of writer and consumes write bytes
	tw.w = func([]byte) (int, error) {
		return 1000, fmt.Errorf("error of writer")
	}
	n, err = w.FlushAll()
	assert.Equal(t, fmt.Errorf("error of writer"), err)
	assert.Equal(t, 3, n)
	assert.Equal(t, 0, w.buf.Len())

}

func TestLineWriter_Flush(t *testing.T) {
	b := bytes.NewBuffer(nil)
	ncall := 0
	tw := &testWriter{
		w: func(p []byte) (int, error) {
			ncall++
			return b.Write(p)
		},
	}
	w := New(false, tw)

	// test that no flush the data if there is line delimiters
	w.buf.WriteString("foo")
	n, err := w.Flush()
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.Equal(t, 0, ncall)

	// test that flush the each lines
	w.buf.WriteString("bar\nbaz\n")
	c := w.buf.Len()
	n, err = w.Flush()
	assert.NoError(t, err)
	assert.Equal(t, c, n)
	assert.Equal(t, 2, ncall)
	// confirm
	assert.Equal(t, "foobar\nbaz\n", b.String())

	// test that returns the return values of writer
	tw.w = func([]byte) (int, error) {
		return 1000, fmt.Errorf("error of writer")
	}
	w.buf.WriteString("qux\n")
	n, err = w.Flush()
	assert.Equal(t, fmt.Errorf("error of writer"), err)
	assert.Equal(t, 4, n)
	assert.Equal(t, 0, w.buf.Len())

	// test for multiline
	b.Reset()
	tw.w = func(p []byte) (int, error) {
		ncall++
		return b.Write(p)
	}
	ncall = 0
	w = New(true, tw)

	// test that no flush the data if there is line delimiters
	w.buf.WriteString("foo")
	n, err = w.Flush()
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.Equal(t, 0, ncall)

	// test that flush the multiline
	w.buf.WriteString("bar\nbaz\n")
	c = w.buf.Len()
	n, err = w.Flush()
	assert.NoError(t, err)
	assert.Equal(t, c, n)
	assert.Equal(t, 1, ncall)
	// confirm
	assert.Equal(t, "foobar\nbaz\n", b.String())

	// test that returns the return values of writer
	tw.w = func([]byte) (int, error) {
		return 0, fmt.Errorf("error of writer")
	}
	w.buf.WriteString("qux\nquux\n")
	n, err = w.Flush()
	assert.Equal(t, fmt.Errorf("error of writer"), err)
	assert.Equal(t, 0, n)
	assert.Equal(t, 9, w.buf.Len())

	// test that returns the return values of writer and consumes write bytes
	tw.w = func([]byte) (int, error) {
		return 1000, fmt.Errorf("error of writer")
	}
	n, err = w.Flush()
	assert.Equal(t, fmt.Errorf("error of writer"), err)
	assert.Equal(t, 9, n)
	assert.Equal(t, 0, w.buf.Len())
}

func TestLineWriter_Write(t *testing.T) {
	b := bytes.NewBuffer(nil)
	ncall := 0
	tw := &testWriter{
		w: func(p []byte) (int, error) {
			ncall++
			return b.Write(p)
		},
	}
	w := New(false, tw)

	// test that write bytes
	foo := []byte("foo")
	n, err := w.Write(foo)
	assert.NoError(t, err)
	assert.Equal(t, len(foo), n)
	assert.Equal(t, 0, ncall)

	// test that flush the buffered data if writes the line delimiter
	barbaz := []byte("bar\nbaz")
	n, err = w.Write(barbaz)
	assert.NoError(t, err)
	assert.Equal(t, len(barbaz), n)
	assert.Equal(t, 1, ncall)
	// confirm
	assert.EqualValues(t, "foobar\n", b.Bytes())

	// test that returns number of bytes written and errors of writer
	b.Reset()
	tw.w = func(p []byte) (int, error) {
		b.Write(p)
		return 1000, fmt.Errorf("error of writer")
	}
	qux := []byte("qux\n")
	n, err = w.Write(qux)
	assert.Equal(t, fmt.Errorf("error of writer"), err)
	assert.Equal(t, len(qux), n)
	assert.EqualValues(t, "bazqux\n", b.Bytes())
}

func TestLineWriter_WriteString(t *testing.T) {
	b := bytes.NewBuffer(nil)
	ncall := 0
	tw := &testWriter{
		w: func(p []byte) (int, error) {
			ncall++
			return b.Write(p)
		},
	}
	w := New(false, tw)

	// test that write bytes
	foo := "foo"
	n, err := w.WriteString(foo)
	assert.NoError(t, err)
	assert.Equal(t, len(foo), n)
	assert.Equal(t, 0, ncall)

	// test that flush the buffered data if writes the line delimiter
	barbaz := "bar\nbaz"
	n, err = w.WriteString(barbaz)
	assert.NoError(t, err)
	assert.Equal(t, len(barbaz), n)
	assert.Equal(t, 1, ncall)
	// confirm
	assert.EqualValues(t, "foobar\n", b.Bytes())

	// test that returns number of bytes written and errors of writer
	b.Reset()
	tw.w = func(p []byte) (int, error) {
		b.Write(p)
		return 1000, fmt.Errorf("error of writer")
	}
	qux := "qux\n"
	n, err = w.WriteString(qux)
	assert.Equal(t, fmt.Errorf("error of writer"), err)
	assert.Equal(t, len(qux), n)
	assert.Equal(t, "bazqux\n", b.String())
}

func TestLineWriter_Printf(t *testing.T) {
	b := bytes.NewBuffer(nil)
	ncall := 0
	tw := &testWriter{
		w: func(p []byte) (int, error) {
			ncall++
			return b.Write(p)
		},
	}
	w := New(false, tw)

	// test that write bytes
	foo := "foo"
	n, err := w.Printf("%s", foo)
	assert.NoError(t, err)
	assert.Equal(t, len(foo), n)
	assert.Equal(t, 0, ncall)

	// test that flush the buffered data if writes the line delimiter
	barbaz := "bar\nbaz"
	n, err = w.Printf("%s", barbaz)
	assert.NoError(t, err)
	assert.Equal(t, len(barbaz), n)
	assert.Equal(t, 1, ncall)
	// confirm
	assert.EqualValues(t, "foobar\n", b.Bytes())

	// test that returns number of bytes written and errors of writer
	b.Reset()
	tw.w = func(p []byte) (int, error) {
		b.Write(p)
		return 1000, fmt.Errorf("error of writer")
	}
	qux := "qux\n"
	n, err = w.Printf("%s", qux)
	assert.Equal(t, fmt.Errorf("error of writer"), err)
	assert.Equal(t, len(qux), n)
	assert.EqualValues(t, "bazqux\n", b.Bytes())
}
