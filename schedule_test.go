// Copyright 2016 Martin Schnabel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sched

import (
	"testing"
	"time"
)

type birthday struct {
	time.Time
}

func (b birthday) Next(after time.Time) time.Time {
	if b.After(after) {
		return b.Time
	}
	t := b.Time.AddDate(after.Year()-b.Year(), 0, 0)
	if !t.After(after) {
		t = t.AddDate(1, 0, 0)
	}
	return t
}

func TestSched(t *testing.T) {
	r1 := birthday{parse("1983-11-07")}
	r2 := birthday{parse("1986-11-07")}
	var s Schedule
	s.Add(r1)
	s.Add(r2)
	expect := []Event{
		{r1, parse("1983-11-07")},
		{r1, parse("1984-11-07")},
		{r1, parse("1985-11-07")},
		{r1, parse("1986-11-07")},
		{r2, parse("1986-11-07")},
		{r2, parse("1987-11-07")},
		{r1, parse("1987-11-07")},
		{r1, parse("1988-11-07")},
		{r2, parse("1988-11-07")},
		{r2, parse("1989-11-07")},
		{r1, parse("1989-11-07")},
		{r1, parse("1990-11-07")},
		{r2, parse("1990-11-07")},
	}
	events := take(len(expect), &s)
	if len(events) != len(expect) {
		t.Fatalf("want %d events got %d", len(expect), len(events))
	}
	for i, ev := range events {
		ex := expect[i]
		if ev.Rule != ex.Rule {
			t.Errorf("expect event %d rule to be %s got %s", i, ex.Rule, ev.Rule)
		}
		if !ev.Time.Equal(ex.Time) {
			t.Errorf("expect event %d rule at %s at %s", i, ex.Time, ev.Time)
		}
	}
}

func take(n int, src *Schedule) []Event {
	res := make([]Event, 0, n)
	for i := 0; i < n; i++ {
		ev := src.Next()
		if ev.Rule == nil {
			break
		}
		res = append(res, ev)
	}
	return res
}

func parse(s string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		t, err = time.Parse("2006-01-02", s)
		if err != nil {
			panic("unknown time format " + s)
		}
	}
	return t
}
