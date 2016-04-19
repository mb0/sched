// Copyright 2016 Martin Schnabel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sched

import "time"

// Next returns the next time after now that is intervals more than start.
// It can be used to implement common repeating rules with a fixed interval.
func Next(start time.Time, interval time.Duration, now time.Time) time.Time {
	if start.After(now) {
		return start
	}
	// the largest representable duration is approximately 290 years
	// so lets adjust our start time to calculate representable one later
	if ydiff := now.Year() - start.Year(); ydiff > 250 {
		start = start.AddDate(ydiff-ydiff%100, 0, 0)
	}
	add := now.Sub(start) % interval
	return now.Add(interval - add)
}

// Interval is simple Rule implementation with start time and interval duration.
type Interval struct {
	Start time.Time
	time.Duration
}

// Next returns the next point in time for this rule after now.
func (i Interval) Next(now time.Time) time.Time {
	return Next(i.Start, i.Duration, now)
}
