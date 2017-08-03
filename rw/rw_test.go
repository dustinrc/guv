package rw_test

import (
	"io"
	"testing"

	"github.com/dustinrc/guv/manager"
	"github.com/dustinrc/guv/rw"
)

type FakeRW struct{}

func (frw FakeRW) Write([]byte) (int, error) { return 0, nil }
func (frw FakeRW) Read([]byte) (int, error)  { return 0, nil }

var goodRates = []struct {
	in, out int
}{
	{1, 1},
	{50, 50},
	{1024, 1024},
	{32, 32},
	{1, 1},
}

var badRates = []int{
	0,
	-1,
	-50,
	-1024,
}

func TestReaderGoodResize(t *testing.T) {
	gr, _ := rw.NewReader(FakeRW{}, 1)

	for _, tt := range goodRates {
		err := gr.Resize(tt.in)
		if err != nil {
			t.Errorf("Unexpected error during reader resize: %v", err)
		}
		size := gr.Size()
		if size != tt.out {
			t.Errorf("Incorrect reader size: expected %d, actual %d", tt.out, size)
		}
	}
}

func TestReaderBadResize(t *testing.T) {
	gr, _ := rw.NewReader(FakeRW{}, 1)

	for _, badSize := range badRates {
		err := gr.Resize(badSize)
		if err == nil {
			t.Errorf("Expected error for bad reader size of %d", badSize)
		}
	}
}

func TestGoodNewReader(t *testing.T) {
	gr, err := rw.NewReader(FakeRW{}, 1)
	size := gr.Size()
	if size != 1 {
		t.Errorf("Incorrect initial reader size: expected 1, actual %d", size)
	}
	if err != nil {
		t.Errorf("Unexpected error for good initial reader size: %v", err)
	}
}

func TestBadNewReader(t *testing.T) {
	gr, err := rw.NewReader(FakeRW{}, -1)
	size := gr.Size()
	if size != 0 {
		t.Errorf("Incorrect initial reader size for bad input: expect 0, actual %d", size)
	}
	if err == nil {
		t.Error("No error for bad initial reader size")
	}
}

func TestReaderImplementsReaderInterface(t *testing.T) {
	var _ io.Reader = (*rw.Reader)(nil)
}

func TestReaderImplementsManageableInterface(t *testing.T) {
	var _ manager.Manageable = (*rw.Reader)(nil)
}

func TestWriterGoodResize(t *testing.T) {
	gw, _ := rw.NewWriter(FakeRW{}, 1)

	for _, tt := range goodRates {
		err := gw.Resize(tt.in)
		if err != nil {
			t.Errorf("Unexpected error during writer resize: %v", err)
		}
		size := gw.Size()
		if size != tt.out {
			t.Errorf("Incorrect writer size: expected %d, actual %d", tt.out, size)
		}
	}
}

func TestWriterBadResize(t *testing.T) {
	gw, _ := rw.NewWriter(FakeRW{}, 1)

	for _, badSize := range badRates {
		err := gw.Resize(badSize)
		if err == nil {
			t.Errorf("Expected error for bad writer size of %d", badSize)
		}
	}
}

func TestGoodNewWriter(t *testing.T) {
	gw, err := rw.NewWriter(FakeRW{}, 1)
	size := gw.Size()
	if size != 1 {
		t.Errorf("Incorrect initial writer size: expected 1, actual %d", size)
	}
	if err != nil {
		t.Errorf("Unexpected error for good initial writer size: %v", err)
	}
}

func TestBadNewWriter(t *testing.T) {
	gw, err := rw.NewWriter(FakeRW{}, -1)
	size := gw.Size()
	if size != 0 {
		t.Errorf("Incorrect initial writer size for bad input: expect 0, actual %d", size)
	}
	if err == nil {
		t.Error("No error for bad initial writer size")
	}
}

func TestWriterImplementsWriterInterface(t *testing.T) {
	var _ io.Writer = (*rw.Writer)(nil)
}

func TestWriterImplementsManageableInterface(t *testing.T) {
	var _ manager.Manageable = (*rw.Writer)(nil)
}
