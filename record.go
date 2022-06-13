// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ttt

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrActiveRecordExists = errors.New("active record already exists")
	ErrNoActiveRecord     = errors.New("no active record exists")
)

type Record struct {
	id    int
	Start *time.Time
	End   *time.Time
	tags  string
}

func (r *Record) Tags() []string {
	return strings.Split(r.tags, ",")
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
		return fmt.Errorf("%w", ErrActiveRecordExists)
	} else {
		rec.Start = &dts
		_, err := t.db.Exec("INSERT INTO records (start) VALUES(?);", rec.Start)
		return err
	}
}

func (t *TimeTrackingDb) EndRecord(dts time.Time) error {
	rec, err := t.GetCurrentRecord()
	if err != nil {
		return err
	}
	if rec.Active() {
		rec.End = &dts
		_, err := t.db.Exec("UPDATE records SET end=? WHERE rowId=?;", rec.End, rec.id)
		return err
	} else {
		return fmt.Errorf("%w", ErrNoActiveRecord)
	}
}

func (t *TimeTrackingDb) GetCurrentRecord() (Record, error) {
	row := t.db.QueryRow("SELECT r.rowId, r.start, r.end FROM records AS r WHERE r.end IS NULL or r.end = '';")
	rec := Record{}
	err := row.Scan(&rec.id, &rec.Start, &rec.End)
	if err == sql.ErrNoRows {
		return rec, nil
	}
	return rec, err
}

func (t *TimeTrackingDb) AddRecord(day time.Time, tags ...string) error {
	//TODO
	return nil
}
