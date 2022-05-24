// Copyright 2016 Martin Schnabel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sched_test

import (
	"fmt"
	"github.com/mb0/sched"
	"sync/atomic"
	"time"
)

func Example_schedScheduler() {
	type action struct {
		sched.Interval
		Do func(time.Time)
	}
	var (
		now   = time.Now()
		done  = make(chan struct{})
		count int32
	)
	a := &action{
		sched.Interval{Start: now, Duration: 10 * time.Millisecond},
		func(t time.Time) {
			fmt.Println("action triggered after", t.Sub(now))
		},
	}
	trigger := func(ev sched.Event) {
		switch a := ev.Rule.(type) {
		case *action:
			a.Do(ev.Time)
		default:
			fmt.Println("unknown rule", a)
		}
		if atomic.AddInt32(&count, 1) >= 3 {
			close(done)
		}
	}
	s := sched.New(now, trigger)
	// equivalent to s.Sched.Add(a); s.Start()
	s.Add(a)
	<-done
	s.Stop()
	// Output:
	// action triggered after 10ms
	// action triggered after 20ms
	// action triggered after 30ms
}
