package manager

import (
	"fmt"
	"time"
)

type Manageable interface {
	Size() (size int)
	Resize(newSize int) (err error)
}

type ManagerCheck func() (newSize int)

type Manager struct {
	Name     string
	resource Manageable
	check    ManagerCheck
	freq     time.Duration
	running  bool
}

func New(resource Manageable, check ManagerCheck, freq time.Duration) *Manager {
	m := &Manager{
		Name:     fmt.Sprintf("Manager[%T %p]", resource, resource),
		resource: resource,
		check:    check,
		freq:     freq,
	}
	return m
}

func (m *Manager) manage(messages chan string, errors chan error) {
	for m.running {
		select {
		case <-time.After(m.freq):
			size := m.resource.Size()
			newSize := m.check()
			if newSize != size {
				messages <- fmt.Sprintf("%s: will resize from %d to %d", m.Name, size, newSize)
				if err := m.resource.Resize(newSize); err != nil {
					errors <- fmt.Errorf("%s: could not resize: %v", m.Name, err)
				}
			}
		}
	}
}

func (m *Manager) Start(msgChanSize int, errChanSize int) (chan string, chan error) {
	m.running = true
	messages := make(chan string, msgChanSize)
	errors := make(chan error, errChanSize)
	go m.manage(messages, errors)
	return messages, errors
}

func (m *Manager) Stop() {
	m.running = false
}
