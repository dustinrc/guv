package manager_test

import (
	"testing"
	"time"

	"github.com/dustinrc/guv/manager"
)

type FakeResource struct{}

func (fr FakeResource) Resize(int) error { return nil }
func (fr FakeResource) Size() int        { return 0 }

func FakeCheck() int { return 0 }

func TestManager(t *testing.T) {
	m := manager.New(FakeResource{}, FakeCheck, time.Minute)
	if m == nil {
		t.Error("Manager did not initialize correctly")
	}
}
