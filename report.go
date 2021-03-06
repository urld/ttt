// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ttt

import (
	"context"
	"database/sql"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func (t *TimeTrackingDb) GetRecords(ctx context.Context, from, to time.Time) (<-chan Record, <-chan error, error) {
	rows, err := t.db.Query("SELECT r.start, r.end, r.absence FROM records AS r WHERE r.Start BETWEEN ? AND ? AND r.End IS NOT NULL ORDER BY r.start;",
		from, to)
	if err != nil {
		return nil, nil, err
	}

	out := make(chan Record)
	ec := make(chan error)
	go func() {
		defer rows.Close()
		defer close(out)
		defer close(ec)

		for rows.Next() {
			var absence sql.NullString
			var record Record
			err = rows.Scan(&record.Start, &record.End, &absence)
			if err != nil {
				ec <- err
			}
			if absence.Valid {
				record.Absence = strings.TrimSpace(absence.String)
			}

			select {
			case out <- record:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, ec, nil
}

func (t *TimeTrackingDb) GetBusinessDays(ctx context.Context, records <-chan Record) (<-chan BusinessDay, <-chan error, error) {
	out := make(chan BusinessDay)
	ec := make(chan error)
	go func() {
		defer close(out)
		defer close(ec)
		var day BusinessDay
		var prevEnd time.Time
		for rec := range records {

			if day.Date != *rec.Start {
				// record belongs to new day
				if day != (BusinessDay{}) {
					// previous day is non-empty and can be sent
					select {
					case out <- day:
					case <-ctx.Done():
						return
					}
				}
				// fill values of new day
				day = BusinessDay{}
				day.Date = *rec.Start
				_, day.ISOWeek = rec.Start.ISOWeek()
				day.WorkHours = t.config.workingHours[day.Date.Weekday()]
			}
			// aggregate values
			day.WorkedHours += rec.End.Sub(*rec.Start)
			if rec.Absence != "" {
				// later absence reasons have precedence
				day.Absence = rec.Absence
			}
			if prevEnd != (time.Time{}) {
				day.BreakHours += rec.Start.Sub(prevEnd)
			}
			prevEnd = *rec.End
		}
		if day != (BusinessDay{}) {
			// last one needs to be sent too
			out <- day
		}
	}()
	return out, ec, nil
}

type BusinessDay struct {
	ISOWeek     int
	Date        time.Time
	WorkedHours time.Duration
	WorkHours   time.Duration
	BreakHours  time.Duration
	Absence     string
}

func (t *TimeTrackingDb) EffWorkedHours(day BusinessDay) time.Duration {
	worked := day.WorkedHours
	if day.WorkedHours > t.config.breakThreshold {
		// apply break deduction
		effDeduction := max(0, t.config.breakDeduction-day.BreakHours)
		worked = max(day.WorkedHours-effDeduction, t.config.breakThreshold)
	}
	if day.Absence != "" {
		// apply absence compensation
		worked = max(day.WorkHours, worked)
	}
	return worked
}

func max(x, y time.Duration) time.Duration {
	if x < y {
		return y
	}
	return x
}
