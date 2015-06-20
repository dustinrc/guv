package guv_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/dustinrc/guv"
)

type fakeJob struct{}

func (j fakeJob) Run() {
	select {
	case <-time.After(time.Microsecond):
		break
	}
}

type fakeJobLong struct{}

func (j fakeJobLong) Run() {
	select {
	case <-time.After(5 * time.Millisecond):
		break
	}
}

var goodSizes = []struct {
	in, out int
}{
	{1, 1},
	{50, 50},
	{1024, 1024},
	{0, runtime.NumCPU()},
	{32, 32},
	{1, 1},
}

var badSizes = []int{
	-1,
	-50,
	-1024,
}

func TestPoolGoodResize(t *testing.T) {
	p, _ := guv.NewPool(0)
	size := p.Size()
	if size != runtime.NumCPU() {
		t.Errorf("Incorrect pool size: expected %d, actual %d.", 1, size)
	}
	knownPrev := size

	for _, tt := range goodSizes {
		prev, err := p.Resize(tt.in)
		if prev != knownPrev {
			t.Errorf("Incorrect previous pool size: expected %d, actual %d.", knownPrev, prev)
		}
		if err != nil {
			t.Errorf("Unexpected error during pool resize: %v", err)
		}
		size := p.Size()
		if size != tt.out {
			t.Errorf("Incorrect pool size: expected %d, actual %d.", tt.out, size)
		}
		knownPrev = size
	}

	p.Wait()
}

func TestPoolBadResize(t *testing.T) {
	p, _ := guv.NewPool(0)

	for _, badSize := range badSizes {
		prev, err := p.Resize(badSize)
		size := p.Size()
		if prev != size {
			t.Errorf("Pool size changed on bad resize: was %d, now %d", prev, size)
		}
		if err == nil {
			t.Errorf("Expected error for bad pool size of %d", badSize)
		}
	}
}

func TestGoodNewPool(t *testing.T) {
	p, err := guv.NewPool(1)
	size := p.Size()
	if size != 1 {
		t.Errorf("Incorrect initial pool size: expected 1, actual %s", size)
	}
	if err != nil {
		t.Errorf("Unexpected error for good initial pool size: %v", err)
	}
}

func TestBadNewPool(t *testing.T) {
	p, err := guv.NewPool(-1)
	size := p.Size()
	if size != 0 {
		t.Error("Incorrect initial pool size for bad input: expect 0, actual %d", size)
	}
	if err == nil {
		t.Error("No error for bad initial pool size")
	}
}
