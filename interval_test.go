// Copyright 2016 Martin Schnabel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sched

import (
	"testing"
	"time"
)

func TestInterval(t *testing.T) {
	iv := Interval{Start: parse("2015-01-12 21:54:21"), Duration: 5 * time.Minute}

	cases := []struct {
		after time.Time
		next  time.Time
	}{
		{parse("0600-01-01 00:00:00"), iv.Start},
		{parse("2000-09-12 12:00:00"), iv.Start},
		{iv.Start, parse("2015-01-12 21:59:21")},
		{parse("2015-06-30 23:59:21"), parse("2015-07-01 00:04:21")},
		{parse("2015-09-12 22:54:20"), parse("2015-09-12 22:54:21")},
		{parse("2015-09-13 21:54:20"), parse("2015-09-13 21:54:21")},
		{parse("2015-10-12 21:54:20"), parse("2015-10-12 21:54:21")},
		{parse("2015-11-12 21:54:20"), parse("2015-11-12 21:54:21")},
		{parse("2016-11-07 00:51:00"), parse("2016-11-07 00:54:21")},
		{parse("2016-11-07 03:51:20"), parse("2016-11-07 03:54:21")},
		{parse("2016-11-07 06:51:40"), parse("2016-11-07 06:54:21")},
		{parse("2016-11-07 09:52:00"), parse("2016-11-07 09:54:21")},
		{parse("2016-11-07 12:52:20"), parse("2016-11-07 12:54:21")},
		{parse("2016-11-07 15:53:40"), parse("2016-11-07 15:54:21")},
		{parse("2016-11-07 18:54:00"), parse("2016-11-07 18:54:21")},
		{parse("2016-11-07 21:54:20"), parse("2016-11-07 21:54:21")},
		{parse("2116-11-07 21:54:20"), parse("2116-11-07 21:54:21")},
		{parse("2216-11-07 21:54:20"), parse("2216-11-07 21:54:21")},
		// following cases would tickle a bug if not for the start adjustment
		{parse("2316-11-07 21:54:20"), parse("2316-11-07 21:54:21")},
		{parse("9316-11-07 21:54:20"), parse("9316-11-07 21:54:21")},
	}

	for i, test := range cases {
		result := iv.Next(test.after)
		if !test.next.Equal(result) {
			t.Errorf("case %d want %s got %s", i, test.next, result)
		}
	}
}
