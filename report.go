// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ttt

import (
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func (t *TimeTrackingDb) GetRecords() ([]Record, error) {
	rows, err := t.db.Query("SELECT r.id, r.start, r.end, r.type_id, t.type FROM records AS r INNER JOIN types AS t ON t.id = r.type_id;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []Record
	for rows.Next() {
		var rec Record
		err = rows.Scan(&rec.id, &rec.Start, &rec.End, &rec.typeId, &rec.Type)
		if err != nil {
			return records, err
		}
		if rec.End != nil {
			rec.Duration = rec.End.Sub(*rec.Start)
			records = append(records, rec)
		}
	}
	return records, nil
}

func (t *TimeTrackingDb) GetWeekAggr() error {
	rows, err := t.db.Query("SELECT strftime('%W', r.start) AS week, SUM(strftime('%s',r.end)-strftime('%s',r.start)) AS duration FROM records AS r WHERE r.Start BETWEEN ? AND ? GROUP BY week ORDER BY week;",
		fromDate(), time.Now())
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var week int
		var duration time.Duration
		err = rows.Scan(&week, &duration)
		if err != nil {
			return err
		}
		fmt.Println(week, duration*time.Second)
	}
	return nil
}

func (t *TimeTrackingDb) GetDayAggr() error {
	rows, err := t.db.Query("SELECT strftime('%W',r.start) as week, date(r.start) AS day, SUM(strftime('%s',r.end)-strftime('%s',r.start)) AS duration FROM records AS r WHERE r.Start BETWEEN ? AND ? GROUP BY day ORDER BY day;",
		fromDate(), time.Now())
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var week int
		var day string
		var duration time.Duration
		err = rows.Scan(&week, &day, &duration)
		if err != nil {
			return err
		}
		fmt.Println(week, day, duration*time.Second)
	}
	return nil
}

func fromDate() time.Time {
	return time.Now().Add(-360 * 24 * time.Hour)
}
