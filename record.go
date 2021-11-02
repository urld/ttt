// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ttt

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ActiveRecordExistsError = errors.New("active record already exists")
	NoActiveRecordError     = errors.New("no active record exists")
)

type Record struct {
	id       int
	Start    *time.Time
	End      *time.Time
	typeId   int
	Type     string
	Duration time.Duration
}

func (r *Record) Active() bool {
	return r.Start != nil && r.End == nil
}

func (t *TimeTrackingDb) StartRecord(dts time.Time) error {
	rec, err := t.GetCurrentRecord()
	if err != nil {
		return err
	}
	if rec.Active() {
		return fmt.Errorf("%w", ActiveRecordExistsError)
	} else {
		dts = t.cleanDts(dts)
		rec.Start = &dts
		_, err := t.db.Exec("INSERT INTO records (type_id, start) VALUES(?, ?);", rec.typeId, rec.Start)
		return err
	}
}

func (t *TimeTrackingDb) EndRecord(dts time.Time) error {
	rec, err := t.GetCurrentRecord()
	if err != nil {
		return err
	}
	if rec.Active() {
		dts = t.cleanDts(dts)
		rec.End = &dts
		_, err := t.db.Exec("UPDATE records SET end=? WHERE id=?;", rec.End, rec.id)
		return err
	} else {
		return fmt.Errorf("%w", NoActiveRecordError)
	}
}

func (t *TimeTrackingDb) GetCurrentRecord() (Record, error) {
	row := t.db.QueryRow("SELECT r.id, r.start, r.end, t.type FROM records AS r INNER JOIN types AS t ON r.type_id = t.id WHERE r.end IS NULL or r.end = '';")
	rec := Record{}
	err := row.Scan(&rec.id, &rec.Start, &rec.End, &rec.Type)
	if err == sql.ErrNoRows {
		return rec, nil
	}
	return rec, err
}

func (t *TimeTrackingDb) cleanDts(dts time.Time) time.Time {
	res, _ := time.ParseDuration(t.config.inputResolution)
	return dts.Round(res)
}
