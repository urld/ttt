// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
	"time"
)

func TestFmtDuration(t *testing.T) {
	d, _ := time.ParseDuration("1h1m50s")

	str := fmtDuration(d, clock)
	if str != " 01:02" {
		t.Errorf("expected: ' 01:02' got:'%s'", str)
	}

	str = fmtDuration(d, hours)
	if str != " 1h 2m" {
		t.Errorf("expected: ' 1h 2m' got:'%s'", str)
	}

	str = fmtDuration(d, decimal)
	if str != " 1.03" {
		t.Errorf("expected: ' 1.03' got:'%s'", str)
	}
}

func TestFmtDurationNegative(t *testing.T) {
	d, _ := time.ParseDuration("-1h1m50s")

	str := fmtDuration(d, clock)
	if str != "-01:02" {
		t.Errorf("expected: '-01:02' got:'%s'", str)
	}

	str = fmtDuration(d, hours)
	if str != "-1h 2m" {
		t.Errorf("expected: '-1h 2m' got:'%s'", str)
	}

	str = fmtDuration(d, decimal)
	if str != "-1.03" {
		t.Errorf("expected: '-1.03' got:'%s'", str)
	}
}
