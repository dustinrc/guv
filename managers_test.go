package guv_test

import (
	"testing"
	"time"

	"github.com/dustinrc/guv"
)

type FakeResource struct{}

func (fr FakeResource) Resize(int) (int, error) { return 0, nil }
func (fr FakeResource) Size() int               { return 0 }

func FakeCheck() int { return 0 }

func TestManager(t *testing.T) {
	m := guv.NewManager(FakeResource{}, FakeCheck, time.Minute)
	if m == nil {
		t.Error("Manager did not initialize correctly")
	}
}
