package guv_test

import (
	"testing"
	"time"

	"github.com/dustinrc/guv"
)

func FakeCheck() int { return 0 }

func TestManager(t *testing.T) {
	m, _ := guv.NewManager(FakeCheck, time.Minute)
	if m == nil {
		t.Error("Manager did not initialize correctly")
	}
}
