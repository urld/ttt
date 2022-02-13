// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ttt

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type TimeTrackingDb struct {
	db       *sql.DB
	filename string
	config   timeTrackingConfig
}

type timeTrackingConfig struct {
	inputResolution time.Duration
	breakThreshold  time.Duration
	breakDeduction  time.Duration
	workingHours    [7]time.Duration
	holidays        string
}

func LoadDb(filename string) (*TimeTrackingDb, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	err = applySchema(db)
	if err != nil {
		return nil, err
	}
	t := TimeTrackingDb{db: db, filename: filename, config: timeTrackingConfig{}}
	t.loadConfig()
	return &t, nil
}

func CreateDb(filename string) (*TimeTrackingDb, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	err = initSchema(db)
	if err != nil {
		return nil, err
	}
	err = applySchema(db)
	if err != nil {
		return nil, err
	}
	t := TimeTrackingDb{db: db, filename: filename, config: timeTrackingConfig{}}
	err = t.loadConfig()
	return &t, err
}

func (t *TimeTrackingDb) loadConfig() error {
	rows, err := t.db.Query("SELECT c.property, c.value FROM config AS c;")
	if err != nil {
		return err
	}
	defer rows.Close()

	var errs []error

	for rows.Next() {
		var property string
		var value string
		err = rows.Scan(&property, &value)
		if err != nil {
			errs = append(errs, err)
		}

		switch property {
		case "input_resolution":
			t.config.inputResolution, errs = parseDuration(value, errs)
		case "break_deduction":
			t.config.breakDeduction, errs = parseDuration(value, errs)
		case "break_threshold":
			t.config.breakThreshold, errs = parseDuration(value, errs)
		case "monday_hours":
			t.config.workingHours[time.Monday], errs = parseDuration(value, errs)
		case "tuesday_hours":
			t.config.workingHours[time.Tuesday], errs = parseDuration(value, errs)
		case "wednesday_hours":
			t.config.workingHours[time.Wednesday], errs = parseDuration(value, errs)
		case "thursday_hours":
			t.config.workingHours[time.Thursday], errs = parseDuration(value, errs)
		case "friday_hours":
			t.config.workingHours[time.Friday], errs = parseDuration(value, errs)
		case "saturdayhours":
			t.config.workingHours[time.Saturday], errs = parseDuration(value, errs)
		case "sunday_hours":
			t.config.workingHours[time.Sunday], errs = parseDuration(value, errs)
		case "holidays":
			t.config.holidays = value
		}
	}
	if len(errs) > 0 {
		return errors.New(fmt.Sprint(errs))
	}
	return nil
}

func parseDuration(value string, errs []error) (time.Duration, []error) {
	value = strings.ReplaceAll(value, " ", "")
	d, err := time.ParseDuration(value)
	if err != nil {
		errs = append(errs, err)
	}
	return d, errs
}

func (t *TimeTrackingDb) Close() error {
	if t.db == nil {
		return nil
	}
	return t.db.Close()
}
