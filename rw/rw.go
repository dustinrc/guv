package rw

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type Reader struct {
	r rwThrottleFunc

	mu   sync.Mutex
	rate int
}

func NewReader(r io.Reader, rate int) (*Reader, error) {
	gr := &Reader{
		r: throttleWrapper(r.Read),
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
	return gr.r(p, gr.Size())
}

type Writer struct {
	w rwThrottleFunc

	mu   sync.Mutex
	rate int
}

func NewWriter(w io.Writer, rate int) (*Writer, error) {
	gw := &Writer{
		w: throttleWrapper(w.Write),
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
	return gw.w(p, gw.rate)
}

type rwFunc func([]byte) (int, error)
type rwThrottleFunc func([]byte, int) (int, error)

func throttleWrapper(f rwFunc) rwThrottleFunc {
	var cooldown time.Duration
	var currRate float64

	return func(p []byte, rate int) (n int, err error) {
		start := time.Now()
		time.Sleep(cooldown)
		n, err = f(p)
		elapsed := time.Since(start)

		currRate = float64(n) / float64(elapsed) * float64(time.Second)
		ratio := currRate / float64(rate)
		if ratio > 25.0 {
			ratio = 25.0
		}
		if ratio > 1.0 {
			cooldown += time.Duration(ratio) * time.Millisecond
		} else {
			cooldown -= time.Duration((1.0-ratio)*1000) * time.Microsecond
		}
		if cooldown < 0 {
			cooldown = 0
		}

		return
	}
}
