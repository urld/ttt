// Copyright (c) 2021, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ttt

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type TimeTrackingDb struct {
	db       *sql.DB
	filename string
	config   timeTrackingConfig
}
type timeTrackingConfig struct {
	inputResolution string
}

func LoadDb(filename string) (TimeTrackingDb, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return TimeTrackingDb{}, err
	}
	err = applySchema(db)
	if err != nil {
		return TimeTrackingDb{}, err
	}
	t := TimeTrackingDb{db, filename, timeTrackingConfig{}}
	t.loadConfig()
	return t, nil
}

func CreateDb(filename string) (TimeTrackingDb, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return TimeTrackingDb{}, err
	}
	err = initSchema(db)
	if err != nil {
		return TimeTrackingDb{}, err
	}
	err = applySchema(db)
	if err != nil {
		return TimeTrackingDb{}, err
	}
	t := TimeTrackingDb{db, filename, timeTrackingConfig{}}
	t.loadConfig()
	return t, nil
}

func (t *TimeTrackingDb) loadConfig() error {
	rows, err := t.db.Query("SELECT c.property, c.value FROM config AS c;")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var property string
		var value string
		err = rows.Scan(&property, &value)
		if err != nil {
			return err
		}
		switch property {
		case "input_resolution":
			t.config.inputResolution = value
		}
	}
	return nil
}

func (t *TimeTrackingDb) Close() error {
	if t.db == nil {
		return nil
	}
	return t.db.Close()
}
