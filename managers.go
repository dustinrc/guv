package guv

import (
	"fmt"
	"log"
	"time"
)

type Manageable interface {
	Size() (size int)
	Resize(newSize int) (previous int, err error)
}

type ManagerCheck func() (newSize int)

type Manager struct {
	Name     string
	resource Manageable
	check    ManagerCheck
	freq     time.Duration
	running  bool
}

func NewManager(resource Manageable, check ManagerCheck, freq time.Duration) *Manager {
	m := &Manager{
		Name:     fmt.Sprintf("Manager %v", resource),
		resource: resource,
		check:    check,
		freq:     freq,
	}
	return m
}

func (m *Manager) Start() {
	m.running = true
	for m.running {
		select {
		case <-time.After(m.freq):
			size := m.resource.Size()
			newSize := m.check()
			if newSize != size {
				log.Printf("%s: will resize from %d to %d", m.Name, size, newSize)
				_, err := m.resource.Resize(m.check())
				if err != nil {
					log.Printf("Manager Error: could not resize: %v", err)
				}
			}
		}
	}
}

func (m *Manager) Stop() {
	m.running = false
}
