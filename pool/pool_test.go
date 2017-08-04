package pool_test

import (
	"testing"
	"time"

	"github.com/dustinrc/guv/manager"
	"github.com/dustinrc/guv/pool"
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
	{32, 32},
	{1, 1},
}

var badSizes = []int{
	0,
	-1,
	-50,
	-1024,
}

func TestPoolGoodResize(t *testing.T) {
	p, _ := pool.New(1)
	size := p.Size()
	if size != 1 {
		t.Errorf("Incorrect pool size: expected %d, actual %d", 1, size)
	}

	for _, tt := range goodSizes {
		err := p.Resize(tt.in)
		if err != nil {
			t.Errorf("Unexpected error during pool resize: %v", err)
		}
		size := p.Size()
		if size != tt.out {
			t.Errorf("Incorrect pool size: expected %d, actual %d", tt.out, size)
		}
	}

	p.Wait()
}

func TestPoolBadResize(t *testing.T) {
	p, _ := pool.New(1)

	for _, badSize := range badSizes {
		err := p.Resize(badSize)
		if err == nil {
			t.Errorf("Expected error for bad pool size of %d", badSize)
		}
	}
}

func TestGoodNewPool(t *testing.T) {
	p, err := pool.New(1)
	size := p.Size()
	if size != 1 {
		t.Errorf("Incorrect initial pool size: expected 1, actual %d", size)
	}
	if err != nil {
		t.Errorf("Unexpected error for good initial pool size: %v", err)
	}
}

func TestBadNewPool(t *testing.T) {
	p, err := pool.New(-1)
	size := p.Size()
	if size != 0 {
		t.Errorf("Incorrect initial pool size for bad input: expect 0, actual %d", size)
	}
	if err == nil {
		t.Error("No error for bad initial pool size")
	}
}

func TestPoolImplementsManageableInterface(t *testing.T) {
	var _ manager.Manageable = (*pool.Pool)(nil)
}
