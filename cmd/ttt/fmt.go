// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"
	"time"
)

func fmtDate(d time.Time) string {
	return d.Format("2006-01-02 Mon")
}

func fmtTime(d time.Time) string {
	return d.Format("Mon 2006-01-02 15:04")
}

func fmtDurationTrim(d time.Duration, style durationFmt) string {
	return strings.TrimSpace(fmtDuration(d, style))
}

func fmtAbsence(a string) string {
	if a != "" {
		return "x"
	}
	return a
}

func fmtDuration(d time.Duration, style durationFmt) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	m := (d - h*time.Hour) / time.Minute
	if m < 0 {
		m = -m
	}
	switch style {
	case clock:
		return fmt.Sprintf("% 03d:%02d", h, m)
	case hours:
		return fmt.Sprintf("% dh %dm", h, m)
	case decimal:
		return fmt.Sprintf("% .2f", float64(d)/float64(time.Hour))
	}
	return d.String()
}

func parseTime(arg string) (time.Time, error) {
	var t time.Time
	var err error
	if len(arg) > 5 {
		t, err = time.Parse("2006-01-02 15:04", arg)
	} else {
		today := time.Now()
		t, err = time.Parse("15:04", arg)
		t = t.AddDate(today.Year(), int(today.Month())-1, today.Day()-1)
	}
	return t, err
}
