// Copyright 2016 Martin Schnabel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sched is used to find the next time event from multiple rules.
package sched

import (
	"container/heap"
	"time"
)

// Rule represents a deterministic sequence of points in time.
// Implementations must be comparable.
type Rule interface {
	// Next returns the next point in time after now.
	// If the returned time is not after now this rule has ended.
	Next(now time.Time) time.Time
}

// Schedule is a time based queue of events sourced from multiple rules.
type Schedule struct {
	Last Event
	queue
}

// NewSchedule returns an empty Schedule beginning at now.
func NewSchedule(now time.Time) *Schedule {
	return &Schedule{Last: Event{Time: now}}
}

// has returns whether rule r is scheduled
func (s *Schedule) has(r Rule) bool {
	for _, ev := range s.queue {
		if ev.Rule == r {
			return true
		}
	}
	return false
}

// Add schedules a rule or returns false.
func (s *Schedule) Add(r Rule) bool {
	// fail if the rule was already scheduled
	if s.has(r) {
		return false
	}
	// next occurance
	next := r.Next(s.Last.Time)
	// fail if rule has ended
	if !next.After(s.Last.Time) {
		return false
	}
	// add the next occurance to the queue
	heap.Push(&s.queue, Event{r, next})
	return true
}

// Remove removes a rule from the schedule or returns false.
func (s *Schedule) Remove(r Rule) bool {
	for i, ev := range s.queue {
		if ev.Rule == r {
			heap.Remove(&s.queue, i)
			return true
		}
	}
	return false
}

// Next returns the next event. The event rule is nil if no event is scheduled.
// The returned event rule is rescheduled or removed from the queue.
func (s *Schedule) Next() Event {
	if len(s.queue) == 0 {
		return Event{}
	}
	s.Last = s.queue[0]
	nt := s.Last.Next(s.Last.Time)
	if nt.After(s.Last.Time) {
		s.queue[0].Time = nt
		heap.Fix(&s.queue, 0)
	} else {
		heap.Pop(&s.queue)
	}
	return s.Last
}

// Event is an occurance of Rule at a specific point in time.
type Event struct {
	Rule
	time.Time
}

// queue is a slice of events implementing heap.Interface
type queue []Event

func (q queue) Len() int           { return len(q) }
func (q queue) Less(i, j int) bool { return q[i].Before(q[j].Time) }
func (q queue) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func (q *queue) Push(x interface{}) {
	*q = append(*q, x.(Event))
}

func (q *queue) Pop() interface{} {
	old := *q
	last := old[len(old)-1]
	*q = old[0 : len(old)-1]
	return last
}
