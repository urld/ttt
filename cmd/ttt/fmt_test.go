// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
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

func TestParseTodayTime(t *testing.T) {
	now := time.Now()
	d, _ := parseTime("10:00")
	fmt.Println(d)
	if d.Year() != now.Year() || d.Month() != now.Month() || d.Day() != now.Day() {
		t.Errorf("expected todays date: '%s' got: '%s'", now, d)
	}
	dfmt := d.Format("15:04:05.999999999Z07:00")
	if dfmt != "10:00:00Z" {
		t.Errorf("expected: '10:00:00Z', got: '%s'", dfmt)
	}
}

func TestParseDateTime(t *testing.T) {
	d, _ := parseTime("2022-03-19 10:00")
	dfmt := d.Format(time.RFC3339Nano)
	if dfmt != "2022-03-19T10:00:00Z" {
		t.Errorf("expected: '2022-03-19T10:00:00Z', got: '%s'", dfmt)
	}
}
