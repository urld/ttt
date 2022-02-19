// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ttt

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func (t *TimeTrackingDb) GetBusinessDays(from, to time.Time) ([]BusinessDay, error) {
	rows, err := t.db.Query("SELECT strftime('%W', r.start) as week_no, date(r.start) AS date, SUM(strftime('%s',r.end)-strftime('%s',r.start)) AS duration FROM records AS r WHERE r.Start BETWEEN ? AND ? AND r.End IS NOT NULL GROUP BY date ORDER BY date;",
		from, to)
	if err != nil {
		return []BusinessDay{}, err
	}
	defer rows.Close()

	var result []BusinessDay
	for rows.Next() {
		var weekNo int
		var date string
		var duration time.Duration
		err = rows.Scan(&weekNo, &date, &duration)
		if err != nil {
			return result, err
		}

		day := BusinessDay{}
		day.ISOWeek = weekNo
		day.Date, _ = time.Parse("2006-01-02", date)
		day.WorkedHours = duration * time.Second
		day.WorkHours = t.config.workingHours[day.Date.Weekday()]
		// TODO: load holiday data
		result = append(result, day)
	}
	return result, nil
}

type BusinessDay struct {
	ISOWeek     int
	Date        time.Time
	WorkedHours time.Duration
	WorkHours   time.Duration
}

func fromDate() time.Time {
	return time.Now().Add(-360 * 24 * time.Hour)
}
