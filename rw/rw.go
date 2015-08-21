package rw

import (
	"fmt"
	"io"
	"sync"
)

type Reader struct {
	r io.Reader

	mu   sync.Mutex
	rate int
}

func NewReader(r io.Reader, rate int) (*Reader, error) {
	gr := &Reader{
		r: r,
	}
	err := gr.Resize(rate)
	return gr, err
}

func (gr *Reader) Resize(rate int) error {
	gr.mu.Lock()
	defer gr.mu.Unlock()

	if rate < 1 {
		return fmt.Errorf("bad rate: %v", rate)
	} else {
		gr.rate = rate
	}
	return nil
}

func (gr *Reader) Size() int {
	gr.mu.Lock()
	defer gr.mu.Unlock()
	return gr.rate
}

func (gr *Reader) Read(p []byte) (n int, err error) {
	return 0, nil
}

type Writer struct {
	w io.Writer

	mu   sync.Mutex
	rate int
}

func NewWriter(w io.Writer, rate int) (*Writer, error) {
	gw := &Writer{
		w: w,
	}
	err := gw.Resize(rate)
	return gw, err
}

func (gw *Writer) Resize(rate int) error {
	gw.mu.Lock()
	defer gw.mu.Unlock()

	if rate < 1 {
		return fmt.Errorf("bad rate: %v", rate)
	} else {
		gw.rate = rate
	}
	return nil
}

func (gw *Writer) Size() int {
	gw.mu.Lock()
	defer gw.mu.Unlock()
	return gw.rate
}

func (gw *Writer) Write(p []byte) (n int, err error) {
	return 0, nil
}
