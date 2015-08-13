package manager_test

import (
	"errors"
	"testing"
	"time"

	"github.com/dustinrc/guv/manager"
)

type FakeResource struct{}

func (fr FakeResource) Resize(int) error { return errors.New("fake error") }
func (fr FakeResource) Size() int        { return 0 }

func FakeCheck() int { return int(time.Now().Unix()) }

func TestManager(t *testing.T) {
	m := manager.New(FakeResource{}, FakeCheck, time.Minute)
	if m == nil {
		t.Error("Manager did not initialize correctly")
	}
}

func TestMessageChannel(t *testing.T) {
	m := manager.New(FakeResource{}, FakeCheck, time.Microsecond)
	messages, _ := m.Start(0, 0)
	select {
	case <-messages:
	case <-time.After(time.Millisecond):
		t.Error("Did not receive expected message")
	}
}

func TestErrorChannel(t *testing.T) {
	m := manager.New(FakeResource{}, FakeCheck, time.Microsecond)
	_, errors := m.Start(1, 0)
	select {
	case <-errors:
	case <-time.After(time.Millisecond):
		t.Error("Did not receive expected error")
	}
}

func TestStop(t *testing.T) {
	m := manager.New(FakeResource{}, FakeCheck, time.Microsecond)
	messages, _ := m.Start(0, 0)
	m.Stop()
	select {
	case <-messages:
		t.Error("Manager did not stop")
	case <-time.After(time.Millisecond):
	}
}
