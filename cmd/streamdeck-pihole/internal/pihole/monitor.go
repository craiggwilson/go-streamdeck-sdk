package pihole

import (
	"sync"
	"time"
)

func newMonitor(ph *PiHole, refreshInterval time.Duration) *Monitor {
	return &Monitor{
		ph:        ph,
		refreshInterval: refreshInterval,
		immediate: make(chan struct{}),
		done:      make(chan struct{}),
	}
}

// Monitor handles watching a Pi-Hole for status changes on a specified interval.
type Monitor struct {
	ph *PiHole
	refreshInterval time.Duration

	current StatusUpdate

	mu        sync.Mutex
	subs      []*subscription
	immediate chan struct{}
	done      chan struct{}
}

// Disable disables the Pi-Hole.
func (m *Monitor) Disable(durationSeconds int) {
	err := m.ph.Disable(durationSeconds)
	m.push(Disabled, time.Now().Add(time.Duration(durationSeconds) * time.Second), err)
}

// Enable enables the Pi-Hole.
func (m *Monitor) Enable() {
	err := m.ph.Enable()
	m.push(Enabled, time.Time{}, err)
}

// RefreshInterval returns the refresh interval used.
func (m *Monitor) RefreshInterval() time.Duration {
	return m.refreshInterval
}

// RefreshIn tells the monitor to refresh itself after the specified duration. It will also refresh itself immediately.
// If duration is 0, the monitor will only be refreshed immediately.
func (m *Monitor) RefreshIn(duration time.Duration) {
	if duration == 0 {
		select {
		case m.immediate <- struct{}{}:
		default:
		}
		return
	}

	go func() {
		<-time.After(duration)
		select {
		case m.immediate <- struct{}{}:
		default:
		}
	}()
}

// Status updates the status immediately and returns the result.
func (m *Monitor) Status() StatusUpdate {
	return m.check()
}

// Stop shuts the monitor down.
func (m *Monitor) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Copy this out because calling unsub messes with Monitor.unsubs.
	subs := make([]*subscription, len(m.subs))
	copy(subs, m.subs)

	for _, sub := range m.subs {
		sub.unsub()
	}

	close(m.done)
	close(m.immediate)
}

// Subscribe returns a channel to watch for updates as well as a function used to unsubscribe.
func (m *Monitor) Subscribe() (<-chan StatusUpdate, func()) {
	m.mu.Lock()
	defer m.mu.Unlock()

	sub := subscription{
		ch: make(chan StatusUpdate),
	}

	m.subs = append(m.subs, &sub)
	return sub.ch, sub.unsub
}

func (m *Monitor) check() StatusUpdate {
	newStatus, err := m.ph.Status()
	if err != nil {
		return m.push(Unknown, time.Time{}, err)
	}

	m.mu.Lock()
	currentStatus := m.current
	m.mu.Unlock()

	if newStatus != currentStatus.Status {
		return m.push(newStatus, currentStatus.DisabledUntil, nil)
	}

	return currentStatus
}

// push handles pushing new status updates to subscriptions and also removes closed
// subscriptions in cooperation with the unsubscribe functions.
func (m *Monitor) push(status Status, disabledUntil time.Time, err error) StatusUpdate {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.current = StatusUpdate {
		Status: status,
		DisabledUntil: disabledUntil,
		Err: err,
	}

	var newSubs []*subscription
	for i := 0; i < len(m.subs); i++ {
		sub := m.subs[i]
		closed := sub.pushIfNotClosed(m.current)
		if sub.closed && newSubs == nil {
			newSubs = make([]*subscription, i, len(m.subs)-1)
			copy(newSubs, m.subs[:i])
		}

		if !closed && newSubs != nil {
			newSubs = append(newSubs, sub)
		}
	}

	return m.current
}

func (m *Monitor) start() {
	go func() {
		ticker := time.NewTicker(m.refreshInterval)
		defer ticker.Stop()
		for {
			select {
			case _, ok := <-m.immediate:
				if ok {
					m.check()
				}
			case <-ticker.C:
				m.check()
			case <-m.done:
				return
			}
		}
	}()
}

type subscription struct {
	mu     sync.Mutex
	ch     chan StatusUpdate
	closed bool
}

// pushIfNotClosed pushes a status update if the sub isn't closed and returns it's closed status.
func (s *subscription) pushIfNotClosed(su StatusUpdate) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.closed {
		select {
		case s.ch <- su:
		default:
		}
	}

	return s.closed
}

func (s *subscription) unsub() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.closed {
		s.closed = true
		close(s.ch)
	}
}
