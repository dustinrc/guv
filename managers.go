package guv

import "time"

type ManagerCheckFunc func() int

type Manager struct {
	pool  *Pool
	check ManagerCheckFunc
	freq  time.Duration
}

func NewManager(check ManagerCheckFunc, frequency time.Duration) (*Manager, error) {
	pool, err := NewPool(0)
	if err != nil {
		return nil, err
	}
	m := &Manager{
		pool:  pool,
		check: check,
		freq:  frequency,
	}
	return m, nil
}
