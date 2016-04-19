// Copyright 2016 Martin Schnabel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sched

import (
	"sync"
	"time"
)

// Scheduler schedules rules and calls its handler with events.
type Scheduler struct {
	sync.Mutex
	Sched   *Schedule
	timer   *time.Timer
	handler func(Event)
}

// New returns a scheduler that calls handler with events.
func New(now time.Time, handler func(Event)) *Scheduler {
	s := &Scheduler{
		Sched:   NewSchedule(now),
		handler: handler,
	}
	return s
}

// Stop stops the scheduler.
func (s *Scheduler) Stop() {
	s.Lock()
	s.timer.Stop()
	s.timer = nil
	s.Unlock()
}

// Add adds rule r to the schedule.
// Adding a rule will automatically start the scheduler.
func (s *Scheduler) Add(r Rule) bool {
	now := time.Now()
	next := r.Next(now)
	// fail if rule has ended
	if !next.After(now) {
		return false
	}
	s.Lock()
	// fail if the rule was already scheduled
	if s.Sched.has(r) {
		s.Unlock()
		return false
	}
	// if the active rule happens after r
	if s.timer != nil && s.Sched.Last.Time.After(next) {
		// stop the timer
		s.timer.Stop()
		s.timer = nil
		// remove the active rule
		s.Sched.Remove(s.Sched.Last.Rule)
		// reset the last schedule time
		s.Sched.Last.Time = now
		// and re-enqueue the active rule
		s.Sched.Add(s.Sched.Last.Rule)
	}

	s.Sched.Add(r)
	// start the timer if not yet active
	if s.timer == nil {
		s.runNext(now)
	}
	s.Unlock()
	return true
}

// Remove removes rule r from the schedule.
func (s *Scheduler) Remove(r Rule) bool {
	s.Lock()
	ok := s.Sched.Remove(r)
	// if the actrive rule is r
	if s.Sched.Last.Rule == r {
		// stop the timer
		if s.timer != nil {
			s.timer.Stop()
			s.timer = nil
		}
		// and start the next rule
		s.runNext(time.Now())
		ok = true
	}
	s.Unlock()
	return ok
}

// runNext triggers all queued events up to now and starts
// a new timer for the event if necessary.
// runNext must be called with the scheduler lock held.
func (s *Scheduler) runNext(now time.Time) {
	ev := s.Sched.Next()
	for ev.Rule != nil && now.After(ev.Time) {
		s.handler(ev)
		ev = s.Sched.Next()
	}
	if ev.Rule != nil {
		s.timer = time.AfterFunc(ev.Time.Sub(now), s.triggerAndRunNext)
	}
}

// triggerAndRunNext will trigger the active event and call runNext
// to handle the next events
func (s *Scheduler) triggerAndRunNext() {
	s.Lock()
	s.handler(s.Sched.Last)
	s.runNext(time.Now())
	s.Unlock()
}
